package client

import (
	"context"
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
	"nhooyr.io/websocket"
)

type GetReportResult struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

type GetReportsResult struct {
	ChainlinkBlob []string `json:"chainlinkBlob"`
}

type User struct {
	Id        string `json:"id"`
	Secret    string `json:"secret" db:"secret"`
	Role      string `json:"role" db:"role"` // 0 = user, 1 = admin
	Disabled  bool   `json:"disabled" db:"disabled"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

type NewReportWSMessage struct {
	FeedId     []byte `json:"feedId"`
	FullReport []byte `json:"report"`
}

type WebsocketConnectQuery struct {
	FeedIds []string `form:"feedIds"`
}

type MercuryServer struct {
	URL       string
	UserId    string
	UserKey   string
	APIClient *resty.Client
}

// Create new mercury server client for userId and userKey that are used for HMAC authentication
func NewMercuryServerClient(url string, userId string, userKey string) *MercuryServer {
	rc := resty.New().SetBaseURL(url)
	return &MercuryServer{
		URL:       url,
		APIClient: rc,
		UserId:    userId,
		UserKey:   userKey,
	}
}

func (s *MercuryServer) DialWS(ctx context.Context, urlParams string) (*websocket.Conn, *http.Response, error) {
	timestamp := genReqTimestamp()
	var path string
	if urlParams != "" {
		path = fmt.Sprintf("/ws%s", urlParams)
	} else {
		path = "/ws"
	}
	hmacSignature := genHmacSignature("GET", path, []byte{}, []byte(s.UserKey), s.UserId, timestamp)
	return websocket.Dial(ctx, fmt.Sprintf("%s%s", s.URL, path), &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Authorization":                    []string{s.UserId},
			"X-Authorization-Timestamp":        []string{timestamp},
			"X-Authorization-Signature-SHA256": []string{hmacSignature},
		},
	})
}

func (s *MercuryServer) CallGet(path string) (map[string]interface{}, *http.Response, error) {
	timestamp := genReqTimestamp()
	hmacSignature := genHmacSignature("GET", path, []byte{}, []byte(s.UserKey), s.UserId, timestamp)
	result := map[string]interface{}{}
	resp, err := s.APIClient.R().
		SetHeader("Authorization", s.UserId).
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
func (s *MercuryServer) AddUser(newUserSecret string, newUserRole string, newUserDisabled bool) (*User, *http.Response, error) {
	request := map[string]interface{}{
		"secret":   newUserSecret,
		"role":     newUserRole,
		"disabled": newUserDisabled,
	}
	result := struct {
		User User
	}{}
	path := "/admin/user"
	timestamp := genReqTimestamp()
	b, _ := json.Marshal(request)
	hmacSignature := genHmacSignature("POST", path, b, []byte(s.UserKey), s.UserId, timestamp)
	resp, err := s.APIClient.R().
		SetHeader("Authorization", s.UserId).
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

// Need admin role
func (s *MercuryServer) GetUsers() ([]User, *http.Response, error) {
	var result []User
	path := "/admin/user"
	timestamp := genReqTimestamp()
	hmacSignature := genHmacSignature("GET", path, []byte{}, []byte(s.UserKey), s.UserId, timestamp)
	resp, err := s.APIClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", s.UserId).
		SetHeader("X-Authorization-Timestamp", timestamp).
		SetHeader("X-Authorization-Signature-SHA256", hmacSignature).
		SetResult(&result).
		Get(path)
	if err != nil {
		return nil, nil, err
	}
	return result, resp.RawResponse, err
}

type FeedIdType string

const (
	StringFeedId FeedIdType = "string"
	HexFeedId    FeedIdType = "hex"
)

func (s *MercuryServer) BulkGetReportsByFeedId(feedId string, afterBlockNumber uint64, limit uint64, feedIdType FeedIdType) (*GetReportsResult, *http.Response, error) {
	result := &GetReportsResult{}
	path := fmt.Sprintf("/client/bulk?afterBlockNumber=%d&limit=%d", afterBlockNumber, limit)
	if feedIdType == StringFeedId {
		path = fmt.Sprintf("%s&feedIDStr=%s", path, feedId)
	} else if feedIdType == HexFeedId {
		path = fmt.Sprintf("%s&feedIDHex=%s", path, feedId)
	} else {
		return nil, nil, fmt.Errorf("%s not supported", feedIdType)
	}
	timestamp := genReqTimestamp()
	hmacSignature := genHmacSignature("GET", path, []byte{}, []byte(s.UserKey), s.UserId, timestamp)
	resp, err := s.APIClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", s.UserId).
		SetHeader("X-Authorization-Timestamp", timestamp).
		SetHeader("X-Authorization-Signature-SHA256", hmacSignature).
		SetResult(&result).
		Get(path)
	if err != nil && resp == nil {
		return nil, nil, err
	}
	if err != nil {
		return nil, resp.RawResponse, err
	}
	return result, resp.RawResponse, err
}

func (s *MercuryServer) GetReportsByFeedId(feedId string, blockNumber uint64, feedIdType FeedIdType) (*GetReportResult, *http.Response, error) {
	result := &GetReportResult{}
	path := fmt.Sprintf("/client?blockNumber=%d", blockNumber)
	if feedIdType == StringFeedId {
		path = fmt.Sprintf("%s&feedIDStr=%s", path, feedId)
	} else if feedIdType == HexFeedId {
		path = fmt.Sprintf("%s&feedIDHex=%s", path, feedId)
	} else {
		return nil, nil, fmt.Errorf("%s not supported", feedIdType)
	}
	timestamp := genReqTimestamp()
	hmacSignature := genHmacSignature("GET", path, []byte{}, []byte(s.UserKey), s.UserId, timestamp)
	resp, err := s.APIClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", s.UserId).
		SetHeader("X-Authorization-Timestamp", timestamp).
		SetHeader("X-Authorization-Signature-SHA256", hmacSignature).
		SetResult(&result).
		Get(path)
	if err != nil && resp == nil {
		return nil, nil, err
	}
	if err != nil {
		return nil, resp.RawResponse, err
	}
	return result, resp.RawResponse, err
}

func genReqTimestamp() string {
	// The timestamp of the request. This is used to prevent replay attacks.
	// The timestamp should be within 5 seconds of the server's time (by default).
	// The server will reject requests with timestamps in the future.
	return strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
}

func genHmacSignature(method string, path string, body []byte, secret []byte, clientId string, timestamp string) string {
	// Get the hash for the body
	bodyHash := sha256.New()
	bodyHash.Write(body)
	bodyHashString := hex.EncodeToString(bodyHash.Sum(nil))

	// Generate the message to be signed
	message := fmt.Sprintf("%s %s %s %s %s", method, path, bodyHashString, clientId, timestamp)
	log.Debug().Msgf("message: %s", message)

	// Generate the signature
	signedMessage := hmac.New(sha256.New, secret)
	signedMessage.Write([]byte(message))
	return hex.EncodeToString(signedMessage.Sum(nil))
}
