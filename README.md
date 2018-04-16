# magnetlink spider
a magnet-link spider in p2p.
一个磁力链接收集器，让你简单快速地收集DHT网络中其他节点下载资源的信息.



#### Install && Usage
    go get -u github.com/bttown/spider
	go build github.com/bttown/spider
	./spider

#### Notice
收集器刚启动的时候需要较长的时间(作者使用阿里云1核1G的机器测试大概需要1到2天:()来和DHT网络中的其他节点通信，当我们的节点被大量其他节点收录时，大量资源就会不请自来了

![snapshot](./snapshot.jpg)

#### Code

```go
package main

import (
	"fmt"
	"github.com/bttown/bloomfilter"
	"github.com/bttown/dht"
	"github.com/bttown/metadata"
	"log"
	"os"
)

var hashFilter = bloomfilter.New(10000000)
var blackList = bloomfilter.New(10000000)

func saveTorrentFile(name string, metadata *metadata.Metadata) {
	f, err := os.Create(name)
	if err != nil {
		return
	}
	defer f.Close()

	f.Write(metadata.Torrent())
}

func main() {
	c := metadata.NewCollector()
	defer c.Close()

	c.OnFinish(func(req *metadata.Request, meta *metadata.Metadata) {
		// 过滤掉重复资源
		if hashFilter.MightContains(req.HashInfo) {
			return
		}
		hashFilter.Put(req.HashInfo)
		magnetLink := fmt.Sprintf("magnet:?xt=urn:btih:%s", req.HashInfo)
		torrentFileName := fmt.Sprintf("torrents/%s.torrent", meta.Name)
		log.Println("[Metadata]", magnetLink, meta.Name)
		saveTorrentFile(torrentFileName, meta)
	})

	c.OnError(func(req *metadata.Request, err error) {
		// 将无法访问的节点地址加入黑名单
		blackList.Put(req.RemoteAddr())
		log.Println("[Error]", err)
	})

	node := dht.NewNode(dht.OptionAddress("0.0.0.0:8662"))
	node.PeerHandler = func(ip string, port int, hashInfo, peerID string) {
		// 过滤掉无法访问的节点
		if blackList.MightContains(fmt.Sprintf("%s:%d", ip, port)) {
			return
		}

		err := c.Get(&metadata.Request{
			IP:       ip,
			Port:     port,
			HashInfo: hashInfo,
			PeerID:   peerID,
		})
		if err != nil {
			panic(err)
		}

	}
	node.Serve()
}

```