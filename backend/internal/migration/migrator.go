package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Migration struct {
	Version string
	Name    string
	SQL     string
}

type Migrator struct {
	db         *sql.DB
	migrations []Migration
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db: db,
	}
}

func (m *Migrator) AddMigration(version, name, sql string) {
	m.migrations = append(m.migrations, Migration{
		Version: version,
		Name:    name,
		SQL:     sql,
	})
}

func (m *Migrator) LoadMigrations(migrationDir string) error {
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		parts := strings.Split(strings.TrimSuffix(file.Name(), ".sql"), "_")
		if len(parts) < 2 {
			continue
		}

		version := parts[0]
		name := strings.Join(parts[1:], "_")

		content, err := os.ReadFile(filepath.Join(migrationDir, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		m.AddMigration(version, name, string(content))
	}

	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	return nil
}

func (m *Migrator) Up() error {
	if err := m.createMigrationTable(); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	for _, migration := range m.migrations {
		applied, err := m.isMigrationApplied(migration.Version)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if applied {
			continue
		}

		if err := m.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		fmt.Printf("Applied migration: %s_%s\n", migration.Version, migration.Name)
	}

	return nil
}

func (m *Migrator) createMigrationTable() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func (m *Migrator) isMigrationApplied(version string) (bool, error) {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = $1", version).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *Migrator) applyMigration(migration Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(migration.SQL); err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO migrations (version, name, applied_at) VALUES ($1, $2, $3)",
		migration.Version, migration.Name, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}
