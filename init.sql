CREATE TABLE IF NOT EXISTS tax_deduction_configs (
    id SERIAL PRIMARY KEY,
    personal DECIMAL(10, 2) DEFAULT 60000.00,
    kreceipt DECIMAL(10, 2) DEFAULT 50000.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tax_deduction_configs (personal, kreceipt) VALUES (60000.00, 50000.000);
