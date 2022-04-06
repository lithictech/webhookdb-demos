up:
	docker-compose up -d
install:
	bundle install
psql:
	psql "postgres://webhookdb_demo:webhookdb_demo@localhost:18015/webhookdb_demo"

demo-app-db-rb:
	bundle exec app-db-rb/app-db.rb
demo-app-fdw-rb:
	bundle exec app-fdw-rb/app-fdw.rb
demo-unittest-rb:
	bundle exec unittest-rb/unittest.rb
