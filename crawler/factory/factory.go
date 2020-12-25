package factory

import (
	"context"
	"net/http"

	"github.com/olivere/elastic/v7"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/extractor/tika"
	"github.com/ipfs-search/ipfs-search/index/elasticsearch"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/protocol/ipfs"
	"github.com/ipfs-search/ipfs-search/queue/amqp"
)

type Factory struct {
	config     *Config
	httpClient *http.Client
	instr      *instr.Instrumentation
}

func (f *Factory) Crawler(ctx context.Context) (*crawler.Crawler, error) {
	queues, err := f.getQueues(ctx)
	if err != nil {
		return nil, err
	}

	indexes, err := f.getIndexes(ctx)
	if err != nil {
		return nil, err
	}

	protocol := ipfs.New(f.config.IPFS, f.httpClient, f.instr)
	extractor := tika.New(f.config.Tika, f.httpClient, protocol, f.instr)

	return crawler.New(f.config.Crawler, indexes, queues, protocol, extractor), nil
}

func getHttpClient() *http.Client {
	// TODO: Get more advanced client with circuit breaking etc. over manual
	// retrying get etc.
	// Ref: https://github.com/gojek/heimdall#creating-a-hystrix-like-circuit-breaker
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

func New(c *Config, i *instr.Instrumentation) *Factory {
	return &Factory{
		config:     c,
		httpClient: getHttpClient(),
		instr:      i,
	}
}

func (f *Factory) getElasticClient() (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(f.config.ElasticSearch.URL),
		elastic.SetHttpClient(f.httpClient),
	)
}

func (f *Factory) getIndexes(ctx context.Context) (*crawler.Indexes, error) {
	esClient, err := f.getElasticClient()
	if err != nil {
		return nil, err
	}

	return &crawler.Indexes{
		Files: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: f.config.Indexes.Files.Name},
		),
		Directories: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: f.config.Indexes.Directories.Name},
		),
		Invalids: elasticsearch.New(
			esClient,
			&elasticsearch.Config{Name: f.config.Indexes.Invalids.Name},
		),
	}, nil
}

func (f *Factory) getQueues(ctx context.Context) (*crawler.Queues, error) {
	amqpConnection, err := amqp.NewConnection(ctx, f.config.AMQP, f.instr)
	if err != nil {
		return nil, err
	}

	fq, err := amqpConnection.NewChannelQueue(ctx, f.config.Queues.Files.Name)
	if err != nil {
		return nil, err
	}

	dq, err := amqpConnection.NewChannelQueue(ctx, f.config.Queues.Directories.Name)
	if err != nil {
		return nil, err
	}

	hq, err := amqpConnection.NewChannelQueue(ctx, f.config.Queues.Hashes.Name)
	if err != nil {
		return nil, err
	}

	return &crawler.Queues{
		Files:       fq,
		Directories: dq,
		Hashes:      hq,
	}, nil
}
