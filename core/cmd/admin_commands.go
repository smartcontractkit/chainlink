package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
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
			Name:   "profile",
			Usage:  "Collects profile metrics from the node.",
			Action: client.Profile,
			Flags: []cli.Flag{
				cli.Uint64Flag{
					Name:  "seconds, s",
					Usage: "duration of profile capture",
					Value: 8,
				},
				cli.StringFlag{
					Name:  "output_dir, o",
					Usage: "output directory of the captured profile",
					Value: "/tmp/",
				},
			},
		},
		{
			Name:   "status",
			Usage:  "Displays the health of various services running inside the node.",
			Action: client.Status,
			Flags:  []cli.Flag{},
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
							Usage:    "email of user to be edited",
							Required: true,
						},
						cli.StringFlag{
							Name:     "new-role, newrole",
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
		NewRole: c.String("new-role"),
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

// Status will display the health of various services
func (cli *Client) Status(c *cli.Context) error {
	resp, err := cli.HTTP.Get("/health?full=1", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &HealthCheckPresenters{})
}

// Profile will collect pprof metrics and store them in a folder.
func (cli *Client) Profile(c *cli.Context) error {
	seconds := c.Uint("seconds")
	baseDir := c.String("output_dir")

	genDir := filepath.Join(baseDir, fmt.Sprintf("debuginfo-%s", time.Now().Format(time.RFC3339)))

	err := os.Mkdir(genDir, 0o755)
	if err != nil {
		return cli.errorOut(err)
	}
	var wgPprof sync.WaitGroup
	vitals := []string{
		"allocs",       // A sampling of all past memory allocations
		"block",        // Stack traces that led to blocking on synchronization primitives
		"cmdline",      // The command line invocation of the current program
		"goroutine",    // Stack traces of all current goroutines
		"heap",         // A sampling of memory allocations of live objects.
		"mutex",        // Stack traces of holders of contended mutexes
		"profile",      // CPU profile.
		"threadcreate", // Stack traces that led to the creation of new OS threads
		"trace",        // A trace of execution of the current program.
	}
	wgPprof.Add(len(vitals))
	cli.Logger.Infof("Collecting profiles: %v", vitals)
	cli.Logger.Infof("writing debug info to %s", genDir)

	errs := make(chan error, len(vitals))
	for _, vt := range vitals {
		go func(vt string) {
			defer wgPprof.Done()
			uri := fmt.Sprintf("/v2/debug/pprof/%s?seconds=%d", vt, seconds)
			resp, err := cli.HTTP.Get(uri)
			if err != nil {
				errs <- fmt.Errorf("error collecting %s: %w", vt, err)
				return
			}
			defer func() {
				if resp.Body != nil {
					resp.Body.Close()
				}
			}()
			if resp.StatusCode == http.StatusUnauthorized {
				errs <- fmt.Errorf("error collecting %s: %w", vt, errUnauthorized)
				return
			}
			if resp.StatusCode == http.StatusBadRequest {
				// best effort to interpret the underlying problem
				pprofVersion := resp.Header.Get("X-Go-Pprof")
				if pprofVersion == "1" {
					b, err := io.ReadAll(resp.Body)
					if err != nil {
						errs <- fmt.Errorf("error collecting %s: %w", vt, errBadRequest)
						return
					}
					respContent := string(b)
					// taken from pprof.Profile https://github.com/golang/go/blob/release-branch.go1.20/src/net/http/pprof/pprof.go#L133
					if strings.Contains(respContent, "profile duration exceeds server's WriteTimeout") {
						errs <- fmt.Errorf("%w: %s", ErrProfileTooLong, respContent)
					} else {
						errs <- fmt.Errorf("error collecting %s: %w: %s", vt, errBadRequest, respContent)
					}
				} else {
					errs <- fmt.Errorf("error collecting %s: %w", vt, errBadRequest)
				}
				return
			}
			// write to file
			f, err := os.Create(filepath.Join(genDir, vt))
			if err != nil {
				errs <- fmt.Errorf("error creating file for %s: %w", vt, err)
				return
			}
			wc := utils.NewDeferableWriteCloser(f)
			defer wc.Close()

			_, err = io.Copy(wc, resp.Body)
			if err != nil {
				errs <- fmt.Errorf("error writing to file for %s: %w", vt, err)
				return
			}
			err = wc.Close()
			if err != nil {
				errs <- fmt.Errorf("error closing file for %s: %w", vt, err)
				return
			}
		}(vt)
	}
	wgPprof.Wait()
	close(errs)
	// Atmost one err is emitted per vital.
	cli.Logger.Infof("collected %d/%d profiles", len(vitals)-len(errs), len(vitals))
	if len(errs) > 0 {
		var merr error
		for err := range errs {
			merr = errors.Join(merr, err)
		}
		return cli.errorOut(fmt.Errorf("profile collection failed:\n%v", merr))
	}
	return nil
}
