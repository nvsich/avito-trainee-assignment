-- employees table
drop index if exists employees_username_idx;

-- items table
drop index if exists items_name_idx;

-- employee_inventory table
drop index if exists employee_inventory_employee_item_idx;

-- transfers table
drop index if exists transfers_from_employee_to_employee_idx;
drop index if exists transfers_to_employee_from_employee_idx;
