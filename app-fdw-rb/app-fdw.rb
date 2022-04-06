#!/usr/bin/env ruby
# frozen_string_literal: true

require "sequel"
require "stripe"

def _mustenv(v)
  e = ENV[v]
  return e if e
  raise "Must set #{v}"
end

database_url = "postgres://webhookdb_demo:webhookdb_demo@localhost:18015/webhookdb_demo"
stripe_key = _mustenv("STRIPE_API_KEY")
whdb_schema = _mustenv("WHDB_SCHEMA").to_sym
whdb_stripe_customers = _mustenv("WHDB_STRIPE_CUSTOMERS").to_sym
table = :app_fdw_rb_customers

Stripe.api_key = stripe_key

DB = Sequel.connect(database_url)
DB.extension :pg_json

DB.create_table?(table) do
  primary_key :id
  text :name
  text :stripe_id
end

name = 'App Demo using FDW'
stripe_customer = Stripe::Customer.create(name: name)
app_customer = DB[table].returning.insert(name: name, stripe_id: stripe_customer.id).first

puts "Created application customer: #{app_customer}"
puts "You can use the Stripe Customer object that was returned at this point."
puts "For future calls though, you will want to query the customer via WebhookDB,"
puts "rather than calling the Stripe API (or keeping in sync in any other way)."
puts ""

found = loop do
  found = DB[Sequel[whdb_schema][whdb_stripe_customers]].where(stripe_id: app_customer[:stripe_id]).select(:pk, :stripe_id, :name, :data).first
  break found if found
  puts "Waiting 1s for Stripe Customer to appear in WebhookDB FDW"
  sleep(1)
end

puts "Found the customer in WebhookDB FDW: #{found}"
