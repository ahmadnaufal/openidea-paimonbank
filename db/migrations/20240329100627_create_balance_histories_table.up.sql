CREATE TABLE IF NOT EXISTS balance_histories (
  id VARCHAR(48) PRIMARY KEY,
  user_id VARCHAR(48) NOT NULL,
  currency VARCHAR(6) NOT NULL,
  balance INTEGER NOT NULL,
  source_bank_account_number VARCHAR(32) NOT NULL, 
  source_bank_name VARCHAR(32) NOT NULL,
  transfer_proof_img_url VARCHAR(128) NOT NULL,
  created_at TIMESTAMP(0) DEFAULT NOW(),
  updated_at TIMESTAMP(0) DEFAULT NOW()
);
