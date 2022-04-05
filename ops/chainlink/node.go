package chainlink

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink-relay/ops/utils"
	"github.com/smartcontractkit/integrations-framework/client"
)

// Node implements the node parameters
type Node struct {
	Name   string
	P2P    client.P2PData
	Config client.ChainlinkConfig
	Call   client.Chainlink
	Keys   NodeKeys
}

// NodeKeys implements the needed keys from the core node
type NodeKeys struct {
	OCR2KeyID             string
	OCR2OnchainPublicKey  string
	OCR2Transmitter       string
	OCR2TransmitterID     string
	OCR2OffchainPublicKey string
	OCR2ConfigPublicKey   string
	P2PID                 string
}

// CredentialsString returns formatted string for node input
func (n *Node) CredentialsString() string {
	return fmt.Sprintf("%s\n%s", n.Config.Email, n.Config.Password)
}

// Health returns if the node is functional or not
func (n *Node) Health() (interface{}, error) {
	return http.Get(n.Config.URL + "/health")
}

// Ready checks when node is ready
func (n *Node) Ready() error {
	msg := utils.LogStatus(fmt.Sprintf("Waiting for health checks on %s", n.Name))
	timeout := 300
	var err error
	time.Sleep(2 * time.Second) // removing this breaks running `up` multiple times...
	for i := 0; i < timeout; i++ {
		_, err = n.Health()
		if err == nil {
			cl, err := client.NewChainlink(&n.Config, http.DefaultClient)
			n.Call = cl
			return msg.Check(err)
		}
		time.Sleep(1 * time.Second)
	}
	return msg.Check(err)
}

// AddBridge adds adapter to CL node
func (n *Node) AddBridge(name, url string) error {
	msg := utils.LogStatus(fmt.Sprintf("Adding %s EA to %s", name, n.Name))

	// check if exists
	_, err := n.Call.ReadBridge(name)
	if err == nil {
		msg.Exists()
		return msg.Check(nil)
	}

	err = n.Call.CreateBridge(&client.BridgeTypeAttributes{
		Name: name,
		URL:  url,
	})
	return msg.Check(err)
}

func (n Node) DeleteAllJobs() error {
	msg := utils.LogStatus("Cleared existing jobs from CL node")

	// get all jobs
	jobs, err := n.Call.ReadJobs()
	for _, j := range jobs.Data {
		// remove job based on ID
		if err := n.Call.DeleteJob(j["id"].(string)); err != nil {
			return msg.Check(err)
		}
	}

	if len(jobs.Data) == 0 {
		fmt.Print(" - No jobs present")
	}
	return msg.Check(err)
}

// TODO: verify does this work for evm and other chains
func (n *Node) GetKeys(chain string) error {
	msg := utils.LogStatus(fmt.Sprintf("Retrieved keys from %s", n.Name))

	ocrKey, err := n.fetchOCR2Key(chain, true)
	if err != nil {
		return msg.Check(err)
	}

	p2pKeys, err := n.Call.ReadP2PKeys()
	if err != nil {
		return msg.Check(err)
	}

	txKey, err := n.fetchTxKey(chain, true)
	if err != nil {
		return msg.Check(err)
	}

	// parse keys into expected format
	n.Keys = NodeKeys{
		OCR2KeyID:             utils.RemovePrefix(ocrKey.ID),
		OCR2OnchainPublicKey:  utils.RemovePrefix(ocrKey.Attributes.OnChainPublicKey),
		OCR2Transmitter:       utils.RemovePrefix(txKey.Attributes.PublicKey),
		OCR2TransmitterID:     utils.RemovePrefix(txKey.ID),
		OCR2OffchainPublicKey: utils.RemovePrefix(ocrKey.Attributes.OffChainPublicKey),
		OCR2ConfigPublicKey:   utils.RemovePrefix(ocrKey.Attributes.ConfigPublicKey),
		P2PID:                 utils.RemovePrefix(p2pKeys.Data[0].Attributes.PeerID),
	}
	n.P2P.PeerID = n.Keys.P2PID // set p2p peerID in the p2p struct

	return msg.Check(err)
}

func (n Node) fetchOCR2Key(chain string, create bool) (ocrKey client.OCR2KeyData, err error) {
	ocrKeys, err := n.Call.ReadOCR2Keys()
	if err != nil {
		return ocrKey, err
	}
	// parse first key that matches chain
	for _, k := range ocrKeys.Data {
		if k.Attributes.ChainType == chain {
			ocrKey = k
			break
		}
	}

	// return if key is found
	if ocrKey != (client.OCR2KeyData{}) {
		return ocrKey, nil
	}

	// return error if key not found and not creating
	if !create {
		return ocrKey, fmt.Errorf("could not find ocr2 key for %s", chain)
	}

	// create key
	if _, err = n.Call.CreateOCR2Key(chain); err != nil {
		return ocrKey, err
	}
	return n.fetchOCR2Key(chain, false) // fetch key again
}

func (n Node) fetchTxKey(chain string, create bool) (client.TxKeyData, error) {
	// check if key exists
	txKeys, err := n.Call.ReadTxKeys(chain) // this doesn't work for evm (uses a different function)
	if err != nil {
		return client.TxKeyData{}, nil
	}

	// return if key is found
	if len(txKeys.Data) > 0 {
		return txKeys.Data[0], nil
	}

	if !create { // if not told to create, then return not found err
		return client.TxKeyData{}, errors.New("no transaction key found")
	}

	// if key is not found, create
	if _, err := n.Call.CreateTxKey(chain); err != nil {
		return client.TxKeyData{}, err
	}

	// after creating fetch again without create
	return n.fetchTxKey(chain, false)
}
