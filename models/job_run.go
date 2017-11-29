package models

import (
	"time"
)

type JobRun struct {
	ID        string    `storm:"id,index,unique"`
	JobID     string    `storm:"index"`
	CreatedAt time.Time `storm:"index"`
}
