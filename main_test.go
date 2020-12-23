package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var buf = &bytes.Buffer{}

func init() {
	writer := io.MultiWriter(os.Stdout, buf)
	logger = log.New(writer, "", log.Lmsgprefix)
}

func expectLog(msg string, buf *bytes.Buffer, t *testing.T) {
	line, _ := buf.ReadString('\n')
	if line[:len(line)-1] != msg {
		t.Fatalf("Expect '%s' but was:\t%s", msg, line)
	}
}

func TestWebhookOpenPR(t *testing.T) {
	defer buf.Reset()

	dat, err := ioutil.ReadFile("test/pr_open.json")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(dat))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	webhook(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "" {
		t.Fatal(err)
	}

	expectLog("POST /webhook", buf, t)
	expectLog("Opened PR. &{TEST example refs/heads/PR 2}", buf, t)
}

func TestWebhookUpdateOpenPR(t *testing.T) {
	defer buf.Reset()

	dat, err := ioutil.ReadFile("test/push_commit.json")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(dat))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	webhook(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "" {
		t.Fatal(err)
	}

	expectLog("POST /webhook", buf, t)
	expectLog("Refs in repo changed.", buf, t)
	expectLog("Updated open PR: &{TEST example refs/heads/PR 2}", buf, t)
}
