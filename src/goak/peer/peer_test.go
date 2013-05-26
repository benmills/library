package peer

import (
	"github.com/bmizerany/pat"
	"github.com/benmills/quiz"

	"net/http/httptest"
	"testing"

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
	test := quiz.Test(t)

	node := testNode()
	defer node.Close()

	statusCode, _ := http_client.HttpRequest("GET", node.URL+"/peers", "")

	test.Expect(statusCode).ToEqual(404)
}

func TestAddPeer(t *testing.T) {
	test := quiz.Test(t)

	node := testNode()
	defer node.Close()

	statusCode, _ := http_client.HttpRequest("PUT", node.URL+"/peers", "peer.url")
	test.Expect(statusCode).ToEqual(201)

	statusCode, body := http_client.HttpRequest("GET", node.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToEqual(`{"peers":["peer.url"]}`)
}

func TestGetMultiplePeers(t *testing.T) {
	test := quiz.Test(t)

	node := testNode()
	defer node.Close()

	http_client.HttpRequest("PUT", node.URL+"/peers", "peer1.url")
	http_client.HttpRequest("PUT", node.URL+"/peers", "peer2.url")

	statusCode, body := http_client.HttpRequest("GET", node.URL+"/peers", "")

	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToEqual(`{"peers":["peer1.url","peer2.url"]}`)
}

func TestAddPeerFailsOnMultipleCalls(t *testing.T) {
	test := quiz.Test(t)

	node := testNode()
	defer node.Close()

	var statusCode int

	statusCode, _ = http_client.HttpRequest("PUT", node.URL+"/peers", "peer.url")
	test.Expect(statusCode).ToEqual(201)

	statusCode, _ = http_client.HttpRequest("PUT", node.URL+"/peers", "peer.url")
	test.Expect(statusCode).ToEqual(409)
}

func TestAddPeerCallsBack(t *testing.T) {
	test := quiz.Test(t)

	nodeA := testNode()
	defer nodeA.Close()
	nodeB := testNode()
	defer nodeB.Close()

	http_client.HttpRequest("PUT", nodeA.URL+"/peers", nodeB.URL)

	var statusCode int
	var body string

	statusCode, body = http_client.HttpRequest("GET", nodeA.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToEqual(`{"peers":["`+nodeB.URL+`"]}`)

	statusCode, body = http_client.HttpRequest("GET", nodeB.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToEqual(`{"peers":["`+nodeA.URL+`"]}`)
}

func TestAddPeerUpdatesExistingPeers(t *testing.T) {
	test := quiz.Test(t)

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
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToContain(nodeB.URL)
	test.Expect(body).ToContain(nodeC.URL)

	statusCode, body = http_client.HttpRequest("GET", nodeB.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToContain(nodeA.URL)
	test.Expect(body).ToContain(nodeC.URL)

	statusCode, body = http_client.HttpRequest("GET", nodeC.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToContain(nodeA.URL)
	test.Expect(body).ToContain(nodeB.URL)
}
