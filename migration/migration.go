package migration

import (
	"embed"

	"github.com/jmoiron/sqlx"
)

type Migration interface {
	// Root Get Root directory
	Root() string
	// Extension Get file extension
	Extension() string
	// Path get full path of file
	Path(name string) string
	// Init initialize database migration
	Init() error
	// Summary get migration summery
	Summary() (Summary, error)
	// StageSummary get migration summery for stage
	StageSummary(stage string) (Summary, error)
	// Up run migration stage
	//
	// pass only to run migrate on certain files only
	Up(stage string, only ...string) ([]string, error)
	// Down rollback migration stage
	//
	// pass only to run rollback on certain files only
	Down(stage string, only ...string) ([]string, error)
}

// NewMigration create new migration instance
func NewMigration(db *sqlx.DB, fs FS, root, ext string) Migration {
	instance := new(migration)
	instance.root = root
	instance.ext = ext
	instance.db = db
	instance.fs = fs
	instance.Init()
	return instance
}

// NewDirMigration create new migration instance from physical dir files
func NewDirMigration(db *sqlx.DB, root, ext string) (Migration, error) {
	if fs, err := NewDirFS(root, ext); err != nil {
		return nil, err
	} else {
		instance := new(migration)
		instance.root = root
		instance.ext = ext
		instance.db = db
		instance.fs = fs
		instance.Init()
		return instance, nil
	}
}

// NewEmbedMigration create new migration instance from embedded files
func NewEmbedMigration(db *sqlx.DB, fs embed.FS, root, ext string) (Migration, error) {
	if eFs, err := NewEmbedFS(fs, root, ext); err != nil {
		return nil, err
	} else {
		instance := new(migration)
		instance.root = root
		instance.ext = ext
		instance.db = db
		instance.fs = eFs
		instance.Init()
		return instance, nil
	}
}
