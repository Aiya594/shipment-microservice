package ports

import "github.com/Aiya594/shipment-microservice/internal/domain/models"

type ShipmentsRepository interface {
	Save(s *models.Shipment) error
	GetByID(id string) (*models.Shipment, error)
	AddEvent(se *models.ShipmentEvent, id string) error
	GetEvents(id string) ([]models.ShipmentEvent, error)
}
