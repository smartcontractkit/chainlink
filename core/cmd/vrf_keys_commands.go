package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"go.uber.org/multierr"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/urfave/cli"
)

type VRFKeyPresenter struct {
	JAID // Include this to overwrite the presenter JAID so it can correctly render the ID in JSON
	presenters.VRFKeyResource
}

// RenderTable implements TableRenderer
func (p *VRFKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"Compressed", "Uncompressed", "Hash"}
	rows := [][]string{p.ToRow()}
	renderList(headers, rows, rt.Writer)
	_, err := rt.Write([]byte("\n"))
	return err
}

func (p *VRFKeyPresenter) ToRow() []string {
	return []string{
		p.Compressed,
		p.Uncompressed,
		p.Hash,
	}
}

type VRFKeyPresenters []VRFKeyPresenter

// RenderTable implements TableRenderer
func (ps VRFKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"Compressed", "Uncompressed", "Hash"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(headers, rows, rt.Writer)
	_, err := rt.Write([]byte("\n"))
	return err
}

// CreateVRFKey creates a key in the VRF keystore, protected by the password in
// the vrf password file provided when starting the chainlink node.
func (cli *Client) CreateVRFKey(c *cli.Context) error {
	resp, err := cli.HTTP.Post("/v2/keys/vrf", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter VRFKeyPresenter
	return cli.renderAPIResponse(resp, &presenter)
}

// ImportVRFKey reads a file into an EncryptedVRFKey in the db
func (cli *Client) ImportVRFKey(c *cli.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the filepath of the key to be imported"))
	}

	oldPasswordFile := c.String("oldpassword")
	if len(oldPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --oldpassword/-p flag"))
	}
	oldPassword, err := ioutil.ReadFile(oldPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.Args().Get(0)
	keyJSON, err := ioutil.ReadFile(filepath)
	if err != nil {
		return cli.errorOut(err)
	}

	normalizedPassword := normalizePassword(string(oldPassword))
	resp, err := cli.HTTP.Post("/v2/keys/vrf/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter VRFKeyPresenter
	return cli.renderAPIResponse(resp, &presenter, "Imported VRF key")
}

// ExportVRFKey saves encrypted copy of VRF key with given public key to
// requested file path.
func (cli *Client) ExportVRFKey(c *cli.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the ID (compressed public key) of the key to export"))
	}

	newPasswordFile := c.String("newpassword")
	if len(newPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --newpassword/-p flag"))
	}
	newPassword, err := ioutil.ReadFile(newPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.String("output")
	if len(filepath) == 0 {
		return cli.errorOut(errors.New("Must specify --output/-o flag"))
	}

	pk, err := getPublicKey(c)
	if err != nil {
		return cli.errorOut(err)
	}

	normalizedPassword := normalizePassword(string(newPassword))
	resp, err := cli.HTTP.Post("/v2/keys/vrf/export/"+pk.String()+"?newpassword="+normalizedPassword, nil)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		errResult, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return cli.errorOut(errors.Errorf("error exporting status code %d error reading body %s", resp.StatusCode, err2))
		}
		return cli.errorOut(errors.Errorf("error exporting status code %d body %s", resp.StatusCode, string(errResult)))
	}

	keyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read response body"))
	}

	err = utils.WriteFileWithMaxPerms(filepath, keyJSON, 0600)
	if err != nil {
		return cli.errorOut(errors.Wrapf(err, "Could not write %v", filepath))
	}

	_, err = os.Stderr.WriteString(fmt.Sprintf("Exported VRF key %s to %s\n", pk.String(), filepath))
	if err != nil {
		return cli.errorOut(err)
	}

	return nil
}

// DeleteVRFKey deletes (hard or soft) the VRF key with given public key from the db
// and memory. V2 jobs referencing the VRF key will be removed if the key is deleted
// (no such protection for the V1 jobs exists).
func (cli *Client) DeleteVRFKey(c *cli.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the key ID (compressed public key) to be deleted"))
	}
	id, err := getPublicKey(c)
	if err != nil {
		return cli.errorOut(err)
	}

	if !confirmAction(c) {
		return nil
	}

	var queryStr string
	if c.Bool("hard") {
		queryStr = "?hard=true"
	}

	resp, err := cli.HTTP.Delete(fmt.Sprintf("/v2/keys/vrf/%s%s", id, queryStr))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter VRFKeyPresenter
	return cli.renderAPIResponse(resp, &presenter, "VRF key deleted")
}

func getPublicKey(c *cli.Context) (secp256k1.PublicKey, error) {
	pkHexString := c.Args().Get(0)
	if pkHexString == "" {
		return secp256k1.PublicKey{}, fmt.Errorf("must specify public key")
	}
	publicKey, err := secp256k1.NewPublicKeyFromHex(pkHexString)
	if err != nil {
		return secp256k1.PublicKey{}, errors.Wrap(err, "failed to parse public key")
	}
	return publicKey, nil
}

// ListKeys Lists the keys in the db
func (cli *Client) ListVRFKeys(c *cli.Context) error {
	resp, err := cli.HTTP.Get("/v2/keys/vrf", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenters VRFKeyPresenters
	return cli.renderAPIResponse(resp, &presenters, "🔑 VRF Keys")
}
