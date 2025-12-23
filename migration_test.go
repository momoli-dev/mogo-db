package database_test

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/momoli-dev/mogo/database"
)

func testMigrationFS() fstest.MapFS {
	return fstest.MapFS{
		"migration/00001_create_migration_test.sql": &fstest.MapFile{
			Data: []byte(`
-- +goose Up
-- +goose StatementBegin

CREATE TABLE migration_test (
	id serial PRIMARY KEY,
	name text NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE migration_test;
-- +goose StatementEnd
`),
		},
		"migration/00002_create_migration_test_2.sql": &fstest.MapFile{
			Data: []byte(`
-- +goose Up
-- +goose StatementBegin
CREATE TABLE migration_test_2 (
	id serial PRIMARY KEY,
	name text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE migration_test_2;
-- +goose StatementEnd
`),
		},
		"migration/00003_add_column_to_migration_test_2.sql": &fstest.MapFile{
			Data: []byte(`
-- +goose Up
-- +goose StatementBegin
ALTER TABLE migration_test_2 ADD COLUMN age integer;

INSERT INTO migration_test_2 (name, age) VALUES ('test', 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE migration_test_2 DROP COLUMN age;
-- +goose StatementEnd
`),
		},
	}
}

func TestMigrateUpAll(t *testing.T) {
	testFS := testMigrationFS()
	RunWithConn(t, func(conn *database.Conn) {
		err := database.MigrateUpAll(conn, testFS)
		require.NoError(t, err)

		rows1, err := conn.Handle().Query("SELECT * FROM migration_test")
		require.NoError(t, err)
		require.NotNil(t, rows1)

		rows2, err := conn.Handle().Query("SELECT * FROM migration_test_2")
		require.NoError(t, err)
		require.NotNil(t, rows2)

		hasColumn := false
		for rows2.Next() {
			var id int
			var name string
			var age int
			err := rows2.Scan(&id, &name, &age)
			require.NoError(t, err)
			require.Equal(t, 0, age)
			hasColumn = true
		}
		require.True(t, hasColumn)
	})
	ResetDB(t)
}

func TestMigrateDownAll(t *testing.T) {
	testFS := testMigrationFS()
	RunWithConn(t, func(conn *database.Conn) {
		err := database.MigrateDownAll(conn, testFS)
		require.NoError(t, err)
	})
	ResetDB(t)
}

func TestMigrateUpThenDown(t *testing.T) {
	testFS := testMigrationFS()
	RunWithConn(t, func(conn *database.Conn) {
		err := database.MigrateUpAll(conn, testFS)
		require.NoError(t, err)

		err = database.MigrateDownAll(conn, testFS)
		require.NoError(t, err)

		_, err = conn.Handle().Query("SELECT * FROM migration_test")
		require.Error(t, err)

		_, err = conn.Handle().Query("SELECT * FROM migration_test_2")
		require.Error(t, err)
	})
	ResetDB(t)
}
