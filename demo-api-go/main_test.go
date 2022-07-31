package main

import (
	"database/sql"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	// One of the benefits of using WebhookDB is that your code does NOT need to take on
	// an additional external dependency. No mocking, no fake services.
	// Instead, you connect to the same Postgres database you are already using in unit tests.
	// For example, if you are using Django or Rails, this would still be your 'application' database,
	// it does not need to be one dedicated to WebhookDB.
	//
	// Note that there is a Docker Compose file in this repo that sets up the database used in this unit test.
	// Running `make test` in this repo also starts up that database.
	db, err := sql.Open("postgres", "postgres://webhookdb_demo:webhookdb_demo@localhost:18015/webhookdb_demo?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	// We need to create the table that we'll be querying.
	// to get this schema, run `webhookdb fixtures stripe_product_v1`.
	// It will spit out the SQL that is used to create the WebhookDB table.
	//
	// When this SQL should be run depends on your database fixturing strategy:
	// - If you create your test database outside of running your tests
	//   (ie, you need to explicitly migrate your test database),
	//   you should put this code in SQL files and run it as part of your test database migration.
	// - If you create your test database as part of a test run,
	//   you should run this at the appropriate time, like once per-suite
	//   (for example in a test class setup or before-all hook).
	//
	// Note that WebhookDB 'normalizes' data into standard types. For example,
	// while Stripe's "created" field is an integer Unix timestamp,
	// WebhookDB denormalizes that into a `timestamptz` field.
	// However the default integer timestamp would be available under `data->>'created'`.
	if _, err := db.Exec(`
DROP TABLE IF EXISTS stripe_product_v1_fixture; 
CREATE TABLE stripe_product_v1_fixture (
  pk bigserial PRIMARY KEY,
  stripe_id text UNIQUE NOT NULL,
  created timestamptz,
  name text,
  package_dimensions text,
  statement_descriptor text,
  unit_label text,
  updated timestamptz,
  data jsonb NOT NULL
)`); err != nil {
		t.Fatal(err)
	}
	// Fixture the test rows. We want to find the 'Match 1' row.
	// Note that we only need to insert the columns and data we care about:
	// in this case, 'data' and its 'description' key.
	// Alternatively, you could store actual API responses and insert them.
	// NOTE: Certain licenses allow the `webhookdb fixtures` call to return SQL that will
	// populate denormalized columns (like 'created', etc) from the resource/'data' payload.
	// This is sort of the WebhookDB sauce, though, so we aren't ready to
	// include it in the standard fixtured schemas yet.
	if _, err := db.Exec(`INSERT INTO public.stripe_product_v1_fixture(stripe_id, data) VALUES ('pr_1', '{"description":"Match 1"}')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO public.stripe_product_v1_fixture(stripe_id, data) VALUES ('pr_2', '{"description":"Nothing 1"}')`); err != nil {
		t.Fatal(err)
	}
	// Run the code and make sure we get back what we expect.
	b, err := RunSearch(db, "stripe_product_v1_fixture", "atch")
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"description": "Match 1"}`
	if strings.TrimSpace(string(b)) != expected {
		t.Fatalf("%s should have equaled %s", string(b), expected)
	}

}
