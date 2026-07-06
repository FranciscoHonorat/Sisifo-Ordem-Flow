package event

import "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"

type OrderAdded struct {
	BaseEvent
	OrderID    valueobject.OrderID
	ProductID  valueobject.ProductID
	Quantity   valueobject.Quantity
	UnitPrice  valueobject.Money
	TotalPrice valueobject.Money
}

func NewOrderAdded(orderID valueobject.OrderID, productID valueobject.ProductID, quantity valueobject.Quantity, unitPrice, totalPrice valueobject.Money) OrderAdded {
	return OrderAdded{
		BaseEvent:  NewBaseEvent("order.item_added", orderID.String()),
		OrderID:    orderID,
		ProductID:  productID,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}
}
