package server

import(
	"github.com/drewolson/testflight"
	"github.com/bmizerany/assert"
	"testing"
)

func withServer(f func(*testflight.Requester)) {
	goakServer := New()
	testflight.WithServer(goakServer.Handler(), f)
}

func TestAddAKey(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		response := r.Put("/mykey", testflight.FORM_ENCODED, "myvalue")
		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, "myvalue", response.Body)
	})
}

func TestFetchKey(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		r.Put("/mykey", testflight.FORM_ENCODED, "myvalue")
		r.Put("/notmykey", testflight.FORM_ENCODED, "notmyvalue")

		response := r.Get("/mykey")
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "myvalue", response.Body)
	})
}

func TestFetchUnknownKey(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		response := r.Get("/badkey")
		assert.Equal(t, 404, response.StatusCode)
	})
}

func TestUpdateKey(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		r.Put("/mykey", testflight.FORM_ENCODED, "myvalue")
		response := r.Put("/mykey", testflight.FORM_ENCODED, "mysecondvalue")
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "mysecondvalue", response.Body)

		response = r.Get("/mykey")
		assert.Equal(t, "mysecondvalue", response.Body)
	})
}
