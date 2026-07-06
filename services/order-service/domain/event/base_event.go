package event

import "time"

type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
	AggregateId() string
}

type BaseEvent struct {
	eventName   string
	occurredAt  time.Time
	aggregateId string
}

func NewBaseEvent(name, aggregateId string) BaseEvent {
	return BaseEvent{
		eventName:   name,
		occurredAt:  time.Now().UTC(),
		aggregateId: aggregateId,
	}
}

func (b BaseEvent) EventName() string     { return b.eventName }
func (b BaseEvent) OccurredAt() time.Time { return b.occurredAt }
func (b BaseEvent) AggregateId() string   { return b.aggregateId }
