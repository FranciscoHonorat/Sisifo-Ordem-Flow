package event

type OrderCancelled struct {
	BaseEvent
	CustomerID string
	Reason     string
}

func NewOrderCancelled(orderID, customerID, reason string) OrderCancelled {
	return OrderCancelled{
		BaseEvent:  NewBaseEvent("order.cancelled", orderID),
		CustomerID: customerID,
		Reason:     reason,
	}
}
