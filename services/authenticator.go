package services

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/models"
	"golang.org/x/crypto/ssh/terminal"
)

func Authenticate(store *Store) {
	password, ok := getExistingPassword(store)
	if ok {
		checkPassword(password)
	} else {
		createPassword(store)
	}
}

func getExistingPassword(store *Store) (models.Password, bool) {
	var passwords []models.Password
	if err := store.All(&passwords); err != nil {
		logger.Fatal(err)
	}

	if len(passwords) == 0 {
		return models.Password{}, false
	} else {
		return passwords[0], true
	}
}

func checkPassword(password models.Password) {
	for {
		if password.Check(promptPassword("Enter Password:")) {
			printGreeting()
			break
		}
	}
}

func createPassword(store *Store) {
	for {
		phrase := promptPassword("New Password:")
		phraseConfirmation := promptPassword("Confirm Password: ")
		if phrase == phraseConfirmation {
			store.AddPassword(models.NewPassword(phrase))
			store.CreateKey(phrase)
			printGreeting()
			break
		} else {
			fmt.Println("Passwords don't match. Please try again.")
		}
	}
}

func promptPassword(prompt string) string {
	fmt.Print(prompt)
	bytePwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println()
	return string(bytePwd)
}

func printGreeting() {
	fmt.Println(`
     _/_/_/  _/                  _/            _/        _/            _/
  _/        _/_/_/      _/_/_/      _/_/_/    _/            _/_/_/    _/  _/
 _/        _/    _/  _/    _/  _/  _/    _/  _/        _/  _/    _/  _/_/
_/        _/    _/  _/    _/  _/  _/    _/  _/        _/  _/    _/  _/  _/
 _/_/_/  _/    _/    _/_/_/  _/  _/    _/  _/_/_/_/  _/  _/    _/  _/    _/
`)
}
