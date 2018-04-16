package scraper

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("_examples/challenge.html")
		if err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Server", "cloudflare-nginx")
		w.WriteHeader(503)
		w.Write(b)
	}))
	defer ts.Close()

	scraper, err := NewTransport(http.DefaultTransport)
	if err != nil {
		t.Fatal(err)
	}

	c := http.Client{
		Transport: scraper,
	}

	// must be
	// https://socket.bittrex.com/cdn-cgi/l/chk_jschl?
	// jschl_vc=5530cb027b54af31d8771d96f7b93817&pass=1523904262.081-6A4JAY%2B1KA&jschl_answer=22.5755820428

	res, err := c.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
}
