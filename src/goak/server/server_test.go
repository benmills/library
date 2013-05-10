package server

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func httpRequest(method string, url string, data string) (int, string) {
	request, _ := http.NewRequest(method, url, strings.NewReader(data))
	client := http.Client{}
	response, _ := client.Do(request)
	rawBody, _ := ioutil.ReadAll(response.Body)

	return response.StatusCode, string(rawBody)
}

type TestNode struct {
	*httptest.Server
	node *Server
}

func testServer() *TestNode {
	goakServer := New()
	return &TestNode{httptest.NewServer(goakServer.Handler()), goakServer}
}

func TestAddAKey(t *testing.T) {
	server := testServer()
	defer server.Close()

	statusCode, body := httpRequest("PUT", server.URL+"/data/mykey", "bar")

	assert.Equal(t, 201, statusCode)
	assert.Equal(t, "bar", body)
}

func TestFetchKey(t *testing.T) {
	server := testServer()
	defer server.Close()

	httpRequest("PUT", server.URL+"/data/mykey", "bar")
	statusCode, body := httpRequest("GET", server.URL+"/data/mykey", "bar")

	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "bar", body)
}

func TestFetchUnknownKey(t *testing.T) {
	server := testServer()
	defer server.Close()

	statusCode, _ := httpRequest("GET", server.URL+"/data/mykey", "bar")

	assert.Equal(t, 404, statusCode)
}

func TestUpdateKey(t *testing.T) {
	server := testServer()
	defer server.Close()

	httpRequest("PUT", server.URL+"/data/mykey", "bar")
	httpRequest("PUT", server.URL+"/data/mykey", "baz")
	statusCode, body := httpRequest("GET", server.URL+"/data/mykey", "")

	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "baz", body)
}

func TestFetchesAcrossNodes(t *testing.T) {
	server1 := testServer()
	defer server1.Close()
	server2 := testServer()
	defer server2.Close()

	httpRequest("PUT", server1.URL+"/peer", server2.URL)
	httpRequest("PUT", server2.URL+"/peer", server1.URL)

	statusCode, _ := httpRequest("PUT", server1.URL+"/data/mykey", "bar")
	assert.Equal(t, 201, statusCode)

	statusCode2, body := httpRequest("GET", server2.URL+"/data/mykey", "")
	assert.Equal(t, 200, statusCode2)
	assert.Equal(t, "bar", body)
}

func TestGetPeerWithNoPeer(t *testing.T) {
	server1 := testServer()
	defer server1.Close()

	statusCode, _ := httpRequest("GET", server1.URL+"/peer", "")
	assert.Equal(t, 404, statusCode)
}

func TestGetPeer(t *testing.T) {
	server1 := testServer()
	defer server1.Close()
	server2 := testServer()
	defer server2.Close()

	httpRequest("PUT", server1.URL+"/peer", server2.URL)

	statusCode, body := httpRequest("GET", server1.URL+"/peer", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, server2.URL, body)
}

func TestAddPeer(t *testing.T) {
	server1 := testServer()
	defer server1.Close()
	server2 := testServer()
	defer server2.Close()

	statusCode, _ := httpRequest("PUT", server1.URL+"/peer", server2.URL)
	assert.Equal(t, 201, statusCode)

	statusCode, body := httpRequest("GET", server1.URL+"/peer", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, server2.URL, body)
}
