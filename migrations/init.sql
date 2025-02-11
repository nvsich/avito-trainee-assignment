-- TODO: это up.sql, todo сделать .down.sql
-- TODO: create table if not exists
create table products
(
    id    uuid primary key,
    name  text not null,
    price int  not null
);

insert into products(id, name, price)
values ('4ba3ad9c-07e2-45d5-9c3f-5c3ffcf2f6a5', 't-shirt', 80),
       ('9d1423c4-f8a6-416c-af24-3b03e8f1594e', 'cup', 20),
       ('2392fe7d-9d34-4d7f-9df0-4d5367ba5db8', 'book', 50),
       ('72d74ee6-00d5-4f5d-b3d8-04c319cd0c4b', 'pen', 10),
       ('4523eaa4-2fc2-4943-a9b4-a71f7a31b099', 'powerbank', 200),
       ('ce64e97d-1a18-48ab-9974-607fbac8f58d', 'hoody', 300),
       ('2cc578ec-8381-4944-a787-05d13fbae770', 'umbrella', 200),
       ('8a3b7185-547c-4008-8e03-990f6cc437ba', 'socks', 10),
       ('64a00672-9833-4fd8-8831-32c775375931', 'wallet', 50),
       ('3d7db05a-035d-4e4e-b29c-a89f6513101c', 'pink-hoody', 500);

create table employees
(
    id           uuid primary key,
    login        text not null unique,
    password_hash text not null,
    balance      int  not null
);

create table employee_inventory
(
    id          uuid primary key,
    employee_id uuid not null,
    product_id  uuid not null,
    amount      int  not null,

    foreign key (employee_id) references employees (id),
    foreign key (product_id) references products (id)
);

-- TODO: подумать надо ли это
-- create table purchases
-- (
--     id          uuid primary key,
--     employee_id uuid      not null,
--     product_id  uuid      not null,
--     date        timestamp not null
-- );

create table transactions
(
    id            uuid primary key,
    from_employee uuid not null,
    to_employee   uuid not null,
    amount        int  not null,

    foreign key (from_employee) references employees (id),
    foreign key (to_employee) references employees (id)
);
