package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"

	cid "github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	namespace "github.com/ipfs/go-datastore/namespace"
	query "github.com/ipfs/go-datastore/query"
	logging "github.com/ipfs/go-log/v2"
	ktds "github.com/ipfs/go-datastore/keytransform"
)

var log = logging.Logger("provider.queue")

// Queue provides a best-effort durability, FIFO interface to the datastore for storing cids
//
// Best-effort durability just means that cids in the process of being provided when a
// crash or shutdown occurs may be in the queue when the node is brought back online
// depending on whether the underlying datastore has synchronous or asynchronous writes.
type Queue struct {
	// used to differentiate queues in datastore
	// e.g. provider vs reprovider
	ctx     context.Context
	ds      datastore.Datastore // Must be threadsafe
	dequeue chan cid.Cid
	enqueue chan cid.Cid
	close   context.CancelFunc
	closed  sync.WaitGroup

	counter uint64
}

// NewQueue creates a queue for cids
func NewQueue(ds datastore.Datastore) *Queue {
	namespaced := namespace.Wrap(ds, datastore.NewKey("/queue"))
	cancelCtx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		ctx:     cancelCtx,
		ds:      namespaced,
		dequeue: make(chan cid.Cid),
		enqueue: make(chan cid.Cid),
		close:   cancel,
	}
	q.closed.Add(1)
	go q.worker()
	return q
}

// Close stops the queue
func (q *Queue) Close() error {
	q.close()
	q.closed.Wait()
	// We don't close dequeue because the provider which consume this get caught in
	// an infinite loop dequeing cid.Undef if we do that.
	// The provider has it's own select on top of dequeue and will handle this by itself.
	return nil
}

// Enqueue puts a cid in the queue
func (q *Queue) Enqueue(cid cid.Cid) error {
	log.Warnf("enqueueing cid: ", cid)
	select {
	case q.enqueue <- cid:
		return nil
	case <-q.ctx.Done():
		return errors.New("failed to enqueue CID: shutting down")
	}
}

// Dequeue returns a channel that if listened to will remove entries from the queue
func (q *Queue) Dequeue() <-chan cid.Cid {
	return q.dequeue
}

// worker run dequeues and enqueues when available.
func (q *Queue) worker() {
	var k datastore.Key = datastore.Key{}
	var c cid.Cid = cid.Undef

	dsImpl, ok := q.ds.(*ktds.Datastore)
	if !ok {
		log.Errorf("queue datastore is not a ktds.Datastore")
		return
		//panic("queue datastore is not a ktds.Datastore")
	}
	dsBatch, err := dsImpl.Batch(q.ctx)
	if err != nil {
		log.Errorf("error creating batch: %s, stopping provider", err)
		return
		//panic( err)
	}

	defer q.closed.Done()
	defer q.close()

	var es []*query.Entry
	esIdx := 0
	batchCommitFlag := false
	for {
		if c == cid.Undef {
			if esIdx >= len(es) {
				es, err = q.getQueueHeadBatch()
				if err != nil {
					log.Errorf("error querying for head of queue: %s, stopping provider", err)
					return
				}
				esIdx = 0
				batchCommitFlag = true
				log.Infof("glin test..., got %d entries from queue", len(es))
			}

			var head *query.Entry = nil
			if esIdx < len(es) {
				head = es[esIdx]
				esIdx++
			}
			//head, err := q.getQueueHead()

			switch {
			//case err != nil:
			//	log.Errorf("error querying for head of queue: %s, stopping provider", err)
			//	return
			case head != nil:
				k = datastore.NewKey(head.Key)
				c, err = cid.Parse(head.Value)
				if err != nil {
					log.Warnf("error parsing queue entry cid with key (%s), removing it from queue: %s", head.Key, err)
					err = q.ds.Delete(q.ctx, k)
					if err != nil {
						log.Errorf("error deleting queue entry with key (%s), due to error (%s), stopping provider", head.Key, err)
						return
					}
					continue
				}
			default:
				c = cid.Undef
			}
		}

		// If c != cid.Undef set dequeue and attempt write, otherwise wait for enqueue
		var dequeue chan cid.Cid
		if c != cid.Undef {
			dequeue = q.dequeue
		}

		select {
		case toQueue := <-q.enqueue:
			keyPath := fmt.Sprintf("%020d/%s", q.counter, c.String())
			q.counter++
			nextKey := datastore.NewKey(keyPath)

			if c == cid.Undef {
				// fast path, skip rereading the datastore if we don't have anything in hand yet
				c = toQueue
				k = nextKey
			}

			if err := q.ds.Put(q.ctx, nextKey, toQueue.Bytes()); err != nil {
				log.Errorf("Failed to enqueue cid: %s", err)
				continue
			}
		case dequeue <- c:
			err := dsBatch.Delete(q.ctx, k)
			//err := q.ds.Delete(q.ctx, k)
			if err != nil {
				log.Errorf("Failed to delete queued cid %s with key %s: %s", c, k, err)
				continue
			}
			if batchCommitFlag {
				if err := dsBatch.Commit(q.ctx); err != nil {
					log.Errorf("Failed to commit batch: %s", err)
					continue
				}
				batchCommitFlag = false
			}

			c = cid.Undef
		case <-q.ctx.Done():
			return
		}
	}
}

func (q *Queue) getQueueHead() (*query.Entry, error) {
	qry := query.Query{Orders: []query.Order{query.OrderByKey{}}, Limit: 1}
	results, err := q.ds.Query(q.ctx, qry)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	r, ok := results.NextSync()
	if !ok {
		return nil, nil
	}

	return &r.Entry, r.Error
}

func (q *Queue) getQueueHeadBatch() ([]*query.Entry, error) {
	limit := 3000
	qry := query.Query{Orders: []query.Order{query.OrderByKey{}}, Limit: limit}
	results, err := q.ds.Query(q.ctx, qry)
	if err != nil {
		return nil, err
	}
	defer results.Close()


	var es []*query.Entry
	for {
		e, ok := results.NextSync()
		if !ok {
			break
		}
		if e.Error != nil {
			return es, e.Error
		}
		es = append(es, &e.Entry)
		limit--
		if limit <= 0 {
			return es, nil
		}
	}
	return es, nil
}

