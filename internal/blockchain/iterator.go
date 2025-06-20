package blockchain

import (
	"blockchain/internal/block"
	httpmap "blockchain/internal/httpMap"
	"fmt"
)

/*block 2 -> block 1 -> genesis*/
type BlockchainIterator struct {
	currentBlockHash string
}

func (bi *BlockchainIterator) Next() (b *block.Block, err error) {
	if len(bi.currentBlockHash) == 0 {
		return
	}

	data, err := httpmap.Load(BlocksFile, bi.currentBlockHash)
	if err != nil {
		return
	}

	b, err = block.StringDeserializeBlock(data)
	if err != nil {
		return
	}

	bi.currentBlockHash = b.StringPrevBlockHash()

	return
}

func ForEach(bc *Blockchain) <-chan *block.Block {
	c := make(chan *block.Block)

	go func() {
		it := bc.Iterator()

		for {
			b, err := it.Next()
			if b == nil {
				if err != nil {
					fmt.Println(err)
				}
				break
			}
			c <- b
		}

		close(c)
	}()

	return c
}
