package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"ops_service/internal/api"
	"ops_service/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// IntegrationTest performs a full flow test:
// 1. Create a pickup point
// 2. Create a receipt
// 3. Add 50 products
// 4. Close the receipt
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	baseURL := "http://localhost:8080"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	moderatorToken, err := getToken(client, baseURL, "moderator")
	require.NoError(t, err, "Failed to get moderator token")

	pickupPoint, err := createPickupPoint(client, baseURL, moderatorToken)
	require.NoError(t, err, "Failed to create pickup point")
	assert.NotZero(t, pickupPoint.ID, "Pickup point ID should not be zero")
	assert.Equal(t, model.CityMoscow, pickupPoint.City, "Pickup point city should be Moscow")

	clientToken, err := getToken(client, baseURL, "client")
	require.NoError(t, err, "Failed to get client token")

	receipt, err := createReceipt(client, baseURL, clientToken, pickupPoint.ID)
	require.NoError(t, err, "Failed to create receipt")
	assert.NotZero(t, receipt.ID, "Receipt ID should not be zero")
	assert.Equal(t, pickupPoint.ID, receipt.PickupPointID, "Receipt pickup point ID should match")
	assert.Equal(t, model.ReceiptStatusInProgress, receipt.Status, "Receipt status should be in progress")

	productTypes := []model.ProductType{
		model.ProductTypeElectronics,
		model.ProductTypeClothing,
		model.ProductTypeFootwear,
	}

	for i := 0; i < 50; i++ {
		productType := productTypes[i%len(productTypes)]
		product, err := addProduct(client, baseURL, clientToken, pickupPoint.ID, productType)
		require.NoError(t, err, "Failed to add product #%d", i+1)
		assert.NotZero(t, product.ID, "Product ID should not be zero")
		assert.Equal(t, productType, product.Type, "Product type should match")
		assert.Equal(t, receipt.ID, product.ReceiptID, "Product receipt ID should match")
	}

	err = closeReceipt(client, baseURL, clientToken, receipt.ID)
	require.NoError(t, err, "Failed to close receipt")

	closedReceipt, err := getReceipt(client, baseURL, clientToken, receipt.ID)
	require.NoError(t, err, "Failed to get closed receipt")
	assert.Equal(t, model.ReceiptStatusClosed, closedReceipt.Status, "Receipt status should be closed")
	assert.Len(t, closedReceipt.Products, 50, "Receipt should have 50 products")
}

func getToken(client *http.Client, baseURL, role string) (string, error) {
	url := fmt.Sprintf("%s/dummyLogin?role=%s", baseURL, role)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response api.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Token, nil
}

func createPickupPoint(client *http.Client, baseURL, token string) (*model.PickupPoint, error) {
	url := fmt.Sprintf("%s/pickup-points", baseURL)

	requestBody, err := json.Marshal(api.CreatePickupPointRequest{
		City: model.CityMoscow,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var pickupPoint model.PickupPoint
	if err := json.NewDecoder(resp.Body).Decode(&pickupPoint); err != nil {
		return nil, err
	}

	return &pickupPoint, nil
}

func createReceipt(client *http.Client, baseURL, token string, pickupPointID int) (*model.Receipt, error) {
	url := fmt.Sprintf("%s/receipts", baseURL)

	requestBody, err := json.Marshal(api.CreateReceiptRequest{
		PickupPointID: pickupPointID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var receipt model.Receipt
	if err := json.NewDecoder(resp.Body).Decode(&receipt); err != nil {
		return nil, err
	}

	return &receipt, nil
}

func addProduct(client *http.Client, baseURL, token string, pickupPointID int, productType model.ProductType) (*model.Product, error) {
	url := fmt.Sprintf("%s/receipts/products?pickup_point_id=%d", baseURL, pickupPointID)

	requestBody, err := json.Marshal(api.AddProductRequest{
		Type: productType,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var product model.Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func closeReceipt(client *http.Client, baseURL, token string, receiptID int) error {
	url := fmt.Sprintf("%s/receipts/%d/close", baseURL, receiptID)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func getReceipt(client *http.Client, baseURL, token string, receiptID int) (*model.Receipt, error) {
	url := fmt.Sprintf("%s/receipts/%d", baseURL, receiptID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var receipt model.Receipt
	if err := json.NewDecoder(resp.Body).Decode(&receipt); err != nil {
		return nil, err
	}

	return &receipt, nil
}
