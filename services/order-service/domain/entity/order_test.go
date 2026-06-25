package entity_test

import (
	"errors"
	"testing"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/domainErrors"
	orderEntity "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/entity"
	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
	"github.com/google/uuid"
)

func TestOrder(t *testing.T) {
	t.Run("Test for NewOrder", func(t *testing.T) {
		tests := []struct {
			name        string
			id          uuid.UUID
			customerID  uuid.UUID
			expectedErr error
		}{
			{"Valid Order", uuid.New(), uuid.New(), nil},
			{"Invalid Order ID", uuid.Nil, uuid.New(), domainErrors.ErrInvalidOrderID},
			{"Invalid Customer ID", uuid.New(), uuid.Nil, domainErrors.ErrInvalidCustomerID},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				orderID, err := valueobject.NewOrderID(tt.id)
				if err != nil {
					if !errors.Is(err, tt.expectedErr) {
						t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
					}
					return
				}
				customerID, err := valueobject.NewCustomerID(tt.customerID)
				if err != nil {
					if !errors.Is(err, tt.expectedErr) {
						t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
					}
					return
				}
				_, err = orderEntity.NewOrder(orderID, customerID)
				if err != tt.expectedErr {
					t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
				}
			})
		}
	})

	t.Run("Test for NewOrderMust", func(t *testing.T) {
		tests := []struct {
			name        string
			id          uuid.UUID
			customerID  uuid.UUID
			expectPanic bool
		}{
			{"Valid Order", uuid.New(), uuid.New(), false},
			{"Invalid Order ID", uuid.Nil, uuid.New(), true},
			{"Invalid Customer ID", uuid.New(), uuid.Nil, true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				orderID, err := valueobject.NewOrderID(tt.id)
				if err != nil && !tt.expectPanic {
					t.Errorf("Unexpected error creating order ID: %v", err)
				}
				customerID, err := valueobject.NewCustomerID(tt.customerID)
				if err != nil && !tt.expectPanic {
					t.Errorf("Unexpected error creating customer ID: %v", err)
				}

				defer func() {
					if r := recover(); r != nil {
						if !tt.expectPanic {
							t.Errorf("Unexpected panic: %v", r)
						}
					} else {
						if tt.expectPanic {
							t.Errorf("Expected panic but did not get one")
						}
					}
				}()

				orderEntity.NewOrderMust(orderID, customerID)
			})
		}
	})

	t.Run("Test for recalculateTotal", func(t *testing.T) {
		order := orderEntity.NewOrderMust(
			valueobject.NewOrderIDMust(uuid.New()),
			valueobject.NewCustomerIDMust(uuid.New()),
		)

		itemUSD := valueobject.NewOrderItemMust(
			valueobject.NewProductIDMust(uuid.New()),
			valueobject.NewMoneyMust(100, "USD"),
			valueobject.NewQuantityMust(1),
		)

		itemBRL := valueobject.NewOrderItemMust(
			valueobject.NewProductIDMust(uuid.New()),
			valueobject.NewMoneyMust(200, "BRL"),
			valueobject.NewQuantityMust(1),
		)

		if err := order.AddItem(itemUSD); err != nil {
			t.Fatalf("Unexpected error adding item: %v", err)
		}

		itemBefore := len(order.Items())
		totalBefore := order.TotalPrice()

		err := order.AddItem(itemBRL)
		if !errors.Is(err, domainErrors.ErrCurrencyMismatch) {
			t.Fatalf("esperava ErrCurrencyMismatch, got: %v", err)
		}

		itemAfter := len(order.Items())
		totalAfter := order.TotalPrice()

		if itemBefore != itemAfter {
			t.Errorf("Expected item count to remain the same after failed addition. Before: %d, After: %d", itemBefore, itemAfter)
		}

		if totalBefore != totalAfter {
			t.Errorf("Expected total price to remain the same after failed addition. Before: %v, After: %v", totalBefore, totalAfter)
		}
	})

	t.Run("Test for AddItem", func(t *testing.T) {
		tests := []struct {
			name        string
			item        valueobject.OrderItem
			expectedErr error
		}{
			{"Valid Item", valueobject.NewOrderItemMust(valueobject.NewProductIDMust(uuid.New()), valueobject.NewMoneyMust(100, "USD"), valueobject.NewQuantityMust(1)), nil},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				orderID, err := valueobject.NewOrderID(uuid.New())
				if err != nil {
					t.Errorf("Unexpected error creating order ID: %v", err)
				}
				customerID, err := valueobject.NewCustomerID(uuid.New())
				if err != nil {
					t.Errorf("Unexpected error creating customer ID: %v", err)
				}
				order := orderEntity.NewOrderMust(orderID, customerID)

				defer func() {
					if r := recover(); r != nil {
						if tt.name == "Valid Item" {
							t.Errorf("Unexpected panic: %v", r)
						}
					} else {
						if tt.name != "Valid Item" {
							t.Errorf("Expected panic but did not get one")
						}
					}
				}()

				order.AddItem(tt.item)
			})
		}
	})

	t.Run("Test for UpdateStatus", func(t *testing.T) {
		tests := []struct {
			name        string
			newStatus   valueobject.OrderStatus
			expectedErr error
		}{
			{"Valid Status Update", valueobject.OrderStatusPaid, nil},
			{"Invalid Status Update", "invalid_status", domainErrors.ErrInvalidOrderStatus},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				orderID, err := valueobject.NewOrderID(uuid.New())
				if err != nil {
					t.Errorf("Unexpected error creating order ID: %v", err)
				}
				customerID, err := valueobject.NewCustomerID(uuid.New())
				if err != nil {
					t.Errorf("Unexpected error creating customer ID: %v", err)
				}
				order := orderEntity.NewOrderMust(orderID, customerID)

				err = order.UpdateStatus(tt.newStatus)
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
				}
			})
		}
	})

	t.Run("Test for MarshalJSON and UnmarshalJSON", func(t *testing.T) {
		tests := []struct {
			name        string
			order       *orderEntity.Order
			expectedErr error
		}{
			{"Valid Order", orderEntity.NewOrderMust(valueobject.NewOrderIDMust(uuid.New()), valueobject.NewCustomerIDMust(uuid.New())), nil},
			{"Invalid Order (empty items)", orderEntity.NewOrderMust(valueobject.NewOrderIDMust(uuid.New()), valueobject.NewCustomerIDMust(uuid.New())), nil},
			{"Invalid Order (invalid status)", func() *orderEntity.Order {
				order := orderEntity.NewOrderMust(valueobject.NewOrderIDMust(uuid.New()), valueobject.NewCustomerIDMust(uuid.New()))
				order.UpdateStatus("invalid_status")
				return order
			}(), domainErrors.ErrInvalidOrderStatus},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				data, err := tt.order.MarshalJSON()
				if err != nil {
					t.Errorf("Unexpected error during MarshalJSON: %v", err)
				}

				var unmarshaledOrder orderEntity.Order
				err = unmarshaledOrder.UnmarshalJSON(data)
				if err != nil {
					t.Errorf("Unexpected error during UnmarshalJSON: %v", err)
				}

				if unmarshaledOrder.ID() != tt.order.ID() || unmarshaledOrder.CustomerID() != tt.order.CustomerID() {
					t.Errorf("Unmarshaled order does not match original order")
				}
			})
		}
	})
}
