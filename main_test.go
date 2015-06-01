package main

import "net/http/httptest"
import "io"
import "io/ioutil"
import "testing"
import "fmt"
import "net/http"
import "os"
import "bytes"
import "net/url"

var (
	server     *httptest.Server
	reader     io.Reader //Ignore this for now
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

func TestTrustedEndpoint(t *testing.T) {
	payload := url.Values{"username": {"data_portal"}}
	b := bytes.NewBufferString(payload.Encode())

	request, _ := http.NewRequest("POST", trustedUrl, b)
	request.Header.Set("X-Triscuits-Auth", "adventuretime")
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
	request, err := http.NewRequest("POST", trustedUrl, nil)
	code, body := send(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	if code != 401 {
		t.Errorf("Expected 401, got a %d", code)
	}

	expectedBody := "nope\n"
	if body != expectedBody {
		t.Errorf("Expected '%s', got a '%s'", expectedBody, body)
	}
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

func TestTrustedEndpointFailsWithoutUsername(t *testing.T) {
	request, err := http.NewRequest("POST", trustedUrl, nil)
	request.Header.Set("X-Triscuits-Auth", "adventuretime")
	code, body := send(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	if code != 400 {
		t.Errorf("Expected 400, got a %d", code)
	}

	expectedBody := "missing parameter 'username'\n"
	if body != expectedBody {
		t.Errorf("Expected '%s', got a '%s'", expectedBody, body)
	}
}
