-- Add migration script here
ALTER TABLE balance_snapshot ADD CONSTRAINT balance_snapshot_account_id_key UNIQUE (account_id);