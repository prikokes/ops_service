package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ops_service/internal/model"
	"time"
)

// PickupPointRepo is a PostgreSQL implementation of PickupPointRepository
type PickupPointRepo struct {
	db *sql.DB
}

// NewPickupPointRepo creates a new PickupPointRepo
func NewPickupPointRepo(db *sql.DB) *PickupPointRepo {
	return &PickupPointRepo{db: db}
}

// Create creates a new pickup point in the database
func (r *PickupPointRepo) Create(ctx context.Context, point *model.PickupPoint) error {
	query := `
		INSERT INTO pickup_points (registered_at, city)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		point.RegisteredAt,
		point.City,
	).Scan(&point.ID)

	return err
}

// GetByID returns a pickup point by ID
func (r *PickupPointRepo) GetByID(ctx context.Context, id int) (*model.PickupPoint, error) {
	query := `
		SELECT id, registered_at, city
		FROM pickup_points
		WHERE id = $1
	`

	point := &model.PickupPoint{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&point.ID,
		&point.RegisteredAt,
		&point.City,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("pickup point not found")
		}
		return nil, err
	}

	return point, nil
}

// List returns a list of pickup points with pagination and date filter
func (r *PickupPointRepo) List(ctx context.Context, page, pageSize int, startDate, endDate *time.Time) ([]*model.PickupPoint, int, error) {
	// First, count total number of pickup points with the given filter
	countQuery := `
		SELECT COUNT(DISTINCT pp.id)
		FROM pickup_points pp
	`
	
	// Add filter if dates are provided
	var args []interface{}
	var whereClause string
	if startDate != nil && endDate != nil {
		whereClause = `
			WHERE EXISTS (
				SELECT 1 FROM receipts r
				WHERE r.pickup_point_id = pp.id
				AND r.created_at BETWEEN $1 AND $2
			)
		`
		args = append(args, startDate, endDate)
	} else if startDate != nil {
		whereClause = `
			WHERE EXISTS (
				SELECT 1 FROM receipts r
				WHERE r.pickup_point_id = pp.id
				AND r.created_at >= $1
			)
		`
		args = append(args, startDate)
	} else if endDate != nil {
		whereClause = `
			WHERE EXISTS (
				SELECT 1 FROM receipts r
				WHERE r.pickup_point_id = pp.id
				AND r.created_at <= $1
			)
		`
		args = append(args, endDate)
	}
	
	countQuery += whereClause
	
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pickup points: %w", err)
	}
	
	// Now get the paginated pickup points
	query := `
		SELECT pp.id, pp.registered_at, pp.city
		FROM pickup_points pp
	` + whereClause + `
		ORDER BY pp.id
		LIMIT $%d OFFSET $%d
	`
	
	// Add pagination parameters
	offset := (page - 1) * pageSize
	argPos := len(args) + 1
	formattedQuery := fmt.Sprintf(query, argPos, argPos+1)
	args = append(args, pageSize, offset)
	
	rows, err := r.db.QueryContext(ctx, formattedQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query pickup points: %w", err)
	}
	defer rows.Close()
	
	var pickupPoints []*model.PickupPoint
	for rows.Next() {
		point := &model.PickupPoint{}
		if err := rows.Scan(&point.ID, &point.RegisteredAt, &point.City); err != nil {
			return nil, 0, fmt.Errorf("failed to scan pickup point: %w", err)
		}
		pickupPoints = append(pickupPoints, point)
	}
	
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over pickup points: %w", err)
	}
	
	return pickupPoints, total, nil
} 