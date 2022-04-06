#!/usr/bin/env ruby
# frozen_string_literal: true

require "sequel"
require "stripe"
require "minitest"
require "minitest/autorun"

def _mustenv(v)
  e = ENV[v]
  return e if e
  raise "Must set #{v}"
end

database_url = "postgres://webhookdb_demo:webhookdb_demo@localhost:18015/webhookdb_demo"
WHDB_STRIPE_CUSTOMERS = _mustenv("WHDB_STRIPE_CUSTOMERS").to_sym
raise "Should have done WHDB_STRIPE_CUSTOMERS=stripe_customer_v1_fixture before running command" unless
  WHDB_STRIPE_CUSTOMERS == :stripe_customer_v1_fixture
CUST_TBL = :unittest_rb_customers
ALERT_TBL = :unittest_rb_alerts

DB = Sequel.connect(database_url)
DB.extension :pg_json

#
# :section:
# Normally this would be done when your test DB is migrated
#
DB.drop_table?(WHDB_STRIPE_CUSTOMERS)
DB.drop_table?(ALERT_TBL)
DB.drop_table?(CUST_TBL)

Dir.glob("unittest-rb/*.schema.sql").each do |schpath|
  DB << File.read(schpath)
end

DB.create_table?(CUST_TBL) do
  primary_key :id
  text :name
  text :stripe_id
end

DB.create_table?(ALERT_TBL) do
  primary_key :id
  text :message
  foreign_key :customer_id, CUST_TBL
end

#
# :section:
# Here is the application code.
# Note that this workload is not practical without WebhookDB-
# you would either need to store the balance yourself,
# or you'd have to make an API call for every customer
# to determine their balance.
#

def alert_negative_balance
  # This would be a good use case for using a FDW so a single query can be used.
  # See app-fdw-rb for an example.
  DB[WHDB_STRIPE_CUSTOMERS].where { balance < 0 }.each do |stripe_cust|
    app_cust = DB[CUST_TBL].where(stripe_id: stripe_cust[:stripe_id]).first
    DB[ALERT_TBL].insert(
      customer_id: app_cust[:id],
      message: "#{app_cust[:name]} has a negative balance of #{stripe_cust[:balance]}"
    )
  end
end

#
# :section:
# And this is what you would have in a unit test.
# Fixture the application data, fixture the webhookdb data,
# then assert the application DB has in it what's expected.
#

class TestAlertNegativeBalances < Minitest::Test
  def test_alerts
    pos_balance_customer = DB[CUST_TBL].returning.insert(name: 'C1', stripe_id: 'cus_1')
    neg_balance_customer = DB[CUST_TBL].returning.insert(name: 'C2', stripe_id: 'cus_2')
    zero_balance_customer = DB[CUST_TBL].returning.insert(name: 'C3', stripe_id: 'cus_3')
    neg_balance_customer2 = DB[CUST_TBL].returning.insert(name: 'C4', stripe_id: 'cus_4')
    DB[WHDB_STRIPE_CUSTOMERS].insert(stripe_id: 'cus_1', balance: 1, data: '{}')
    DB[WHDB_STRIPE_CUSTOMERS].insert(stripe_id: 'cus_2', balance: -1, data: '{}')
    DB[WHDB_STRIPE_CUSTOMERS].insert(stripe_id: 'cus_3', balance: 0, data: '{}')
    DB[WHDB_STRIPE_CUSTOMERS].insert(stripe_id: 'cus_4', balance: -5, data: '{}')

    alert_negative_balance

    alerts = DB[ALERT_TBL].order(:id).all
    assert_equal(2, alerts.count, "Should have had 2 alerts: #{alerts}")
    assert_equal(alerts[0][:message], "C2 has a negative balance of -1")
    assert_equal(alerts[1][:message], "C4 has a negative balance of -5")
  end
end
