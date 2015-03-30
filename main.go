package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

var jenkinsUrl = strings.TrimSuffix(os.Getenv("JENKINS_URL"), "/")
var proxyPort = os.Getenv("PROXY_PORT")

func httpError(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}

type Commit struct {
	Branch string
}

type BBMessage struct {
	Commits []Commit
}

func handler(w http.ResponseWriter, r *http.Request) {
	d, _ := httputil.DumpRequest(r, true)
	log.Println(string(d))
	rawurl := string(r.URL.Path)
	rawurl = strings.TrimPrefix(strings.TrimSuffix(rawurl, "/"), "/")
	itemsInPath := strings.Split(rawurl, "/")
	if len(itemsInPath) != 3 {
		httpError(w, "access <job>/<branch>[,<branch>]/<auth_token>*")
		return
	}
	job, branches, authToken := itemsInPath[0], itemsInPath[1], itemsInPath[2]

	var payloadDict BBMessage

	payload := make([]byte, r.ContentLength)
	_, err := r.Body.Read(payload)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	log.Println(string(payload))

	err = json.Unmarshal(payload, &payloadDict)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	// the following is a bit lame, but I prefer the idea of sets to lists
	testedBranches := make(map[string]struct{})
	for _, b := range strings.Split(branches, ",") {
		testedBranches[b] = struct{}{}
	}
	commitedBranches := make(map[string]struct{})
	for _, commit := range payloadDict.Commits {
		commitedBranches[commit.Branch] = struct{}{}
	}
	log.Println("branches for which we should run tests:", testedBranches)
	log.Println("branches touched in the push for this message:", commitedBranches)
	runJob := false

	// intersection of keys in the two maps
	for b, _ := range commitedBranches {
		if _, shouldBeTested := testedBranches[b]; shouldBeTested {
			runJob = true
			break
		}
	}
	if runJob {
		log.Println("we will attempt to run the job")
		client := &http.Client{}
		newUrl := jenkinsUrl + "/job/" + job + "/build?token=" + authToken
		newRequest, err := http.NewRequest("POST", newUrl, bytes.NewBuffer(payload))
		resp, err := client.Do(newRequest)
		if err != nil {
			httpError(w, "error while forwarding the request: "+err.Error())
			return
		}
		defer resp.Body.Close()

	} else {
		log.Println("this push will not run the job")
	}
}

func main() {
	if len(jenkinsUrl) == 0 || len(proxyPort) == 0 {
		log.Fatal("You must define PROXY_PORT and JENKINS_URL env vars.")
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+proxyPort, nil)
}
