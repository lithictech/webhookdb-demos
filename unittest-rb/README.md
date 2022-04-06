# Unit testing with WebhookDB

You will need a Ruby installation with bundler, Docker,
and some way to run Postgres commands (ie `psql`) to run this demo.

From the repo root, run:

    make up
    make install

You can unit test using any of the WebhookDB integration approaches-
the details are basically the same. We will use the 'separate DB'
integration via [`app-db-rb`](https://github.com/lithictech/webhookdb-demos/tree/main/app-db-rb) for this example.

The basic idea is that we create tables in our test database,
insert into those, and query against them.
Using the same database for all tests avoids complexity
around multiple databases for testing,
though of course you can do this if desired.

To generate the table schema, run this from your [WebhookDB terminal](https://webhookdb.com/terminal/):

```
rob@lithic.tech/rob_lithic_demo > webhookdb fixtures stripe_customer_v1
CREATE TABLE stripe_customer_v1_fixture (
  pk bigserial PRIMARY KEY,
  "stripe_id" text UNIQUE NOT NULL,
  "balance" integer ,
  "created" timestamptz ,
  "email" text ,
  "name" text ,
  "phone" text ,
  "updated" timestamptz ,
  data jsonb NOT NULL
);
CREATE INDEX IF NOT EXISTS balance_idx ON stripe_customer_v1_fixture ("balance");
CREATE INDEX IF NOT EXISTS created_idx ON stripe_customer_v1_fixture ("created");
CREATE INDEX IF NOT EXISTS email_idx ON stripe_customer_v1_fixture ("email");
CREATE INDEX IF NOT EXISTS phone_idx ON stripe_customer_v1_fixture ("phone");
CREATE INDEX IF NOT EXISTS updated_idx ON stripe_customer_v1_fixture ("updated");
```

You would take that output and run it as part of your migrations against your test database.
For our demo, we've put the contents into
a [`.sql`](https://github.com/lithictech/webhookdb-demos/blob/main/unittest-rb/stripe_customer_v1.schema.sql) file
and run it before the test.

Our example unit test checks that we create an alert for all customers
who have a negative balance on their stripe customer. 

You can run the demo:

    WHDB_STRIPE_CUSTOMERS=stripe_customer_v1_fixture make demo-unittest-rb

Note that in order to run the demo, we need the customers table configured.
But in this case, we point it at the fixtures table that was created by our sql.
Another option would be to rename our production table so it and the fixtured table
have the same name- we follow this approach when integrating WebhookDB
with an ORM, as [in this example](https://github.com/lithictech/webhookdb-demos/tree/main/app-fdw-rb#view-backed-orm-models).
