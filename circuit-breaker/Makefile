clean:
	docker compose down --volumes --remove-orphans
	rm -rf localstack

down:
	docker compose down

up:
	docker compose up -d --build

db-drop:
	docker compose exec -T postgres bash -lc 'PGPASSWORD="password" psql -U user -d mydb -c "DROP TABLE orders;"'

db-create:
	docker compose exec -T postgres bash -lc 'PGPASSWORD="password" psql -U user -d mydb -c "CREATE TABLE IF NOT EXISTS orders (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, amount DECIMAL(10, 2) NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"'

db-query:
	docker compose exec -T postgres bash -lc 'PGPASSWORD="password" psql -U user -d mydb -c "SELECT * FROM orders ORDER BY id DESC LIMIT 100;"'

db-delete:
	docker compose exec -T postgres bash -lc 'PGPASSWORD="password" sql="DELETE FROM orders"; \
	if [ -n "$(WHERE)" ]; then sql="$$sql WHERE $(WHERE)"; fi; \
	PGPASSWORD="password" psql -U user -d mydb -c "$$sql"; \
	PGPASSWORD="password" psql -U user -d mydb -c "SELECT COUNT(*) AS remaining FROM orders";'

sqs-receive:
	AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy AWS_REGION=sa-east-1 aws sqs receive-message --endpoint-url=http://localhost:4566 --queue-url http://localhost:4566/000000000000/orders-queue --region sa-east-1

sqs-send:
	@count=$${COUNT:-1}; \
	if [ $${count} -le 0 ]; then \
		echo "COUNT must be a positive integer"; \
		exit 1; \
	fi; \
	while [ $${count} -gt 0 ]; do \
		AWS_PAGER="" AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy AWS_REGION=sa-east-1 aws --no-cli-pager sqs send-message \
			--endpoint-url=http://localhost:4566 \
			--queue-url=http://localhost:4566/000000000000/orders-queue \
			--message-body '{"name": "John Due", "amount": 100.00}' \
			--region sa-east-1; \
		count=$$((count-1)); \
	done
