package peer

import (
	"github.com/bmizerany/pat"
	"github.com/benmills/quiz"

	"net/http/httptest"
	"testing"

	"goak/httpclient"
)

func testNode() *httptest.Server {
	m := pat.New()
	goakPeer := New("localhost:someport")
	goakPeer.Handler(m)
	httpServer := httptest.NewServer(m)
	goakPeer.SetURL(httpServer.URL)

	return httpServer
}

func TestGetPeerWithNoPeer(t *testing.T) {
	test := quiz.Test(t)

	node := testNode()
	defer node.Close()

	statusCode, _ := httpclient.Get(node.URL+"/peers", "")

	test.Expect(statusCode).ToEqual(404)
}

func TestAddPeer(t *testing.T) {
	test := quiz.Test(t)

	node := testNode()
	defer node.Close()

	statusCode, _ := httpclient.Put(node.URL+"/peers", "peer.url")
	test.Expect(statusCode).ToEqual(201)

	statusCode, body := httpclient.Get(node.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToEqual(`{"peers":["peer.url"]}`)
}

func TestGetMultiplePeers(t *testing.T) {
	test := quiz.Test(t)

	node := testNode()
	defer node.Close()

	httpclient.Put(node.URL+"/peers", "peer1.url")
	httpclient.Put(node.URL+"/peers", "peer2.url")

	statusCode, body := httpclient.Get(node.URL+"/peers", "")

	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToEqual(`{"peers":["peer1.url","peer2.url"]}`)
}

func TestAddPeerFailsOnMultipleCalls(t *testing.T) {
	test := quiz.Test(t)

	node := testNode()
	defer node.Close()

	var statusCode int

	statusCode, _ = httpclient.Put(node.URL+"/peers", "peer.url")
	test.Expect(statusCode).ToEqual(201)

	statusCode, _ = httpclient.Put(node.URL+"/peers", "peer.url")
	test.Expect(statusCode).ToEqual(409)
}

func TestJoinCallsBack(t *testing.T) {
	test := quiz.Test(t)

	nodeA := testNode()
	defer nodeA.Close()
	nodeB := testNode()
	defer nodeB.Close()

	httpclient.Put(nodeA.URL+"/peers/join", nodeB.URL)

	var statusCode int
	var body string

	statusCode, body = httpclient.Get(nodeA.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToEqual(`{"peers":["`+nodeB.URL+`"]}`)

	statusCode, body = httpclient.Get(nodeB.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToEqual(`{"peers":["`+nodeA.URL+`"]}`)
}

func TestJoinUpdatesExistingPeers(t *testing.T) {
	test := quiz.Test(t)

	nodeA := testNode()
	defer nodeA.Close()
	nodeB := testNode()
	defer nodeB.Close()
	nodeC := testNode()
	defer nodeC.Close()

	httpclient.Put(nodeA.URL+"/peers/join", nodeB.URL)
	httpclient.Put(nodeA.URL+"/peers/join", nodeC.URL)

	var statusCode int
	var body string

	statusCode, body = httpclient.Get(nodeA.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToContain(nodeB.URL)
	test.Expect(body).ToContain(nodeC.URL)

	statusCode, body = httpclient.Get(nodeB.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToContain(nodeA.URL)
	test.Expect(body).ToContain(nodeC.URL)

	statusCode, body = httpclient.Get(nodeC.URL+"/peers", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToContain(nodeA.URL)
	test.Expect(body).ToContain(nodeB.URL)
}

func TestNodeStats(t *testing.T) {
	test := quiz.Test(t)

	node := testNode()
	defer node.Close()

	statusCode, body := httpclient.Get(node.URL+"/stats", "")

	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToContain(`"vnodeCount":1024`)
	test.Expect(body).ToContain(`"vnodeSize":4194303`)
	test.Expect(body).ToContain(`"vnodeStart":0`)
	test.Expect(body).ToContain(`"ring":["`+node.URL+`"]`)
}

func TestNodeSetRing(t *testing.T) {
	test := quiz.Test(t)

	node := testNode()
	defer node.Close()

	statusCode, _ := httpclient.Put(node.URL+"/ring", `{"ring":["`+node.URL+`","b","c"]}`)
	test.Expect(statusCode).ToEqual(201)

	_, body := httpclient.Get(node.URL+"/stats", "")
	test.Expect(body).ToContain(`"ring":["`+node.URL+`","b","c"]`)
}

func TestAddNodeUpdatesRing(t *testing.T) {
	test := quiz.Test(t)
	var statusCode int
	var body string

	nodeA := testNode()
	defer nodeA.Close()

	nodeB := testNode()
	defer nodeB.Close()

	httpclient.Put(nodeA.URL+"/peers/join", nodeB.URL)

	statusCode, body = httpclient.Get(nodeA.URL+"/stats", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToContain(`"vnodeCount":512`)
	test.Expect(body).ToContain(`"vnodeStart":0`)

	statusCode, body = httpclient.Get(nodeB.URL+"/stats", "")
	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToContain(`"vnodeCount":512`)
	test.Expect(body).ToContain(`"vnodeStart":2147483137`)
}
