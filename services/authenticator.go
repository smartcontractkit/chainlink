package services

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/store"
	"golang.org/x/crypto/ssh/terminal"
)

func Authenticate(store *store.Store) {
	if store.KeyStore.HasAccounts() {
		checkPassword(store)
	} else {
		createAccount(store)
	}
}

func checkPassword(store *store.Store) {
	for {
		phrase := promptPassword("Enter Password:")
		if err := store.KeyStore.Unlock(phrase); err != nil {
			fmt.Printf(err.Error())
		} else {
			printGreeting()
			break
		}
	}
}

func createAccount(store *store.Store) {
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
	var rval string
	withTerminalResetter(func() {
		fmt.Print(prompt)
		bytePwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println()
		rval = string(bytePwd)
	})
	return rval
}

// Explicitly reset terminal state in the event of a signal (CTRL+C)
// to ensure typed characters are echoed in terminal:
// https://groups.google.com/forum/#!topic/Golang-nuts/kTVAbtee9UA
func withTerminalResetter(f func()) {
	initialTermState, err := terminal.GetState(syscall.Stdin)
	if err != nil {
		logger.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		terminal.Restore(syscall.Stdin, initialTermState)
		os.Exit(1)
	}()

	f()
	signal.Stop(c)
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
