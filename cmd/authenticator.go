package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"golang.org/x/crypto/ssh/terminal"
)

type Authenticator interface {
	Authenticate(*store.Store, string)
}

type TerminalAuthenticator struct {
	Exiter func(int)
}

// Authenticate checks to see if there are accounts present in
// the KeyStore, and if there are none, a new account will be created
// by prompting for a password. If there are accounts present, the
// account which is unlocked by the given password will be used.
func (k TerminalAuthenticator) Authenticate(store *store.Store, pwd string) {
	if len(pwd) != 0 {
		k.authenticateWithPwd(store, pwd)
	} else {
		k.authenticationPrompt(store)
	}
}

func (k TerminalAuthenticator) authenticationPrompt(store *store.Store) {
	if store.KeyStore.HasAccounts() {
		promptAndCheckPassword(store)
	} else {
		createAccount(store)
	}
}

func (k TerminalAuthenticator) authenticateWithPwd(store *store.Store, pwd string) {
	if !store.KeyStore.HasAccounts() {
		fmt.Println("Cannot authenticate with password because there are no accounts")
		k.Exiter(1)
	} else if err := checkPassword(store, pwd); err != nil {
		k.Exiter(1)
	}
}

func checkPassword(store *store.Store, phrase string) error {
	if err := store.KeyStore.Unlock(phrase); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		printGreeting()
		return nil
	}
}

func promptAndCheckPassword(store *store.Store) {
	for {
		phrase := promptPassword("Enter Password:")
		if checkPassword(store, phrase) == nil {
			break
		}
	}
}

func createAccount(store *store.Store) {
	for {
		phrase := promptPassword("New Password: ")
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
