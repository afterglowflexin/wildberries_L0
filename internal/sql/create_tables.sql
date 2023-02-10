CREATE TABLE orders (
  order_id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  order_data jsonb
);

INSERT INTO orders (order_data) VALUES ('{"Red":"Red"}'), ('{"Blue":"Blue"}');

SELECT * FROM orders;