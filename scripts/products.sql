--Create Tables

CREATE TABLE IF NOT EXISTS products
(
    id uuid NOT NULL,
    name character varying(255),
    escription text,
    image character varying(255),
    price numeric(10,2),
    sku character varying(255),
    updated_at timestamp without time zone,
    created_at timestamp without time zone,
    CONSTRAINT products_pkey PRIMARY KEY (id)
);

GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE products TO appuser;
ALTER TABLE IF EXISTS products OWNER to appuser;

CREATE TABLE IF NOT EXISTS users
(
    id uuid NOT NULL,
    username character varying(255),
    passhash character varying(255),
    mobile character varying(50),
    name character varying(255),
    email character varying(255),
    updated_at timestamp without time zone,
    created_at timestamp without time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE users TO appuser;
ALTER TABLE IF EXISTS users OWNER to appuser;


--Create Functions

CREATE OR REPLACE FUNCTION get_product(productid uuid)
    RETURNS SETOF products 
    LANGUAGE 'sql'
    COST 100
    VOLATILE PARALLEL UNSAFE
    ROWS 1000

AS $BODY$
SELECT id, name, description, image, price, sku, updated_at, created_at
FROM products WHERE id=productId
LIMIT 1
$BODY$;

ALTER FUNCTION get_product(uuid)  OWNER TO appuser;


CREATE OR REPLACE FUNCTION get_products(offst integer, lmt integer)
    RETURNS SETOF products 
    LANGUAGE 'sql'
    COST 100
    VOLATILE PARALLEL UNSAFE
    ROWS 1000

AS $BODY$
SELECT id, name, description, image, price, sku, updated_at, created_at
FROM products
ORDER BY "name", id
OFFSET offst LIMIT lmt;
$BODY$;

ALTER FUNCTION get_products(integer, integer) OWNER TO appuser;

CREATE OR REPLACE FUNCTION get_user(userid uuid)
    RETURNS SETOF users 
    LANGUAGE 'sql'
    COST 100
    VOLATILE PARALLEL UNSAFE
    ROWS 1000

AS $BODY$
SELECT id, username, passhash, mobile, name, email, updated_at, created_at
FROM users WHERE id=userId
LIMIT 1
$BODY$;

ALTER FUNCTION get_user(uuid) OWNER TO appuser;


CREATE OR REPLACE FUNCTION get_user_by_username(
	user_name character varying)
    RETURNS SETOF users 
    LANGUAGE 'sql'
    COST 100
    VOLATILE PARALLEL UNSAFE
    ROWS 1000

AS $BODY$
SELECT id, username, passhash, mobile, name, email, updated_at, created_at
FROM users WHERE LOWER(username)=LOWER(user_name)
LIMIT 1
$BODY$;

ALTER FUNCTION get_user_by_username(character varying) OWNER TO appuser;

--Create Procedures

CREATE OR REPLACE PROCEDURE products_insert(
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
	
    INSERT INTO products (id, name, description, price, image, sku, updated_at, created_at)
    SELECT gen_random_uuid(),
           p_name,
           p_description,
           p_price,
           p_image,
           p_sku,
		   p_create_date,
		   p_create_date
    RETURNING id INTO next_id;

    COMMIT;

END;
$BODY$;
ALTER PROCEDURE products_insert(text, text, money, text, text, timestamp without time zone, uuid) OWNER TO appuser;



CREATE OR REPLACE PROCEDURE products_update(
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
  sku = product_sku,
  updated_at = NOW()
  WHERE id = product_id;
END;
$BODY$;
ALTER PROCEDURE products_update(uuid, text, text, money, text, text) OWNER TO appuser;


CREATE OR REPLACE PROCEDURE products_update_image(IN product_id uuid, IN product_image text)
LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
  UPDATE products SET image = product_image, updated_at = NOW() WHERE id = product_id;
END;
$BODY$;
ALTER PROCEDURE products_update_image(uuid, text) OWNER TO appuser;


CREATE OR REPLACE PROCEDURE products_update_price(
	IN product_id uuid,
	IN product_price money)
LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
  UPDATE products SET price = product_price, updated_at = NOW() WHERE id = product_id;
END;
$BODY$;
ALTER PROCEDURE products_update_price(uuid, money) OWNER TO appuser;


CREATE OR REPLACE PROCEDURE users_insert(
	IN p_username text,
	IN p_passhash text,
	IN p_mobile text,
	IN p_name text,
	IN p_email text,
	IN p_create_date timestamp without time zone,
	INOUT next_id uuid)
LANGUAGE 'plpgsql'
AS $BODY$

BEGIN
	
    INSERT INTO Users (id, username, passhash, mobile, name, email, updated_at, created_at)
    SELECT gen_random_uuid(),
           p_username,
           p_passhash,
           p_mobile,
           p_name,
           p_email,
		   p_create_date,
		   p_create_date
    RETURNING id INTO next_id;

    COMMIT;

END;
$BODY$;
ALTER PROCEDURE users_insert(text, text, text, text, text, timestamp without time zone, uuid) OWNER TO appuser;



--DUMMY DATA

INSERT INTO products(id, name, description, image, price, sku, updated_at, created_at)
VALUES 
('48dd8c7a-9ac1-4263-88e4-bb01b5e29001','iPhone 16 Pro Max','The latest iPhone 16 Pro Max offers everything a premium flagship smartphone should, including a brilliant 6.9-inch AMOLED display for all the media consumption -- and mobile productivity, of course. This year''s model also looks and feels different than any prior Pro Max devices due to its thinner bezels, larger screen, and addition of the Camera Control button, a physical switch that lets you quickly open the camera and snap photos without ever touching the screen.','iphone16-pro-max.png',1257.00,'iphone-16-pro','2024-12-04 20:35:22.617887','2024-12-04 20:35:22.617887'),
('063d0ff7-e17e-4957-8d92-a988caeda8a1','Samsung Galaxy S24 Ultra','Samsung''s Galaxy S24 line was among the first smartphones to go all-in on AI this year, and the S24 Ultra, the most premium of the three, is the best Android phone today. The new Galaxy AI model embedded in the device brings a host of generative capabilities, including real-time phone call translations, the ability to circle an object on screen to perform an image-based Google search, AI-assisted photo editing and transcriptions, and even a Chat Assist feature for figuring out how to phrase a message in different tones.','galaxys24.png',1167.00,'samsung-galaxy-s24-ultra','2024-12-04 20:35:22.617887','2024-12-04 20:35:22.617887'),
('c3e4bc7a-4a0d-4881-87fc-4c6d0b30037d','Google Pixel 9 Pro XL','When it comes to camera performance, you really can''t go wrong with any of the flagship devices from the big three (Apple, Samsung, and Google). Depending on your preference for color temperature and feature set, you may lean towards one manufacturer over the other. But more often than not, Google''s Pixel camera system satisfies most users, and the latest Pixel 9 Pro (and Pro XL) remains a champion for instant capturing and post-processing.','pixel9pro-xl.png',1367.00,'pixel-9-pro-xl','2024-12-04 20:35:22.617887','2024-12-04 20:35:22.617887'),
('838435c6-073f-48d8-9741-8b01268c212c','CMF Phone 1 by Nothing','The best cheap phone you can buy today is the CMF Phone 1. Starting at $239, the Phone 1 has several features going for it that put it above devices that cost hundreds of dollars more, such as the ability to manually replace the back cover, screw in accessories (including a kickstand, wallet slot, and more), and insert a MicroSD card for expanded storage.','cmfphone1.png',237.00,'CMF-Phone-1','2024-12-04 20:35:22.617887','2024-12-04 20:35:22.617887'),
('d97c1376-ecec-4e92-a444-06bbc43231c6','OnePlus Open','The number of foldable phones on the market has never been higher, thanks to the collective effort of just about every manufacturer, including Google with its Pixel 9 Pro Fold, Motorola with its Razr lineup, and OnePlus with the OnePlus Open. While Samsung has held the reins of the best foldable honor for years, I''m giving the top spot right now to the OnePlus Open. ','oneplus-open.png',785.00,'oneplus-open-2024','2024-12-04 20:35:22.617887','2024-12-04 20:35:22.617887'),
('1e024fa3-f6f5-41ca-ba3a-ec00d6c10e55','Samsung Galaxy Z Flip 6','The new Galaxy Z Flip 6, unveiled at Samsung Unpacked in July, packs a ton of character and features into a tiny clamshell that pays homage to flip phones of the past.','galaxy-z-flip-6.png',1111.00,'Galaxy-Z-Flip-6-25','2024-12-04 20:35:22.617887','2024-12-04 20:35:22.617887'),
('2d248bb4-e831-44b1-8595-446d460cc511','OnePlus 12','OnePlus has had its ups and downs over the past four years, pivoting from value-driven smartphones to ultra-premium and then back to square one with last year''s OnePlus 11. This year, it''s doubling down on its value-driven flagships.','oneplus-12.png',1049.99,'OnePlus-12','2024-12-04 21:09:04.230299','2024-12-04 20:59:07.717535'),
('6369403b-4c58-4ae9-89bd-a7884e4e6b66','Xiaomi 14T Pro','MediaTek Dimensity 9300+','657454ds01c.png',1299.00,'xiaomi-14t-pro','2024-12-04 21:09:35.806982','2024-12-04 20:35:22.617887');