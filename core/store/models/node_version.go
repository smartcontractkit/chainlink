package models

import "time"

type NodeVersion struct {
	Version   string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"index"`
}

func NewNodeVersion(version string) NodeVersion {
	return NodeVersion{
		Version:   version,
		CreatedAt: time.Now(),
	}
}
