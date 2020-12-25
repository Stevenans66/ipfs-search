package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/streadway/amqp"
)

// TODO: Consider getting rid of any references to other modules here.

type Worker struct {
	crawler    *crawler.Crawler
	deliveries <-chan amqp.Delivery
	*instr.Instrumentation
}

func New(crawler *crawler.Crawler, deliveries <-chan amqp.Delivery, instr *instr.Instrumentation) (*Worker, error) {
	return &Worker{
		crawler,
		deliveries,
		instr,
	}, nil
}

func (w *Worker) crawlDelivery(ctx context.Context, d amqp.Delivery) error {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{},
	}

	if err := json.Unmarshal(d.Body, r); err != nil {
		return err
	}

	if !r.IsValid() {
		return fmt.Errorf("Invalid resource: %v", r)
	}

	log.Printf("Crawling: %v\n", r)

	return w.crawler.Crawl(ctx, r)
}

func (w *Worker) work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-w.deliveries:
			if !ok {
				// This is a fatal error; it should never happen - crash the program!
				panic("unexpected channel close")
			}
			if err := w.crawlDelivery(ctx, d); err != nil {
				shouldRetry := crawler.IsTemporaryErr(err)

				if err := d.Reject(shouldRetry); err != nil {
					log.Printf("Reject error %s\n", d.Body)
					// span.RecordError(ctx, err)
				}
				log.Printf("Error '%s' in delivery '%s'", err, d.Body)
				// span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
			} else {
				if err := d.Ack(false); err != nil {
					log.Printf("Ack error %s\n", d.Body)

					// span.RecordError(ctx, err)
				}
				log.Printf("Done crawling: %s\n", d.Body)
			}
		}
	}
}

func (w *Worker) Go(ctx context.Context, n uint) {
	var i uint
	for i = 0; i < n; i++ {
		go w.work(ctx)
	}
}
