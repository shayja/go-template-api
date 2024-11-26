
CREATE DATABASE shop
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LOCALE_PROVIDER = 'libc'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

CREATE ROLE appuser WITH
	LOGIN
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	NOREPLICATION
	NOBYPASSRLS
	CONNECTION LIMIT -1
	PASSWORD 'xxxxxx';
COMMENT ON ROLE appuser IS 'User that can exec commands/procedures';


GRANT appuser TO pg_read_all_data, pg_write_all_data;

CREATE TABLE IF NOT EXISTS public.products
(
    id uuid primary key NOT NULL,
    name character varying(255) COLLATE pg_catalog."default",
    description character varying(255) COLLATE pg_catalog."default",
    image character varying(255) COLLATE pg_catalog."default",
    price numeric(10,2),
    sku character varying(255) COLLATE pg_catalog."default",
    create_date timestamp without time zone,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)
REVOKE ALL ON TABLE public.products FROM appuser;
GRANT DELETE, INSERT, SELECT, UPDATE ON TABLE public.products TO appuser;


CREATE OR REPLACE PROCEDURE public.products_insert(
	IN p_name text,
	IN p_description text,
	IN p_price money,
	IN p_image text,
	IN p_sku text,
	IN p_create_date timestamp without time zone,
	INOUT next_id uuid)
LANGUAGE 'plpgsql'
AS $BODY$

BEGIN
	
    INSERT INTO products (id, name, description, price, image, sku, create_date)
    SELECT gen_random_uuid(),
           p_name,
           p_description,
           p_price,
           p_image,
           p_sku,
		   p_create_date
    RETURNING id INTO next_id;

    COMMIT;

END;
$BODY$;

CREATE OR REPLACE PROCEDURE public.products_update(
	IN product_id uuid,
	IN product_name text,
	IN product_description text,
	IN product_price money,
	IN product_image text,
	IN product_sku text)
LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
  UPDATE products 
  SET name = product_name,
  description = product_description,
  price = product_price,
  image = product_image,
  sku = product_sku 
  WHERE id = product_id;
END;
$BODY$;

CREATE OR REPLACE PROCEDURE public.products_update_image(
	IN product_id uuid,
	IN product_image text)
LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
  UPDATE products SET image = product_image WHERE id = product_id;
END;
$BODY$;

CREATE OR REPLACE PROCEDURE public.products_update_price(
	IN product_id uuid,
	IN product_price money)
LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
  UPDATE products SET price = product_price WHERE id = product_id;
END;
$BODY$;

CREATE OR REPLACE FUNCTION public.get_product(
	productid uuid)
    RETURNS SETOF products 
    LANGUAGE 'sql'
    COST 100
    VOLATILE PARALLEL UNSAFE
    ROWS 1000

AS $BODY$
SELECT id, name, description, image, price, sku, create_date
FROM products WHERE id=productId
LIMIT 1
$BODY$;


CREATE OR REPLACE FUNCTION public.get_products(
	offst integer,
	lmt integer)
    RETURNS SETOF products 
    LANGUAGE 'sql'
    COST 100
    VOLATILE PARALLEL UNSAFE
    ROWS 1000

AS $BODY$
SELECT id, name, description, image, price, sku, create_date /*, row_number() OVER (ORDER BY "name", id) AS rn*/
FROM products
ORDER BY "name", id
OFFSET offst LIMIT lmt;
$BODY$;


