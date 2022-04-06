# Integrating WebhookDB using a Foreign Data Wrapper

You will need a Ruby installation with bundler, Docker,
and some way to run Postgres commands (ie `psql`) to run this demo.

From the repo root, run:

    make up
    make install

We will need to set up WebhookDB's database as a FDW into our Postgres.
From your WebhookDB terminal, generate the FDW details:

```
rob@lithic.tech/rob_lithic_demo > webhookdb db fdw

CREATE EXTENSION IF NOT EXISTS postgres_fdw;
DROP SERVER IF EXISTS webhookdb_remote CASCADE;
CREATE SERVER webhookdb_remote
  FOREIGN DATA WRAPPER postgres_fdw
  OPTIONS (host 'rob_lithic_demo.db.webhookdb.com', port '5432', dbname 'adb14a2b6c2c2a58db1549', fetch_size '50000');

CREATE USER MAPPING FOR CURRENT_USER
  SERVER webhookdb_remote
  OPTIONS (user 'aro14afc01a3b145a83a8e', password 'a14a93733100441a9dee');

CREATE SCHEMA IF NOT EXISTS webhookdb_remote;
IMPORT FOREIGN SCHEMA public
  FROM SERVER webhookdb_remote
  INTO webhookdb_remote;

CREATE SCHEMA IF NOT EXISTS webhookdb;

CREATE MATERIALIZED VIEW IF NOT EXISTS webhookdb.stripe_customer_v1 AS SELECT * FROM webhookdb_remote.stripe_customer_v1_8f82;
```

Copy the output, then:

    make psql

Paste the output to import WebhookDB into the database.

You now have an FDW set up. Note that when you create new WebhookDB integrations,
you can recreate the FDW using this same process.
This will make sure new tables show up as foreign tables in your database.

You can customize the output via the CLI, such as excluding the materialized view.
The materialized view eliminates any FDW-related performance concerns
but comes at the cost of potentially stale data, so is mostly useful for analytics.

You will need the Stripe secret key for your test account.
Get it from [https://dashboard.stripe.com/test/apikeys](https://dashboard.stripe.com/test/apikeys).

    export STRIPE_API_KEY=sk_test_verylongstring123

We also need to configure the Stripe schema and table for our app:

    export WHDB_SCHEMA=webhookdb_remote
    export WHDB_STRIPE_CUSTOMERS=stripe_customer_v1_8f82

Then you can run the demo:

    make demo-app-fdw-rb

If anything is not configured, the demo will error.

## View-backed ORM Models

To use view-backed models, you'll want to use a consistent, hard-coded table name
rather than the auto-generated one we provide.
You can run this (using your own integration id, use `webhookdb integrations list` to see all):

    webhookdb db rename-table svi_6i8987gxug15z0mbhi0dhwns0 --new-name=stripe_customers_v1

And rename your integration to `stripe_customers` or whatever.
Then in your ORM, you would have something like this in Ruby:

```ruby
# Sequel
class StripeCustomer < Sequel::Model(Sequel[:webhookdb_remote][:stripe_customers])
end

# ActiveRecord
class StripeCustomer < ActiveRecord::Base
  self.table_name = 'webhookdb_remote.stripe_customers'
end
```

Or line this in Python:

```py
# SQLAlchemy
Base = declarative_base()
class StripeCustomer(Base):
    __tablename__ = "stripe_customers"
    __table_args__ = {"schema": "webhookdb_remote"}

# Django
class StripeCustomer(models.Model):
    class Meta:
        db_table = '"webhookdb_remote\".\"stripe_customers"'
```