package main

import "net/http/httptest"
import "io/ioutil"
import "testing"
import "fmt"
import "net/http"
import "os"
import "bytes"
import "net/url"
import "encoding/base64"

var (
	server     *httptest.Server
	trustedUrl string
)

func init() {
	//mock tableau server
	tableau_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World. U are secured.")
	}))

	os.Setenv("TRISCUITS_HMAC", "adventuretime")
	os.Setenv("TABLEAU_URL", fmt.Sprintf(tableau_server.URL))
	server = httptest.NewServer(triscuits()) //Creating new server with the user handlers

	trustedUrl = fmt.Sprintf("%s/trusted", server.URL)
}

func send(request *http.Request) (int, string) {
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return res.StatusCode, string(body)
}

func TestDecodingHMAC(t *testing.T) {
	expected_hmac := DecodeHMAC(CalculateHMAC("hello world"))
	decoded_message := DecodeHMAC(base64.StdEncoding.EncodeToString(expected_hmac))
	if bytes.Compare(expected_hmac, decoded_message) != 0 {
		t.Errorf(fmt.Sprintf("expected to decode the message, but got %x instead", decoded_message))
	}
}

func TestCalculatingHMAC(t *testing.T) {
	expected_hmac := DecodeHMAC("mWBFMDhNoDfb9rAXjfaPM1IQMbOjitBk+tS6A6P0kTI=")
	calculated_hmac := DecodeHMAC(CalculateHMAC("hello world"))
	if bytes.Compare(expected_hmac, calculated_hmac) != 0 {
		t.Errorf("expected %v, got %v", expected_hmac, calculated_hmac)
	}
}

// func TestCheckHMAC(t *testing.T) {
// 	message := "hello world"
// 	message_hmac := CalculateHMAC(message)
// 	if CheckHMAC(message, message_hmac) {
// 		t.Errorf("failed hmac comparison. this is weird, right?")
// 	}
// }

func TestTrustedEndpoint(t *testing.T) {
	payload := url.Values{"username": {"data_portal"}}
	b := bytes.NewBufferString(payload.Encode())

	request, _ := http.NewRequest("POST", trustedUrl, b)
	request.Header.Set("X-Triscuits-Auth", CalculateHMAC("data_portal"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	code, body := send(request)

	if code != 200 {
		t.Errorf("Expected 200, got a %d", code)
	}

	expectedBody := "Hello World. U are secured.\n"
	if body != expectedBody {
		t.Errorf("Expected '%s', got a '%s'", expectedBody, body)
	}
}

func TestTrustedEndpointFailsWithoutHeader(t *testing.T) {
	payload := url.Values{"username": {"data_portal"}}
	b := bytes.NewBufferString(payload.Encode())

	request, _ := http.NewRequest("POST", trustedUrl, b)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	code, body := send(request)

	if code != 401 {
		t.Errorf("Expected 401, got a %d", code)
	}

	expectedBody := "nope\n"
	if body != expectedBody {
		t.Errorf("Expected '%s', got a '%s'", expectedBody, body)
	}
}

func TestTrustedEndpointFailsWithoutUsername(t *testing.T) {
	request, _ := http.NewRequest("POST", trustedUrl, nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	code, body := send(request)

	if code != 400 {
		t.Errorf("Expected 400, got a %d", code)
	}

	expectedBody := "missing parameter 'username'\n"
	if body != expectedBody {
		t.Errorf("Expected '%s', got a '%s'", expectedBody, body)
	}
}
