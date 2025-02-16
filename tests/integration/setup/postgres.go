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
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPostgres(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	networkName := "test-network"
	network, err := createNetwork(ctx, networkName)
	require.NoError(t, err)

	container, dsn, err := startPostgresContainer(ctx, networkName)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	dsnForMigrations := "postgres://postgres:password@postgres:5432/shop_test?sslmode=disable"
	require.NoError(t, runMigrations(ctx, dsnForMigrations, networkName))

	cleanup := func() {
		pool.Close()
		_ = container.Terminate(ctx)
		_ = network.Remove(ctx)
	}

	return pool, cleanup
}

func createNetwork(ctx context.Context, networkName string) (testcontainers.Network, error) {
	return testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
		},
	})
}

func startPostgresContainer(ctx context.Context, networkName string) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "shop_test",
		},
		WaitingFor: wait.ForSQL("5432/tcp", "pgx", func(host string, port nat.Port) string {
			return fmt.Sprintf("postgres://postgres:password@%s:%s/shop_test?sslmode=disable", host, port.Port())
		}),
		Networks:       []string{networkName},
		NetworkAliases: map[string][]string{networkName: {"postgres"}},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", err
	}

	dsn := fmt.Sprintf("postgres://postgres:password@%s:%s/shop_test?sslmode=disable", host, port.Port())
	return container, dsn, nil
}

func runMigrations(ctx context.Context, dsn, networkName string) error {
	migrationsPath := getMigrationsPath()

	migrateReq := testcontainers.ContainerRequest{
		Image: "migrate/migrate:v4.15.2",
		Mounts: testcontainers.Mounts(
			testcontainers.ContainerMount{
				Source: testcontainers.GenericBindMountSource{HostPath: migrationsPath},
				Target: "/migrations",
			},
		),
		WaitingFor: wait.ForExit().WithExitTimeout(30 * time.Second),
		Cmd:        []string{"-path", "/migrations", "-database", dsn, "up"},
		Networks:   []string{networkName},
	}

	migrateContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: migrateReq,
		Started:          true,
	})
	if err != nil {
		return err
	}
	defer migrateContainer.Terminate(ctx)

	return nil
}

func getMigrationsPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(wd, "../../..", "migrations")
}
