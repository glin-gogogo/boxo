package mdutils

import (
	bsrv "github.com/glin-gogogo/boxo/blockservice"
	blockstore "github.com/glin-gogogo/boxo/blockstore"
	offline "github.com/glin-gogogo/boxo/exchange/offline"
	dag "github.com/glin-gogogo/boxo/ipld/merkledag"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	ipld "github.com/ipfs/go-ipld-format"
)

// Mock returns a new thread-safe, mock DAGService.
func Mock() ipld.DAGService {
	return dag.NewDAGService(Bserv())
}

// Bserv returns a new, thread-safe, mock BlockService.
func Bserv() bsrv.BlockService {
	bstore := blockstore.NewBlockstore(dssync.MutexWrap(ds.NewMapDatastore()))
	return bsrv.New(bstore, offline.Exchange(bstore))
}
