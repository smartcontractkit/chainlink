// Package client enables interaction with APIs of test components like the mockserver and Chainlink nodes
package nodeclient

import (
	"crypto/tls"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const (
	// ChainlinkKeyPassword used to encrypt exported keys
	ChainlinkKeyPassword string = "twochains"
	// NodeURL string for logging
	NodeURL string = "Node URL"
)

var (
	// OneLINK representation of a single LINK token
	OneLINK           = big.NewFloat(1e18)
	mapKeyTypeToChain = map[string]string{
		"evm":      "eTHKeys",
		"solana":   "encryptedSolanaKeys",
		"starknet": "encryptedStarkNetKeys",
	}
)

type ChainlinkClient struct {
	APIClient         *resty.Client
	Config            *ChainlinkConfig
	pageSize          int
	primaryEthAddress string
	ethAddresses      []string
	l                 zerolog.Logger
}

// NewChainlinkClient creates a new Chainlink model using a provided config
func NewChainlinkClient(c *ChainlinkConfig, logger zerolog.Logger) (*ChainlinkClient, error) {
	rc, err := initRestyClient(c.URL, c.Email, c.Password, c.Headers, c.HTTPTimeout)
	if err != nil {
		return nil, err
	}
	return &ChainlinkClient{
		Config:    c,
		APIClient: rc,
		pageSize:  25,
		l:         logger,
	}, nil
}

func initRestyClient(url string, email string, password string, headers map[string]string, timeout *time.Duration) (*resty.Client, error) {
	isDebug := os.Getenv("RESTY_DEBUG") == "true"
	// G402 - TODO: certificates
	//nolint
	rc := resty.New().SetBaseURL(url).SetHeaders(headers).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).SetDebug(isDebug)
	if timeout != nil {
		rc.SetTimeout(*timeout)
	}
	session := &Session{Email: email, Password: password}
	// Retry the connection on boot up, sometimes pods can still be starting up and not ready to accept connections
	var resp *resty.Response
	var err error
	retryCount := 20
	for i := 0; i < retryCount; i++ {
		resp, err = rc.R().SetBody(session).Post("/sessions")
		if err != nil {
			log.Warn().Err(err).Str("URL", url).Interface("Session Details", session).Msg("Error connecting to Chainlink node, retrying")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("error connecting to chainlink node after %d attempts: %w", retryCount, err)
	}
	rc.SetCookies(resp.Cookies())
	log.Debug().Str("URL", url).Msg("Connected to Chainlink node")
	return rc, nil
}

// URL Chainlink instance http url
func (c *ChainlinkClient) URL() string {
	return c.Config.URL
}

func (c *ChainlinkClient) WithRetryCount(retryCount int) *ChainlinkClient {
	c.APIClient.SetRetryCount(retryCount)
	return c
}

// Health returns all statuses health info
func (c *ChainlinkClient) Health() (*HealthResponse, *http.Response, error) {
	respBody := &HealthResponse{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Requesting health data")
	resp, err := c.APIClient.R().
		SetResult(&respBody).
		Get("/health")
	if err != nil {
		return nil, nil, err
	}
	return respBody, resp.RawResponse, err
}

// CreateJobRaw creates a Chainlink job based on the provided spec string
func (c *ChainlinkClient) CreateJobRaw(spec string) (*Job, *http.Response, error) {
	job := &Job{}
	c.l.Info().Str("Node URL", c.Config.URL).Msg("Creating Job")
	c.l.Trace().Str("Node URL", c.Config.URL).Str("Job Body", spec).Msg("Creating Job")
	resp, err := c.APIClient.R().
		SetBody(&JobForm{
			TOML: spec,
		}).
		SetResult(&job).
		Post("/v2/jobs")
	if err != nil {
		return nil, nil, err
	}
	return job, resp.RawResponse, err
}

// MustCreateJob creates a Chainlink job based on the provided spec struct and returns error if
// the request is unsuccessful
func (c *ChainlinkClient) MustCreateJob(spec JobSpec) (*Job, error) {
	job, resp, err := c.CreateJob(spec)
	if err != nil {
		return nil, err
	}
	return job, VerifyStatusCodeWithResponse(resp, http.StatusOK)
}

func (c *ChainlinkClient) GetConfig() ChainlinkConfig {
	return *c.Config
}

// CreateJob creates a Chainlink job based on the provided spec struct
func (c *ChainlinkClient) CreateJob(spec JobSpec) (*Job, *resty.Response, error) {
	job := &Job{}
	specString, err := spec.String()
	if err != nil {
		return nil, nil, err
	}
	c.l.Info().Str("Node URL", c.Config.URL).Str("Type", spec.Type()).Msg("Creating Job")
	c.l.Trace().Str("Node URL", c.Config.URL).Str("Type", spec.Type()).Str("Spec", specString).Msg("Creating Job")
	resp, err := c.APIClient.R().
		SetBody(&JobForm{
			TOML: specString,
		}).
		SetResult(&job).
		Post("/v2/jobs")
	if err != nil {
		return nil, nil, err
	}
	return job, resp, err
}

// ReadJobs reads all jobs from the Chainlink node
func (c *ChainlinkClient) ReadJobs() (*ResponseSlice, *http.Response, error) {
	specObj := &ResponseSlice{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Getting Jobs")
	resp, err := c.APIClient.R().
		SetResult(&specObj).
		Get("/v2/jobs")
	if err != nil {
		return nil, nil, err
	}
	return specObj, resp.RawResponse, err
}

// ReadJob reads a job with the provided ID from the Chainlink node
func (c *ChainlinkClient) ReadJob(id string) (*Response, *http.Response, error) {
	specObj := &Response{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("ID", id).Msg("Reading Job")
	resp, err := c.APIClient.R().
		SetResult(&specObj).
		SetPathParams(map[string]string{
			"id": id,
		}).
		Get("/v2/jobs/{id}")
	if err != nil {
		return nil, nil, err
	}
	return specObj, resp.RawResponse, err
}

// MustDeleteJob deletes a job with a provided ID from the Chainlink node and returns error if
// the request is unsuccessful
func (c *ChainlinkClient) MustDeleteJob(id string) error {
	resp, err := c.DeleteJob(id)
	if err != nil {
		return err
	}
	return VerifyStatusCode(resp.StatusCode, http.StatusNoContent)
}

// DeleteJob deletes a job with a provided ID from the Chainlink node
func (c *ChainlinkClient) DeleteJob(id string) (*http.Response, error) {
	c.l.Info().Str(NodeURL, c.Config.URL).Str("ID", id).Msg("Deleting Job")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"id": id,
		}).
		Delete("/v2/jobs/{id}")
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, err
}

// CreateSpec creates a job spec on the Chainlink node
func (c *ChainlinkClient) CreateSpec(spec string) (*Spec, *http.Response, error) {
	s := &Spec{}
	r := strings.NewReplacer("\n", "", " ", "", "\\", "") // Makes it more compact and readable for logging
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Spec", r.Replace(spec)).Msg("Creating Spec")
	resp, err := c.APIClient.R().
		SetBody([]byte(spec)).
		SetResult(&s).
		Post("/v2/specs")
	if err != nil {
		return nil, nil, err
	}
	return s, resp.RawResponse, err
}

// ReadSpec reads a job spec with the provided ID on the Chainlink node
func (c *ChainlinkClient) ReadSpec(id string) (*Response, *http.Response, error) {
	specObj := &Response{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("ID", id).Msg("Reading Spec")
	resp, err := c.APIClient.R().
		SetResult(&specObj).
		SetPathParams(map[string]string{
			"id": id,
		}).
		Get("/v2/specs/{id}")
	if err != nil {
		return nil, nil, err
	}
	return specObj, resp.RawResponse, err
}

// MustReadRunsByJob attempts to read all runs for a job and returns error if
// the request is unsuccessful
func (c *ChainlinkClient) MustReadRunsByJob(jobID string) (*JobRunsResponse, error) {
	runsObj, resp, err := c.ReadRunsByJob(jobID)
	if err != nil {
		return nil, err
	}
	return runsObj, VerifyStatusCode(resp.StatusCode, http.StatusOK)
}

// ReadRunsByJob reads all runs for a job
func (c *ChainlinkClient) ReadRunsByJob(jobID string) (*JobRunsResponse, *http.Response, error) {
	runsObj := &JobRunsResponse{}
	c.l.Debug().Str(NodeURL, c.Config.URL).Str("JobID", jobID).Msg("Reading runs for a job")
	resp, err := c.APIClient.R().
		SetResult(&runsObj).
		SetPathParams(map[string]string{
			"jobID": jobID,
		}).
		Get("/v2/jobs/{jobID}/runs")
	if err != nil {
		return nil, nil, err
	}
	return runsObj, resp.RawResponse, err
}

// DeleteSpec deletes a job spec with the provided ID from the Chainlink node
func (c *ChainlinkClient) DeleteSpec(id string) (*http.Response, error) {
	c.l.Info().Str(NodeURL, c.Config.URL).Str("ID", id).Msg("Deleting Spec")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"id": id,
		}).
		Delete("/v2/specs/{id}")
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, err
}

// MustCreateBridge creates a bridge on the Chainlink node based on the provided attributes and returns error if
// the request is unsuccessful
func (c *ChainlinkClient) MustCreateBridge(bta *BridgeTypeAttributes) error {
	c.l.Debug().Str(NodeURL, c.Config.URL).Str("Name", bta.Name).Msg("Creating Bridge")
	resp, err := c.CreateBridge(bta)
	if err != nil {
		return err
	}
	return VerifyStatusCode(resp.StatusCode, http.StatusOK)
}

func (c *ChainlinkClient) CreateBridge(bta *BridgeTypeAttributes) (*http.Response, error) {
	c.l.Debug().Str(NodeURL, c.Config.URL).Str("Name", bta.Name).Msg("Creating Bridge")
	resp, err := c.APIClient.R().
		SetBody(bta).
		Post("/v2/bridge_types")
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, err
}

// ReadBridge reads a bridge from the Chainlink node based on the provided name
func (c *ChainlinkClient) ReadBridge(name string) (*BridgeType, *http.Response, error) {
	bt := BridgeType{}
	c.l.Debug().Str(NodeURL, c.Config.URL).Str("Name", name).Msg("Reading Bridge")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"name": name,
		}).
		SetResult(&bt).
		Get("/v2/bridge_types/{name}")
	if err != nil {
		return nil, nil, err
	}
	return &bt, resp.RawResponse, err
}

