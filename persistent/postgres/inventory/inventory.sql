-- product table
create table public.product (
    sku text not null,
    stock integer not null,
    constraint product_pk primary key (sku)
)

-- order table
create table public."order" (
    id uuid not null,
    sku text not null references product(sku),
    quantity integer not null,
    "status" text not null,
    constraint order_pk primary key (id)
)
