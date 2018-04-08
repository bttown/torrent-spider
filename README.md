# spider
a bittorrent spider

![snapshot](./snapshot.jpg)

#### Install
    go get -u github.com/bttown/spider

#### Usage
```go
package main

import (
	"github.com/bttown/dht"
	"github.com/bttown/metadata"
	"github.com/bttown/bloomfilter"
	"log"
	"fmt"
)

var filter = bloomfilter.New(10000)

func main() {
	c := metadata.NewCollector()
	c.OnFinish(func(req *metadata.Request, meta *metadata.Metadata) {
		magnetLink := fmt.Sprintf("magnet:?xt=urn:btih:%s", req.HashInfo)
		log.Println("[onMetadata]", magnetLink, meta.Name)
		filter.Put(req.HashInfo)
	})


	node := dht.NewNode(dht.OptionAddress("0.0.0.0:8661"))
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

```