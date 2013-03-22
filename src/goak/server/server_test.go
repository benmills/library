package server

import (
	"github.com/bmizerany/assert"
	"github.com/drewolson/testflight"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func withServer(f func(*testflight.Requester)) {
	goakServer := New()
	testflight.WithServer(goakServer.Handler(), f)
}

func TestAddAKey(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		response := r.Put("/data/mykey", testflight.FORM_ENCODED, "myvalue")
		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, "myvalue", response.Body)
	})
}

func TestFetchKey(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		r.Put("/data/mykey", testflight.FORM_ENCODED, "myvalue")
		r.Put("/data/notmykey", testflight.FORM_ENCODED, "notmyvalue")

		response := r.Get("/data/mykey")
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "myvalue", response.Body)
	})
}

func TestFetchUnknownKey(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		response := r.Get("/data/badkey")
		assert.Equal(t, 404, response.StatusCode)
	})
}

func TestUpdateKey(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		r.Put("/data/mykey", testflight.FORM_ENCODED, "myvalue")
		response := r.Put("/data/mykey", testflight.FORM_ENCODED, "mysecondvalue")
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "mysecondvalue", response.Body)

		response = r.Get("/data/mykey")
		assert.Equal(t, "mysecondvalue", response.Body)
	})
}

func TestFetchesAcrossNodes(t *testing.T) {
	node1 := New()
	node2 := New()
	node1.AddPeer("http://localhost:9191")
	node2.AddPeer("http://localhost:9090")
	go http.ListenAndServe(":9090", node1.Handler())
	go http.ListenAndServe(":9191", node2.Handler())
	request, _ := http.NewRequest("PUT", "http://localhost:9090/data/mykey", strings.NewReader("bar"))
	client := http.Client{}
	response, _ := client.Do(request)
	assert.Equal(t, 201, response.StatusCode)

	response, _ = http.Get("http://localhost:9191/data/mykey")
	assert.Equal(t, 200, response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "bar", string(body))
}
