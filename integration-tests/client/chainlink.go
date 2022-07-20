// Package client enables interaction with APIs of test components like the mockserver and Chainlink nodes
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	chainlinkChart "github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"golang.org/x/sync/errgroup"
)

// OneLINK representation of a single LINK token
var OneLINK = big.NewFloat(1e18)

// Chainlink interface that enables interactions with a chainlink node
type Chainlink interface {
	URL() string
	CreateJob(spec JobSpec) (*Job, error)
	CreateJobRaw(spec string) (*Job, error)
	ReadJobs() (*ResponseSlice, error)
	ReadJob(id string) (*Response, error)
	DeleteJob(id string) error

	CreateSpec(spec string) (*Spec, error)
	ReadSpec(id string) (*Response, error)
	DeleteSpec(id string) error

	CreateBridge(bta *BridgeTypeAttributes) error
	ReadBridge(name string) (*BridgeType, error)
	DeleteBridge(name string) error

	ReadRunsByJob(jobID string) (*JobRunsResponse, error)

	CreateOCRKey() (*OCRKey, error)
	ReadOCRKeys() (*OCRKeys, error)
	DeleteOCRKey(id string) error

	CreateOCR2Key(chain string) (*OCR2Key, error)
	ReadOCR2Keys() (*OCR2Keys, error)
	DeleteOCR2Key(id string) error

	CreateP2PKey() (*P2PKey, error)
	ReadP2PKeys() (*P2PKeys, error)
	DeleteP2PKey(id int) error

	ReadETHKeys() (*ETHKeys, error)
	ReadPrimaryETHKey() (*ETHKeyData, error)
	PrimaryEthAddress() (string, error)
	UpdateEthKeyMaxGasPriceGWei(keyId string, gwei int) (*ETHKey, error)

	CreateTxKey(chain string) (*TxKey, error)
	ReadTxKeys(chain string) (*TxKeys, error)
	DeleteTxKey(chain, id string) error

	ReadTransactionAttempts() (*ctfClient.TransactionsData, error)
	ReadTransactions() (*ctfClient.TransactionsData, error)
	SendNativeToken(amount *big.Int, fromAddress, toAddress string) (ctfClient.TransactionData, error)

	CreateVRFKey() (*VRFKey, error)
	ReadVRFKeys() (*VRFKeys, error)
	ExportVRFKey(keyId string) (*VRFExportKey, error)
	ImportVRFKey(vrfExportKey *VRFExportKey) (*VRFKey, error)

	CreateCSAKey() (*CSAKey, error)
	ReadCSAKeys() (*CSAKeys, error)

	CreateEI(eia *EIAttributes) (*EIKeyCreate, error)
	ReadEIs() (*EIKeys, error)
	DeleteEI(name string) error

	CreateTerraChain(node *TerraChainAttributes) (*TerraChainCreate, error)
	CreateTerraNode(node *TerraNodeAttributes) (*TerraNodeCreate, error)

	CreateSolanaChain(node *SolanaChainAttributes) (*SolanaChainCreate, error)
	CreateSolanaNode(node *SolanaNodeAttributes) (*SolanaNodeCreate, error)

	RemoteIP() string
	SetSessionCookie() error

	SetPageSize(size int)

	Profile(profileTime time.Duration, profileFunction func(Chainlink)) (*ctfClient.ChainlinkProfileResults, error)

	// SetClient is used for testing
	SetClient(client *http.Client)
}

type chainlink struct {
	*ctfClient.BasicHTTPClient
	Config            *ChainlinkConfig
	pageSize          int
	primaryEthAddress string
}

// NewChainlink creates a new chainlink model using a provided config
func NewChainlink(c *ChainlinkConfig, httpClient *http.Client) (Chainlink, error) {
	cl := &chainlink{
		Config:          c,
		BasicHTTPClient: ctfClient.NewBasicHTTPClient(httpClient, c.URL),
		pageSize:        25,
	}
	return cl, cl.SetSessionCookie()
}

// URL chainlink instance http url
func (c *chainlink) URL() string {
	return c.Config.URL
}

// CreateJobRaw creates a Chainlink job based on the provided spec string
func (c *chainlink) CreateJobRaw(spec string) (*Job, error) {
	job := &Job{}
	log.Info().Str("Node URL", c.Config.URL).Str("Job Body", spec).Msg("Creating Job")
	err := c.do(http.MethodPost, "/v2/jobs", &JobForm{
		TOML: spec,
	}, &job, http.StatusOK)
	return job, err
}

