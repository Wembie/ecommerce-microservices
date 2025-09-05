package models_test

import (
	"testing"

	"ecommerce.orders.manager/internal/models"
	"ecommerce.orders.manager/internal/utils"

	"github.com/stretchr/testify/assert"
	
)

func TestNewPaginatedResponse(t *testing.T) {
	tests := []struct {
		name          string
		items         []string
		page          int
		size          int
		total         int
		expectedPages int
		expectNext    *int
		expectPrev    *int
	}{
		{
			name:          "first page with next",
			items:         []string{"a", "b"},
			page:          0,
			size:          2,
			total:         5,
			expectedPages: 3,
			expectNext:    utils.IntPtr(1),
			expectPrev:    nil,
		},
		{
			name:          "middle page with prev and next",
			items:         []string{"c", "d"},
			page:          1,
			size:          2,
			total:         5,
			expectedPages: 3,
			expectNext:    utils.IntPtr(2),
			expectPrev:    utils.IntPtr(0),
		},
		{
			name:          "last page with prev only",
			items:         []string{"e"},
			page:          2,
			size:          2,
			total:         5,
			expectedPages: 3,
			expectNext:    nil,
			expectPrev:    utils.IntPtr(1),
		},
		{
			name:          "no items still has one page",
			items:         []string{},
			page:          0,
			size:          10,
			total:         0,
			expectedPages: 1,
			expectNext:    nil,
			expectPrev:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := models.NewPaginatedResponse(tt.items, tt.page, tt.size, tt.total)

			assert.Equal(t, tt.items, resp.Items)
			assert.Equal(t, tt.page, resp.Page)
			assert.Equal(t, tt.size, resp.Size)
			assert.Equal(t, tt.total, resp.Total)
			assert.Equal(t, tt.expectedPages, resp.Pages)

			if tt.expectNext != nil {
				assert.NotNil(t, resp.NextPage)
				assert.Equal(t, *tt.expectNext, *resp.NextPage)
			} else {
				assert.Nil(t, resp.NextPage)
			}

			if tt.expectPrev != nil {
				assert.NotNil(t, resp.PreviousPage)
				assert.Equal(t, *tt.expectPrev, *resp.PreviousPage)
			} else {
				assert.Nil(t, resp.PreviousPage)
			}
		})
	}
}