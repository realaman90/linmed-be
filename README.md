# linmed-be

-- Users and Authentication
CREATE TABLE roles (
    role_id INT PRIMARY KEY AUTO_INCREMENT,
    role_name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE permissions (
    permission_id INT PRIMARY KEY AUTO_INCREMENT,
    permission_name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role_permissions (
    role_id INT,
    permission_id INT,
    FOREIGN KEY (role_id) REFERENCES roles(role_id),
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE users (
    user_id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role_id INT,
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

-- Facilities and Locations
CREATE TABLE facilities (
    facility_id INT PRIMARY KEY AUTO_INCREMENT,
    facility_name VARCHAR(255) NOT NULL,
    facility_type VARCHAR(50), -- hospital, clinic, etc.
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    contact_person VARCHAR(100),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE floors (
    floor_id INT PRIMARY KEY AUTO_INCREMENT,
    facility_id INT,
    floor_number VARCHAR(10),
    floor_name VARCHAR(100),
    floor_plan_url VARCHAR(255),
    FOREIGN KEY (facility_id) REFERENCES facilities(facility_id)
);

CREATE TABLE stations (
    station_id INT PRIMARY KEY AUTO_INCREMENT,
    floor_id INT,
    station_name VARCHAR(100),
    station_type VARCHAR(50),
    location_x DECIMAL(10,2), -- X coordinate on floor plan
    location_y DECIMAL(10,2), -- Y coordinate on floor plan
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (floor_id) REFERENCES floors(floor_id)
);

-- Products and Inventory
CREATE TABLE product_categories (
    category_id INT PRIMARY KEY AUTO_INCREMENT,
    category_name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_category_id INT,
    FOREIGN KEY (parent_category_id) REFERENCES product_categories(category_id)
);

CREATE TABLE products (
    product_id INT PRIMARY KEY AUTO_INCREMENT,
    category_id INT,
    product_name VARCHAR(255) NOT NULL,
    description TEXT,
    manufacturer VARCHAR(100),
    model_number VARCHAR(100),
    sku VARCHAR(50) UNIQUE,
    service_interval_days INT, -- Number of days between required services
    lifetime_months INT, -- Expected lifetime in months
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES product_categories(category_id)
);

CREATE TABLE station_products (
    station_product_id INT PRIMARY KEY AUTO_INCREMENT,
    station_id INT,
    product_id INT,
    serial_number VARCHAR(100),
    installation_date DATE,
    expiration_date DATE,
    next_inspection_date DATE,
    next_service_date DATE,
    status VARCHAR(20),
    notes TEXT,
    FOREIGN KEY (station_id) REFERENCES stations(station_id),
    FOREIGN KEY (product_id) REFERENCES products(product_id)
);

-- Maintenance and Services
CREATE TABLE service_records (
    service_id INT PRIMARY KEY AUTO_INCREMENT,
    station_product_id INT,
    service_type VARCHAR(50), -- installation, inspection, maintenance, replacement
    service_date DATE,
    performed_by INT, -- references users
    description TEXT,
    next_service_date DATE,
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (station_product_id) REFERENCES station_products(station_product_id),
    FOREIGN KEY (performed_by) REFERENCES users(user_id)
);

-- Alerts and Notifications
CREATE TABLE alert_types (
    alert_type_id INT PRIMARY KEY AUTO_INCREMENT,
    alert_name VARCHAR(100),
    description TEXT,
    severity VARCHAR(20) -- high, medium, low
);

CREATE TABLE alerts (
    alert_id INT PRIMARY KEY AUTO_INCREMENT,
    alert_type_id INT,
    station_product_id INT,
    alert_message TEXT,
    status VARCHAR(20), -- active, acknowledged, resolved
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP,
    resolved_by INT,
    FOREIGN KEY (alert_type_id) REFERENCES alert_types(alert_type_id),
    FOREIGN KEY (station_product_id) REFERENCES station_products(station_product_id),
    FOREIGN KEY (resolved_by) REFERENCES users(user_id)
);

-- Orders and Inventory Management
CREATE TABLE orders (
    order_id INT PRIMARY KEY AUTO_INCREMENT,
    facility_id INT,
    ordered_by INT,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20),
    total_amount DECIMAL(10,2),
    notes TEXT,
    FOREIGN KEY (facility_id) REFERENCES facilities(facility_id),
    FOREIGN KEY (ordered_by) REFERENCES users(user_id)
);

CREATE TABLE order_items (
    order_item_id INT PRIMARY KEY AUTO_INCREMENT,
    order_id INT,
    product_id INT,
    quantity INT,
    unit_price DECIMAL(10,2),
    FOREIGN KEY (order_id) REFERENCES orders(order_id),
    FOREIGN KEY (product_id) REFERENCES products(product_id)
);

-- Dashboard Preferences
CREATE TABLE dashboard_preferences (
    user_id INT,
    widget_id VARCHAR(50),
    widget_order INT,
    is_visible BOOLEAN DEFAULT true,
    settings JSON,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    PRIMARY KEY (user_id, widget_id)
);