// ReadBridges reads bridges from the Chainlink node
func (c *ChainlinkClient) ReadBridges() (*Bridges, *resty.Response, error) {
	result := &Bridges{}
	c.l.Debug().Str(NodeURL, c.Config.URL).Msg("Getting all bridges")
	resp, err := c.APIClient.R().
		SetResult(&result).
		Get("/v2/bridge_types")
	if err != nil {
		return nil, nil, err
	}
	return result, resp, err
}

// DeleteBridge deletes a bridge on the Chainlink node based on the provided name
func (c *ChainlinkClient) DeleteBridge(name string) (*http.Response, error) {
	c.l.Debug().Str(NodeURL, c.Config.URL).Str("Name", name).Msg("Deleting Bridge")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"name": name,
		}).
		Delete("/v2/bridge_types/{name}")
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, err
}

// CreateOCRKey creates an OCRKey on the Chainlink node
func (c *ChainlinkClient) CreateOCRKey() (*OCRKey, *http.Response, error) {
	ocrKey := &OCRKey{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Creating OCR Key")
	resp, err := c.APIClient.R().
		SetResult(ocrKey).
		Post("/v2/keys/ocr")
	if err != nil {
		return nil, nil, err
	}
	return ocrKey, resp.RawResponse, err
}

// MustReadOCRKeys reads all OCRKeys from the Chainlink node and returns error if
// the request is unsuccessful
func (c *ChainlinkClient) MustReadOCRKeys() (*OCRKeys, error) {
	ocrKeys := &OCRKeys{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading OCR Keys")
	resp, err := c.APIClient.R().
		SetResult(ocrKeys).
		Get("/v2/keys/ocr")
	if err != nil {
		return nil, err
	}
	err = VerifyStatusCode(resp.StatusCode(), http.StatusOK)
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
func (c *ChainlinkClient) DeleteOCRKey(id string) (*http.Response, error) {
	c.l.Info().Str(NodeURL, c.Config.URL).Str("ID", id).Msg("Deleting OCR Key")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"id": id,
		}).
		Delete("/v2/keys/ocr/{id}")
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, err
}

// CreateOCR2Key creates an OCR2Key on the Chainlink node
func (c *ChainlinkClient) CreateOCR2Key(chain string) (*OCR2Key, *http.Response, error) {
	ocr2Key := &OCR2Key{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Creating OCR2 Key")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"chain": chain,
		}).
		SetResult(ocr2Key).
		Post("/v2/keys/ocr2/{chain}")
	if err != nil {
		return nil, nil, err
	}
	return ocr2Key, resp.RawResponse, err
}

