package main

import (
	"fmt"
	"github.com/bttown/dht"
	"github.com/bttown/metadata"
	"log"
)

var (
	collectorQueriesBufferSize = 5000
	collectorMaxPendingQueries = 2000
)

var (
	// DHT 节点
	node = dht.NewNode(dht.OptionAddress("0.0.0.0:8662"))
	// 种子信息获取器
	collector = metadata.NewCollector(metadata.Options{
		QueriesBufferSize: collectorQueriesBufferSize,
		MaxPendingQueries: collectorMaxPendingQueries,
	})
)

func main() {
	collector.OnFinish(func(req metadata.Request, torrent metadata.Torrent) {
		magnetLink := fmt.Sprintf("magnet:?xt=urn:btih:%s", req.HashInfo)
		log.Println("[Metadata]", magnetLink, torrent.Info.Name)
	})
	defer collector.Close()

	node.PeerHandler = func(ip string, port int, hashInfo, peerID string) {
		if err := collector.Get(&metadata.Request{
			IP:       ip,
			Port:     port,
			HashInfo: hashInfo,
			PeerID:   peerID,
		}); err != nil {
			panic(err)
		}

	}
	node.Serve()
}
