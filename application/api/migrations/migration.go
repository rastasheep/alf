package migrations

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Migration struct {
	Version int
	Scripts []string
}

var migrations = []*Migration{
	{
		Version: 1,
		Scripts: []string{
			`create table executions (
				id      serial primary key,
				query   text not null,

				created_at timestamp not null default current_timestamp
			)`,
		},
	},
	{
		Version: 2,
		Scripts: []string{
			`create table templates (
				id      serial primary key,
				query   text not null,

				updated_at timestamp not null default current_timestamp,
				created_at timestamp not null default current_timestamp
			)`,
		},
	},
}

func createVersionTable(db *sqlx.DB) error {
	_, err := db.Exec(`
	do $$ begin
		create table if not exists versions (
			version int not null unique,
			updated_at timestamp not null default current_timestamp
		);
		if not exists (select 1 from versions where version = 0) then
			insert into versions (version) values (0);
		end if;
	end; $$;`)

	if err != nil {
		return fmt.Errorf("Failed to create Version table: %v", err)
	}
	return nil
}

func applyMigration(db *sqlx.DB, migration *Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, q := range migration.Scripts {
		if _, err := tx.Exec(q); err != nil {
			return fmt.Errorf("Migration to version %d failed at step %d: %v", migration.Version, i+1, err)
		}
	}

	_, err = tx.Exec(`insert into versions (version) values ($1)`, migration.Version)
	if err != nil {
		return fmt.Errorf("Version update failed: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Migration to version %d failed: %v", migration.Version, err)
	}

	return nil
}

func Exec(db *sqlx.DB) error {
	if err := createVersionTable(db); err != nil {
		return err
	}

	version := 0
	if err := db.QueryRow(`select max(version) from versions`).Scan(&version); err != nil {
		return err
	}

	for _, migration := range migrations {
		if migration.Version > version {
			err := applyMigration(db, migration)
			if err != nil {
				return err
			}
			version = migration.Version
		}
	}
	return nil
}
