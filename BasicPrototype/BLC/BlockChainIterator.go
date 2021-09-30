/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/24 11:13 上午
 * Description:
 *
 */
package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

/**
 * @Author: sundaohan
 * @Description: 迭代器
 */
type BlockChainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

/**
 * @Author: sundaohan
 * @Description: 返回迭代器对象
 * @receiver bc
 * @return *BlockChainIterator
 */
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{
		bc.Tip,
		bc.DB,
	}
}

func (bci *BlockChainIterator) Next() *Block {
	var block *Block
	err := bci.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			currentBlockBytes := b.Get(bci.CurrentHash)
			//获取到当前currentHash的区块
			block = DeserializeBlock(currentBlockBytes)
			//更新迭代器里面的CurrentHash
			bci.CurrentHash = block.PreBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block
}
