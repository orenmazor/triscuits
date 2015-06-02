package main

import "os"
import "fmt"
import "github.com/gorilla/mux"
import "errors"
import "net/http"
import "io/ioutil"
import "net/url"
import "github.com/orenmazor/hmaclib"

func main() {
	fmt.Println("listening on 31337. bring it!")
	http.Handle("/", triscuits())
	http.ListenAndServe("0.0.0.0:31337", nil)
}

func triscuits() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/trusted", TicketRequestHandler)
	return r
}

func getUsername(req *http.Request) (string, error) {
	req.ParseForm()

	if len(req.Form["username"]) == 1 {
		return req.Form["username"][0], nil
	} else {
		return "", errors.New("missing parameter 'username'")
	}
}

func TicketRequestHandler(w http.ResponseWriter, req *http.Request) {
	username, err := getUsername(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if authorized_request(username, req.Header["X-Triscuits-Auth"]) {
		ticket := generate_ticket(username)
		fmt.Fprint(w, ticket)
	} else {
		http.Error(w, "nope", http.StatusUnauthorized)
	}
}

func authorized_request(username string, headers []string) bool {
	if len(headers) == 0 {
		return false
	}

	auth_header := headers[0]
	return hmaclib.CheckHMAC([]byte(username), auth_header, []byte(os.Getenv("TRISCUITS_HMAC")))
}

func generate_ticket(user string) string {
	resp, _ := http.PostForm(os.Getenv("TABLEAU_URL"), url.Values{"username": {user}})
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
