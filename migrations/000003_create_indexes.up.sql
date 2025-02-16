-- employees table
create unique index if not exists employees_username_idx on employees (username);

-- items table
create unique index if not exists items_name_idx on items (name);

-- employee_inventory table
create unique index if not exists employee_inventory_employee_item_idx on employee_inventory (employee_id, item_id);

-- transfers table
create unique index if not exists transfers_to_employee_from_employee_idx on transfers (to_employee, from_employee);
create unique index if not exists transfers_from_employee_to_employee_idx on transfers (from_employee, to_employee);