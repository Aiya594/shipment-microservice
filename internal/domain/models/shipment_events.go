package models

import "time"

type ShipmentEvent struct {
	Status Status
	Time   time.Time
}
