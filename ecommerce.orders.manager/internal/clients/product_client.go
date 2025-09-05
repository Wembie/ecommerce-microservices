package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ecommerce.orders.manager/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProductClient interface {
	GetProduct(ctx context.Context, productID uuid.UUID) (*models.ProductInfo, error)
	UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error
}

type productClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewProductClient(baseURL string, logger *zap.Logger) ProductClient {
	return &productClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

func (c *productClient) GetProduct(ctx context.Context, productID uuid.UUID) (*models.ProductInfo, error) {
	c.logger.Info("Getting product info", zap.String("product_id", productID.String()))

	url := fmt.Sprintf("%s/products/%s", c.baseURL, productID.String())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.logger.Error("Failed to create request", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to make request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Unexpected status code", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("product service returned status %d", resp.StatusCode)
	}

	var product models.ProductInfo
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		c.logger.Error("Failed to decode response", zap.Error(err))
		return nil, err
	}

	return &product, nil
}

func (c *productClient) UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error {
	c.logger.Info("Updating product stock", zap.String("product_id", productID.String()), zap.Int("stock", stock))

	url := fmt.Sprintf("%s/products/%s/stock", c.baseURL, productID.String())

	updateReq := map[string]int{"stock": stock}
	jsonData, err := json.Marshal(updateReq)
	if err != nil {
		c.logger.Error("Failed to marshal request", zap.Error(err))
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create request", zap.Error(err))
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to make request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Unexpected status code", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("product service returned status %d", resp.StatusCode)
	}

	return nil
}