// CreateJob creates a Chainlink job based on the provided spec struct
func (c *chainlink) CreateJob(spec JobSpec) (*Job, error) {
	job := &Job{}
	specString, err := spec.String()
	if err != nil {
		return nil, err
	}
	log.Info().Str("Node URL", c.Config.URL).Str("Type", spec.Type()).Msg("Creating Job")
	err = c.do(http.MethodPost, "/v2/jobs", &JobForm{
		TOML: specString,
	}, &job, http.StatusOK)
	return job, err
}

// ReadJobs reads all jobs from the Chainlink node
func (c *chainlink) ReadJobs() (*ResponseSlice, error) {
	specObj := &ResponseSlice{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Getting Jobs")
	return specObj, c.do(http.MethodGet, "/v2/jobs", nil, specObj, http.StatusOK)
}

// ReadJob reads a job with the provided ID from the Chainlink node
func (c *chainlink) ReadJob(id string) (*Response, error) {
	specObj := &Response{}
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Reading Job")
	return specObj, c.do(http.MethodGet, fmt.Sprintf("/v2/jobs/%s", id), nil, specObj, http.StatusOK)
}

// DeleteJob deletes a job with a provided ID from the Chainlink node
func (c *chainlink) DeleteJob(id string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Deleting Job")
	return c.do(http.MethodDelete, fmt.Sprintf("/v2/jobs/%s", id), nil, nil, http.StatusNoContent)
}

// CreateSpec creates a job spec on the Chainlink node
func (c *chainlink) CreateSpec(spec string) (*Spec, error) {
	s := &Spec{}
	r := strings.NewReplacer("\n", "", " ", "", "\\", "") // Makes it more compact and readable for logging
	log.Info().Str("Node URL", c.Config.URL).Str("Spec", r.Replace(spec)).Msg("Creating Spec")
	return s, c.doRaw(http.MethodPost, "/v2/specs", []byte(spec), s, http.StatusOK)
}

// ReadSpec reads a job spec with the provided ID on the Chainlink node
func (c *chainlink) ReadSpec(id string) (*Response, error) {
	specObj := &Response{}
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Reading Spec")
	return specObj, c.do(http.MethodGet, fmt.Sprintf("/v2/specs/%s", id), nil, specObj, http.StatusOK)
}

// ReadRunsByJob reads all runs for a job
func (c *chainlink) ReadRunsByJob(jobID string) (*JobRunsResponse, error) {
	runsObj := &JobRunsResponse{}
	log.Debug().Str("Node URL", c.Config.URL).Str("JobID", jobID).Msg("Reading runs for a job")
	return runsObj, c.do(http.MethodGet, fmt.Sprintf("/v2/jobs/%s/runs", jobID), nil, runsObj, http.StatusOK)
}

// DeleteSpec deletes a job spec with the provided ID from the Chainlink node
func (c *chainlink) DeleteSpec(id string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Deleting Spec")
	return c.do(http.MethodDelete, fmt.Sprintf("/v2/specs/%s", id), nil, nil, http.StatusNoContent)
}

// CreateBridge creates a bridge on the Chainlink node based on the provided attributes
func (c *chainlink) CreateBridge(bta *BridgeTypeAttributes) error {
	log.Info().Str("Node URL", c.Config.URL).Str("Name", bta.Name).Msg("Creating Bridge")
	return c.do(http.MethodPost, "/v2/bridge_types", bta, nil, http.StatusOK)
}

// ReadBridge reads a bridge from the Chainlink node based on the provided name
func (c *chainlink) ReadBridge(name string) (*BridgeType, error) {
	bt := BridgeType{}
	log.Info().Str("Node URL", c.Config.URL).Str("Name", name).Msg("Reading Bridge")
	return &bt, c.do(http.MethodGet, fmt.Sprintf("/v2/bridge_types/%s", name), nil, &bt, http.StatusOK)
}

// DeleteBridge deletes a bridge on the Chainlink node based on the provided name
func (c *chainlink) DeleteBridge(name string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("Name", name).Msg("Deleting Bridge")
	return c.do(http.MethodDelete, fmt.Sprintf("/v2/bridge_types/%s", name), nil, nil, http.StatusOK)
}

// CreateOCRKey creates an OCRKey on the Chainlink node
func (c *chainlink) CreateOCRKey() (*OCRKey, error) {
	ocrKey := &OCRKey{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating OCR Key")
	return ocrKey, c.do(http.MethodPost, "/v2/keys/ocr", nil, ocrKey, http.StatusOK)
}

// ReadOCRKeys reads all OCRKeys from the Chainlink node
func (c *chainlink) ReadOCRKeys() (*OCRKeys, error) {
	ocrKeys := &OCRKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading OCR Keys")
	err := c.do(http.MethodGet, "/v2/keys/ocr", nil, ocrKeys, http.StatusOK)
	for index := range ocrKeys.Data {
		ocrKeys.Data[index].Attributes.ConfigPublicKey = strings.TrimPrefix(
			ocrKeys.Data[index].Attributes.ConfigPublicKey, "ocrcfg_")
		ocrKeys.Data[index].Attributes.OffChainPublicKey = strings.TrimPrefix(
			ocrKeys.Data[index].Attributes.OffChainPublicKey, "ocroff_")
		ocrKeys.Data[index].Attributes.OnChainSigningAddress = strings.TrimPrefix(
			ocrKeys.Data[index].Attributes.OnChainSigningAddress, "ocrsad_")
	}
	return ocrKeys, err
}

// DeleteOCRKey deletes an OCRKey based on the provided ID
func (c *chainlink) DeleteOCRKey(id string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Deleting OCR Key")
	err := c.do(http.MethodDelete, fmt.Sprintf("/v2/keys/ocr/%s", id), nil, nil, http.StatusOK)
	return err
}

// CreateOCR2Key creates an OCR2Key on the Chainlink node
func (c *chainlink) CreateOCR2Key(chain string) (*OCR2Key, error) {
	ocr2Key := &OCR2Key{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating OCR2 Key")
	err := c.do(http.MethodPost, fmt.Sprintf("/v2/keys/ocr2/%s", chain), nil, ocr2Key, http.StatusOK)
	return ocr2Key, err
}

// ReadOCR2Keys reads all OCR2Keys from the Chainlink node
func (c *chainlink) ReadOCR2Keys() (*OCR2Keys, error) {
	ocr2Keys := &OCR2Keys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading OCR2 Keys")
	err := c.do(http.MethodGet, "/v2/keys/ocr2", nil, ocr2Keys, http.StatusOK)
	return ocr2Keys, err
}

// DeleteOCR2Key deletes an OCR2Key based on the provided ID
func (c *chainlink) DeleteOCR2Key(id string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Deleting OCR2 Key")
	err := c.do(http.MethodDelete, fmt.Sprintf("/v2/keys/ocr2/%s", id), nil, nil, http.StatusOK)
	return err
}

// CreateP2PKey creates an P2PKey on the Chainlink node
func (c *chainlink) CreateP2PKey() (*P2PKey, error) {
	p2pKey := &P2PKey{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating P2P Key")
	err := c.do(http.MethodPost, "/v2/keys/p2p", nil, p2pKey, http.StatusOK)
	return p2pKey, err
}

// ReadP2PKeys reads all P2PKeys from the Chainlink node
func (c *chainlink) ReadP2PKeys() (*P2PKeys, error) {
	p2pKeys := &P2PKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading P2P Keys")
	err := c.do(http.MethodGet, "/v2/keys/p2p", nil, p2pKeys, http.StatusOK)
	if len(p2pKeys.Data) == 0 {
		err = fmt.Errorf("Found no P2P Keys on the chainlink node. Node URL: %s", c.Config.URL)
		log.Err(err).Msg("Error getting P2P keys")
		return nil, err
	}
	for index := range p2pKeys.Data {
		p2pKeys.Data[index].Attributes.PeerID = strings.TrimPrefix(p2pKeys.Data[index].Attributes.PeerID, "p2p_")
	}
	return p2pKeys, err
}

// DeleteP2PKey deletes a P2PKey on the Chainlink node based on the provided ID
func (c *chainlink) DeleteP2PKey(id int) error {
	log.Info().Str("Node URL", c.Config.URL).Int("ID", id).Msg("Deleting P2P Key")
	err := c.do(http.MethodDelete, fmt.Sprintf("/v2/keys/p2p/%d", id), nil, nil, http.StatusOK)
	return err
}

// ReadETHKeys reads all ETH keys from the Chainlink node
func (c *chainlink) ReadETHKeys() (*ETHKeys, error) {
	ethKeys := &ETHKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading ETH Keys")
	err := c.do(http.MethodGet, "/v2/keys/eth", nil, ethKeys, http.StatusOK)
	if len(ethKeys.Data) == 0 {
		log.Warn().Str("Node URL", c.Config.URL).Msg("Found no ETH Keys on the node")
	}
	return ethKeys, err
}

// UpdateEthKeyMaxGasPriceGWei updates the maxGasPriceGWei for an eth key
func (c *chainlink) UpdateEthKeyMaxGasPriceGWei(keyId string, gWei int) (*ETHKey, error) {
	ethKey := &ETHKey{}
	log.Info().Str("Node URL", c.Config.URL).Str("ID", keyId).Int("maxGasPriceGWei", gWei).Msg("Update maxGasPriceGWei for eth key")
	err := c.do(http.MethodPut, fmt.Sprintf("/v2/keys/eth/%s?maxGasPriceGWei=%d", keyId, gWei), nil, ethKey, http.StatusOK)
	return ethKey, err
}

// ReadPrimaryETHKey reads updated information about the chainlink's primary ETH key
func (c *chainlink) ReadPrimaryETHKey() (*ETHKeyData, error) {
	ethKeys, err := c.ReadETHKeys()
	if err != nil {
		return nil, err
	}
	if len(ethKeys.Data) == 0 {
		return nil, fmt.Errorf("Error retrieving primary eth key on node %s: No ETH keys present", c.URL())
	}
	return &ethKeys.Data[0], nil
}

// PrimaryEthAddress returns the primary ETH address for the chainlink node
func (c *chainlink) PrimaryEthAddress() (string, error) {
	if c.primaryEthAddress == "" {
		ethKeys, err := c.ReadETHKeys()
		if err != nil {
			return "", err
		}
		c.primaryEthAddress = ethKeys.Data[0].Attributes.Address
	}
	return c.primaryEthAddress, nil
}

// CreateTxKey creates a tx key on the Chainlink node
func (c *chainlink) CreateTxKey(chain string) (*TxKey, error) {
	txKey := &TxKey{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating Tx Key")
	err := c.do(http.MethodPost, fmt.Sprintf("/v2/keys/%s", chain), nil, txKey, http.StatusOK)
	return txKey, err
}

// ReadTxKeys reads all tx keys from the Chainlink node
func (c *chainlink) ReadTxKeys(chain string) (*TxKeys, error) {
	txKeys := &TxKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading Tx Keys")
	err := c.do(http.MethodGet, fmt.Sprintf("/v2/keys/%s", chain), nil, txKeys, http.StatusOK)
	return txKeys, err
}

// DeleteTxKey deletes an tx key based on the provided ID
func (c *chainlink) DeleteTxKey(chain string, id string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("ID", id).Msg("Deleting Tx Key")
	err := c.do(http.MethodDelete, fmt.Sprintf("/v2/keys/%s/%s", chain, id), nil, nil, http.StatusOK)
	return err
}

// ReadTransactionAttempts reads all transaction attempts on the chainlink node
func (c *chainlink) ReadTransactionAttempts() (*ctfClient.TransactionsData, error) {
	txsData := &ctfClient.TransactionsData{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading Transaction Attempts")
	err := c.do(http.MethodGet, "/v2/tx_attempts", nil, txsData, http.StatusOK)
	return txsData, err
}

// ReadTransactions reads all transactions made by the chainlink node
func (c *chainlink) ReadTransactions() (*ctfClient.TransactionsData, error) {
	txsData := &ctfClient.TransactionsData{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading Transactions")
	return txsData, c.do(http.MethodGet, "/v2/transactions", nil, txsData, http.StatusOK)
}

// SendNativeToken sends native token (ETH usually) of a specified amount from one of its addresses to the target address
// WARNING: The txdata object that chainlink sends back is almost always blank.
func (c *chainlink) SendNativeToken(amount *big.Int, fromAddress, toAddress string) (ctfClient.TransactionData, error) {
	request := SendEtherRequest{
		DestinationAddress: toAddress,
		FromAddress:        fromAddress,
		Amount:             amount.String(),
		AllowHigherAmounts: true,
	}
	txData := SingleTransactionDataWrapper{}
	err := c.do(http.MethodPost, "/v2/transfers", request, txData, http.StatusOK)
	log.Info().
		Str("Node URL", c.Config.URL).
		Str("From", fromAddress).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Msg("Sending Native Token")
	return txData.Data, err
}

// ReadVRFKeys reads all VRF keys from the Chainlink node
func (c *chainlink) ReadVRFKeys() (*VRFKeys, error) {
	vrfKeys := &VRFKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading VRF Keys")
	err := c.do(http.MethodGet, "/v2/keys/vrf", nil, vrfKeys, http.StatusOK)
	if len(vrfKeys.Data) == 0 {
		log.Warn().Str("Node URL", c.Config.URL).Msg("Found no VRF Keys on the node")
	}
	return vrfKeys, err
}

// CreateVRFKey creates a VRF key on the Chainlink node
func (c *chainlink) CreateVRFKey() (*VRFKey, error) {
	vrfKey := &VRFKey{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating VRF Key")
	return vrfKey, c.do(http.MethodPost, "/v2/keys/vrf", nil, vrfKey, http.StatusOK)
}

// ExportVRFKey exports a vrf key by key id
func (c *chainlink) ExportVRFKey(keyId string) (*VRFExportKey, error) {
	vrfExportKey := &VRFExportKey{}
	log.Info().Str("Node URL", c.Config.URL).Str("ID", keyId).Msg("Exporting VRF Key")
	err := c.do(http.MethodPost, fmt.Sprintf("/v2/keys/vrf/export/%s", keyId), nil, vrfExportKey, http.StatusOK)
	return vrfExportKey, err
}

// ImportVRFKey import vrf key
func (c *chainlink) ImportVRFKey(vrfExportKey *VRFExportKey) (*VRFKey, error) {
	vrfKey := &VRFKey{}
	log.Info().Str("Node URL", c.Config.URL).Str("ID", vrfExportKey.VrfKey.Address).Msg("Importing VRF Key")
	err := c.do(http.MethodPost, "/v2/keys/vrf/import", vrfExportKey, vrfKey, http.StatusOK)
	return vrfKey, err
}

// CreateCSAKey creates a CSA key on the Chainlink node, only 1 CSA key per noe
func (c *chainlink) CreateCSAKey() (*CSAKey, error) {
	csaKey := &CSAKey{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Creating CSA Key")
	return csaKey, c.do(http.MethodPost, "/v2/keys/csa", nil, csaKey, http.StatusOK)
}

// ReadCSAKeys reads CSA keys from the Chainlink node
func (c *chainlink) ReadCSAKeys() (*CSAKeys, error) {
	csaKeys := &CSAKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading CSA Keys")
	err := c.do(http.MethodGet, "/v2/keys/csa", nil, csaKeys, http.StatusOK)
	if len(csaKeys.Data) == 0 {
		log.Warn().Str("Node URL", c.Config.URL).Msg("Found no CSA Keys on the node")
	}
	return csaKeys, err
}

// CreateEI creates an EI on the Chainlink node based on the provided attributes and returns the respective secrets
func (c *chainlink) CreateEI(eia *EIAttributes) (*EIKeyCreate, error) {
	ei := EIKeyCreate{}
	log.Info().Str("Node URL", c.Config.URL).Str("Name", eia.Name).Msg("Creating External Initiator")
	return &ei, c.do(http.MethodPost, "/v2/external_initiators", eia, &ei, http.StatusCreated)
}

// ReadEIs reads all of the configured EIs from the chainlink node
func (c *chainlink) ReadEIs() (*EIKeys, error) {
	ei := EIKeys{}
	log.Info().Str("Node URL", c.Config.URL).Msg("Reading EI Keys")
	return &ei, c.do(http.MethodGet, "/v2/external_initiators", nil, &ei, http.StatusOK)
}

// DeleteEI deletes an external initiator in the Chainlink node based on the provided name
func (c *chainlink) DeleteEI(name string) error {
	log.Info().Str("Node URL", c.Config.URL).Str("Name", name).Msg("Deleting EI")
	return c.do(http.MethodDelete, fmt.Sprintf("/v2/external_initiators/%s", name), nil, nil, http.StatusNoContent)
}

// CreateTerraChain creates a terra chain
func (c *chainlink) CreateTerraChain(chain *TerraChainAttributes) (*TerraChainCreate, error) {
	response := TerraChainCreate{}
	log.Info().Str("Node URL", c.Config.URL).Str("Chain ID", chain.ChainID).Msg("Creating Terra Chain")
	return &response, c.do(http.MethodPost, "/v2/chains/terra", chain, &response, http.StatusCreated)
}

// CreateTerraNode creates a terra node
func (c *chainlink) CreateTerraNode(node *TerraNodeAttributes) (*TerraNodeCreate, error) {
	response := TerraNodeCreate{}
	log.Info().Str("Node URL", c.Config.URL).Str("Name", node.Name).Msg("Creating Terra Node")
	return &response, c.do(http.MethodPost, "/v2/nodes/terra", node, &response, http.StatusOK)
}

// CreateSolana creates a solana chain
func (c *chainlink) CreateSolanaChain(chain *SolanaChainAttributes) (*SolanaChainCreate, error) {
	response := SolanaChainCreate{}
	log.Info().Str("Node URL", c.Config.URL).Str("Chain ID", chain.ChainID).Msg("Creating Solana Chain")
	return &response, c.do(http.MethodPost, "/v2/chains/solana", chain, &response, http.StatusCreated)
}

// CreateSolanaNode creates a solana node
func (c *chainlink) CreateSolanaNode(node *SolanaNodeAttributes) (*SolanaNodeCreate, error) {
	response := SolanaNodeCreate{}
	log.Info().Str("Node URL", c.Config.URL).Str("Name", node.Name).Msg("Creating Solana Node")
	return &response, c.do(http.MethodPost, "/v2/nodes/solana", node, &response, http.StatusOK)
}

// RemoteIP retrieves the inter-cluster IP of the chainlink node, for use with inter-node communications
func (c *chainlink) RemoteIP() string {
	return c.Config.RemoteIP
}

// SetSessionCookie authenticates against the Chainlink node and stores the cookie in client state
func (c *chainlink) SetSessionCookie() error {
	session := &Session{Email: c.Config.Email, Password: c.Config.Password}
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/sessions", c.Config.URL),
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(
			"error while reading response: %v\nURL: %s\nresponse received: %s",
			err,
			c.Config.URL,
			string(b),
		)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf(
			"status code of %d was returned when trying to get a session\nURL: %s\nresponse received: %s",
			resp.StatusCode,
			c.Config.URL,
			b,
		)
	}
	if len(resp.Cookies()) == 0 {
		return fmt.Errorf("no cookie was returned after getting a session")
	}
	c.BasicHTTPClient.Cookies = resp.Cookies()

	sessionFound := false
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "clsession" {
			sessionFound = true
		}
	}
	if !sessionFound {
		return fmt.Errorf("chainlink: session cookie wasn't returned on login")
	}
	return nil
}

// Profile starts a profile session on the Chainlink node for a pre-determined length, then runs the provided function
// to profile it.
func (c *chainlink) Profile(profileTime time.Duration, profileFunction func(Chainlink)) (*ctfClient.ChainlinkProfileResults, error) {
	profileSeconds := int(profileTime.Seconds())
	profileResults := ctfClient.NewBlankChainlinkProfileResults()
	profileErrorGroup := new(errgroup.Group)
	var profileExecutedGroup sync.WaitGroup
	log.Info().Int("Seconds to Profile", profileSeconds).Str("Node URL", c.Config.URL).Msg("Starting Node PPROF session")
	for _, rep := range profileResults.Reports {
		profileExecutedGroup.Add(1)
		profileReport := rep
		// The profile function returns with the profile results after the profile time frame has concluded
		// e.g. a profile API call of 5 seconds will start profiling, wait for 5 seconds, then send back results
		profileErrorGroup.Go(func() error {
			log.Debug().Str("Type", profileReport.Type).Msg("PROFILING")
			uri := fmt.Sprintf("/v2/debug/pprof/%s?seconds=%d", profileReport.Type, profileSeconds)
			profileExecutedGroup.Done()
			rawBytes, err := c.doRawBytes(http.MethodGet, uri, nil, nil, http.StatusOK)
			if err != nil {
				return err
			}
			log.Debug().Str("Type", profileReport.Type).Msg("DONE PROFILING")
			profileReport.Data = rawBytes
			return nil
		})
	}
	// Wait for the profiling to actually get triggered on the node before running the function to profile
	// An imperfect solution, but an effective one.
	profileExecutedGroup.Wait()

	funcStart := time.Now()
	// Feed this chainlink node into the profiling function
	profileFunction(c)
	actualRunTime := time.Since(funcStart)
	actualSeconds := int(actualRunTime.Seconds())

	if actualSeconds > profileSeconds {
		log.Warn().
			Int("Actual Seconds", actualSeconds).
			Int("Profile Seconds", profileSeconds).
			Msg("Your profile function took longer than expected to run, increase profileTime")
	} else if actualSeconds < profileSeconds && actualSeconds > 0 {
		log.Warn().
			Int("Actual Seconds", actualSeconds).
			Int("Profile Seconds", profileSeconds).
			Msg("Your profile function took shorter than expected to run, you can decrease profileTime")
	}
	profileResults.ActualRunSeconds = actualSeconds
	profileResults.ScheduledProfileSeconds = profileSeconds
	return profileResults, profileErrorGroup.Wait() // Wait for all the results of the profiled function to come in
}

// SetPageSize globally sets the page
func (c *chainlink) SetPageSize(size int) {
	c.pageSize = size
}

// SetClient overrides the http client, used for mocking out the Chainlink server for unit testing
func (c *chainlink) SetClient(client *http.Client) {
	c.HttpClient = client
}

func (c *chainlink) doRawBytes(
	method,
	endpoint string,
	body []byte, obj interface{},
	expectedStatusCode int,
) ([]byte, error) {
	client := c.HttpClient

	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s%s", c.Config.URL, endpoint),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	for _, cookie := range c.Cookies {
		req.AddCookie(cookie)
	}

	q := req.URL.Query()
	q.Add("size", fmt.Sprint(c.pageSize))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"error while reading response: %v\nURL: %s\nresponse received: %s",
			err,
			c.Config.URL,
			string(b),
		)
	}
	if resp.StatusCode == http.StatusNotFound {
		return b, ctfClient.ErrNotFound
	} else if resp.StatusCode == http.StatusUnprocessableEntity {
		return b, ctfClient.ErrUnprocessableEntity
	} else if resp.StatusCode != expectedStatusCode {
		return b, fmt.Errorf(
			"unexpected response code, got %d, expected %d\nURL: %s\nresponse received: %s",
			resp.StatusCode,
			expectedStatusCode,
			c.Config.URL,
			string(b),
		)
	}
	return b, err
}

func (c *chainlink) doRaw(
	method,
	endpoint string,
	body []byte, obj interface{},
	expectedStatusCode int,
) error {
	respBody, err := c.doRawBytes(method, endpoint, body, obj, expectedStatusCode)
	if obj == nil || err != nil {
		return err
	}

	err = json.Unmarshal(respBody, &obj)
	if err != nil {
		return fmt.Errorf(
			"error while unmarshaling response to JSON: %v\nURL: %s\nresponse received: %s",
			err,
			c.Config.URL,
			string(respBody),
		)
	}
	return err
}

func (c *chainlink) do(
	method,
	endpoint string,
	body interface{},
	obj interface{},
	expectedStatusCode int,
) error {
	b, err := json.Marshal(body)
	if body != nil && err != nil {
		return err
	}
	return c.doRaw(method, endpoint, b, obj, expectedStatusCode)
}

// ConnectChainlinkNodes creates new chainlink clients
func ConnectChainlinkNodes(e *environment.Environment) ([]Chainlink, error) {
	var clients []Chainlink
	localURLs := e.URLs[chainlinkChart.NodesLocalURLsKey]
	internalURLs := e.URLs[chainlinkChart.NodesInternalURLsKey]
	for i, localURL := range localURLs {
		internalHost := parseHostname(internalURLs[i])
		c, err := NewChainlink(&ChainlinkConfig{
			URL:      localURL,
			Email:    "notreal@fakeemail.ch",
			Password: "fj293fbBnlQ!f9vNs",
			RemoteIP: internalHost,
		}, http.DefaultClient)
		clients = append(clients, c)
		if err != nil {
			return nil, err
		}
	}
	return clients, nil
}

func parseHostname(s string) string {
	r := regexp.MustCompile(`://(?P<Host>.*):`)
	return r.FindStringSubmatch(s)[1]
}
