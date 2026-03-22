package app

import (
	"fmt"
	"time"

	"github.com/Aiya594/shipment-microservice/internal/domain/models"
	"github.com/Aiya594/shipment-microservice/internal/ports"
)

type ShipmentService struct {
	repo ports.ShipmentsRepository
}

func NewShipmentService(repo ports.ShipmentsRepository) *ShipmentService {
	return &ShipmentService{repo: repo}
}

func (s *ShipmentService) CreateShipment(origin, destination, driver, unit string, cost, revenue float64) (*models.Shipment, error) {
	if origin == "" || destination == "" || driver == "" {
		return nil, fmt.Errorf("origin, destination, and driver are required:%w", ErrBadCredentials)
	}
	if cost < 0 || revenue < 0 {
		return nil, fmt.Errorf("cost and revenue must be non-negative:%w", ErrBadCredentials)
	}

	shipment := models.NewShipment(origin, destination, driver, unit, cost, revenue)
	err := s.repo.Save(shipment)
	if err != nil {
		return nil, fmt.Errorf("failed to save shipment: %w", err)
	}
	return shipment, nil
}

func (s *ShipmentService) GetShipment(id string) (*models.Shipment, error) {
	if id == "" {
		return nil, fmt.Errorf("shipment ID is required:%w", ErrBadCredentials)
	}

	shipment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	events, err := s.repo.GetEvents(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment events: %w", err)
	}
	shipment.Events = events

	return shipment, nil
}

func (s *ShipmentService) AddShipmentEvent(shipmentID string, status models.Status) error {
	if shipmentID == "" {
		return fmt.Errorf("shipment ID is required:%w", ErrBadCredentials)
	}

	//get current shipment to validate transition
	shipment, err := s.GetShipment(shipmentID)
	if err != nil {
		return fmt.Errorf("failed to get shipment for status update: %w", err)
	}

	//validate status transition
	if !shipment.ValidStatusTransition(status) {
		return fmt.Errorf("invalid status transition from %s to %s", shipment.Status, status)
	}

	//create event
	event := &models.ShipmentEvent{
		Status: status,
		Time:   time.Now(),
	}

	//save event
	err = s.repo.AddEvent(event, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to add shipment event: %w", err)
	}

	return nil
}

func (s *ShipmentService) GetShipmentEvents(shipmentID string) ([]models.ShipmentEvent, error) {
	if shipmentID == "" {
		return nil, fmt.Errorf("shipment ID is required:%w", ErrBadCredentials)
	}

	events, err := s.repo.GetEvents(shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment events: %w", err)
	}

	return events, nil
}