// ReadOCR2Keys reads all OCR2Keys from the Chainlink node
func (c *ChainlinkClient) ReadOCR2Keys() (*OCR2Keys, *http.Response, error) {
	ocr2Keys := &OCR2Keys{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading OCR2 Keys")
	resp, err := c.APIClient.R().
		SetResult(ocr2Keys).
		Get("/v2/keys/ocr2")
	return ocr2Keys, resp.RawResponse, err
}

// MustReadOCR2Keys reads all OCR2Keys from the Chainlink node returns err if response not 200
func (c *ChainlinkClient) MustReadOCR2Keys() (*OCR2Keys, error) {
	ocr2Keys := &OCR2Keys{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading OCR2 Keys")
	resp, err := c.APIClient.R().
		SetResult(ocr2Keys).
		Get("/v2/keys/ocr2")
	if err != nil {
		return nil, err
	}
	err = VerifyStatusCode(resp.StatusCode(), http.StatusOK)
	return ocr2Keys, err
}

// DeleteOCR2Key deletes an OCR2Key based on the provided ID
func (c *ChainlinkClient) DeleteOCR2Key(id string) (*http.Response, error) {
	c.l.Info().Str(NodeURL, c.Config.URL).Str("ID", id).Msg("Deleting OCR2 Key")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"id": id,
		}).
		Delete("/v2/keys/ocr2/{id}")
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, err
}

// CreateP2PKey creates an P2PKey on the Chainlink node
func (c *ChainlinkClient) CreateP2PKey() (*P2PKey, *http.Response, error) {
	p2pKey := &P2PKey{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Creating P2P Key")
	resp, err := c.APIClient.R().
		SetResult(p2pKey).
		Post("/v2/keys/p2p")
	if err != nil {
		return nil, nil, err
	}
	return p2pKey, resp.RawResponse, err
}

// MustReadP2PKeys reads all P2PKeys from the Chainlink node and returns error if
// the request is unsuccessful
func (c *ChainlinkClient) MustReadP2PKeys() (*P2PKeys, error) {
	p2pKeys := &P2PKeys{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading P2P Keys")
	resp, err := c.APIClient.R().
		SetResult(p2pKeys).
		Get("/v2/keys/p2p")
	if err != nil {
		return nil, err
	}
	err = VerifyStatusCode(resp.StatusCode(), http.StatusOK)
	if len(p2pKeys.Data) == 0 {
		err = fmt.Errorf("Found no P2P Keys on the Chainlink node. Node URL: %s", c.Config.URL)
		c.l.Err(err).Msg("Error getting P2P keys")
		return nil, err
	}
	for index := range p2pKeys.Data {
		p2pKeys.Data[index].Attributes.PeerID = strings.TrimPrefix(p2pKeys.Data[index].Attributes.PeerID, "p2p_")
	}
	return p2pKeys, err
}

// DeleteP2PKey deletes a P2PKey on the Chainlink node based on the provided ID
func (c *ChainlinkClient) DeleteP2PKey(id int) (*http.Response, error) {
	c.l.Info().Str(NodeURL, c.Config.URL).Int("ID", id).Msg("Deleting P2P Key")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"id": strconv.Itoa(id),
		}).
		Delete("/v2/keys/p2p/{id}")
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, err
}

