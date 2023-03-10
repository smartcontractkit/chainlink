package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type GetReportsResult struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

type User struct {
	Id        string `json:"id"`
	Secret    string `json:"secret" db:"secret"`
	Role      string `json:"role" db:"role"` // 0 = user, 1 = admin
	Disabled  bool   `json:"disabled" db:"disabled"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

type MercuryServer struct {
	URL       string
	APIClient *resty.Client
}

func NewMercuryServer(url string) *MercuryServer {
	rc := resty.New().SetBaseURL(url)
	return &MercuryServer{
		URL:       url,
		APIClient: rc,
	}
}

func (ms *MercuryServer) CallGet(path string, userId string, userSecret string) (map[string]interface{}, *http.Response, error) {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
	hmacSignature := generateHmacSignature("GET", path, []byte{}, []byte(userSecret), userId, timestamp)
	result := map[string]interface{}{}
	resp, err := ms.APIClient.R().
		SetHeader("Authorization", userId).
		SetHeader("X-Authorization-Timestamp", timestamp).
		SetHeader("X-Authorization-Signature-SHA256", hmacSignature).
		SetResult(&result).
		Get(path)
	if err != nil {
		return nil, nil, err
	}
	return result, resp.RawResponse, nil
}

// Add new user with "admin" or "user" role
func (ms *MercuryServer) AddUser(adminId string, adminSecret string, newUserSecret string, newUserRole string, newUserDisabled bool) (*User, *http.Response, error) {
	request := map[string]interface{}{
		"secret":   newUserSecret,
		"role":     newUserRole,
		"disabled": newUserDisabled,
	}
	result := struct {
		User User
	}{}
	path := "/admin/user"
	timestamp := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
	b, _ := json.Marshal(request)
	hmacSignature := generateHmacSignature("POST", path, b, []byte(adminSecret), adminId, timestamp)
	resp, err := ms.APIClient.R().
		SetHeader("Authorization", adminId).
		SetHeader("X-Authorization-Timestamp", timestamp).
		SetHeader("X-Authorization-Signature-SHA256", hmacSignature).
		SetBody(request).
		SetResult(result).
		Post(path)
	if err != nil {
		return nil, nil, err
	}
	return &result.User, resp.RawResponse, err
}

func (ms *MercuryServer) GetUsers(adminId string, adminSecret string) (*[]User, *http.Response, error) {
	var result []User
	path := "/admin/user"
	timestamp := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
	hmacSignature := generateHmacSignature("GET", path, []byte{}, []byte(adminSecret), adminId, timestamp)
	resp, err := ms.APIClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", adminId).
		SetHeader("X-Authorization-Timestamp", timestamp).
		SetHeader("X-Authorization-Signature-SHA256", hmacSignature).
		SetResult(&result).
		Get(path)
	if err != nil {
		return nil, nil, err
	}
	return &result, resp.RawResponse, err
}

func (ms *MercuryServer) GetReports(userId string, userSecret string, feedIDStr string, blockNumber uint64) (*GetReportsResult, *http.Response, error) {
	result := &GetReportsResult{}
	path := fmt.Sprintf("/client?feedIDStr=%s&L2Blocknumber=%d", feedIDStr, blockNumber)
	timestamp := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
	hmacSignature := generateHmacSignature("GET", path, []byte{}, []byte(userSecret), userId, timestamp)
	resp, err := ms.APIClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", userId).
		SetHeader("X-Authorization-Timestamp", timestamp).
		SetHeader("X-Authorization-Signature-SHA256", hmacSignature).
		SetResult(&result).
		Get(path)
	if err != nil {
		return nil, resp.RawResponse, err
	}
	return result, resp.RawResponse, err
}

func generateHmacSignature(method string, path string, body []byte, secret []byte, clientId string, timestamp string) string {
	// Get the hash for the body
	bodyHash := sha256.New()
	bodyHash.Write(body)
	bodyHashString := hex.EncodeToString(bodyHash.Sum(nil))

	// Generate the message to be signed
	message := fmt.Sprintf("%s %s %s %s %s", method, path, bodyHashString, clientId, timestamp)
	log.Info().Msgf("message: %s", message)

	// Generate the signature
	signedMessage := hmac.New(sha256.New, secret)
	signedMessage.Write([]byte(message))
	return hex.EncodeToString(signedMessage.Sum(nil))
}
