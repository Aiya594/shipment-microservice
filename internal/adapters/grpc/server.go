package grpcAdapter

import (
	"context"

	"github.com/Aiya594/shipment-microservice/internal/app"
	"github.com/Aiya594/shipment-microservice/internal/domain/models"
	"github.com/Aiya594/shipment-microservice/proto"
)

type ShipmentServer struct {
	proto.UnimplementedShipmentServiceServer
	service *app.ShipmentService
}

func NewShipmentServer(service *app.ShipmentService) *ShipmentServer {
	return &ShipmentServer{service: service}
}

func (s *ShipmentServer) CreateShipment(ctx context.Context, req *proto.CreateShipmentRequest) (*proto.CreateShipmentResponse, error) {
	shipment, err := s.service.CreateShipment(
		req.Origin,
		req.Destination,
		req.Driver,
		req.Unit,
		req.Cost,
		req.DriverRevenue,
	)
	if err != nil {
		return nil, err
	}

	return &proto.CreateShipmentResponse{
		Id:        shipment.ID,
		Reference: shipment.Reference,
		Status:    string(shipment.Status),
	}, nil
}

func (s *ShipmentServer) GetShipment(ctx context.Context, req *proto.GetShipmentRequest) (*proto.GetShipmentResponse, error) {
	shipment, err := s.service.GetShipment(req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.GetShipmentResponse{
		Id:            shipment.ID,
		Reference:     shipment.Reference,
		Status:        string(shipment.Status),
		Origin:        shipment.Origin,
		Destination:   shipment.Destination,
		Unit:          shipment.Unit,
		Cost:          shipment.Cost,
		Driver:        shipment.Driver,
		DriverRevenue: shipment.DriverRevenue,
	}, nil
}

func (s *ShipmentServer) AddShipmentEvent(ctx context.Context, req *proto.AddShipmentEventRequest) (*proto.AddShipmentEventResponse, error) {
	status := models.Status(req.Status)
	err := s.service.AddShipmentEvent(req.ShipmentId, status)
	if err != nil {
		return &proto.AddShipmentEventResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &proto.AddShipmentEventResponse{
		Success: true,
		Message: "Event added successfully",
	}, nil
}

func (s *ShipmentServer) GetShipmentEvents(ctx context.Context, req *proto.GetShipmentEventsRequest) (*proto.GetShipmentEventsResponse, error) {
	events, err := s.service.GetShipmentEvents(req.ShipmentId)
	if err != nil {
		return nil, err
	}

	protoEvents := make([]*proto.ShipmentEvent, len(events))
	for i, event := range events {
		protoEvents[i] = &proto.ShipmentEvent{
			Status: string(event.Status),
			Time:   event.Time.String(),
		}
	}

	return &proto.GetShipmentEventsResponse{
		Events: protoEvents,
	}, nil
}
