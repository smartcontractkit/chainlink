package cmd

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

type UserInitializer interface {
	Initialize(store *store.Store) (models.User, error)
}

type terminalUserInitializer struct {
	prompter Prompter
}

func NewTerminalUserInitializer() UserInitializer {
	return &terminalUserInitializer{prompter: NewTerminalPrompter()}
}

func (t *terminalUserInitializer) Initialize(store *store.Store) (models.User, error) {
	if user, err := store.FindUser(); err == nil {
		return user, err
	}

	for {
		email := t.prompter.Prompt("Enter API Email: ")
		pwd := t.prompter.PasswordPrompt("Enter API Password: ")
		user, err := models.NewUser(email, pwd)
		if err != nil {
			fmt.Println("Error creating API user: ", err)
			continue
		}
		if err = store.Save(&user); err != nil {
			fmt.Println("Error creating API user: ", err)
		}
		return user, err
	}
}
