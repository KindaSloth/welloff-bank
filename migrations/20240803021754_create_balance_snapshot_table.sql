-- Add migration script here
CREATE TABLE "balance_snapshot" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
  account_id UUID NOT NULL,
  balance DECIMAL(15, 2) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES "account"(id)
);
