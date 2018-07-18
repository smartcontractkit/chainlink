package cmd

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// APIInitializer is the interface used to create the API User credentials
// needed to access the API.
type APIInitializer interface {
	Initialize(store *store.Store) (models.User, error)
}

type terminalAPIInitializer struct {
	prompter Prompter
}

// NewTerminalAPIInitializer creates a concrete instance of APIInitializer
// that uses the terminal to solicit credentials from the user.
func NewTerminalAPIInitializer() APIInitializer {
	return &terminalAPIInitializer{prompter: NewTerminalPrompter()}
}

// Initialize uses the terminal to get credentials from the user that it then saves in the
// store.
func (t *terminalAPIInitializer) Initialize(store *store.Store) (models.User, error) {
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
