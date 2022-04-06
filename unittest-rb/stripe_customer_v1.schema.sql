CREATE TABLE stripe_customer_v1_fixture
(
    pk          bigserial PRIMARY KEY,
    "stripe_id" text UNIQUE NOT NULL,
    "balance"   integer,
    "created"   timestamptz,
    "email"     text,
    "name"      text,
    "phone"     text,
    "updated"   timestamptz,
    data        jsonb       NOT NULL
);
CREATE INDEX IF NOT EXISTS balance_idx ON stripe_customer_v1_fixture ("balance");
CREATE INDEX IF NOT EXISTS created_idx ON stripe_customer_v1_fixture ("created");
CREATE INDEX IF NOT EXISTS email_idx ON stripe_customer_v1_fixture ("email");
CREATE INDEX IF NOT EXISTS phone_idx ON stripe_customer_v1_fixture ("phone");
CREATE INDEX IF NOT EXISTS updated_idx ON stripe_customer_v1_fixture ("updated");
