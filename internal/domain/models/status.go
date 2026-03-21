package models

type Status string

const (
	Pending   Status = "pending"
	Picked    Status = "picked_up"
	InTransit Status = "in_transit"
	Delivered Status = "delivered"
	Cancelled Status = "cancelled"
)

var statusTransitions = map[Status][]Status{
	Pending:   {Picked, Cancelled},
	Picked:    {InTransit},
	InTransit: {Delivered},
	Delivered: {Cancelled},
}
