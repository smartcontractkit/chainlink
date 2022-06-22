package cmd

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"
)

type AdminUsersPresenter struct {
	Email     string
	Role      string
	CreatedAt time.Time
	TokenKey  null.String
	UpdatedAt time.Time
}

func (p *AdminUsersPresenter) ToRow() []string {
	hasToken := "false"
	if p.TokenKey.Valid {
		hasToken = "true"
	}
	return []string{
		p.Email,
		p.Role,
		hasToken,
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
}

var adminUsersTableHeaders = []string{"Email", "Role", "Has API token", "Created At", "Updated at"}

// RenderTable implements TableRenderer
func (p *AdminUsersPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{p.ToRow()}

	renderList(adminUsersTableHeaders, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

type AdminUsersPresenters []AdminUsersPresenter

// RenderTable implements TableRenderer
func (ps AdminUsersPresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(adminUsersTableHeaders, rows, rt.Writer)

	return nil
}

// ListUsers renders the active account address with its ETH & LINK balance
func (cli *Client) ListUsers(c *cli.Context) (err error) {
	// TODO: Andrew - STUBBED
	return nil
}

// CreateUser creates a new user by prompting for email, password, and role
func (cli *Client) CreateUser(c *cli.Context) (err error) {
	createUrl := url.URL{
		Path: "/v2/users",
	}
	query := createUrl.Query()

	if c.IsSet("email") {
		// TODO: Andrew - from cli prompter, p.prompter.Prompt("Enter email: ") - see core/cmd/client.go
		query.Set("email", c.String("email"))
	}
	if c.IsSet("password") {
		// TODO: Andrew - from cli prompter, p.prompter.Prompt("Enter email: ") - see core/cmd/client.go
		query.Set("password", c.String("password"))
	}
	if c.IsSet("role") {
		// TODO: Andrew - from cli prompter, p.prompter.Prompt("Enter email: ") - see core/cmd/client.go
		query.Set("role", c.String("role"))
	}

	createUrl.RawQuery = query.Encode()
	resp, err := cli.HTTP.Post(createUrl.String(), nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &AdminUsersPresenter{}, "User created successfully.")
}

// EditUser can change a user's email, password, and role
func (cli *Client) EditUser(c *cli.Context) (err error) {
	createUrl := url.URL{
		Path: "/v2/users",
	}
	query := createUrl.Query()

	if c.IsSet("email") {
		// TODO: Andrew - from cli prompter, p.prompter.Prompt("Enter email: ") - see core/cmd/client.go
		query.Set("email", c.String("email"))
	}
	if c.IsSet("password") {
		// TODO: Andrew - from cli prompter, p.prompter.Prompt("Enter email: ") - see core/cmd/client.go
		query.Set("password", c.String("password"))
	}
	if c.IsSet("role") {
		// TODO: Andrew - from cli prompter, p.prompter.Prompt("Enter email: ") - see core/cmd/client.go
		query.Set("role", c.String("role"))
	}

	createUrl.RawQuery = query.Encode()
	resp, err := cli.HTTP.Patch(createUrl.String(), nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &AdminUsersPresenter{}, "User updated successfully.")
}

// DeleteUser can change a user's email, password, and role
func (cli *Client) DeleteUser(c *cli.Context) (err error) {
	createUrl := url.URL{
		Path: "/v2/users",
	}
	query := createUrl.Query()

	if c.IsSet("email") {
		// TODO: Andrew - from cli prompter, p.prompter.Prompt("Enter email: ") - see core/cmd/client.go
		query.Set("email", c.String("email"))
	}

	createUrl.RawQuery = query.Encode()
	resp, err := cli.HTTP.Delete(createUrl.String())
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &AdminUsersPresenter{}, "User updated successfully.")
}
