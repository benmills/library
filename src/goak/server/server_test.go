package server

import (
	"github.com/bmizerany/assert"
	"net/http/httptest"
	"testing"

	"goak/http_client"
)

type TestNode struct {
	*httptest.Server
	node *Server
}

func testServer() *TestNode {
	goakServer := New()
	httpServer := httptest.NewServer(goakServer.Handler())
	goakServer.SetURL(httpServer.URL)

	return &TestNode{httpServer, goakServer}
}

func TestAddAKey(t *testing.T) {
	server := testServer()
	defer server.Close()

	statusCode, body := http_client.HttpRequest("PUT", server.URL+"/data/mykey", "bar")

	assert.Equal(t, 201, statusCode)
	assert.Equal(t, "bar", body)
}

func TestFetchKey(t *testing.T) {
	server := testServer()
	defer server.Close()

	http_client.HttpRequest("PUT", server.URL+"/data/mykey", "bar")
	statusCode, body := http_client.HttpRequest("GET", server.URL+"/data/mykey", "bar")

	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "bar", body)
}

func TestFetchUnknownKey(t *testing.T) {
	server := testServer()
	defer server.Close()

	statusCode, _ := http_client.HttpRequest("GET", server.URL+"/data/mykey", "bar")

	assert.Equal(t, 404, statusCode)
}

func TestUpdateKey(t *testing.T) {
	server := testServer()
	defer server.Close()

	http_client.HttpRequest("PUT", server.URL+"/data/mykey", "bar")
	http_client.HttpRequest("PUT", server.URL+"/data/mykey", "baz")
	statusCode, body := http_client.HttpRequest("GET", server.URL+"/data/mykey", "")

	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "baz", body)
}

func TestFetchesAcrossNodes(t *testing.T) {
	server1 := testServer()
	defer server1.Close()
	server2 := testServer()
	defer server2.Close()

	http_client.HttpRequest("PUT", server1.URL+"/peers", server2.URL)

	statusCode, _ := http_client.HttpRequest("PUT", server1.URL+"/data/mykey", "bar")
	assert.Equal(t, 201, statusCode)

	statusCode2, body := http_client.HttpRequest("GET", server2.URL+"/data/mykey", "")
	assert.Equal(t, 200, statusCode2)
	assert.Equal(t, "bar", body)
}
