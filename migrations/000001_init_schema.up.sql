create table if not exists items
(
    id    uuid primary key,
    name  text not null unique,
    price int  not null
);

create table if not exists employees
(
    id            uuid primary key,
    username      text not null unique,
    password_hash text not null,
    balance       int  not null
);

create table if not exists employee_inventory
(
    id          uuid primary key,
    employee_id uuid not null,
    item_id     uuid not null,
    amount      int  not null,

    foreign key (employee_id) references employees (id),
    foreign key (item_id) references items (id)
);

create table if not exists transfers
(
    id            uuid primary key,
    from_employee uuid not null,
    to_employee   uuid not null,
    amount        int  not null,

    foreign key (from_employee) references employees (id),
    foreign key (to_employee) references employees (id)
);
