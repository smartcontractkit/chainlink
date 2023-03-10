package mercuryserver

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	mercuryserverhelm "github.com/smartcontractkit/chainlink-env/pkg/helm/mercury-server"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/test-go/testify/require"
)

type csaKey struct {
	NodeName    string `json:"nodeName"`
	NodeAddress string `json:"nodeAddress"`
	PublicKey   string `json:"publicKey"`
}

type oracle struct {
	Id                    string   `json:"id"`
	Website               string   `json:"website"`
	Name                  string   `json:"name"`
	Status                string   `json:"status"`
	NodeAddress           []string `json:"nodeAddress"`
	OracleAddress         string   `json:"oracleAddress"`
	CsaKeys               []csaKey `json:"csaKeys"`
	Ocr2ConfigPublicKey   []string `json:"ocr2ConfigPublicKey"`
	Ocr2OffchainPublicKey []string `json:"ocr2OffchainPublicKey"`
	Ocr2OnchainPublicKey  []string `json:"ocr2OnchainPublicKey"`
}

// Build config with nodes for Mercury server
func BuildRpcNodesJsonConf(t *testing.T, chainlinkNodes []*client.Chainlink) ([]byte, error) {
	var msRpcNodesConf []*oracle
	for i, chainlinkNode := range chainlinkNodes {
		nodeName := fmt.Sprint(i)
		nodeAddress, err := chainlinkNode.PrimaryEthAddress()
		require.NoError(t, err)
		csaKeys, _, err := chainlinkNode.ReadCSAKeys()
		require.NoError(t, err)
		csaPubKey := csaKeys.Data[0].Attributes.PublicKey
		ocr2Keys, resp, err := chainlinkNode.ReadOCR2Keys()
		_ = ocr2Keys
		_ = resp
		require.NoError(t, err)
		var ocr2Config client.OCR2KeyAttributes
		for _, key := range ocr2Keys.Data {
			if key.Attributes.ChainType == string(chaintype.EVM) {
				ocr2Config = key.Attributes
				break
			}
		}
		ocr2ConfigPublicKey := strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_")
		ocr2OffchainPublicKey := strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_")
		ocr2OnchainPublicKey := strings.TrimPrefix(ocr2Config.OnChainPublicKey, "ocr2on_evm_")

		node := &oracle{
			Id:            fmt.Sprint(i),
			Name:          nodeName,
			Status:        "active",
			NodeAddress:   []string{nodeAddress},
			OracleAddress: "0x0000000000000000000000000000000000000000",
			CsaKeys: []csaKey{
				{
					NodeName:    nodeName,
					NodeAddress: nodeAddress,
					PublicKey:   csaPubKey,
				},
			},
			Ocr2ConfigPublicKey:   []string{ocr2ConfigPublicKey},
			Ocr2OffchainPublicKey: []string{ocr2OffchainPublicKey},
			Ocr2OnchainPublicKey:  []string{ocr2OnchainPublicKey},
		}
		msRpcNodesConf = append(msRpcNodesConf, node)
	}
	return json.Marshal(msRpcNodesConf)
}

func SetupMercuryServer(
	t *testing.T,
	testEnv *environment.Environment,
	dbSettings map[string]interface{},
	serverSettings map[string]interface{},
	adminId string,
	adminEncryptedKey string,
) ed25519.PublicKey {
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnv)
	require.NoError(t, err, "Error connecting to Chainlink nodes")

	rpcNodesJsonConf, _ := BuildRpcNodesJsonConf(t, chainlinkNodes)
	log.Info().Msgf("RPC nodes conf for mercury server: %s", rpcNodesJsonConf)

	// Generate keys for Mercury RPC server
	// rpcPrivKey, rpcPubKey, err := generateEd25519Keys()
	rpcPubKey, rpcPrivKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	initDbSql, err := buildInitialDbSql(adminId, adminEncryptedKey)
	require.NoError(t, err)
	log.Info().Msgf("Initialize mercury server db with:\n%s", initDbSql)

	settings := map[string]interface{}{
		"image": map[string]interface{}{
			"repository": os.Getenv("MERCURY_SERVER_IMAGE"),
			"tag":        os.Getenv("MERCURY_SERVER_TAG"),
		},
		"postgresql": map[string]interface{}{
			"enabled": true,
		},
		"qa": map[string]interface{}{
			"rpcPrivateKey": hex.EncodeToString(rpcPrivKey),
			"enabled":       true,
			"initDbSql":     initDbSql,
		},
		"rpcNodesConf": string(rpcNodesJsonConf),
		"prometheus":   "true",
	}

	if dbSettings != nil {
		settings["db"] = dbSettings
	}
	if serverSettings != nil {
		settings["resources"] = serverSettings
	}

	testEnv.AddHelm(mercuryserverhelm.New(settings)).Run()

	return rpcPubKey
}

func buildInitialDbSql(adminId string, adminEncryptedKey string) (string, error) {
	data := struct {
		UserId       string
		UserRole     string
		EncryptedKey string
	}{
		UserId:       adminId,
		UserRole:     "admin",
		EncryptedKey: adminEncryptedKey,
	}

	// Get file path to the sql
	_, filename, _, _ := runtime.Caller(0)
	tmplPath := path.Join(path.Dir(filename), "/mercury_db_init_sql_template")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
