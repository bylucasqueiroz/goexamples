# DROP
docker compose exec -T postgres bash -lc 'PGPASSWORD="password" psql -U user -d mydb -c "DROP TABLE orders;"'

# CREATE
docker compose exec -T postgres bash -lc 'PGPASSWORD="password" psql -U user -d mydb -c "CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);"'

