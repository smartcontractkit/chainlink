package internal

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Manifest is a list of migrations
type Manifest struct {
	Migrations []Migration `json:"migrations"`
}

// Migration is a single migration
type Migration struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Direction string `json:"direction"`
}

func Load(path string) (Manifest, error) {
	m := Manifest{
		Migrations: []Migration{},
	}
	// Load migrations from git
	r, err := git.PlainOpen(path)
	if err != nil {
		return m, err
	}
	cIter, err := r.Log(&git.LogOptions{Order: git.LogOrderCommitterTime, PathFilter: func(path string) bool {
		return strings.Contains(path, "migrate/migrations") || strings.Contains(path, "migrate/template")
	}})
	if err != nil {
		return m, err
	}
	// ... just iterates over the commits, printing it
	n := 0
	err = cIter.ForEach(func(c *object.Commit) error {
		n++
		fiter, ferr := c.Files()
		if ferr != nil {
			return ferr
		}
		ferr = fiter.ForEach(func(f *object.File) error {

			if strings.HasSuffix(f.Name, "sql") {
				m.Migrations = append(m.Migrations, Migration{
					ID:   c.Hash.String(),
					Name: f.Name,
				})
				fmt.Printf("commit, name: %s, %s\n", c.Hash, f.Name)
			}
			return nil
		})
		if n > 50 {
			fmt.Println("breaking")
			return nil
		}
		return ferr
	})

	return m, err
}
