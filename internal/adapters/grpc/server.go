package grpcAdapter

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Aiya594/shipment-microservice/internal/adapters/postgres"
	"github.com/Aiya594/shipment-microservice/internal/app"
	"github.com/Aiya594/shipment-microservice/internal/domain/models"
	"github.com/Aiya594/shipment-microservice/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShipmentServer struct {
	proto.UnimplementedShipmentServiceServer
	service *app.ShipmentService
	logger  *slog.Logger
}

func NewShipmentServer(service *app.ShipmentService, logger *slog.Logger) *ShipmentServer {
	return &ShipmentServer{service: service, logger: logger}
}

func (s *ShipmentServer) CreateShipment(ctx context.Context, req *proto.CreateShipmentRequest) (*proto.CreateShipmentResponse, error) {
	s.logger.Info("gRPC CreateShipment called",
		"origin", req.Origin,
		"destination", req.Destination,
		"driver", req.Driver,
		"unit", req.Unit,
		"cost", req.Cost,
		"driver_revenue", req.DriverRevenue,
	)

	shipment, err := s.service.CreateShipment(
		req.Origin,
		req.Destination,
		req.Driver,
		req.Unit,
		req.Cost,
		req.DriverRevenue,
	)
	if err != nil {
		s.logger.Error("gRPC CreateShipment failed",
			"error", err,
			"origin", req.Origin,
			"destination", req.Destination,
			"driver", req.Driver,
		)

		if errors.Is(err, app.ErrInvalidArgument) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("gRPC CreateShipment succeeded",
		"shipment_id", shipment.ID,
		"reference", shipment.Reference,
		"status", shipment.Status,
	)

	return &proto.CreateShipmentResponse{
		Id:        shipment.ID,
		Reference: shipment.Reference,
		Status:    string(shipment.Status),
	}, nil
}

func (s *ShipmentServer) GetShipment(ctx context.Context, req *proto.GetShipmentRequest) (*proto.GetShipmentResponse, error) {
	s.logger.Info("gRPC GetShipment called", "shipment_id", req.Id)

	shipment, err := s.service.GetShipment(req.Id)
	if err != nil {
		s.logger.Error("gRPC GetShipment failed",
			"shipment_id", req.Id,
			"error", err,
		)
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("gRPC GetShipment succeeded",
		"shipment_id", shipment.ID,
		"status", shipment.Status,
	)

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
	s.logger.Info("gRPC AddShipmentEvent called",
		"shipment_id", req.ShipmentId,
		"status", req.Status,
	)
	stat := models.Status(req.Status)
	err := s.service.AddShipmentEvent(req.ShipmentId, stat)
	if err != nil {
		s.logger.Error("gRPC AddShipmentEvent failed",
			"shipment_id", req.ShipmentId,
			"status", req.Status,
			"error", err,
		)

		if errors.Is(err, postgres.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else if errors.Is(err, models.ErrInvalidTransition) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("gRPC AddShipmentEvent succeeded",
		"shipment_id", req.ShipmentId,
		"status", req.Status,
	)

	return &proto.AddShipmentEventResponse{
		Success: true,
		Message: "Event added successfully",
	}, nil
}

func (s *ShipmentServer) GetShipmentEvents(ctx context.Context, req *proto.GetShipmentEventsRequest) (*proto.GetShipmentEventsResponse, error) {
	s.logger.Info("gRPC GetShipmentEvents called", "shipment_id", req.ShipmentId)

	events, err := s.service.GetShipmentEvents(req.ShipmentId)
	if err != nil {
		s.logger.Error("gRPC GetShipmentEvents failed",
			"shipment_id", req.ShipmentId,
			"error", err,
		)
		if errors.Is(err, app.ErrInvalidArgument) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoEvents := make([]*proto.ShipmentEvent, len(events))
	for i, event := range events {
		protoEvents[i] = &proto.ShipmentEvent{
			Status: string(event.Status),
			Time:   event.Time.String(),
		}
	}

	s.logger.Info("gRPC GetShipmentEvents succeeded",
		"shipment_id", req.ShipmentId,
		"events_count", len(events),
	)

	return &proto.GetShipmentEventsResponse{
		Events: protoEvents,
	}, nil
}
