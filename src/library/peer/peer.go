package peer

import (
	"encoding/json"
	"github.com/bmizerany/pat"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"library/hashring"
	"library/httpclient"
)

type Peer struct {
	Peers []string
	url string
	ring *hashring.Ring
	node *hashring.Node
	logger *log.Logger
}

func New(url string, logger *log.Logger) *Peer {
	ring := hashring.New()
	node := ring.AddNode(url)

	return &Peer{
		Peers: []string{},
		url: url,
		ring: ring,
		node: node,
		logger: logger,
	}
}

func (peer *Peer) PeerAddressForKey(key string) string {
	return peer.ring.NodeForKey(key).GetName()
}

func (peer *Peer) URL() string {
	return peer.url
}

func (peer *Peer) SetURL(url string) {
	peer.url = url
	peer.node.SetName(url)
}

func (peer *Peer) addPeer(newPeer string) {
	peer.Peers = append(peer.Peers, newPeer)
}

func (peer *Peer) join(newPeer string) {
	for _, p := range peer.Peers {
		httpclient.Put(p+"/peers", newPeer)
		httpclient.Put(newPeer+"/peers", p)
	}

	peer.addPeer(newPeer)
	peer.ring.AddNode(newPeer)
	httpclient.Put(newPeer+"/peers", peer.url)

	for _, p := range peer.Peers {
		nodes := httpclient.JsonData{
			"ring": peer.ring.GetNodes(),
		}
		httpclient.Put(p+"/ring", nodes.Encode())
	}
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
	m.Get("/stats", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		stats := httpclient.JsonData{
			"ring": peer.ring.GetNodes(),
			"vnodeCount": peer.node.VnodeCount(),
			"vnodeSize": peer.node.VnodeSize(),
			"vnodeStart": peer.node.VnodeStart(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		io.WriteString(w, stats.Encode())
	}))

	m.Put("/ring", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		body, _ := ioutil.ReadAll(request.Body)
		data := map[string][]string{}
		json.Unmarshal(body, &data)

		peer.ring.SetNodes(data["ring"])
		peer.node = peer.ring.Get(peer.url)

		w.WriteHeader(201)
	}))

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

	m.Put("/peers/join", http.HandlerFunc(func (w http.ResponseWriter, request *http.Request) {
		body, _ := ioutil.ReadAll(request.Body)
		newPeerURL := string(body)

		if peer.peerExists(newPeerURL) {
			w.WriteHeader(409)
		} else {
			peer.join(newPeerURL)
			w.WriteHeader(201)
		}
	}))

	m.Get("/peers", http.HandlerFunc(func (w http.ResponseWriter, request *http.Request) {
		if peer.HasPeer() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			io.WriteString(w, httpclient.JsonData{"peers": peer.Peers}.Encode())
		} else {
			w.WriteHeader(404)
		}
	}))
}
