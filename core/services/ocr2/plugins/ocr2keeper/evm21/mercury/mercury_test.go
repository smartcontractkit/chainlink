package mercury

import (
	"testing"
)

func TestGenerateHMACFn(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		path     string
		body     []byte
		clientId string
		secret   string
		ts       int64
		expected string
	}{
		{
			name:     "generate hmac function",
			method:   "GET",
			path:     "/example",
			body:     []byte(""),
			clientId: "yourClientId",
			secret:   "yourSecret",
			ts:       1234567890,
			expected: "17b0bb6b14f7b48ef9d24f941ff8f33ad2d5e94ac343380be02c2f1ca32fdbd8",
		},
		{
			name:     "generate hmac function with non-empty body",
			method:   "POST",
			path:     "/api",
			body:     []byte("request body"),
			clientId: "anotherClientId",
			secret:   "anotherSecret",
			ts:       1597534567,
			expected: "d326c168c50c996e271d6b3b4c97944db01163994090f73fcf4fd42f23f06bbb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateHMACFn(tc.method, tc.path, tc.body, tc.clientId, tc.secret, tc.ts)

			if result != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, result)
			}
		})
	}
}
