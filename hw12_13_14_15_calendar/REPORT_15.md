"user=hw12user password=hw12user host='127.0.0.1' database=hw12calendar search_path=hw12calendar

PG_DSN="host=172.18.0.1 port=5437 dbname=hw15 user=hw15user password=hw15user " make migrate-sql-up

make migrate-rmq-up: 

make compose-up