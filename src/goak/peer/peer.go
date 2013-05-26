package peer

import (
	"github.com/bmizerany/pat"
	"io"
	"io/ioutil"
	"net/http"

	"goak/httpclient"
)

type Peer struct {
	Peers []string
	url string
}

func New() *Peer {
	return &Peer{
		Peers: []string{},
		url: "",
	}
}

func (peer *Peer) SetURL(url string) {
	peer.url = url
}

func (peer *Peer) addPeer(newPeer string) {
	for _, p := range peer.Peers {
		go httpclient.Put(p+"/peers", newPeer)
	}

	peer.Peers = append(peer.Peers, newPeer)
	httpclient.Put(newPeer+"/peers", peer.url)
}

func (peer *Peer) HasPeer() bool {
	return len(peer.Peers) > 0
}

func (peer *Peer) peerExists(query string) bool {
	for _, p := range peer.Peers {
		if p == query {
			return true
		}
	}

	return false
}

func (peer *Peer) Handler(m *pat.PatternServeMux)  {
	m.Put("/peers", http.HandlerFunc(func (w http.ResponseWriter, request *http.Request) {
		body, _ := ioutil.ReadAll(request.Body)
		newPeerURL := string(body)

		if peer.peerExists(newPeerURL) {
			w.WriteHeader(409)
		} else {
			peer.addPeer(newPeerURL)
			w.WriteHeader(201)
		}
	}))

	m.Get("/peers", http.HandlerFunc(func (w http.ResponseWriter, request *http.Request) {
		if peer.HasPeer() {
			w.WriteHeader(200)
			io.WriteString(w, httpclient.JsonData{"peers": peer.Peers}.Encode())
		} else {
			w.WriteHeader(404)
		}
	}))
}
