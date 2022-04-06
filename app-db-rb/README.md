# Integrating WebhookDB using a separate connection

You will need a Ruby installation with bundler plus Docker to run this demo.

From the repo root, run:

    make up
    make install

You will need the Stripe secret key for your test account.
Get it from [https://dashboard.stripe.com/test/apikeys](https://dashboard.stripe.com/test/apikeys).

    export STRIPE_API_KEY=sk_test_verylongstring123

We also need to configure the WebhookDB database connection string and Stripe table:

    export WHDB_DATABASE_URL=postgres://aro14afc01a3b145a83a8e:a14a93733100441a9dee@rob_lithic_demo.db.webhookdb.com:5432/adb14a2b6c2c2a58db1549
    export WHDB_STRIPE_CUSTOMERS=stripe_customer_v1_8f82

Then you can run the demo:

    make demo-app-db-rb

If anything is not configured, the demo will error.