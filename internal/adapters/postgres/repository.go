package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Aiya594/shipment-microservice/internal/domain/models"
	"github.com/Aiya594/shipment-microservice/internal/ports"
)

type ShipmentsPostgres struct {
	db *sql.DB
}

func NewShipmentRepo(db *sql.DB) ports.ShipmentsRepository {
	return &ShipmentsPostgres{db: db}
}

// insert shipment info into shipments and initial event(status Pending) in events
func (sr *ShipmentsPostgres) Save(s *models.Shipment) error {
	tx, err := sr.db.Begin()
	if err != nil {
		return fmt.Errorf("couldnt begin transaction:%w", err)
	}
	defer tx.Rollback()
	_, err = tx.Exec(`INSERT INTO shipments (id,reference,status,origin,destination,unit,cost,driver,driver_revenue)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		s.ID, s.Reference, s.Status,
		s.Origin, s.Destination,
		s.Unit, s.Cost,
		s.Driver, s.DriverRevenue)
	if err != nil {
		return fmt.Errorf("couldnt save shipment information:%w", err)
	}

	_, err = tx.Exec(`INSERT INTO events(shipment_id, status)
	VALUES($1,$2)`, s.ID, s.Status)
	if err != nil {
		return fmt.Errorf("couldnt save shipment event information:%w", err)
	}

	return tx.Commit()
}

// get shipment by its id
func (sr *ShipmentsPostgres) GetByID(id string) (*models.Shipment, error) {
	row := sr.db.QueryRow(`SELECT id,reference,status,origin,destination,unit,cost,driver,driver_revenue
	FROM shipments 
	WHERE id=$1`, id)
	var s models.Shipment
	err := row.Scan(
		&s.ID,
		&s.Reference,
		&s.Status,
		&s.Origin,
		&s.Destination,
		&s.Unit,
		&s.Cost,
		&s.Driver,
		&s.DriverRevenue,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("couldnt get info for id=%s:%w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("couldnt get info for id=%s:%w", id, err)
	}

	return &s, nil
}

// add event info into events and update status of shipment in shipments
func (sr *ShipmentsPostgres) AddEvent(se *models.ShipmentEvent, id string) error {
	tx, err := sr.db.Begin()
	if err != nil {
		return fmt.Errorf("couldnt begin transaction:%w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO events(shipment_id, status,time)
	VALUES($1,$2,$3)`, id, se.Status, se.Time)
	if err != nil {
		return fmt.Errorf("couldnt save event info:%w", err)
	}

	_, err = tx.Exec(`UPDATE shipments SET status=$1 WHERE id=$2`, se.Status, id)
	if err != nil {
		return fmt.Errorf("couldnt update shipment status:%w", err)
	}

	return tx.Commit()
}

func (sr *ShipmentsPostgres) GetEvents(id string) ([]models.ShipmentEvent, error) {
	rows, err := sr.db.Query(`
        SELECT status, time
        FROM events 
        WHERE shipment_id = $1
    `, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer rows.Close()
	var events []models.ShipmentEvent
	for rows.Next() {
		var e models.ShipmentEvent
		err := rows.Scan(&e.Status, &e.Time)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}
	return events, nil
}
