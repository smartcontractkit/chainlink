package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/smartcontractkit/chainlink-go/web"
	"github.com/urfave/cli"
)

type Client struct {
	io.Writer
}

func (self *Client) PrettyPrintJSON(v interface{}) error {
	b, err := utils.FormatJSON(v)
	if err != nil {
		return err
	}
	if _, err = self.Write(b); err != nil {
		return err
	}
	return nil
}

func (self *Client) cliError(err error) error {
	if err != nil {
		self.Write([]byte(err.Error()))
	}
	return err
}

func (self *Client) RunNode(c *cli.Context) error {
	cl := services.NewApplication(store.NewConfig())
	services.Authenticate(cl.Store)
	r := web.Router(cl)

	if err := cl.Start(); err != nil {
		logger.Fatal(err)
	}
	defer cl.Stop()
	logger.Fatal(r.Run())
	return nil
}

func (self *Client) ShowJob(c *cli.Context) error {
	cfg := store.NewConfig()
	if !c.Args().Present() {
		return self.cliError(fmt.Errorf("Must pass the job id to be shown"))
	}
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		"http://localhost:8080/jobs/"+c.Args().First(),
	)
	if err != nil {
		return self.cliError(err)
	}
	defer resp.Body.Close()
	var job models.Job
	err = deserializeResponse(resp, &job)
	if err != nil {
		return self.cliError(err)
	}
	return self.cliError(self.PrettyPrintJSON(job))
}

func deserializeResponse(resp *http.Response, dst interface{}) error {
	if resp.StatusCode >= 300 {
		return fmt.Errorf(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &dst); err != nil {
		return err
	}
	return nil
}

func (self *Client) GetJobs(c *cli.Context) error {
	cfg := store.NewConfig()
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		"http://localhost:8080/jobs",
	)
	if err != nil {
		return self.cliError(err)
	}
	defer resp.Body.Close()

	var jobs []models.Job
	err = deserializeResponse(resp, &jobs)
	if err != nil {
		return self.cliError(err)
	}
	return self.cliError(self.PrettyPrintJSON(jobs))
}
