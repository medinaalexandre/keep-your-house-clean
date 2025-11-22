package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const schemaMigrationsTable = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    executed_at TIMESTAMP NOT NULL DEFAULT NOW()
);
`

type Migrator struct {
	db         *sql.DB
	migrations embed.FS
}

func NewMigrator(db *sql.DB, migrations embed.FS) *Migrator {
	return &Migrator{
		db:         db,
		migrations: migrations,
	}
}

func (m *Migrator) Run() error {
	if err := m.createSchemaMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	executedMigrations, err := m.getExecutedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	migrationFiles, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	executedCount := 0
	skippedCount := 0

	for _, migrationFile := range migrationFiles {
		version := m.extractVersion(migrationFile)
		if _, executed := executedMigrations[version]; executed {
			log.Printf("Migration %s already executed, skipping", version)
			skippedCount++
			continue
		}

		log.Printf("Executing migration %s: %s", version, migrationFile)
		if err := m.executeMigration(migrationFile, version); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationFile, err)
		}
		executedCount++
		log.Printf("Migration %s executed successfully", version)
	}

	if executedCount > 0 {
		log.Printf("Migrations executed: %d new, %d skipped", executedCount, skippedCount)
	} else {
		log.Printf("All migrations up to date: %d migrations already executed", skippedCount)
	}

	return nil
}

func (m *Migrator) createSchemaMigrationsTable() error {
	_, err := m.db.Exec(schemaMigrationsTable)
	return err
}

func (m *Migrator) getExecutedMigrations() (map[string]bool, error) {
	rows, err := m.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		executed[version] = true
	}

	return executed, rows.Err()
}

func (m *Migrator) getMigrationFiles() ([]string, error) {
	var files []string

	err := fs.WalkDir(m.migrations, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".sql" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		versionI := m.extractVersionNumber(files[i])
		versionJ := m.extractVersionNumber(files[j])
		return versionI < versionJ
	})

	return files, nil
}

func (m *Migrator) extractVersion(filename string) string {
	base := filepath.Base(filename)
	parts := strings.Split(base, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return base
}

func (m *Migrator) extractVersionNumber(filename string) int {
	version := m.extractVersion(filename)
	num, err := strconv.Atoi(version)
	if err != nil {
		return 0
	}
	return num
}

func (m *Migrator) executeMigration(filename, version string) error {
	content, err := m.migrations.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (version) VALUES ($1)",
		version,
	); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

