package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"ecommerce.products.manager/internal/models"
)

func TestProduct_JSONSerialization(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	desc := "sample description"

	product := models.Product{
		ID:          id,
		Name:        "Laptop",
		Description: &desc,
		Price:       999.99,
		Stock:       10,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	data, err := json.Marshal(product)
	assert.NoError(t, err)

	var decoded models.Product
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, product.ID, decoded.ID)
	assert.Equal(t, product.Name, decoded.Name)
	assert.Equal(t, *product.Description, *decoded.Description)
	assert.Equal(t, product.Price, decoded.Price)
	assert.Equal(t, product.Stock, decoded.Stock)
}

func TestCreateProductRequest_JSON(t *testing.T) {
	desc := "gaming laptop"
	req := models.CreateProductRequest{
		Name:        "Laptop",
		Description: &desc,
		Price:       1500.50,
		Stock:       5,
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var decoded models.CreateProductRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, req.Name, decoded.Name)
	assert.Equal(t, *req.Description, *decoded.Description)
	assert.Equal(t, req.Price, decoded.Price)
	assert.Equal(t, req.Stock, decoded.Stock)
}

func TestUpdateProductRequest_OptionalFields(t *testing.T) {
	id := uuid.New()
	name := "Updated Laptop"
	price := 1200.75
	stock := 20

	req := models.UpdateProductRequest{
		ID:        id,
		Name:      &name,
		Price:     &price,
		Stock:     &stock,
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, id, req.ID)
	assert.NotNil(t, req.Name)
	assert.Equal(t, "Updated Laptop", *req.Name)
	assert.NotNil(t, req.Price)
	assert.Equal(t, 1200.75, *req.Price)
	assert.NotNil(t, req.Stock)
	assert.Equal(t, 20, *req.Stock)
}

func TestDeleteProductResponse_JSON(t *testing.T) {
	resp := models.DeleteProductResponse{
		Success: true,
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded models.DeleteProductResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.True(t, decoded.Success)
}

func TestUpdateStockRequest_JSON(t *testing.T) {
	id := uuid.New()
	req := models.UpdateStockRequest{
		ID:    id,
		Stock: 50,
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, `"stock":50`)
	assert.NotContains(t, jsonStr, `"id"`)

	var decoded models.UpdateStockRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, 50, decoded.Stock)

	assert.Equal(t, uuid.Nil, decoded.ID)
}

func TestSearchProductsRequest_JSON(t *testing.T) {
	name := "Phone"
	stock := 15

	req := models.SearchProductsRequest{
		Name:  &name,
		Stock: &stock,
		Page:  2,
		Size:  20,
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var decoded models.SearchProductsRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, *req.Name, *decoded.Name)
	assert.Equal(t, *req.Stock, *decoded.Stock)
	assert.Equal(t, req.Page, decoded.Page)
	assert.Equal(t, req.Size, decoded.Size)
}
