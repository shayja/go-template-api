## go-template-api

## Authors

- [Shay Jacoby](https://github.com/shayja)

## Description

A Starter Template API CRUD project using Go, Gin and PostgreSQL.

- Go : https://go.dev/doc/
- Gin : https://gin-gonic.com/docs/
- PostgreSQL : https://www.postgresql.org/docs/

To start the project, open Terminal:
go run ./cmd/main.go

using docker compose:
docker compose up --build

To Shot down:
docker-compose down --remove-orphans --volumes

Database:

1. Create a new Postgres DB, choose "shop" as database name, create a new user "appuser" and set a login password.
2. Exceute the sql script from project path /scripts on "shop" database.
3. Set your DB credentials in .env.local file and rename the file to '.env'.
4. Add new .env file in project root with this content, change the password configuration value:

# Database settings:

DB_HOST="localhost"
DB_USER="appuser"
DB_PASSWORD="<YOUR_PASSWORD>"
DB_NAME="shop"
DB_PORT=5432

configure the Postgres admin user credentials:
PGADMIN_DEFAULT_EMAIL="your@admmin.email.here"
PGADMIN_DEFAULT_PASSWORD="<<PGADMIN_ADMIN_PASSWORD>>"

App endpoints:

GET
/api/v1/product

Get all products with paging

example req:
curl --location 'http://localhost:8080/api/v1/product?page=1' \
--data ''

GET
/api/v1/product/:id

Get a product by id

example:
curl --location 'http://localhost:8080/api/v1/product/3954d2d4-94cf-44f8-a237-fa905773cffd' \
--data ''

POST
/api/v1/product

Create a new product

example:
curl --location 'http://localhost:8080/api/v1/product' \
--header 'Content-Type: application/json' \
--data '{"name": "Iphone 15","description": "The latest iphone from Apple", "image": "ihone.png","price": 89.99, "sku": "AABBCC2233"}'

PUT
/api/v1/product/:id

Update an existing product by id

example:
curl --location --request PUT 'http://localhost:8080/api/v1/product/3954d2d4-94cf-44f8-a237-fa905773cffd' \
--header 'Content-Type: application/json' \
--data '{"name": "Samsung Galaxy S22","description": "The latest phone from Samsung", "image": "samsung.png","price": 88.99, "sku": "XXYYZZ2233"}'

PATCH
/api/v1/product/:id

Update product price by id

curl --location --request PATCH 'http://localhost:8080/api/v1/product/3954d2d4-94cf-44f8-a237-fa905773cffd' \
--header 'Content-Type: application/json' \
--data '{"price": 1189.99}'

DELETE
/api/v1/product/:id

Delete an existing product by id

example:
curl --location --request DELETE 'http://localhost:8080/api/v1/product/5' \
--header 'Content-Type: application/json' \
--data ''
