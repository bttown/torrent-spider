# spider
a p2p magnet-link spider
一个磁力链接收集器，让你简单快速地收集DHT网络中其他节点下载资源的信息.



#### Install && Usage
    go get -u github.com/bttown/spider
	go build github.com/bttown/spider
	./spider

#### Notice
收集器刚启动的时候需要较长的时间来和DHT网络中的其他节点通信，当我们的节点被大量其他节点收录时，大量资源就会不请自来了

![snapshot](./snapshot.jpg)

#### Code

```go
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

```