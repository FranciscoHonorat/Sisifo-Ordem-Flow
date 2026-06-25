package valueobject

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	productID ProductID
	unitPrice Money
	quantity  Quantity
	createdAt time.Time
}

func NewOrderItem(productID ProductID, unitPrice Money, quantity Quantity) (OrderItem, error) {
	return OrderItem{
		productID: productID,
		unitPrice: unitPrice,
		quantity:  quantity,
		createdAt: time.Now(),
	}, nil
}

func NewOrderItemMust(productID ProductID, unitPrice Money, quantity Quantity) OrderItem {
	i, err := NewOrderItem(productID, unitPrice, quantity)
	if err != nil {
		panic(err)
	}
	return i
}

func (i OrderItem) SubTotal() (Money, error) {
	subTotal, err := i.unitPrice.Multiply(i.quantity.Value())
	if err != nil {
		return Money{}, fmt.Errorf("subtotal: %w", err)
	}
	return subTotal, nil
}

func (i OrderItem) ProductID() ProductID {
	return i.productID
}

func (i OrderItem) UnitPrice() Money {
	return i.unitPrice
}

func (i OrderItem) Quantity() Quantity {
	return i.quantity
}

func (i OrderItem) CreatedAt() time.Time {
	return i.createdAt
}

func (i OrderItem) MarshalJSON() ([]byte, error) {
	auxOrderItem := struct {
		ProductID string `json:"product_id"`
		UnitPrice Money  `json:"unit_price"`
		Quantity  int64  `json:"quantity"`
	}{
		ProductID: i.productID.String(),
		UnitPrice: i.unitPrice,
		Quantity:  i.quantity.Value(),
	}
	return json.Marshal(auxOrderItem)
}

func (i *OrderItem) UnmarshalJSON(data []byte) error {
	var orderItem struct {
		ProductID string `json:"product_id"`
		UnitPrice Money  `json:"unit_price"`
		Quantity  int64  `json:"quantity"`
	}

	if err := json.Unmarshal(data, &orderItem); err != nil {
		return err
	}

	parseUUID, err := uuid.Parse(orderItem.ProductID)
	if err != nil {
		return fmt.Errorf("unmarshal productID: %w", err)
	}

	productID, err := NewProductID(parseUUID)
	if err != nil {
		return err
	}
	quantity, err := NewQuantity(orderItem.Quantity)
	if err != nil {
		return err
	}
	i.productID = productID
	i.unitPrice = orderItem.UnitPrice
	i.quantity = quantity
	i.createdAt = time.Now()
	return nil
}
