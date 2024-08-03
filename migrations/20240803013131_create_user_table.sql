-- Add migration script here
CREATE TABLE "user" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
