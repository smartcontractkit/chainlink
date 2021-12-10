package cmd

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewTerminalPrompter,
	NewPasswordPrompter,
	NewChangePasswordPrompter,
	NewFileSessionRequestBuilder,
	NewPromptingSessionRequestBuilder,
	NewPromptingAPIInitializer,
	NewSessionCookieAuthenticator,
	NewClient,
)
