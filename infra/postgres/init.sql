-- ===================================================
-- CWCU SIAKAD + E-Learning - PostgreSQL Init Script
-- ===================================================

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE cwcu_db TO cwcu_user;
