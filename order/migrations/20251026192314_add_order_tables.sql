-- +goose Up
create table orders (
	id uuid primary key default gen_random_uuid(),
	user_id uuid not null,
	total_price numeric(18,2) not null,
	transaction_uuid uuid,
	payment_method varchar,
	status varchar not null
);

create table order_items (
	id uuid primary key default gen_random_uuid(),
	order_id uuid not null,
	part_uuid uuid not null,

	constraint fk_order_items_order
		foreign key (order_id)
		references orders(id)
);

-- +goose Down
drop table order_items;
drop table orders;
