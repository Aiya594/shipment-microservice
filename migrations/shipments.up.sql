CREATE TABLE shipments (
    id VARCHAR(50) PRIMARY KEY,
    reference VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    origin VARCHAR(100) NOT NULL,
    destination VARCHAR(100) NOT NULL,
    unit VARCHAR(50),
    cost DECIMAL(10,2),
    driver VARCHAR(100),
    driver_revenue DECIMAL(10,2)
);