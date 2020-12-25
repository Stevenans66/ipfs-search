package factory

import (
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/extractor/tika"
	"github.com/ipfs-search/ipfs-search/protocol/ipfs"
	"github.com/ipfs-search/ipfs-search/queue/amqp"
)

type Index struct {
	Name string
}

type Queue struct {
	Name string
}

type Indexes struct {
	Files       Index
	Directories Index
	Invalids    Index
}

type Queues struct {
	Files       Queue
	Directories Queue
	Hashes      Queue
}

type ElasticSearch struct {
	URL string
}

type Config struct {
	Indexes
	Queues
	ElasticSearch
	Tika    *tika.Config
	IPFS    *ipfs.Config
	AMQP    *amqp.Config
	Crawler *crawler.Config
}
