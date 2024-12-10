# go-template-api

## Authors

- [Shay Jacoby](https://github.com/shayja)

## Description

A Starter Template API CRUD project using Go, Gin and PostgreSQL.

- Go : https://go.dev/doc/
- Gin : https://gin-gonic.com/docs/
- PostgreSQL : https://www.postgresql.org/docs/

If running the app locally, ensure that PostgreSQL is installed on your machine. Alternatively, if you run the app using Docker Compose, everything will be automatically set up within the container.

Database:

1. Create a new PostgreSQL database named "shop".
2. Add a new db user called "appuser" and assign a login password.
3. Execute the SQL script located in the /migrations directory of the project on the "shop" database.
4. Update your database credentials in the .env.local file, then rename it to .env. Ensure the file is placed in the root folder.
5. Adjust the configuration values to match the details of your "appuser" and the database root admin user.

# Database settings:

DB_HOST="localhost"
DB_USER="appuser"
DB_PASSWORD="<YOUR_PASSWORD>"
DB_NAME="shop"
DB_PORT=5432

configure the Postgres admin user credentials:
PGADMIN_DEFAULT_EMAIL="your@admmin.email.here"
PGADMIN_DEFAULT_PASSWORD="<<PGADMIN_ADMIN_PASSWORD>>"

To start the app, open Terminal:
go run ./cmd/main.go

To start using docker compose:
docker compose up --build

To stop the container:
docker-compose down --remove-orphans --volumes

## App endpoints:

**GET**
/api/v1/product

Get all products with paging

example req:
curl --location 'http://localhost:8080/api/v1/product?page=1' \
--data ''

**GET**
/api/v1/product/:id

Get a product by id

example:
curl --location 'http://localhost:8080/api/v1/product/3954d2d4-94cf-44f8-a237-fa905773cffd' \
--data ''

**POST**
/api/v1/product

Create a new product

example:
curl --location 'http://localhost:8080/api/v1/product' \
--header 'Content-Type: application/json' \
--data '{"name": "Iphone 15","description": "The latest iphone from Apple", "image": "ihone.png","price": 89.99, "sku": "AABBCC2233"}'

**PUT**
/api/v1/product/:id

Update an existing product by id

example:
curl --location --request PUT 'http://localhost:8080/api/v1/product/3954d2d4-94cf-44f8-a237-fa905773cffd' \
--header 'Content-Type: application/json' \
--data '{"name": "Samsung Galaxy S22","description": "The latest phone from Samsung", "image": "samsung.png","price": 88.99, "sku": "XXYYZZ2233"}'

**PATCH**
/api/v1/product/:id

Update product price by id

curl --location --request PATCH 'http://localhost:8080/api/v1/product/3954d2d4-94cf-44f8-a237-fa905773cffd' \
--header 'Content-Type: application/json' \
--data '{"price": 1189.99}'

**DELETE**
/api/v1/product/:id

Delete an existing product by id

example:
curl --location --request DELETE 'http://localhost:8080/api/v1/product/5' \
--header 'Content-Type: application/json' \
--data ''
