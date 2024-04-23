CREATE TABLE IF NOT EXISTS TAX_ALLOWANCECONFIG (
  id SERIAL PRIMARY KEY,
  allowance_type VARCHAR(30) NOT NULL,
  allowance_amount DECIMAL(10,2) NOT NULL
);

INSERT INTO TAX_ALLOWANCECONFIG (allowance_type, allowance_amount) VALUES
 ('personal', 60000),
 ('k-receipt', 50000);