# Example API

This is an example app using Stripe to create Products
and [WebhookDB](https://webhookdb.com) to query them
based on their description, which is impossible using the Stripe API.

There is also a test suite showing how to write unit tests against
code depending on your WebhookDB database.

The API has two endpoints:

- An endpoint to create a product in Stripe. It just proxies along what it receives in the request.
- An endpoint to search products in Stripe based on their description.

You can try out this demo by calling the following. You may need to use a different product name:

`curl -X POST -d 'name=My Awesome Product&description=Great product' https://webhookdb-demo-api-go.herokuapp.com/create`

And then to query:

`curl https://webhookdb-demo-api-go.herokuapp.com/search?q=Great`

Read through the comments in `main.go` and `main_test.go` to understand what is going on
and how to implement this in your own application.

## Local Setup

You should be able to run `main test` to run unit tests.

To run the server yourself, you need to set up your WebhookDB Stripe integration,
and then set some environment variables.

```
$ webhookdb auth login

$ webhookdb integrations create stripe_product_v1
You are about to start reflecting Stripe Product info into webhookdb.
We've made an endpoint available for Stripe Product webhooks:

https://api.webhookdb.com/v1/service_integrations/svi_somelongid

From your Stripe Dashboard, go to Developers -> Webhooks -> Add Endpoint.
Use the URL above, and choose all of the following events:
  product.created
  product.deleted
  product.updated
Then click Add Endpoint.

The page for the webhook will have a 'Signing Secret' section.
Reveal it, then copy the secret (it will start with `whsec_`).
      
Paste or type your secret here: 

Great! WebhookDB is now listening for Stripe Product webhooks.
You can query the table through your organization's Postgres connection string:

  psql postgres://arousername:apassword@rob_demo_lithic_tech_org.db.webhookdb.com:5432/dbname
  > SELECT * FROM stripe_product_v1_00ca

You can also run a query through the CLI:

  webhookdb db sql "SELECT * FROM stripe_product_v1_00ca"
  
In order to backfill existing Stripe Products, run this from a shell:

  webhookdb backfill svi_eeswudrwhr1yw68lrj7q5h7b1
```

You need the database connection string and table name printed out in that response,
plus your Stripe API key. Set that in your environment:

```bash
export WEBHOOKDB_URL=postgres://arousername:apassword@rob_demo_lithic_tech_org.db.webhookdb.com:5432/dbname
export WEBHOOKDB_TABLE=stripe_product_v1_00ca
export STRIPE_API_KEY=sk_test_your_stripe_test_key
```

Then `make run`, and you can `curl` to see the demo in action:

```
curl -X POST -d 'name=Test 1&description=First Test Product' http://localhost:18018/create
curl "http://localhost:18018/search?q=Test"
```

It's worth noting that the Stripe webhook is not hitting your local service;
it is hitting WebhookDB's servers. Your application is talking to
one of the WebhookDB database servers to get the Stripe data. 
