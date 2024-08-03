-- Add migration script here
CREATE TYPE transaction_kind AS ENUM ('deposit', 'withdrawal', 'transfer', 'refund');

CREATE TABLE "transaction" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    kind transaction_kind NOT NULL,
    from_account_id UUID,
    to_account_id UUID,
    amount DECIMAL(15, 2) NOT NULL,
    date_issued TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    related_transaction_id UUID, -- if this is a refund, this is the transaction that was refunded

    CONSTRAINT fk_from_account FOREIGN KEY(from_account_id) REFERENCES "account"(id),
    CONSTRAINT fk_to_account FOREIGN KEY(to_account_id) REFERENCES "account"(id),
    CONSTRAINT fk_related_transaction FOREIGN KEY(related_transaction_id) REFERENCES "transaction"(id)
);
