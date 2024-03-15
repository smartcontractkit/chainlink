package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initVRFKeysSubCmd(s *Shell) cli.Command {
	return cli.Command{
		Name:  "vrf",
		Usage: "Remote commands for administering the node's vrf keys",
		Subcommands: cli.Commands{
			{
				Name:   "create",
				Usage:  "Create a VRF key",
				Action: s.CreateVRFKey,
			},
			{
				Name:  "import",
				Usage: "Import VRF key from keyfile",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "old-password, oldpassword, p",
						Usage: "`FILE` containing the password used to encrypt the key in the JSON file",
					},
				},
				Action: s.ImportVRFKey,
			},
			{
				Name:  "export",
				Usage: "Export VRF key to keyfile",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "new-password, newpassword, p",
						Usage: "`FILE` containing the password to encrypt the key (required)",
					},
					cli.StringFlag{
						Name:  "output, o",
						Usage: "`FILE` where the JSON file will be saved (required)",
					},
				},
				Action: s.ExportVRFKey,
			},
			{
				Name: "delete",
				Usage: "Archive or delete VRF key from memory and the database, if present. " +
					"Note that jobs referencing the removed key will also be removed.",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "publicKey, pk"},
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "skip the confirmation prompt",
					},
					cli.BoolFlag{
						Name:  "hard",
						Usage: "hard-delete the key instead of archiving (irreversible!)",
					},
				},
				Action: s.DeleteVRFKey,
			},
			{
				Name: "list", Usage: "List the VRF keys",
				Action: s.ListVRFKeys,
			},
		},
	}
}

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
func (s *Shell) CreateVRFKey(_ *cli.Context) error {
	resp, err := s.HTTP.Post(s.ctx(), "/v2/keys/vrf", nil)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter VRFKeyPresenter
	return s.renderAPIResponse(resp, &presenter)
}

// ImportVRFKey reads a file into an EncryptedVRFKey in the db
func (s *Shell) ImportVRFKey(c *cli.Context) error {
	if !c.Args().Present() {
		return s.errorOut(errors.New("Must pass the filepath of the key to be imported"))
	}

	oldPasswordFile := c.String("old-password")
	if len(oldPasswordFile) == 0 {
		return s.errorOut(errors.New("Must specify --old-password/-p flag"))
	}
	oldPassword, err := os.ReadFile(oldPasswordFile)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.Args().Get(0)
	keyJSON, err := os.ReadFile(filepath)
	if err != nil {
		return s.errorOut(err)
	}

	normalizedPassword := normalizePassword(string(oldPassword))
	resp, err := s.HTTP.Post(s.ctx(), "/v2/keys/vrf/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter VRFKeyPresenter
	return s.renderAPIResponse(resp, &presenter, "Imported VRF key")
}

// ExportVRFKey saves encrypted copy of VRF key with given public key to
// requested file path.
func (s *Shell) ExportVRFKey(c *cli.Context) error {
	if !c.Args().Present() {
		return s.errorOut(errors.New("Must pass the ID (compressed public key) of the key to export"))
	}

	newPasswordFile := c.String("new-password")
	if len(newPasswordFile) == 0 {
		return s.errorOut(errors.New("Must specify --new-password/-p flag"))
	}
	newPassword, err := os.ReadFile(newPasswordFile)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.String("output")
	if len(filepath) == 0 {
		return s.errorOut(errors.New("Must specify --output/-o flag"))
	}

	pk, err := getPublicKey(c)
	if err != nil {
		return s.errorOut(err)
	}

	normalizedPassword := normalizePassword(string(newPassword))
	resp, err := s.HTTP.Post(s.ctx(), "/v2/keys/vrf/export/"+pk.String()+"?newpassword="+normalizedPassword, nil)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return s.errorOut(fmt.Errorf("error exporting: %w", httpError(resp)))
	}

	keyJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not read response body"))
	}

	err = utils.WriteFileWithMaxPerms(filepath, keyJSON, 0o600)
	if err != nil {
		return s.errorOut(errors.Wrapf(err, "Could not write %v", filepath))
	}

	_, err = os.Stderr.WriteString(fmt.Sprintf("Exported VRF key %s to %s\n", pk.String(), filepath))
	if err != nil {
		return s.errorOut(err)
	}

	return nil
}

// DeleteVRFKey deletes (hard or soft) the VRF key with given public key from the db
// and memory. V2 jobs referencing the VRF key will be removed if the key is deleted
// (no such protection for the V1 jobs exists).
func (s *Shell) DeleteVRFKey(c *cli.Context) error {
	if !c.Args().Present() {
		return s.errorOut(errors.New("Must pass the key ID (compressed public key) to be deleted"))
	}
	id, err := getPublicKey(c)
	if err != nil {
		return s.errorOut(err)
	}

	if !confirmAction(c) {
		return nil
	}

	var queryStr string
	if c.Bool("hard") {
		queryStr = "?hard=true"
	}

	resp, err := s.HTTP.Delete(s.ctx(), fmt.Sprintf("/v2/keys/vrf/%s%s", id, queryStr))
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter VRFKeyPresenter
	return s.renderAPIResponse(resp, &presenter, "VRF key deleted")
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
func (s *Shell) ListVRFKeys(_ *cli.Context) error {
	resp, err := s.HTTP.Get(s.ctx(), "/v2/keys/vrf", nil)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenters VRFKeyPresenters
	return s.renderAPIResponse(resp, &presenters, "🔑 VRF Keys")
}
