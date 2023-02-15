package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func initAdminSubCmds(client *Client) []cli.Command {
	return []cli.Command{
		{
			Name:   "chpass",
			Usage:  "Change your API password remotely",
			Action: client.ChangePassword,
		},
		{
			Name:   "login",
			Usage:  "Login to remote client by creating a session cookie",
			Action: client.RemoteLogin,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "text file holding the API email and password needed to create a session cookie",
				},
				cli.BoolFlag{
					Name:  "bypass-version-check",
					Usage: "Bypass versioning check for compatibility of remote node",
				},
			},
		},
		{
			Name:   "logout",
			Usage:  "Delete any local sessions",
			Action: client.Logout,
		},
		{
			Name:  "users",
			Usage: "Create, edit permissions, or delete API users",
			Subcommands: cli.Commands{
				{
					Name:   "list",
					Usage:  "Lists all API users and their roles",
					Action: client.ListUsers,
				},
				{
					Name:   "create",
					Usage:  "Create a new API user",
					Action: client.CreateUser,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "email",
							Usage:    "Email of new user to create",
							Required: true,
						},
						cli.StringFlag{
							Name:     "role",
							Usage:    "Permission level of new user. Options: 'admin', 'edit', 'run', 'view'.",
							Required: true,
						},
					},
				},
				{
					Name:   "chrole",
					Usage:  "Changes an API user's role",
					Action: client.ChangeRole,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "email",
							Usage:    "email of user to be editted",
							Required: true,
						},
						cli.StringFlag{
							Name:     "newrole",
							Usage:    "new permission level role to set for user. Options: 'admin', 'edit', 'run', 'view'.",
							Required: true,
						},
					},
				},
				{
					Name:   "delete",
					Usage:  "Delete an API user",
					Action: client.DeleteUser,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "email",
							Usage:    "Email of API user to delete",
							Required: true,
						},
					},
				},
			},
		},
	}
}

type AdminUsersPresenter struct {
	JAID
	presenters.UserResource
}

var adminUsersTableHeaders = []string{"Email", "Role", "Has API token", "Created at", "Updated at"}

func (p *AdminUsersPresenter) ToRow() []string {
	row := []string{
		p.ID,
		string(p.Role),
		p.HasActiveApiToken,
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
	return row
}

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

	if _, err := rt.Write([]byte("Users\n")); err != nil {
		return err
	}
	renderList(adminUsersTableHeaders, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

// ListUsers renders all API users and their roles
func (cli *Client) ListUsers(c *cli.Context) (err error) {
	resp, err := cli.HTTP.Get("/v2/users/", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &AdminUsersPresenters{})
}

// CreateUser creates a new user by prompting for email, password, and role
func (cli *Client) CreateUser(c *cli.Context) (err error) {
	resp, err := cli.HTTP.Get("/v2/users/", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var links jsonapi.Links
	var users AdminUsersPresenters
	if err := cli.deserializeAPIResponse(resp, &users, &links); err != nil {
		return cli.errorOut(err)
	}
	for _, user := range users {
		if strings.EqualFold(user.Email, c.String("email")) {
			return cli.errorOut(fmt.Errorf("user with email %s already exists", user.Email))
		}
	}

	fmt.Println("Password of new user:")
	pwd := cli.PasswordPrompter.Prompt()

	request := struct {
		Email    string `json:"email"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}{
		Email:    c.String("email"),
		Role:     c.String("role"),
		Password: pwd,
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)
	response, err := cli.HTTP.Post("/v2/users", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(response, &AdminUsersPresenter{}, "Successfully created new API user")
}

// ChangeRole can change a user's role
func (cli *Client) ChangeRole(c *cli.Context) (err error) {
	request := struct {
		Email   string `json:"email"`
		NewRole string `json:"newRole"`
	}{
		Email:   c.String("email"),
		NewRole: c.String("newrole"),
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)
	response, err := cli.HTTP.Patch("/v2/users", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(response, &AdminUsersPresenter{}, "Successfully updated API user")
}

// DeleteUser deletes an API user by email
func (cli *Client) DeleteUser(c *cli.Context) (err error) {
	email := c.String("email")
	if email == "" {
		return cli.errorOut(errors.New("email flag is empty, must specify an email"))
	}

	response, err := cli.HTTP.Delete(fmt.Sprintf("/v2/users/%s", email))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(response, &AdminUsersPresenter{}, "Successfully deleted API user")
}