// MustReadETHKeys reads all ETH keys from the Chainlink node and returns error if
// the request is unsuccessful
func (c *ChainlinkClient) MustReadETHKeys() (*ETHKeys, error) {
	ethKeys := &ETHKeys{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading ETH Keys")
	resp, err := c.APIClient.R().
		SetResult(ethKeys).
		Get("/v2/keys/eth")
	if err != nil {
		return nil, err
	}
	err = VerifyStatusCode(resp.StatusCode(), http.StatusOK)
	if len(ethKeys.Data) == 0 {
		c.l.Warn().Str(NodeURL, c.Config.URL).Msg("Found no ETH Keys on the node")
	}
	return ethKeys, err
}

// UpdateEthKeyMaxGasPriceGWei updates the maxGasPriceGWei for an eth key
func (c *ChainlinkClient) UpdateEthKeyMaxGasPriceGWei(keyId string, gWei int) (*ETHKey, *http.Response, error) {
	ethKey := &ETHKey{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("ID", keyId).Int("maxGasPriceGWei", gWei).Msg("Update maxGasPriceGWei for eth key")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"keyId": keyId,
		}).
		SetQueryParams(map[string]string{
			"maxGasPriceGWei": strconv.Itoa(gWei),
		}).
		SetResult(ethKey).
		Put("/v2/keys/eth/{keyId}")
	if err != nil {
		return nil, nil, err
	}
	err = VerifyStatusCode(resp.StatusCode(), http.StatusOK)
	if err != nil {
		return nil, nil, err
	}
	return ethKey, resp.RawResponse, err
}

// ReadPrimaryETHKey reads updated information about the Chainlink's primary ETH key
func (c *ChainlinkClient) ReadPrimaryETHKey() (*ETHKeyData, error) {
	ethKeys, err := c.MustReadETHKeys()
	if err != nil {
		return nil, err
	}
	if len(ethKeys.Data) == 0 {
		return nil, fmt.Errorf("Error retrieving primary eth key on node %s: No ETH keys present", c.URL())
	}
	return &ethKeys.Data[0], nil
}

// ReadETHKeyAtIndex reads updated information about the Chainlink's ETH key at given index
func (c *ChainlinkClient) ReadETHKeyAtIndex(keyIndex int) (*ETHKeyData, error) {
	ethKeys, err := c.MustReadETHKeys()
	if err != nil {
		return nil, err
	}
	if len(ethKeys.Data) == 0 {
		return nil, fmt.Errorf("Error retrieving primary eth key on node %s: No ETH keys present", c.URL())
	}
	return &ethKeys.Data[keyIndex], nil
}

// PrimaryEthAddress returns the primary ETH address for the Chainlink node
func (c *ChainlinkClient) PrimaryEthAddress() (string, error) {
	if c.primaryEthAddress == "" {
		ethKeys, err := c.MustReadETHKeys()
		if err != nil {
			return "", err
		}
		c.primaryEthAddress = ethKeys.Data[0].Attributes.Address
	}
	return c.primaryEthAddress, nil
}

// EthAddresses returns the ETH addresses for the Chainlink node
func (c *ChainlinkClient) EthAddresses() ([]string, error) {
	if len(c.ethAddresses) == 0 {
		ethKeys, err := c.MustReadETHKeys()
		c.ethAddresses = make([]string, len(ethKeys.Data))
		if err != nil {
			return make([]string, 0), err
		}
		for index, data := range ethKeys.Data {
			c.ethAddresses[index] = data.Attributes.Address
		}
	}
	return c.ethAddresses, nil
}

// EthAddresses returns the ETH addresses of the Chainlink node for a specific chain id
func (c *ChainlinkClient) EthAddressesForChain(chainId string) ([]string, error) {
	var ethAddresses []string
	ethKeys, err := c.MustReadETHKeys()
	if err != nil {
		return nil, err
	}
	for _, ethKey := range ethKeys.Data {
		if ethKey.Attributes.ChainID == chainId {
			ethAddresses = append(ethAddresses, ethKey.Attributes.Address)
		}
	}
	return ethAddresses, nil
}

// PrimaryEthAddressForChain returns the primary ETH address for the Chainlink node for mentioned chain
func (c *ChainlinkClient) PrimaryEthAddressForChain(chainId string) (string, error) {
	ethKeys, err := c.MustReadETHKeys()
	if err != nil {
		return "", err
	}
	for _, ethKey := range ethKeys.Data {
		if ethKey.Attributes.ChainID == chainId {
			return ethKey.Attributes.Address, nil
		}
	}
	return "", nil
}

// ExportEVMKeys exports Chainlink private EVM keys
func (c *ChainlinkClient) ExportEVMKeys() ([]*ExportedEVMKey, error) {
	exportedKeys := make([]*ExportedEVMKey, 0)
	keys, err := c.MustReadETHKeys()
	if err != nil {
		return nil, err
	}
	for _, key := range keys.Data {
		if key.Attributes.ETHBalance != "0" {
			exportedKey := &ExportedEVMKey{}
			_, err := c.APIClient.R().
				SetResult(exportedKey).
				SetPathParam("keyAddress", key.Attributes.Address).
				SetQueryParam("newpassword", ChainlinkKeyPassword).
				Post("/v2/keys/eth/export/{keyAddress}")
			if err != nil {
				return nil, err
			}
			exportedKeys = append(exportedKeys, exportedKey)
		}
	}
	c.l.Info().
		Str(NodeURL, c.Config.URL).
		Str("Password", ChainlinkKeyPassword).
		Msg("Exported EVM Keys")
	return exportedKeys, nil
}

