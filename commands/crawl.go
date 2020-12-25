package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/extractor/tika"
	"github.com/ipfs-search/ipfs-search/index/elasticsearch"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/protocol/ipfs"
	"github.com/ipfs-search/ipfs-search/queue/amqp"
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs-search/ipfs-search/worker"
	samqp "github.com/streadway/amqp"

	"log"
	// "go.opentelemetry.io/otel/api/trace"
	// "go.opentelemetry.io/otel/codes"
)

type consumeChans struct {
	Files       <-chan samqp.Delivery
	Directories <-chan samqp.Delivery
	Hashes      <-chan samqp.Delivery
}

func getConsumeChans(ctx context.Context, cfg *config.Config, instrumentation *instr.Instrumentation) (*consumeChans, error) {
	var c consumeChans

	queues, err := getQueues(ctx, cfg, instrumentation)
	if err != nil {
		return nil, err
	}

	c.Files, err = queues.Files.Consume(ctx)
	if err != nil {
		return nil, err
	}

	c.Directories, err = queues.Directories.Consume(ctx)
	if err != nil {
		return nil, err
	}

	c.Hashes, err = queues.Hashes.Consume(ctx)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func makeWorkers(ctx context.Context, consumeChan <-chan samqp.Delivery, c *crawler.Crawler, n uint) {
	var i uint
	for i = 0; i < n; i++ {
		go work(ctx, consumeChan, c)
	}
}

// Crawl configures and initializes crawling
func Crawl(ctx context.Context, cfg *config.Config) error {
	instFlusher, err := instr.Install("ipfs-crawler")
	if err != nil {
		log.Fatal(err)
	}
	defer instFlusher()

	instrumentation := instr.New()
	tracer := instrumentation.Tracer

	ctx, span := tracer.Start(ctx, "commands.Crawl")
	defer span.End()

	f := factory.New()
	c, err := f.Crawler(ctx)
	if err != nil {
		return err
	}

	consumeChans, err := getConsumeChans(ctx, cfg, instrumentation)

	makeWorkers(ctx, consumeChans.Files, c, cfg.Workers.FileWorkers)
	makeWorkers(ctx, consumeChans.Hashes, c, cfg.Workers.HashWorkers)
	makeWorkers(ctx, consumeChans.Directories, c, cfg.Workers.DirectoryWorkers)

	// Context closure or panic is the only way to stop crawling
	<-ctx.Done()

	return ctx.Err()
}
