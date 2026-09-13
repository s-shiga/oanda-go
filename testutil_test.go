package oanda

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
)

// Tests must inject a fake HTTPClient. Even an accidentally unconfigured client
// cannot reach the API, regardless of credentials in the environment.
func TestMain(m *testing.M) {
	http.DefaultTransport = offlineTransport{}
	os.Exit(m.Run())
}

type offlineTransport struct{}

func (offlineTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("network disabled in tests: %s %s", req.Method, req.URL.Redacted())
}

const testAccountPath = "/v3/accounts/101-001-1234567-001"

type endpointTest struct {
	name     string
	method   string
	path     string
	query    string
	body     string
	status   int
	response string
	// want overrides the expected decoded JSON when a service unwraps a response.
	want string
	call func(*Client) (any, error)
}

func runEndpointTests(t *testing.T, cases []endpointTest) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status := tc.status
			if status == 0 {
				status = http.StatusOK
			}
			fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(status, tc.response)}}
			result, err := tc.call(newFakeClient(fake))
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if len(fake.requests) != 1 {
				t.Fatalf("got %d requests, want 1", len(fake.requests))
			}
			req := fake.requests[0]
			if req.Method != tc.method || req.URL.Path != tc.path {
				t.Errorf("request = %s %s, want %s %s", req.Method, req.URL.Path, tc.method, tc.path)
			}
			wantQuery, err := url.ParseQuery(tc.query)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(req.URL.Query(), wantQuery) {
				t.Errorf("query = %v, want %v", req.URL.Query(), wantQuery)
			}
			if got := req.Header.Get("Authorization"); got != "Bearer test-key" {
				t.Errorf("Authorization = %q", got)
			}
			if got := req.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q", got)
			}
			if tc.body == "" {
				if fake.bodies[0] != "" {
					t.Errorf("unexpected request body: %s", fake.bodies[0])
				}
			} else {
				got, want := decodeTestJSON(t, fake.bodies[0]), decodeTestJSON(t, tc.body)
				if !reflect.DeepEqual(got, want) {
					t.Errorf("body = %s, want %s", fake.bodies[0], tc.body)
				}
			}
			encoded, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			want := tc.want
			if want == "" {
				want = tc.response
			}
			// Check every fixture field; additional zero-valued Go fields are allowed.
			assertJSONFields(t, "$", decodeTestJSON(t, string(encoded)), decodeTestJSON(t, want))
		})
	}
}

func decodeTestJSON(t *testing.T, value string) any {
	t.Helper()
	var result any
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		t.Fatalf("invalid test JSON %q: %v", value, err)
	}
	return result
}

func assertJSONFields(t *testing.T, path string, got, want any) {
	t.Helper()
	switch expected := want.(type) {
	case map[string]any:
		actual, ok := got.(map[string]any)
		if !ok {
			t.Errorf("%s = %#v, want object", path, got)
			return
		}
		for key, value := range expected {
			field, exists := actual[key]
			if !exists {
				t.Errorf("%s.%s missing from decoded response", path, key)
				continue
			}
			assertJSONFields(t, path+"."+key, field, value)
		}
	case []any:
		actual, ok := got.([]any)
		if !ok || len(actual) != len(expected) {
			t.Errorf("%s = %#v, want array of length %d", path, got, len(expected))
			return
		}
		for i, value := range expected {
			assertJSONFields(t, fmt.Sprintf("%s[%d]", path, i), actual[i], value)
		}
	default:
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s = %#v, want %#v", path, got, want)
		}
	}
}
