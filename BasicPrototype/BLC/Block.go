/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/19 11:20 下午
 * Description:
 *
 */
package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

/**
 * @Author: sundaohan
 * @Description: 区块结构
 */
type Block struct {
	//区块高度
	Height int64
	//上一个区块哈希
	PreBlockHash []byte
	//交易数据
	Txs []*Transaction
	//时间戳
	TimeStamp int64
	//哈希
	Hash []byte
	//Nonce
	Nonce int64
}

/**
 * @Author: sundaohan
 * @Description: 创建新区块
 * @param data
 * @param height
 * @param prevBlockHash
 * @return *Block
 */
func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {
	//创建区块
	block := &Block{
		Height:       height,
		PreBlockHash: prevBlockHash,
		Txs:          txs,
		TimeStamp:    time.Now().Unix(),
		Hash:         nil,
		Nonce:        0,
	}
	//设置hash
	//block.SetHash()
	//工作量证明算法,返回有效的Hash和Nonce
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	fmt.Println()
	return block
}

/**
 * @Author: sundaohan
 * @Description: 创建创世区块的方法
 * @param data
 */
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs,
		1,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	)
}

/**
 * @Author: sundaohan
 * @Description: 将区块序列化成字节数组
 * @receiver b
 * @return []byte
 */
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

/**
 * @Author: sundaohan
 * @Description: 将字节数组反序列化
 * @param blockBytes
 * @return *Block
 */
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

/**
 * @Author: sundaohan
 * @Description: 将txs转换成字节数组
 * @receiver b
 * @return []byte
 */
func (b *Block) HashTransactions() []byte {

	//var txHashes [][]byte
	//var txHash [32]byte
	//for _, tx := range b.Txs{
	//	txHashes = append(txHashes, tx.TxHash)
	//}
	//txHash = sha256.Sum256(bytes.Join(txHashes,[]byte{}))
	//return txHash[:]
	var transactions [][]byte
	for _, tx := range b.Txs {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)
	return mTree.RootNode.Data
}