// ExportEVMKeysForChain exports Chainlink private EVM keys for a particular chain
func (c *ChainlinkClient) ExportEVMKeysForChain(chainid string) ([]*ExportedEVMKey, error) {
	exportedKeys := make([]*ExportedEVMKey, 0)
	keys, err := c.MustReadETHKeys()
	if err != nil {
		return nil, err
	}
	for _, key := range keys.Data {
		if key.Attributes.ETHBalance != "0" && key.Attributes.ChainID == chainid {
			exportedKey := &ExportedEVMKey{}
			_, err := c.APIClient.R().
				SetResult(exportedKey).
				SetPathParam("keyAddress", key.Attributes.Address).
				SetQueryParam("newpassword", ChainlinkKeyPassword).
				Post("/v2/keys/eth/export/{keyAddress}")
			if err != nil {
				return nil, err
			}
			exportedKeys = append(exportedKeys, exportedKey)
		}
	}
	c.l.Info().
		Str(NodeURL, c.Config.URL).
		Str("Password", ChainlinkKeyPassword).
		Msg("Exported EVM Keys")
	return exportedKeys, nil
}

// CreateTxKey creates a tx key on the Chainlink node
func (c *ChainlinkClient) CreateTxKey(chain string, chainId string) (*TxKey, *http.Response, error) {
	txKey := &TxKey{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Creating Tx Key")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"chain": chain,
		}).
		SetQueryParam("evmChainID", chainId).
		SetResult(txKey).
		Post("/v2/keys/{chain}")
	if err != nil {
		return nil, nil, err
	}
	return txKey, resp.RawResponse, err
}

// ReadTxKeys reads all tx keys from the Chainlink node
func (c *ChainlinkClient) ReadTxKeys(chain string) (*TxKeys, *http.Response, error) {
	txKeys := &TxKeys{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading Tx Keys")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"chain": chain,
		}).
		SetResult(txKeys).
		Get("/v2/keys/{chain}")
	if err != nil {
		return nil, nil, err
	}
	return txKeys, resp.RawResponse, err
}

// DeleteTxKey deletes an tx key based on the provided ID
func (c *ChainlinkClient) DeleteTxKey(chain string, id string) (*http.Response, error) {
	c.l.Info().Str(NodeURL, c.Config.URL).Str("ID", id).Msg("Deleting Tx Key")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"chain": chain,
			"id":    id,
		}).
		Delete("/v2/keys/{chain}/{id}")
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, err
}

