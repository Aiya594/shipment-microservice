# Shipment Tracking Microservice

A gRPC-based microservice for managing shipments and tracking their status changes throughout the logistics lifecycle.

# Overview

This service models a simplified Transportation Management System (TMS) component responsible for:

- Creating shipments
- Retrieving shipment details
- Tracking shipment status changes
- Storing and retrieving shipment event history

Each shipment progresses through a lifecycle, and every status change is recorded as an event.

## Architecture

This service follows Clean Architecture principles with clear separation of concerns:

- **Domain Layer** (`internal/domain/`): Contains business logic, entities, and rules
- **Application Layer** (`internal/app/`): Contains use cases and application logic
- **Infrastructure Layer** (`internal/adapters/`): Contains external concerns (gRPC, database)
- **Ports** (`internal/ports/`): Defines interfaces for external dependencies

### Architecture Flow
```
gRPC -> Application Layer -> Domain Layer -> Repository Interface -> Infrastructure
```

- gRPC handlers call use cases
- Use cases orchestrate domain logic
- Domain contains pure business rules
- Infrastructure implements persistence

## Features

- Create shipments with origin, destination, driver, and cost information
- Retrieve shipment details
- Track shipment status changes through valid state transitions
- Maintain audit trail of all status changes
- gRPC API for efficient communication


## Design Decisions

### Clean Architecture
- Domain logic is independent of infrastructure
- Easy to swap implementations (e.g., different databases, protocols)

### Status Transitions
- Enforced at the domain level to ensure business rule compliance
- Invalid transitions are rejected with clear error messages
- Status history is maintained for audit purposes

### Database Schema
- Separate tables for shipments and events
- Events table maintains complete audit trail
- Foreign key constraints ensure data integrity
- PostgreSQL is used as the primary database
- Configuration is provided via environment variables

### gRPC Protocol
- Efficient binary protocol suitable for microservices
- Strongly typed contracts via Protocol Buffers
- Support for streaming if needed in the future

## Assumptions

- Shipment reference numbers are generated using format: `REF-DDMMYYYY-XXXX`
- All monetary values are in a single currency
- Status transitions are validated but no complex business rules beyond valid states
- Database connection is handled via environment variables
- Service runs on a single port for all gRPC endpoints

## Shipment Lifecycle
Statuses:

- **Pending**: Initial status when shipment is created
- **Picked Up**: Shipment has been collected
- **In Transit**: Shipment is on the way to destination
- **Delivered**: Shipment has reached its destination
- **Cancelled**: Shipment has been cancelled (can be done from **pending** and **delivered** states)

Default lifecycle:
```
pending → picked_up → in_transit → delivered
```

Business Rules:
- Shipment starts with **pending**
- Status transitions must be valid
- Invalid transitions are rejected

- Each status change:
   - creates an event
   - updates current shipment status
- Duplicate or illogical updates are prevented

# gRPC API

Defined using Protocol Buffers.

Main RPC Methods:
- CreateShipment
- GetShipment
- AddShipmentEvent
- GetShipmentEvents

## How to Run
1. Clone repository
```shell
  git clone https://github.com/Aiya594/shipment-microservice
  cd shipment-microservice
```
   
2. Set environment variables

Create .env as given in `.env.example`

3. Run database migrations (see `migrations/` directory)
4. Run the service
```shell
  go run main.go
```

## Testing

Due to time constraints, automated tests were not implemented.

The system is structured to allow easy testing of core business logic:
- domain layer is independent of infrastructure
- repository interfaces enable mocking

Planned test coverage includes:
- shipment creation
- validation of status transitions
- rejection of invalid transitions

## Future Enhancements

- Add authentication/authorization
- Implement event-driven architecture with message queues
- Add metrics and monitoring
- Implement rate limiting
- Add Docker containerization
