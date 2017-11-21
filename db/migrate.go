package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/gormigrate.v1"
	"log"
)

func migrate() {
	db.LogMode(true)

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "201711211530",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(`
					CREATE TABLE jobs (
						id bigint NOT NULL,
						schedule text
					);
				`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("jobs").Error
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	log.Printf("Migration did run successfully")
}
