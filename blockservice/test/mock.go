package bstest

import (
	testinstance "github.com/glin-gogogo/boxo/bitswap/testinstance"
	tn "github.com/glin-gogogo/boxo/bitswap/testnet"
	"github.com/glin-gogogo/boxo/blockservice"
	mockrouting "github.com/glin-gogogo/boxo/routing/mock"
	delay "github.com/ipfs/go-ipfs-delay"
)

// Mocks returns |n| connected mock Blockservices
func Mocks(n int, opts ...blockservice.Option) []blockservice.BlockService {
	net := tn.VirtualNetwork(delay.Fixed(0))
	routing := mockrouting.NewServer()
	sg := testinstance.NewTestInstanceGenerator(net, routing, nil, nil)
	instances := sg.Instances(n)

	var servs []blockservice.BlockService
	for _, i := range instances {
		servs = append(servs, blockservice.New(i.Blockstore,
			i.Exchange, opts...))
	}
	return servs
}
