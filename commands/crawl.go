package commands

import (
	"context"

	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/protocol/ipfs"
	"github.com/ipfs-search/ipfs-search/extractor/tika"
	"github.com/ipfs-search/ipfs-search/queue/amqp"
	"github.com/ipfs-search/ipfs-search/index/elasticsearch"


	// "github.com/ipfs-search/ipfs-search/worker"
	// "golang.org/x/sync/errgroup"
	"log"
	// "go.opentelemetry.io/otel/api/trace"
	// "go.opentelemetry.io/otel/codes"
)

func getIndexes(ctx context.Context, config.Indexes) (crawler.Indexes, err) {
	indexes := crawler.Indexes{
	Files:       elasticsearch.New(),
	Directories:       elasticsearch.New(),
	Invalids:       elasticsearch.New(),
		}
}

func getQueues(ctx context.Context, config.Queues) (*crawler.Queues, err) {
	amqpConnection, err = amqp.NewConnection(ctx, cfg.AMQP.URL, instrumentation)
	if err != nil {
		return nil, err
	}

	fq, err := amqpConnection.NewChannelQueue(ctx, cfg.Queues.Files.Name)
	if err != nil {
		return nil, err
	}

	dq, err := amqpConnection.NewChannelQueue(ctx, cfg.Queues.Directories.Name)
	if err != nil {
		return nil, err
	}

	hq, err := amqpConnection.NewChannelQueue(ctx, cfg.Queues.Hashes.Name)
	if err != nil {
		return nil, err
	}

	return &crawler.Queues{
		Files: fq,
		Directories: dq,
		Hashes: hq,
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

	queues := getQueues(ctx, cfg.Queues)
	indexes := getIndexes(ctx, cfg.Indexes)

	protocol := ipfs.New(

	)

	extractor := tika.New(

	)

	crawler := crawler.New(cfg.CrawlerConfig(), indexes, queues, protocol, extractor)


func New(config *Config, indexes Indexes, queues Queues, protocol protocol.Protocol, extractor extractor.Extractor) *Crawler {
	return &Crawler{
		config,
		indexes,
		queues,
		protocol,
		extractor,
	}
}
	// TODO: Rewrite me!
	return nil
}

// 	errc := make(chan error, 1)

// 	// Create error group and context
// 	// The derived Context is canceled the first time a function passed to Go returns a non-nil error or the
// 	// first time Wait returns, whichever occurs first.
// 	errg, ctx := errgroup.WithContext(ctx)

// 	startWorkers := func(ctx context.Context, cfg *config.Config, errc chan<- error) error {
// 		ctx, span := tracer.Start(ctx, "commands.startWorkers")
// 		defer span.End()

// 		factory, err := factory.New(ctx, cfg.FactoryConfig(), instrumentation, errc)
// 		if err != nil {
// 			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
// 			return err
// 		}

// 		hashGroup := worker.Group{
// 			Count:   cfg.Crawler.HashWorkers,
// 			Wait:    cfg.Crawler.HashWait,
// 			Factory: factory.NewHashWorker,
// 		}
// 		fileGroup := worker.Group{
// 			Count:   cfg.Crawler.FileWorkers,
// 			Wait:    cfg.Crawler.FileWait,
// 			Factory: factory.NewFileWorker,
// 		}

// 		// Start work loop
// 		errg.Go(func() error {
// 			ctx, span := tracer.Start(ctx, "commands.hashWorkers")
// 			defer span.End()
// 			err := hashGroup.Work(ctx)
// 			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
// 			return err

// 		})
// 		errg.Go(func() error {
// 			ctx, span := tracer.Start(ctx, "commands.fileWorkers")
// 			defer span.End()
// 			err := fileGroup.Work(ctx)
// 			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
// 			return err
// 		})

// 		return nil
// 	}

// 	if err := startWorkers(ctx, cfg, errc); err != nil {
// 		return err
// 	}

// 	log.Printf("Workers started")
// 	span.AddEvent(ctx, "workers-started")

// 	// Log errors, wait for context to cancel
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			// Errorgroup context closed (parent or error ocurred).
// 			err := ctx.Err()
// 			span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))

// 			log.Printf("Shutting down: %s", err)
// 			log.Print("Waiting for workers to finish")

// 			// Wait blocks until all function calls from the Go method
// 			// have returned, then returns the first non-nil error (if any) from them.
// 			err = errg.Wait()
// 			log.Printf("Error group finished: %s", err)
// 			return err
// 		case err := <-errc:
// 			// Log errors
// 			log.Printf("%T: %v", err, err)
// 		}
// 	}
// }
