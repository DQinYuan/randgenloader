package randgenloader

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func SessionTest(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		panic("session error")
	}

	r.ParseForm()

	log.Printf("len: %d", len(r.Form["t"]))

	if session.IsNew {
		w.Write([]byte("do not have session"))
	}
}

func TestSession(t *testing.T) {
	http.HandleFunc("/one", http.HandlerFunc(SessionTest))

	http.ListenAndServe(":9090", nil)
}

func TestMashal(t *testing.T) {
	r := &resultStruct{"haha", "'testname' required"}
	bytes, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))
}