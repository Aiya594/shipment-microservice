package models

import (
	"time"

	"github.com/Aiya594/shipment-microservice/internal/domain/utils"
	"github.com/google/uuid"
)

type Shipment struct {
	ID            string
	Reference     string
	Status        Status
	Origin        string
	Destination   string
	Unit          string
	Cost          float64
	Driver        string
	DriverRevenue float64
	Events        []ShipmentEvent
}

func NewShipment(origin, destination, driver, unit string, cost, revenue float64) *Shipment {
	id := uuid.New()
	ref := utils.GenerateReferenceNum()
	event := ShipmentEvent{Status: Pending, Time: time.Now()}

	return &Shipment{
		ID:            id.String(),
		Reference:     ref,
		Status:        Pending,
		Origin:        origin,
		Destination:   destination,
		Unit:          unit,
		Cost:          cost,
		Driver:        driver,
		DriverRevenue: revenue,
		Events:        []ShipmentEvent{event}}
}

func (s *Shipment) ValidStatusTransition(to Status) bool {
	stats, ok := statusTransitions[s.Status]
	if !ok {
		return false
	}
	for _, st := range stats {
		if to == st {
			return true
		}
	}
	return false
}

func (s *Shipment) UpdateStatus(to Status) error {
	valid := s.ValidStatusTransition(to)
	if !valid {
		return ErrInvalidTransition
	}

	event := ShipmentEvent{Status: to, Time: time.Now()}
	s.Status = to
	s.Events = append(s.Events, event)
	return nil
}
