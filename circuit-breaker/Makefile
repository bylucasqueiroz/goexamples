clean:
	docker compose down --volumes --remove-orphans
	rm -rf localstack

down:
	docker compose down

up:
	docker compose up -d --build

# Query last 100 rows from events
events-query:
	docker compose exec -T postgres bash -lc 'PGPASSWORD="password" psql -U user -d mydb -c "SELECT * FROM orders ORDER BY id DESC LIMIT 100;"'

# Delete rows from events. Optionally pass WHERE, e.g.:
# make events-delete WHERE="created_at < '\''2026-01-01'\''"
events-delete:
	docker compose exec -T postgres bash -lc 'PGPASSWORD="password" sql="DELETE FROM orders"; \
	if [ -n "$(WHERE)" ]; then sql="$$sql WHERE $(WHERE)"; fi; \
	PGPASSWORD="password" psql -U user -d mydb -c "$$sql"; \
	PGPASSWORD="password" psql -U user -d mydb -c "SELECT COUNT(*) AS remaining FROM orders";'

sqs-receive:
	AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy AWS_REGION=sa-east-1 aws sqs receive-message --endpoint-url=http://localhost:4566 --queue-url http://localhost:4566/000000000000/orders-queue --region sa-east-1

sqs-send:
	@count=$${COUNT:-1}; \
	while [ $$count -gt 0 ]; do \
		AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy AWS_REGION=sa-east-1 aws sqs send-message \
			--endpoint-url=http://localhost:4566 \
			--queue-url http://localhost:4566/000000000000/orders-queue \
			--message-body '{"name": "John Due", "amount": 100.00}' \
			--region sa-east-1; \
		count=$$((count-1)); \
	done