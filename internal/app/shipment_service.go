package app

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Aiya594/shipment-microservice/internal/domain/models"
	"github.com/Aiya594/shipment-microservice/internal/ports"
)

type ShipmentService struct {
	repo   ports.ShipmentsRepository
	logger *slog.Logger
}

func NewShipmentService(repo ports.ShipmentsRepository, logger *slog.Logger) *ShipmentService {
	return &ShipmentService{repo: repo, logger: logger}
}

func (s *ShipmentService) CreateShipment(origin, destination, driver, unit string, cost, revenue float64) (*models.Shipment, error) {
	if origin == "" || destination == "" || driver == "" {
		s.logger.Error("couldnt create shipment", "error", ErrInvalidArgument, "origin", origin, "destination", destination, "driver", driver)
		return nil, fmt.Errorf("origin, destination, and driver are required:%w", ErrInvalidArgument)
	}
	if cost < 0 || revenue < 0 {
		s.logger.Error("couldnt create shipment", "error", ErrInvalidArgument, "cost", cost, "revenue", revenue)
		return nil, fmt.Errorf("cost and revenue must be non-negative:%w", ErrInvalidArgument)
	}

	shipment := models.NewShipment(origin, destination, driver, unit, cost, revenue)
	err := s.repo.Save(shipment)
	if err != nil {
		s.logger.Error("couldnt save shipment", "error", err.Error())
		return nil, fmt.Errorf("failed to save shipment: %w", err)
	}
	s.logger.Info("shipment created", "ID", shipment.ID, "reference", shipment.Reference, "status", shipment.Status)
	return shipment, nil
}

func (s *ShipmentService) GetShipment(id string) (*models.Shipment, error) {
	if id == "" {
		s.logger.Error("get shipment failed: empty id", "error", ErrInvalidArgument)
		return nil, fmt.Errorf("shipment ID is required:%w", ErrInvalidArgument)
	}

	shipment, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get shipment by id",
			"shipment_id", id,
			"error", err)
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	events, err := s.repo.GetEvents(id)
	if err != nil {
		s.logger.Error("failed to get events for shipment",
			"shipment_id", id,
			"error", err)
		return nil, fmt.Errorf("failed to get shipment events: %w", err)
	}
	shipment.Events = events
	s.logger.Info("shipment retrieved successfully",
		"shipment_id", id,
		"status", shipment.Status,
		"events_count", len(events))

	return shipment, nil
}

func (s *ShipmentService) AddShipmentEvent(shipmentID string, status models.Status) error {
	if shipmentID == "" {
		s.logger.Error("add shipment event failed: empty shipment id",
			"error", ErrInvalidArgument)
		return fmt.Errorf("shipment ID is required:%w", ErrInvalidArgument)
	}

	//get current shipment to validate transition
	shipment, err := s.GetShipment(shipmentID)
	if err != nil {
		s.logger.Error("failed to get shipment for event",
			"shipment_id", shipmentID,
			"status", status,
			"error", err)
		return fmt.Errorf("failed to get shipment for status update: %w", err)
	}

	//validate status transition
	if !shipment.ValidStatusTransition(status) {
		s.logger.Warn("invalid status transition attempted",
			"shipment_id", shipmentID,
			"from_status", shipment.Status,
			"to_status", status)
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
		s.logger.Error("failed to add event to repository",
			"shipment_id", shipmentID,
			"status", status,
			"error", err)
		return fmt.Errorf("failed to add shipment event: %w", err)
	}

	s.logger.Info("shipment event added",
		"shipment_id", shipmentID,
		"new_status", status,
		"old_status", shipment.Status,
		"event_time", event.Time)
	return nil
}

func (s *ShipmentService) GetShipmentEvents(shipmentID string) ([]models.ShipmentEvent, error) {
	if shipmentID == "" {
		s.logger.Error("get shipment events failed: empty shipment id",
			"error", ErrInvalidArgument)
		return nil, fmt.Errorf("shipment ID is required:%w", ErrInvalidArgument)
	}

	events, err := s.repo.GetEvents(shipmentID)
	if err != nil {
		s.logger.Error("failed to get events from repository",
			"shipment_id", shipmentID,
			"error", err)
		return nil, fmt.Errorf("failed to get shipment events: %w", err)
	}
	s.logger.Info("shipment events retrieved",
		"shipment_id", shipmentID,
		"events_count", len(events))

	return events, nil
}
