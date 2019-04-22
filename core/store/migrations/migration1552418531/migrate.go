package migration1552418531

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
)

// Migrate creates a new bridge_types table with the correct primary key
// because sqlite does not allow you to modify the primary key
// after table creation.
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&initiator{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&jobSpec{}).Error; err != nil {
		return err
	}
	return tx.AutoMigrate(&jobRun{}).Error
}

type jobSpec struct {
	ID        string    `json:"id,omitempty" gorm:"primary_key;not null"`
	DeletedAt null.Time `json:"-" gorm:"index"`
}

type jobRun struct {
	ID        string    `json:"id" gorm:"primary_key;not null"`
	DeletedAt null.Time `json:"-" gorm:"index"`
}

type initiator struct {
	ID        uint      `json:"id" gorm:"primary_key;auto_increment"`
	DeletedAt null.Time `json:"=" gorm:"index"`
}
