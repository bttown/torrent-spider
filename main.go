package main

import (
	"fmt"
	"github.com/bttown/bloomfilter"
	"github.com/bttown/dht"
	"github.com/bttown/metadata"
)

var filter = bloomfilter.New(10000)

func main() {
	c := metadata.NewCollector()
	c.OnFinish(func(req *metadata.Request, meta *metadata.Metadata) {
		magnetLink := fmt.Sprintf("magnet:?xt=urn:btih:%s", req.HashInfo)
		fmt.Println("[onMetadata]", magnetLink, meta.Name)
		filter.Put(req.HashInfo)
	})

	node := dht.NewNode(dht.OptionNodeID(dht.RANDOM), dht.OptionAddress("0.0.0.0:8661"))
	node.PeerHandler = func(ip string, port int, hashInfo, peerID string) {
		if filter.MightContains(hashInfo) {
			return
		}

		c.Get(&metadata.Request{
			IP:       ip,
			Port:     port,
			HashInfo: hashInfo,
			PeerID:   peerID,
		})
	}
	node.Serve()
}
