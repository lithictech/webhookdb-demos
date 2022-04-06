# webhookdb-demos

Demo code for how to use and integrate webhookdb.com into your application, analytics, and infrastructure.

## Setup

For all of these demos, you will need a WebhookDB account, and a service integration.

We'll use a Stripe test account for these demos, since it's easy to get started
and their API is quite easy to work with.

If you don't have an account, register for one at [stripe.com](https://dashboard.stripe.com/register).
We'll be doing everything in 'Test Mode', so you won't need to create real resources.

For your WebhookDB account, you can [run the CLI in your browser](https://webhookdb.com/terminal)
and create an account right from there:

    > webhookdb auth login

We can create an organization just for these tests.
Organization slugs must be unique, so you'll have to create your own:

    rob@lithic.tech/lithictech > webhookdb org create rob-lithic-demo
    Organization created with identifier 'rob_lithic_demo'.
    Use `webhookdb org invite` to invite members to rob-lithic-demo.
    Organization rob-lithic-demo (rob_lithic_demo) is now active.

Finally, let's create a service integration for Stripe customers,
so we can sync them down.
**NOTE**: Make sure you set up a *development* webhook.
The URL will be [https://dashboard.stripe.com/test/webhooks/create](https://dashboard.stripe.com/test/webhooks/create),
notice 'test' in the path.

```
rob@lithic.tech/rob_lithic_demo > webhookdb integrations create stripe_customer_v1

You are about to start reflecting Stripe Customer info into webhookdb.
We've made an endpoint available for Stripe Customer webhooks:

https://api.webhookdb.com/v1/service_integrations/svi_6i8987gxug15z0mbhi0dhwns0

From your Stripe Dashboard, go to Developers -> Webhooks -> Add Endpoint.
Use the URL above, and choose all of the Customer events.
Then click Add Endpoint.

The page for the webhook will have a 'Signing Secret' section.
Reveal it, then copy the secret (it will start with `whsec_`).
      

Paste or type your secret here:
```

Do as the prompt says (use a 'test' webhook as explained above),
and enter the signing secret:

```
Great! WebhookDB is now listening for Stripe Customer webhooks.
You can query the table through your organization's Postgres connection string:

  psql postgres://aro14afc01a3b145a83a8e:a14a93733100441a9dee@rob_lithic_demo.db.webhookdb.com:5432/adb14a2b6c2c2a58db1549
  > SELECT * FROM stripe_customer_v1_8f82

You can also run a query through the CLI:

  webhookdb db sql "SELECT * FROM stripe_customer_v1_8f82"
  
In order to backfill existing Stripe Customers, run this from a shell:

  webhookdb backfill svi_6i8987gxug15z0mbhi0dhwns0
```

If you have not used WebhookDB before, I would suggest you run a query as above.
The database connection is for a real Stripe test account and WebhookDB org
so may contain real data.

## Integrating WebhookDB using a separate connection

The easiest route to integrate WebhookDB into your application is
by creating a new connection to its database.

See the `app-db-rb` folder for example code.

The downside of this approach is that it requires a separate database connection,
and cannot be used for `JOIN` since it is another DB.

## Integrating WebhookDB using Foreign Data Wrappers

If you are using WebhookDB for application data on your critical path/web requests,
we suggest using a Foreign Data Wrapper (FDW) so you can `JOIN` data in WebhookDB
with your own application tables.

See the `app-fdw-rb` folder for example code.

The downside of this approach is that some analytics queries may be slower against
a FDW than real tables. For analytics queries, you can still use a separate connection string.

## Integrating WebhookDB into your database directly

For Enterprise-level customers, we can use a database user you provision
to write data directly into your database. For most query patterns, this does not provide much
of a benefit over the other options, but it is available if needed.

## Unit testing with WebhookDB

When we say "WebhookDB is built to be used by your application,"
we also mean "very easy to test with."

Specifically, WebhookDB is much easier to test with when compared to API mocking;
instead, you can insert DB objects that are queried by your test,
rather than dealing with HTTP mocks to endpoints.

See the `unittest-rb` folder for example code. 
