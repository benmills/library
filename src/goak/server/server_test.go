package server

import (
	"github.com/bmizerany/assert"
	"net/http/httptest"
	"testing"
	"strings"
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

	httpRequest("PUT", server1.URL+"/peers", server2.URL)

	statusCode, _ := httpRequest("PUT", server1.URL+"/data/mykey", "bar")
	assert.Equal(t, 201, statusCode)

	statusCode2, body := httpRequest("GET", server2.URL+"/data/mykey", "")
	assert.Equal(t, 200, statusCode2)
	assert.Equal(t, "bar", body)
}

func TestGetPeerWithNoPeer(t *testing.T) {
	server1 := testServer()
	defer server1.Close()

	statusCode, _ := httpRequest("GET", server1.URL+"/peers", "")
	assert.Equal(t, 404, statusCode)
}

func TestGetPeer(t *testing.T) {
	server1 := testServer()
	defer server1.Close()

	httpRequest("PUT", server1.URL+"/peers", "peer.url")

	statusCode, body := httpRequest("GET", server1.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, `{"peers":["peer.url"]}`, body)
}

func TestGetMultiplePeers(t *testing.T) {
	server1 := testServer()
	defer server1.Close()

	httpRequest("PUT", server1.URL+"/peers", "peer1.url")
	httpRequest("PUT", server1.URL+"/peers", "peer2.url")

	statusCode, body := httpRequest("GET", server1.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, `{"peers":["peer1.url","peer2.url"]}`, body)
}

func TestAddPeer(t *testing.T) {
	server1 := testServer()
	defer server1.Close()

	statusCode, _ := httpRequest("PUT", server1.URL+"/peers", "peer.url")
	assert.Equal(t, 201, statusCode)

	statusCode, body := httpRequest("GET", server1.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, `{"peers":["peer.url"]}`, body)
}

func TestAddPeerFailsOnMultipleCalls(t *testing.T) {
	server1 := testServer()
	defer server1.Close()

	var statusCode int

	statusCode, _ = httpRequest("PUT", server1.URL+"/peers", "peer.url")
	assert.Equal(t, 201, statusCode)

	statusCode, _ = httpRequest("PUT", server1.URL+"/peers", "peer.url")
	assert.Equal(t, 409, statusCode)
}

func TestAddPeerCallsBack(t *testing.T) {
	server1 := testServer()
	defer server1.Close()
	server2 := testServer()
	defer server2.Close()

	httpRequest("PUT", server1.URL+"/peers", server2.URL)

	var statusCode int
	var body string

	statusCode, body = httpRequest("GET", server1.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, `{"peers":["`+server2.URL+`"]}`, body)

	statusCode, body = httpRequest("GET", server2.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, `{"peers":["`+server1.URL+`"]}`, body)
}

func TestAddPeerUpdatesExistingPeers(t *testing.T) {
	serverA := testServer()
	defer serverA.Close()
	serverB := testServer()
	defer serverB.Close()
	serverC := testServer()
	defer serverC.Close()

	httpRequest("PUT", serverA.URL+"/peers", serverB.URL)
	httpRequest("PUT", serverA.URL+"/peers", serverC.URL)

	var statusCode int
	var body string

	statusCode, body = httpRequest("GET", serverA.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, true, strings.Contains(body, serverB.URL))
	assert.Equal(t, true, strings.Contains(body, serverC.URL))

	statusCode, body = httpRequest("GET", serverB.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, true, strings.Contains(body, serverA.URL))
	assert.Equal(t, true, strings.Contains(body, serverC.URL))

	statusCode, body = httpRequest("GET", serverC.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, true, strings.Contains(body, serverA.URL))
	assert.Equal(t, true, strings.Contains(body, serverB.URL))
}
