CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(48) PRIMARY KEY,
  name VARCHAR(52) NOT NULL,
  email VARCHAR(64) UNIQUE, 
  password VARCHAR(256) NOT NULL,
  created_at TIMESTAMP(0) DEFAULT NOW(),
  updated_at TIMESTAMP(0) DEFAULT NOW()
);