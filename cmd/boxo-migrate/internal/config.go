package migrate

import (
	"encoding/json"
	"fmt"
	"io"
)

type Config struct {
	ImportPaths map[string]string
	Modules     []string
}

var DefaultConfig = Config{
	ImportPaths: map[string]string{
		"github.com/ipfs/go-bitswap":                     "github.com/glin-gogogo/boxo/bitswap",
		"github.com/ipfs/go-ipfs-files":                  "github.com/glin-gogogo/boxo/files",
		"github.com/ipfs/tar-utils":                      "github.com/glin-gogogo/boxo/tar",
		"github.com/ipfs/interface-go-ipfs-core":         "github.com/glin-gogogo/boxo/coreiface",
		"github.com/ipfs/go-unixfs":                      "github.com/glin-gogogo/boxo/ipld/unixfs",
		"github.com/ipfs/go-pinning-service-http-client": "github.com/glin-gogogo/boxo/pinning/remote/client",
		"github.com/ipfs/go-path":                        "github.com/glin-gogogo/boxo/path",
		"github.com/ipfs/go-namesys":                     "github.com/glin-gogogo/boxo/namesys",
		"github.com/ipfs/go-mfs":                         "github.com/glin-gogogo/boxo/mfs",
		"github.com/ipfs/go-ipfs-provider":               "github.com/glin-gogogo/boxo/provider",
		"github.com/ipfs/go-ipfs-pinner":                 "github.com/glin-gogogo/boxo/pinning/pinner",
		"github.com/ipfs/go-ipfs-keystore":               "github.com/glin-gogogo/boxo/keystore",
		"github.com/ipfs/go-filestore":                   "github.com/glin-gogogo/boxo/filestore",
		"github.com/ipfs/go-ipns":                        "github.com/glin-gogogo/boxo/ipns",
		"github.com/ipfs/go-blockservice":                "github.com/glin-gogogo/boxo/blockservice",
		"github.com/ipfs/go-ipfs-chunker":                "github.com/glin-gogogo/boxo/chunker",
		"github.com/ipfs/go-fetcher":                     "github.com/glin-gogogo/boxo/fetcher",
		"github.com/ipfs/go-ipfs-blockstore":             "github.com/glin-gogogo/boxo/blockstore",
		"github.com/ipfs/go-ipfs-posinfo":                "github.com/glin-gogogo/boxo/filestore/posinfo",
		"github.com/ipfs/go-ipfs-util":                   "github.com/glin-gogogo/boxo/util",
		"github.com/ipfs/go-ipfs-ds-help":                "github.com/glin-gogogo/boxo/datastore/dshelp",
		"github.com/ipfs/go-verifcid":                    "github.com/glin-gogogo/boxo/verifcid",
		"github.com/ipfs/go-ipfs-exchange-offline":       "github.com/glin-gogogo/boxo/exchange/offline",
		"github.com/ipfs/go-ipfs-routing":                "github.com/glin-gogogo/boxo/routing",
		"github.com/ipfs/go-ipfs-exchange-interface":     "github.com/glin-gogogo/boxo/exchange",
		"github.com/ipfs/go-merkledag":                   "github.com/glin-gogogo/boxo/ipld/merkledag",
		"github.com/boxo/ipld/car":                       "github.com/ipld/go-car",

		// Pre Boxo rename
		"github.com/ipfs/go-libipfs/gateway":               "github.com/glin-gogogo/boxo/gateway",
		"github.com/ipfs/go-libipfs/bitswap":               "github.com/glin-gogogo/boxo/bitswap",
		"github.com/ipfs/go-libipfs/files":                 "github.com/glin-gogogo/boxo/files",
		"github.com/ipfs/go-libipfs/tar":                   "github.com/glin-gogogo/boxo/tar",
		"github.com/ipfs/go-libipfs/coreiface":             "github.com/glin-gogogo/boxo/coreiface",
		"github.com/ipfs/go-libipfs/unixfs":                "github.com/glin-gogogo/boxo/ipld/unixfs",
		"github.com/ipfs/go-libipfs/pinning/remote/client": "github.com/glin-gogogo/boxo/pinning/remote/client",
		"github.com/ipfs/go-libipfs/path":                  "github.com/glin-gogogo/boxo/path",
		"github.com/ipfs/go-libipfs/namesys":               "github.com/glin-gogogo/boxo/namesys",
		"github.com/ipfs/go-libipfs/mfs":                   "github.com/glin-gogogo/boxo/mfs",
		"github.com/ipfs/go-libipfs/provider":              "github.com/glin-gogogo/boxo/provider",
		"github.com/ipfs/go-libipfs/pinning/pinner":        "github.com/glin-gogogo/boxo/pinning/pinner",
		"github.com/ipfs/go-libipfs/keystore":              "github.com/glin-gogogo/boxo/keystore",
		"github.com/ipfs/go-libipfs/filestore":             "github.com/glin-gogogo/boxo/filestore",
		"github.com/ipfs/go-libipfs/ipns":                  "github.com/glin-gogogo/boxo/ipns",
		"github.com/ipfs/go-libipfs/blockservice":          "github.com/glin-gogogo/boxo/blockservice",
		"github.com/ipfs/go-libipfs/chunker":               "github.com/glin-gogogo/boxo/chunker",
		"github.com/ipfs/go-libipfs/fetcher":               "github.com/glin-gogogo/boxo/fetcher",
		"github.com/ipfs/go-libipfs/blockstore":            "github.com/glin-gogogo/boxo/blockstore",
		"github.com/ipfs/go-libipfs/filestore/posinfo":     "github.com/glin-gogogo/boxo/filestore/posinfo",
		"github.com/ipfs/go-libipfs/util":                  "github.com/glin-gogogo/boxo/util",
		"github.com/ipfs/go-libipfs/datastore/dshelp":      "github.com/glin-gogogo/boxo/datastore/dshelp",
		"github.com/ipfs/go-libipfs/verifcid":              "github.com/glin-gogogo/boxo/verifcid",
		"github.com/ipfs/go-libipfs/exchange/offline":      "github.com/glin-gogogo/boxo/exchange/offline",
		"github.com/ipfs/go-libipfs/routing":               "github.com/glin-gogogo/boxo/routing",
		"github.com/ipfs/go-libipfs/exchange":              "github.com/glin-gogogo/boxo/exchange",

		// Unmigrated things
		"github.com/ipfs/go-libipfs/blocks":  "github.com/ipfs/go-block-format",
		"github.com/glin-gogogo/boxo/blocks": "github.com/ipfs/go-block-format",
	},
	Modules: []string{
		"github.com/ipfs/go-bitswap",
		"github.com/ipfs/go-ipfs-files",
		"github.com/ipfs/tar-utils",
		"gihtub.com/ipfs/go-block-format",
		"github.com/ipfs/interface-go-ipfs-core",
		"github.com/ipfs/go-unixfs",
		"github.com/ipfs/go-pinning-service-http-client",
		"github.com/ipfs/go-path",
		"github.com/ipfs/go-namesys",
		"github.com/ipfs/go-mfs",
		"github.com/ipfs/go-ipfs-provider",
		"github.com/ipfs/go-ipfs-pinner",
		"github.com/ipfs/go-ipfs-keystore",
		"github.com/ipfs/go-filestore",
		"github.com/ipfs/go-ipns",
		"github.com/ipfs/go-blockservice",
		"github.com/ipfs/go-ipfs-chunker",
		"github.com/ipfs/go-fetcher",
		"github.com/ipfs/go-ipfs-blockstore",
		"github.com/ipfs/go-ipfs-posinfo",
		"github.com/ipfs/go-ipfs-util",
		"github.com/ipfs/go-ipfs-ds-help",
		"github.com/ipfs/go-verifcid",
		"github.com/ipfs/go-ipfs-exchange-offline",
		"github.com/ipfs/go-ipfs-routing",
		"github.com/ipfs/go-ipfs-exchange-interface",
		"github.com/ipfs/go-libipfs",
	},
}

func ReadConfig(r io.Reader) (Config, error) {
	var config Config
	err := json.NewDecoder(r).Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("reading and decoding config: %w", err)
	}
	return config, nil
}
