package main

import "fmt"
import "os"
import "github.com/gorilla/mux"
import "net/http"
import "io/ioutil"
import "net/url"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/trusted", TicketRequestHandler)
	fmt.Println("listening on 31337. bring it!")
	http.Handle("/", r)
	http.ListenAndServe(":31337", nil)
}

func TicketRequestHandler(w http.ResponseWriter, req *http.Request) {
	if authorized_request(req.Header["X-Triscuits-Auth"]) {
		req.ParseForm()

		ticket := generate_ticket(req.Form["username"][0])
		fmt.Fprint(w, ticket)
	} else {
		http.Error(w, "nope", http.StatusUnauthorized)
	}
}

func authorized_request(headers []string) bool {
	if len(headers) == 0 {
		return false
	}

	auth_header := headers[0]
	expected_hmac := os.Getenv("TRISCUITS_HMAC")
	return expected_hmac == auth_header
}

func generate_ticket(user string) string {
	resp, _ := http.PostForm(os.Getenv("TABLEAU_URL"), url.Values{"username": {user}})
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