// MustReadTransactionAttempts reads all transaction attempts on the Chainlink node
// and returns error if the request is unsuccessful
func (c *ChainlinkClient) MustReadTransactionAttempts() (*TransactionsData, error) {
	txsData := &TransactionsData{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading Transaction Attempts")
	resp, err := c.APIClient.R().
		SetResult(txsData).
		Get("/v2/tx_attempts")
	if err != nil {
		return nil, err
	}
	err = VerifyStatusCode(resp.StatusCode(), http.StatusOK)
	return txsData, err
}

// ReadTransactions reads all transactions made by the Chainlink node
func (c *ChainlinkClient) ReadTransactions() (*TransactionsData, *http.Response, error) {
	txsData := &TransactionsData{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading Transactions")
	resp, err := c.APIClient.R().
		SetResult(txsData).
		Get("/v2/transactions")
	if err != nil {
		return nil, nil, err
	}
	return txsData, resp.RawResponse, err
}

// MustSendNativeToken sends native token (ETH usually) of a specified amount from one of its addresses to the target address
// and returns error if the request is unsuccessful
// WARNING: The txdata object that Chainlink sends back is almost always blank.
func (c *ChainlinkClient) MustSendNativeToken(amount *big.Int, fromAddress, toAddress string) (TransactionData, error) {
	request := SendEtherRequest{
		DestinationAddress: toAddress,
		FromAddress:        fromAddress,
		Amount:             amount.String(),
		AllowHigherAmounts: true,
	}
	txData := SingleTransactionDataWrapper{}
	resp, err := c.APIClient.R().
		SetBody(request).
		SetResult(txData).
		Post("/v2/transfers")

	c.l.Info().
		Str(NodeURL, c.Config.URL).
		Str("From", fromAddress).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Msg("Sending Native Token")
	if err == nil {
		err = VerifyStatusCode(resp.StatusCode(), http.StatusOK)
	}

	return txData.Data, err
}

// ReadVRFKeys reads all VRF keys from the Chainlink node
func (c *ChainlinkClient) ReadVRFKeys() (*VRFKeys, *http.Response, error) {
	vrfKeys := &VRFKeys{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading VRF Keys")
	resp, err := c.APIClient.R().
		SetResult(vrfKeys).
		Get("/v2/keys/vrf")
	if err != nil {
		return nil, nil, err
	}
	if len(vrfKeys.Data) == 0 {
		c.l.Warn().Str(NodeURL, c.Config.URL).Msg("Found no VRF Keys on the node")
	}
	return vrfKeys, resp.RawResponse, err
}

// MustCreateVRFKey creates a VRF key on the Chainlink node
// and returns error if the request is unsuccessful
func (c *ChainlinkClient) MustCreateVRFKey() (*VRFKey, error) {
	vrfKey := &VRFKey{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Creating VRF Key")
	resp, err := c.APIClient.R().
		SetResult(vrfKey).
		Post("/v2/keys/vrf")
	if err == nil {
		err = VerifyStatusCode(resp.StatusCode(), http.StatusOK)
	}
	return vrfKey, err
}

// ExportVRFKey exports a vrf key by key id
func (c *ChainlinkClient) ExportVRFKey(keyId string) (*VRFExportKey, *http.Response, error) {
	vrfExportKey := &VRFExportKey{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("ID", keyId).Msg("Exporting VRF Key")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"keyId": keyId,
		}).
		SetResult(vrfExportKey).
		Post("/v2/keys/vrf/export/{keyId}")
	if err != nil {
		return nil, nil, err
	}
	return vrfExportKey, resp.RawResponse, err
}

// ImportVRFKey import vrf key
func (c *ChainlinkClient) ImportVRFKey(vrfExportKey *VRFExportKey) (*VRFKey, *http.Response, error) {
	vrfKey := &VRFKey{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("ID", vrfExportKey.VrfKey.Address).Msg("Importing VRF Key")
	resp, err := c.APIClient.R().
		SetBody(vrfExportKey).
		SetResult(vrfKey).
		Post("/v2/keys/vrf/import")
	if err != nil {
		return nil, nil, err
	}
	return vrfKey, resp.RawResponse, err
}

// CreateCSAKey creates a CSA key on the Chainlink node, only 1 CSA key per noe
func (c *ChainlinkClient) CreateCSAKey() (*CSAKey, *http.Response, error) {
	csaKey := &CSAKey{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Creating CSA Key")
	resp, err := c.APIClient.R().
		SetResult(csaKey).
		Post("/v2/keys/csa")
	if err != nil {
		return nil, nil, err
	}
	return csaKey, resp.RawResponse, err
}

func (c *ChainlinkClient) MustReadCSAKeys() (*CSAKeys, *resty.Response, error) {
	csaKeys, res, err := c.ReadCSAKeys()
	if err != nil {
		return nil, res, err
	}
	return csaKeys, res, VerifyStatusCodeWithResponse(res, http.StatusOK)
}

// ReadCSAKeys reads CSA keys from the Chainlink node
func (c *ChainlinkClient) ReadCSAKeys() (*CSAKeys, *resty.Response, error) {
	csaKeys := &CSAKeys{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading CSA Keys")
	resp, err := c.APIClient.R().
		SetResult(csaKeys).
		Get("/v2/keys/csa")
	if len(csaKeys.Data) == 0 {
		c.l.Warn().Str(NodeURL, c.Config.URL).Msg("Found no CSA Keys on the node")
	}
	if err != nil {
		return nil, nil, err
	}
	return csaKeys, resp, err
}

// CreateEI creates an EI on the Chainlink node based on the provided attributes and returns the respective secrets
func (c *ChainlinkClient) CreateEI(eia *EIAttributes) (*EIKeyCreate, *http.Response, error) {
	ei := EIKeyCreate{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Name", eia.Name).Msg("Creating External Initiator")
	resp, err := c.APIClient.R().
		SetBody(eia).
		SetResult(&ei).
		Post("/v2/external_initiators")
	if err != nil {
		return nil, nil, err
	}
	return &ei, resp.RawResponse, err
}

// ReadEIs reads all of the configured EIs from the Chainlink node
func (c *ChainlinkClient) ReadEIs() (*EIKeys, *http.Response, error) {
	ei := EIKeys{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading EI Keys")
	resp, err := c.APIClient.R().
		SetResult(&ei).
		Get("/v2/external_initiators")
	if err != nil {
		return nil, nil, err
	}
	return &ei, resp.RawResponse, err
}

// DeleteEI deletes an external initiator in the Chainlink node based on the provided name
func (c *ChainlinkClient) DeleteEI(name string) (*http.Response, error) {
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Name", name).Msg("Deleting EI")
	resp, err := c.APIClient.R().
		SetPathParams(map[string]string{
			"name": name,
		}).
		Delete("/v2/external_initiators/{name}")
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, err
}

// CreateCosmosChain creates a cosmos chain
func (c *ChainlinkClient) CreateCosmosChain(chain *CosmosChainAttributes) (*CosmosChainCreate, *http.Response, error) {
	response := CosmosChainCreate{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Chain ID", chain.ChainID).Msg("Creating Cosmos Chain")
	resp, err := c.APIClient.R().
		SetBody(chain).
		SetResult(&response).
		Post("/v2/chains/cosmos")
	if err != nil {
		return nil, nil, err
	}
	return &response, resp.RawResponse, err
}

// CreateCosmosNode creates a cosmos node
func (c *ChainlinkClient) CreateCosmosNode(node *CosmosNodeAttributes) (*CosmosNodeCreate, *http.Response, error) {
	response := CosmosNodeCreate{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Name", node.Name).Msg("Creating Cosmos Node")
	resp, err := c.APIClient.R().
		SetBody(node).
		SetResult(&response).
		Post("/v2/nodes/cosmos")
	if err != nil {
		return nil, nil, err
	}
	return &response, resp.RawResponse, err
}

// CreateSolanaChain creates a solana chain
func (c *ChainlinkClient) CreateSolanaChain(chain *SolanaChainAttributes) (*SolanaChainCreate, *http.Response, error) {
	response := SolanaChainCreate{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Chain ID", chain.ChainID).Msg("Creating Solana Chain")
	resp, err := c.APIClient.R().
		SetBody(chain).
		SetResult(&response).
		Post("/v2/chains/solana")
	if err != nil {
		return nil, nil, err
	}
	return &response, resp.RawResponse, err
}

// CreateSolanaNode creates a solana node
func (c *ChainlinkClient) CreateSolanaNode(node *SolanaNodeAttributes) (*SolanaNodeCreate, *http.Response, error) {
	response := SolanaNodeCreate{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Name", node.Name).Msg("Creating Solana Node")
	resp, err := c.APIClient.R().
		SetBody(node).
		SetResult(&response).
		Post("/v2/nodes/solana")
	if err != nil {
		return nil, nil, err
	}
	return &response, resp.RawResponse, err
}

// CreateStarkNetChain creates a starknet chain
func (c *ChainlinkClient) CreateStarkNetChain(chain *StarkNetChainAttributes) (*StarkNetChainCreate, *http.Response, error) {
	response := StarkNetChainCreate{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Chain ID", chain.ChainID).Msg("Creating StarkNet Chain")
	resp, err := c.APIClient.R().
		SetBody(chain).
		SetResult(&response).
		Post("/v2/chains/starknet")
	if err != nil {
		return nil, nil, err
	}
	return &response, resp.RawResponse, err
}

// CreateStarkNetNode creates a starknet node
func (c *ChainlinkClient) CreateStarkNetNode(node *StarkNetNodeAttributes) (*StarkNetNodeCreate, *http.Response, error) {
	response := StarkNetNodeCreate{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Name", node.Name).Msg("Creating StarkNet Node")
	resp, err := c.APIClient.R().
		SetBody(node).
		SetResult(&response).
		Post("/v2/nodes/starknet")
	if err != nil {
		return nil, nil, err
	}
	return &response, resp.RawResponse, err
}

// CreateTronChain creates a tron chain
func (c *ChainlinkClient) CreateTronChain(chain *TronChainAttributes) (*TronChainCreate, *http.Response, error) {
	response := TronChainCreate{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Chain ID", chain.ChainID).Msg("Creating Tron Chain")
	resp, err := c.APIClient.R().
		SetBody(chain).
		SetResult(&response).
		Post("/v2/chains/tron")
	if err != nil {
		return nil, nil, err
	}
	return &response, resp.RawResponse, err
}

// CreateTronNode creates a tron node
func (c *ChainlinkClient) CreateTronNode(node *TronNodeAttributes) (*TronNodeCreate, *http.Response, error) {
	response := TronNodeCreate{}
	c.l.Info().Str(NodeURL, c.Config.URL).Str("Name", node.Name).Msg("Creating Tron Node")
	resp, err := c.APIClient.R().
		SetBody(node).
		SetResult(&response).
		Post("/v2/nodes/tron")
	if err != nil {
		return nil, nil, err
	}
	return &response, resp.RawResponse, err
}

// InternalIP retrieves the inter-cluster IP of the Chainlink node, for use with inter-node communications
func (c *ChainlinkClient) InternalIP() string {
	return c.Config.InternalIP
}

// Profile starts a profile session on the Chainlink node for a pre-determined length, then runs the provided function
// to profile it.
func (c *ChainlinkClient) Profile(profileTime time.Duration, profileFunction func(*ChainlinkClient)) (*ChainlinkProfileResults, error) {
	profileSeconds := int(profileTime.Seconds())
	profileResults := NewBlankChainlinkProfileResults()
	profileErrorGroup := new(errgroup.Group)
	var profileExecutedGroup sync.WaitGroup
	c.l.Info().Int("Seconds to Profile", profileSeconds).Str(NodeURL, c.Config.URL).Msg("Starting Node PPROF session")
	for _, rep := range profileResults.Reports {
		profileExecutedGroup.Add(1)
		profileReport := rep
		// The profile function returns with the profile results after the profile time frame has concluded
		// e.g. a profile API call of 5 seconds will start profiling, wait for 5 seconds, then send back results
		profileErrorGroup.Go(func() error {
			c.l.Debug().Str("Type", profileReport.Type).Msg("PROFILING")
			profileExecutedGroup.Done()
			resp, err := c.APIClient.R().
				SetPathParams(map[string]string{
					"reportType": profileReport.Type,
				}).
				SetQueryParams(map[string]string{
					"seconds": strconv.Itoa(profileSeconds),
				}).
				Get("/v2/debug/pprof/{reportType}")
			if err != nil {
				return err
			}
			err = VerifyStatusCode(resp.StatusCode(), http.StatusOK)
			if err != nil {
				return err
			}
			c.l.Debug().Str("Type", profileReport.Type).Msg("DONE PROFILING")
			profileReport.Data = resp.Body()
			return err
		})
	}
	// Wait for the profiling to actually get triggered on the node before running the function to profile
	// An imperfect solution, but an effective one.
	profileExecutedGroup.Wait()

	funcStart := time.Now()
	// Feed this Chainlink node into the profiling function
	profileFunction(c)
	actualRunTime := time.Since(funcStart)
	actualSeconds := int(actualRunTime.Seconds())

	if actualSeconds > profileSeconds {
		c.l.Warn().
			Int("Actual Seconds", actualSeconds).
			Int("Profile Seconds", profileSeconds).
			Msg("Your profile function took longer than expected to run, increase profileTime")
	} else if actualSeconds < profileSeconds && actualSeconds > 0 {
		c.l.Warn().
			Int("Actual Seconds", actualSeconds).
			Int("Profile Seconds", profileSeconds).
			Msg("Your profile function took shorter than expected to run, you can decrease profileTime")
	}
	profileResults.ActualRunSeconds = actualSeconds
	profileResults.ScheduledProfileSeconds = profileSeconds
	return profileResults, profileErrorGroup.Wait() // Wait for all the results of the profiled function to come in
}

// SetPageSize globally sets the page
func (c *ChainlinkClient) SetPageSize(size int) {
	c.pageSize = size
}

// VerifyStatusCode verifies the status code of the response. Favor VerifyStatusCodeWithResponse over this for better errors
func VerifyStatusCode(actStatusCd, expStatusCd int) error {
	if actStatusCd != expStatusCd {
		return fmt.Errorf(
			"unexpected response code, got %d, expected %d",
			actStatusCd,
			expStatusCd,
		)
	}
	return nil
}

// VerifyStatusCodeWithResponse verifies the status code of the response and returns the response as part of the error.
// Favor this over VerifyStatusCode
func VerifyStatusCodeWithResponse(res *resty.Response, expStatusCd int) error {
	actStatusCd := res.RawResponse.StatusCode
	if actStatusCd != expStatusCd {
		return fmt.Errorf(
			"unexpected response code, got %d, expected %d, response: %s",
			actStatusCd,
			expStatusCd,
			res.Body(),
		)
	}
	return nil
}

func CreateNodeKeysBundle(nodes []*ChainlinkClient, chainName string, chainId string) ([]NodeKeysBundle, []*CLNodesWithKeys, error) {
	nkb := make([]NodeKeysBundle, 0)
	var clNodes []*CLNodesWithKeys
	for _, n := range nodes {
		p2pkeys, err := n.MustReadP2PKeys()
		if err != nil {
			return nil, nil, err
		}
		if len(p2pkeys.Data) == 0 {
			return nil, nil, fmt.Errorf("found no P2P Keys on the Chainlink node. Node URL: %s", n.URL())
		}
		peerID := p2pkeys.Data[0].Attributes.PeerID
		// If there is already a txkey present for the chain skip creating a new one
		// otherwise the test logic will need multiple key management (like funding all the keys,
		// for ocr scenarios adding all keys to ocr config)
		var txKey *TxKey
		txKeys, _, err := n.ReadTxKeys(chainName)
		if err != nil {
			return nil, nil, err
		}
		if _, ok := mapKeyTypeToChain[chainName]; ok {
			for _, key := range txKeys.Data {
				if key.Type == mapKeyTypeToChain[chainName] {
					txKey = &TxKey{Data: key}
					break
				}
			}
		}
		// if no txkey is found for the chain, create a new one
		if txKey == nil {
			txKey, _, err = n.CreateTxKey(chainName, chainId)
			if err != nil {
				return nil, nil, err
			}
		}
		keys, _, err := n.ReadOCR2Keys()
		if err != nil {
			return nil, nil, err
		}
		var ocrKey *OCR2Key
		for _, key := range keys.Data {
			if key.Attributes.ChainType == chainName {
				ocrKey = &OCR2Key{Data: key}
				break
			}
		}

		if ocrKey == nil {
			return nil, nil, fmt.Errorf("no OCR key found for chain %s", chainName)
		}
		ethAddress, err := n.PrimaryEthAddressForChain(chainId)
		if err != nil {
			return nil, nil, err
		}
		bundle := NodeKeysBundle{
			PeerID:     peerID,
			OCR2Key:    *ocrKey,
			TXKey:      *txKey,
			P2PKeys:    *p2pkeys,
			EthAddress: ethAddress,
		}
		nkb = append(nkb, bundle)
		clNodes = append(clNodes, &CLNodesWithKeys{Node: n, KeysBundle: bundle})
	}

	return nkb, clNodes, nil
}

// TrackForwarder track forwarder address in db.
func (c *ChainlinkClient) TrackForwarder(chainID *big.Int, address common.Address) (*Forwarder, *http.Response, error) {
	response := &Forwarder{}
	request := ForwarderAttributes{
		ChainID: chainID.String(),
		Address: address.Hex(),
	}
	c.l.Debug().Str(NodeURL, c.Config.URL).
		Str("Forwarder address", (address).Hex()).
		Str("Chain ID", chainID.String()).
		Msg("Track forwarder")
	resp, err := c.APIClient.R().
		SetBody(request).
		SetResult(response).
		Post("/v2/nodes/evm/forwarders/track")
	if err != nil {
		return nil, nil, err
	}
	err = VerifyStatusCode(resp.StatusCode(), http.StatusCreated)
	if err != nil {
		return nil, nil, err
	}

	return response, resp.RawResponse, err
}

// GetForwarders get list of tracked forwarders
func (c *ChainlinkClient) GetForwarders() (*Forwarders, *http.Response, error) {
	response := &Forwarders{}
	c.l.Info().Str(NodeURL, c.Config.URL).Msg("Reading Tracked Forwarders")
	resp, err := c.APIClient.R().
		SetResult(response).
		Get("/v2/nodes/evm/forwarders")
	if err != nil {
		return nil, nil, err
	}
	err = VerifyStatusCode(resp.StatusCode(), http.StatusOK)
	if err != nil {
		return nil, nil, err
	}
	return response, resp.RawResponse, err
}

// Replays log poller from block number
func (c *ChainlinkClient) ReplayLogPollerFromBlock(fromBlock, evmChainID int64) (*ReplayResponse, *http.Response, error) {
	specObj := &ReplayResponse{}
	c.l.Info().Str(NodeURL, c.Config.URL).Int64("From block", fromBlock).Int64("EVM chain ID", evmChainID).Msg("Replaying Log Poller from block")
	resp, err := c.APIClient.R().
		SetResult(&specObj).
		SetQueryParams(map[string]string{
			"evmChainID": strconv.FormatInt(evmChainID, 10),
		}).
		SetPathParams(map[string]string{
			"fromBlock": strconv.FormatInt(fromBlock, 10),
		}).
		Post("/v2/replay_from_block/{fromBlock}")
	if err != nil {
		return nil, nil, err
	}

	return specObj, resp.RawResponse, err
}
