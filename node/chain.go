package node

import (
	"encoding/hex"
	"fmt"

	"github.com/IvanTarjan/P2PgRPC/proto"
	"github.com/IvanTarjan/P2PgRPC/types"
)

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{headers: make([]*proto.Header, 0)}
}

func (list *HeaderList) Add(h *proto.Header){
	list.headers = append(list.headers, h)
}

func (list *HeaderList) Height() int {
	return len(list.headers)-1
}

func (list *HeaderList) Len() int {
	return len(list.headers)
}

func (list *HeaderList) Get(index int) (*proto.Header){
	if index > list.Height(){
		panic("index out of bounds")
	}
	return list.headers[index]
}

type Chain struct {
	blockStore BlockStorer
	headers *HeaderList
}

func NewChain(bs BlockStorer) *Chain {
	return &Chain{
		blockStore: bs,
		headers: NewHeaderList(),
	}
}

func (c *Chain) Height() int {
	return c.headers.Height()
}

func (c *Chain) AddBlock(b *proto.Block) error {
	// Add header to list of headers.
	c.headers.Add(b.Header)
	// validation
	return c.blockStore.Put(b)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error){
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error){
	if c.Height() < height{
		return nil, fmt.Errorf("given height [%d] too high - height [%d]", height, c.Height())
	}
	header := c.headers.Get(height)
	return c.GetBlockByHash(types.HashHeader(header))
}