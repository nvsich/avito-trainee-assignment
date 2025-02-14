package setup

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

// TODO: ПЕРЕДЕЛАТЬ, ЧТОБЫ БРАТЬ ИЗ МИГРАЦИЙ!!!
var initScript = `
create table items
(
    id    uuid primary key,
    name  text not null,
    price int  not null
);

insert into items(id, name, price)
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
    id            uuid primary key,
    username      text not null unique,
    password_hash text not null,
    balance       int  not null
);

create table employee_inventory
(
    id          uuid primary key,
    employee_id uuid not null,
    item_id     uuid not null,
    amount      int  not null,

    foreign key (employee_id) references employees (id),
    foreign key (item_id) references items (id)
);

create table transfers
(
    id            uuid primary key,
    from_employee uuid not null,
    to_employee   uuid not null,
    amount        int  not null,

    foreign key (from_employee) references employees (id),
    foreign key (to_employee) references employees (id)
);
`

func TestPostgres(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "shop_test",
		},
		WaitingFor: wait.ForSQL("5432/tcp", "pgx",
			func(host string, port nat.Port) string {
				portStr := port.Port()
				return fmt.Sprintf("postgres://postgres:password@%s:%s/shop_test?sslmode=disable", host, portStr)
			}),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://postgres:password@%s:%s/shop_test?sslmode=disable", host, port.Port())

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, initScript)
	require.NoError(t, err)

	cleanup := func() {
		pool.Close()
		_ = container.Terminate(ctx)
	}

	return pool, cleanup
}
