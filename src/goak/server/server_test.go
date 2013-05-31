package server

import (
	"github.com/benmills/quiz"
	"net/http/httptest"
	"testing"

	"goak/httpclient"
)

type TestNode struct {
	*httptest.Server
	node *Server
}

func testServer() *TestNode {
	goakServer := New("localhost:someport")
	httpServer := httptest.NewServer(goakServer.Handler())
	goakServer.SetURL(httpServer.URL)

	return &TestNode{httpServer, goakServer}
}

func TestAddAKey(t *testing.T) {
	test := quiz.Test(t)

	server := testServer()
	defer server.Close()

	statusCode, body := httpclient.Put(server.URL+"/data/mykey", "bar")

	test.Expect(statusCode).ToEqual(201)
	test.Expect(body).ToEqual("bar")
}

func TestFetchKey(t *testing.T) {
	test := quiz.Test(t)

	server := testServer()
	defer server.Close()

	httpclient.Put(server.URL+"/data/mykey", "bar")
	statusCode, body := httpclient.Get(server.URL+"/data/mykey", "bar")

	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToEqual("bar")
}

func TestFetchUnknownKey(t *testing.T) {
	test := quiz.Test(t)

	server := testServer()
	defer server.Close()

	statusCode, _ := httpclient.Get(server.URL+"/data/mykey", "bar")

	test.Expect(statusCode).ToEqual(404)
}

func TestUpdateKey(t *testing.T) {
	test := quiz.Test(t)

	server := testServer()
	defer server.Close()

	httpclient.Put(server.URL+"/data/mykey", "bar")
	httpclient.Put(server.URL+"/data/mykey", "baz")
	statusCode, body := httpclient.Get(server.URL+"/data/mykey", "")

	test.Expect(statusCode).ToEqual(200)
	test.Expect(body).ToEqual("baz")
}

func TestFetchesAcrossNodes(t *testing.T) {
	test := quiz.Test(t)

	server1 := testServer()
	defer server1.Close()
	server2 := testServer()
	defer server2.Close()

	httpclient.Put(server1.URL+"/peers", server2.URL)

	statusCode, _ := httpclient.Put(server1.URL+"/data/mykey", "bar")
	test.Expect(statusCode).ToEqual(201)

	statusCode2, body := httpclient.Get(server2.URL+"/data/mykey", "")
	test.Expect(statusCode2).ToEqual(200)
	test.Expect(body).ToEqual("bar")
}
