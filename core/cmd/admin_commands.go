package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type AdminUsersPresenter struct {
	JAID
	presenters.UserResource
}

var adminUsersTableHeaders = []string{"Email", "Role", "Has API token", "Created At", "Updated at"}

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
	response, err := cli.HTTP.Delete(fmt.Sprintf("/v2/users/%s", c.String("email")))
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
