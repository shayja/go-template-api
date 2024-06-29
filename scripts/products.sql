
CREATE TABLE IF NOT EXISTS public.products
(
    id integer NOT NULL DEFAULT nextval('products_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default",
    description character varying(255) COLLATE pg_catalog."default",
    image character varying(255) COLLATE pg_catalog."default",
    price numeric(10,2),
    sku character varying(255) COLLATE pg_catalog."default",
    CONSTRAINT products_pkey PRIMARY KEY (id)
)
