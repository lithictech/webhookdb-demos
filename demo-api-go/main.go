package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Start by making sure some environment variables are set up.
	stripeApiKey := os.Getenv("STRIPE_API_KEY")
	webhookdbUrl := os.Getenv("WEBHOOKDB_URL")
	webhookdbTable := os.Getenv("WEBHOOKDB_TABLE")
	if stripeApiKey == "" {
		log.Fatal("Must set STRIPE_API_KEY")
	}
	if !strings.HasPrefix(stripeApiKey, "sk_test_") {
		log.Fatal("This app only works with your Stripe private test key (sk_test_ prefix)")
	}
	if webhookdbUrl == "" {
		log.Fatal("Must set WEBHOOKDB_URL")
	}
	if webhookdbTable == "" {
		log.Fatal("Must set WEBHOOKDB_TABLE")
	}
	db, err := sql.Open("postgres", webhookdbUrl)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}

	// Add both endpoints.
	http.HandleFunc("/create", Create(stripeApiKey))
	http.HandleFunc("/search", Search(db, webhookdbTable))
	port := os.Getenv("PORT")
	if port == "" {
		port = "18018"
		fmt.Printf("Run curl against http://localhost:%s\n", port)
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

var client = http.Client{}

func Create(stripeApiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Here we are just proxying the request to the Stripe API.
		// Note that WebhookDB is focused on querying, not mutations.
		// You still must use the Stripe SDK (or, as we recommend and do here, just using HTTP directly)
		// to change data in the API.
		//
		// Alternatively, you can create products directly in the Stripe dashboard
		// rather than using the API.
		req, err := http.NewRequest("POST", "https://api.stripe.com/v1/products", r.Body)
		if handleErr(w, err) {
			return
		}
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(stripeApiKey+":")))
		req.Header.Add("Content-Type", r.Header.Get("Content-Type"))
		resp, err := client.Do(req)
		if handleErr(w, err) {
			return
		}
		b, err := ioutil.ReadAll(resp.Body)
		if handleErr(w, err) {
			return
		}
		w.Write(b)
	}
}

func Search(db *sql.DB, productsTable string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		// Break this out so we can test it more easily for the sake of illustration;
		// in production code you probably want to also test via the API!
		b, err := RunSearch(db, productsTable, q)
		if handleErr(w, err) {
			return
		}
		w.Write(b)
	}
}

func RunSearch(db *sql.DB, productsTable, q string) ([]byte, error) {
	// Query the database for products with a description matching the given query.
	//
	// First turn 'Great' and '*Great*' into '%GREAT%' for an ILIKE comparison.
	q = strings.ReplaceAll(q, "*", "%")
	q = strings.Trim(q, "%")
	q = "%" + q + "%"
	// Now issue the query. The 'description' column isn't denormalized by default,
	// so it needs to be queried through the 'data' column, which stores the full resource JSON.
	//
	// Note that it is actually impossible to do sort of search via the Stripe API.
	// If you wanted to look through product descriptions,
	// you'd need to paginate the entire API.
	// Instead, this query executes very quickly.
	// It could be made quicker by adding an index to the searched column/JSONB field.
	rows, err := db.Query(fmt.Sprintf("SELECT data FROM %s WHERE data->>'description' ILIKE $1", productsTable), q)
	if err != nil {
		return nil, err
	}
	w := bytes.NewBuffer(nil)
	var jsoncol string
	for rows.Next() {
		if err := rows.Scan(&jsoncol); err != nil {
			return nil, err
		}
		w.Write([]byte(jsoncol))
		w.Write([]byte("\n"))
	}
	return w.Bytes(), nil
}

func handleErr(w http.ResponseWriter, e error) bool {
	if e == nil {
		return false
	}
	w.Write([]byte(fmt.Sprintf("Error: %v", e)))
	w.WriteHeader(500)
	return true
}
