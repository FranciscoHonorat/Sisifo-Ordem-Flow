package valueobject_test

import (
	"testing"

	order "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
	"github.com/google/uuid"
)

func TestOrderItem(t *testing.T) {
	t.Run("Test for NewOrderItem", func(t *testing.T) {
		tests := []struct {
			name          string
			productID     order.ProductID
			unitPrice     order.Money
			quantity      order.Quantity
			expectedError error
		}{
			{"Valid OrderItem", order.NewProductIDMust(uuid.New()), order.NewMoneyMust(100, "USD"), order.NewQuantityMust(2), nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				p, err := order.NewOrderItem(tt.productID, tt.unitPrice, tt.quantity)
				if err != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				}

				if err == nil {
					if p.ProductID() != tt.productID {
						t.Errorf("Expected ProductID: %v, got: %v", tt.productID, p.ProductID())
					}
					if p.UnitPrice() != tt.unitPrice {
						t.Errorf("Expected UnitPrice: %v, got: %v", tt.unitPrice, p.UnitPrice())
					}
					if p.Quantity() != tt.quantity {
						t.Errorf("Expected Quantity: %v, got: %v", tt.quantity, p.Quantity())
					}
				}
			})
		}
	})

	t.Run("Test for subTotal", func(t *testing.T) {
		tests := []struct {
			name          string
			productID     order.ProductID
			unitPrice     order.Money
			quantity      order.Quantity
			expectedTotal order.Money
		}{
			{"Valid OrderItem", order.NewProductIDMust(uuid.New()), order.NewMoneyMust(100, "USD"), order.NewQuantityMust(2), order.NewMoneyMust(200, "USD")},
			{"Valid OrderItem with different currency", order.NewProductIDMust(uuid.New()), order.NewMoneyMust(50, "BRL"), order.NewQuantityMust(3), order.NewMoneyMust(150, "BRL")},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				p, err := order.NewOrderItem(tt.productID, tt.unitPrice, tt.quantity)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				total, err := p.SubTotal()
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if total != tt.expectedTotal {
					t.Errorf("Expected SubTotal: %v, got: %v", tt.expectedTotal, total)
				}
			})
		}
	})

	t.Run("Test for MarshalJSON", func(t *testing.T) {
		pid := order.NewProductIDMust(uuid.New())
		tests := []struct {
			name         string
			productID    order.ProductID
			unitPrice    order.Money
			quantity     order.Quantity
			expectedJSON string
		}{
			{"Valid OrderItem", pid, order.NewMoneyMust(100, "USD"), order.NewQuantityMust(2), `{"product_id":"` + pid.String() + `","unit_price":{"amount":100,"currency":"USD"},"quantity":2}`},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				p, err := order.NewOrderItem(tt.productID, tt.unitPrice, tt.quantity)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				jsonData, err := p.MarshalJSON()
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if string(jsonData) != tt.expectedJSON {
					t.Errorf("Expected JSON: %v, got: %v", tt.expectedJSON, string(jsonData))
				}
			})
		}
	})

	t.Run("Test for UnmarshalJSON", func(t *testing.T) {
		pid := order.NewProductIDMust(uuid.New())
		tests := []struct {
			name         string
			jsonData     string
			expectedItem order.OrderItem
			expectError  bool
		}{
			{
				"Valid OrderItem",
				`{"product_id":"` + pid.String() + `","unit_price":{"amount":100,"currency":"USD"},"quantity":2}`,
				order.NewOrderItemMust(pid, order.NewMoneyMust(100, "USD"), order.NewQuantityMust(2)),
				false,
			},
			{"Invalid JSON", `{"product_id":"invalid-uuid",...}`, order.OrderItem{}, true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var item order.OrderItem
				err := item.UnmarshalJSON([]byte(tt.jsonData))
				if (err != nil) != tt.expectError {
					t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
				}
				if !tt.expectError && (item.ProductID() != tt.expectedItem.ProductID() || item.UnitPrice() != tt.expectedItem.UnitPrice() || item.Quantity() != tt.expectedItem.Quantity()) {
					t.Errorf("Expected OrderItem: %+v, got: %+v", tt.expectedItem, &item)
				}
			})
		}
	})
}
