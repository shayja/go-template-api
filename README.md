## go-template-api

## Authors

- [Shay Jacoby](https://github.com/shayja)

## Description

A Starter Template API CRUD project using Go, Gin and PostgreSQL.

- Go : https://go.dev/doc/
- Gin : https://gin-gonic.com/docs/
- PostgreSQL : https://www.postgresql.org/docs/

To start the project, open Terminal:
go run main.go

Database:

1. Create a new Postgres DB
2. Exceute the sql script from project path /scripts
3. Set your DB credentials in .env file.

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
curl --location 'http://localhost:8080/api/v1/product/1' \
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
curl --location --request PUT 'http://localhost:8080/api/v1/product/2' \
--header 'Content-Type: application/json' \
--data '{"name": "Samsung Galaxy S22","description": "The latest phone from Samsung", "image": "samsung.png","price": 88.99, "sku": "XXYYZZ2233"}'

PATCH
/api/v1/product/:id

Update product price by id

curl --location --request PATCH 'http://localhost:8080/api/v1/product/2' \
--header 'Content-Type: application/json' \
--data '{"price": 1189.99}'

DELETE
/api/v1/product/:id

Delete an existing product by id

example:
curl --location --request DELETE 'http://localhost:8080/api/v1/product/5' \
--header 'Content-Type: application/json' \
--data ''
