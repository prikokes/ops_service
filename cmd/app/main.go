package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"ops_service/internal/api"
	"ops_service/internal/auth"
	"ops_service/internal/config"
	"ops_service/internal/db"
	grpcserver "ops_service/internal/grpc"
	"ops_service/internal/metrics"
	"ops_service/internal/middleware"
	"ops_service/internal/model"
	"ops_service/internal/repository"
	"ops_service/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize db: %s", err.Error())
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Failed to run migrations: %s", err.Error())
	}

	repos := repository.NewRepository(database)

	appMetrics := metrics.NewMetrics("ops_service")
	prometheusHandler := metrics.NewPrometheusHandler(appMetrics, cfg.Prometheus)

	tokenService := auth.NewTokenService(cfg.JWT.Secret, cfg.JWT.TTL)
	userService := service.NewUserService(repos.User)
	pickupPointService := service.NewPickupPointService(repos.PickupPoint)
	receiptService := service.NewReceiptService(repos.Receipt, repos.PickupPoint, repos.Product)

	router := gin.Default()
	router.Use(middleware.MetricsMiddleware(appMetrics))

	authHandler := api.NewAuthHandler(userService, tokenService)
	pickupPointHandler := api.NewPickupPointHandler(pickupPointService)
	receiptHandler := api.NewReceiptHandler(receiptService)

	authMiddleware := middleware.AuthMiddleware(tokenService)
	moderatorMiddleware := middleware.RoleMiddleware(model.RoleModerator)
	clientMiddleware := middleware.RoleMiddleware(model.RoleClient)

	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.GET("/dummyLogin", authHandler.DummyLogin)

	pickupPointRoutes := router.Group("/pickup-points")
	{
		pickupPointRoutes.Use(authMiddleware)

		pickupPointRoutes.POST("", moderatorMiddleware, func(c *gin.Context) {
			pickupPointHandler.Create(c)
			if c.Writer.Status() == http.StatusCreated {
				appMetrics.IncPickupPointsTotal()
			}
		})

		pickupPointRoutes.GET("/:id", pickupPointHandler.Get)

		pickupPointRoutes.GET("", pickupPointHandler.List)
	}

	receiptRoutes := router.Group("/receipts")
	{
		receiptRoutes.Use(authMiddleware)

		// Create receipt (client only)
		receiptRoutes.POST("", clientMiddleware, func(c *gin.Context) {
			receiptHandler.Create(c)
			if c.Writer.Status() == http.StatusCreated {
				appMetrics.IncReceiptsTotal()
			}
		})

		receiptRoutes.GET("/:id", receiptHandler.Get)

		receiptRoutes.POST("/:id/close", clientMiddleware, receiptHandler.Close)

		productRoutes := receiptRoutes.Group("/products")
		{
			productRoutes.POST("", clientMiddleware, func(c *gin.Context) {
				receiptHandler.AddProduct(c)
				if c.Writer.Status() == http.StatusCreated {
					appMetrics.IncProductsTotal()
				}
			})

			productRoutes.DELETE("", clientMiddleware, receiptHandler.DeleteLastProduct)
		}
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	grpcServer := grpcserver.NewServer(repos.PickupPoint)

	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %s", err)
		}
	}()

	go func() {
		log.Printf("Starting gRPC server on port %d", cfg.GRPC.Port)
		if err := grpcServer.Run(cfg.GRPC); err != nil {
			log.Fatalf("Failed to start gRPC server: %s", err)
		}
	}()

	go func() {
		log.Printf("Starting Prometheus server on port %d", cfg.Prometheus.Port)
		if err := prometheusHandler.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start Prometheus server: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server forced to shutdown: %s", err)
	}

	if err := prometheusHandler.Shutdown(ctx); err != nil {
		log.Fatalf("Prometheus server forced to shutdown: %s", err)
	}

	grpcServer.Shutdown()

	log.Println("Servers exited")
}
