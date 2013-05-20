package peer

import (
	"github.com/bmizerany/pat"
	"github.com/bmizerany/assert"
	"net/http/httptest"
	"testing"
	"strings"

	"goak/http_client"
)

func testNode() *httptest.Server {
	m := pat.New()
	goakPeer := New()
	goakPeer.Handler(m)
	httpServer := httptest.NewServer(m)
	goakPeer.SetURL(httpServer.URL)

	return httpServer
}

func TestGetPeerWithNoPeer(t *testing.T) {
	node := testNode()
	defer node.Close()

	statusCode, _ := http_client.HttpRequest("GET", node.URL+"/peers", "")
	assert.Equal(t, 404, statusCode)
}

func TestAddPeer(t *testing.T) {
	node := testNode()
	defer node.Close()

	statusCode, _ := http_client.HttpRequest("PUT", node.URL+"/peers", "peer.url")
	assert.Equal(t, 201, statusCode)

	statusCode, body := http_client.HttpRequest("GET", node.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, `{"peers":["peer.url"]}`, body)
}

func TestGetMultiplePeers(t *testing.T) {
	node := testNode()
	defer node.Close()

	http_client.HttpRequest("PUT", node.URL+"/peers", "peer1.url")
	http_client.HttpRequest("PUT", node.URL+"/peers", "peer2.url")

	statusCode, body := http_client.HttpRequest("GET", node.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, `{"peers":["peer1.url","peer2.url"]}`, body)
}

func TestAddPeerFailsOnMultipleCalls(t *testing.T) {
	node := testNode()
	defer node.Close()

	var statusCode int

	statusCode, _ = http_client.HttpRequest("PUT", node.URL+"/peers", "peer.url")
	assert.Equal(t, 201, statusCode)

	statusCode, _ = http_client.HttpRequest("PUT", node.URL+"/peers", "peer.url")
	assert.Equal(t, 409, statusCode)
}

func TestAddPeerCallsBack(t *testing.T) {
	nodeA := testNode()
	defer nodeA.Close()
	nodeB := testNode()
	defer nodeB.Close()

	http_client.HttpRequest("PUT", nodeA.URL+"/peers", nodeB.URL)

	var statusCode int
	var body string

	statusCode, body = http_client.HttpRequest("GET", nodeA.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, `{"peers":["`+nodeB.URL+`"]}`, body)

	statusCode, body = http_client.HttpRequest("GET", nodeB.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, `{"peers":["`+nodeA.URL+`"]}`, body)
}

func TestAddPeerUpdatesExistingPeers(t *testing.T) {
	nodeA := testNode()
	defer nodeA.Close()
	nodeB := testNode()
	defer nodeB.Close()
	nodeC := testNode()
	defer nodeC.Close()

	http_client.HttpRequest("PUT", nodeA.URL+"/peers", nodeB.URL)
	http_client.HttpRequest("PUT", nodeA.URL+"/peers", nodeC.URL)

	var statusCode int
	var body string

	statusCode, body = http_client.HttpRequest("GET", nodeA.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, true, strings.Contains(body, nodeB.URL))
	assert.Equal(t, true, strings.Contains(body, nodeC.URL))

	statusCode, body = http_client.HttpRequest("GET", nodeB.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, true, strings.Contains(body, nodeA.URL))
	assert.Equal(t, true, strings.Contains(body, nodeC.URL))

	statusCode, body = http_client.HttpRequest("GET", nodeC.URL+"/peers", "")
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, true, strings.Contains(body, nodeA.URL))
	assert.Equal(t, true, strings.Contains(body, nodeB.URL))
}
