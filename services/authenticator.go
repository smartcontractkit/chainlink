package services

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-go/logger"
	"golang.org/x/crypto/ssh/terminal"
)

func Authenticate(store *Store) {
	if store.KeyStore.HasAccounts() {
		checkPassword(store)
	} else {
		createAccount(store)
	}
}

func checkPassword(store *Store) {
	for {
		ok := true
		phrase := promptPassword("Enter Password:")
		for _, account := range store.KeyStore.Accounts() {
			err := store.KeyStore.Unlock(account, phrase)
			if err != nil {
				fmt.Printf("Invalid password for account: %s\nPlease try again.\n\n", account.Address.Hex())
				ok = false
			}
		}
		if ok {
			printGreeting()
			break
		}
	}
}

func createAccount(store *Store) {
	for {
		phrase := promptPassword("New Password:")
		phraseConfirmation := promptPassword("Confirm Password: ")
		if phrase == phraseConfirmation {
			_, err := store.KeyStore.NewAccount(phrase)
			if err != nil {
				logger.Fatal(err)
			}
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
