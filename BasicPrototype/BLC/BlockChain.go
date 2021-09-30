/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/19 11:21 下午
 * Description:
 *
 */
package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"sdhChain/BasicPrototype/Wallet"
	"strconv"
	"time"
)

const (
	//数据库名字
	dbName = "blockchain_%s.db"
	//表名
	blockTableName = "blocks"
)

/**
 * @Author: sundaohan
 * @Description: 区块链对象
 */
type BlockChain struct {
	Tip []byte //最新的区块的Hash
	DB  *bolt.DB
}

/**
 * @Author: sundaohan
 * @Description: 判断数据库是否存在
 * @return bool
 */
func DBExists(dbName string) bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

/**
 * @Author: sundaohan
 * @Description: 创建带有创世区块的区块链
 * @return *BlockChain
 */
func CreateBlockChainWithGenesisBlock(address string, nodeID string) *BlockChain {
	//格式化数据库名字
	dbName := fmt.Sprintf(dbName, nodeID)
	//判断数据库是否存在
	if DBExists(dbName) {
		fmt.Println("创世区块已存在!")
		os.Exit(1)
	}
	fmt.Println("正在创建创世区块......")
	// 打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	var blockHash []byte
	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			//创建创世区块
			txCoinBase := NewCoinBaseTransaction(address)
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinBase})
			//blockHash = genesisBlock.Hash
			//将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			blockHash = genesisBlock.Hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

	//返回区块链对象
	return &BlockChain{
		blockHash,
		db,
	}

}

/**
 * @Author: sundaohan
 * @Description: 创建新区块加入到区块链中
 * @receiver bc
 * @param data
 * @param height
 * @param preHash
 */
func (bc *BlockChain) AddBlockToBlockChain(txs []*Transaction) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 获取表
		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			// 获取当前最新区块
			blockBytes := b.Get(bc.Tip)
			block := DeserializeBlock(blockBytes)
			// 创建新区块
			newBlock := NewBlock(txs, block.Height+1, block.Hash)
			// 将区块序列化存储到数据库中
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 更新数据库"l"对应的hash
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			// 更新blockchain的tip
			bc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

}

/**
 * @Author: sundaohan
 * @Description: 遍历输出所有区块
 * @receiver bc
 */
func (bc *BlockChain) PrintChain() {

	blockChainIterator := bc.Iterator()

	for {
		block := blockChainIterator.Next()
		fmt.Printf("Height : %d, PrevBlockHash: %x, TimeStamp: %s, Hash : %x, Nonce : %d, Txs:\n",
			block.Height,
			block.PreBlockHash,
			time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"),
			block.Hash,
			block.Nonce)
		for _, tx := range block.Txs {
			fmt.Printf("TxHash: %x\n", tx.TxHash)

			for _, in := range tx.Vins {
				fmt.Printf("Vins.TxHash: %x\n", in.TxHash)
				fmt.Printf("Vins.Vout: %d\n", in.Vout)
				fmt.Printf("Vins.PublicKey: %x\n", in.PublicKey)
			}

			for _, out := range tx.Vouts {
				fmt.Printf("Vouts.Value :%d\n", out.Value)
				fmt.Printf("Vouts.Ripemd160Hash: %x\n", out.Ripemd160Hash)
			}

		}
		fmt.Println("-------------------------------------------------------")
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

}

/**
 * @Author: sundaohan
 * @Description: 获取区块链对象
 * @return *BlockChain
 */
func GetBlockChainObject(nodeID string) *BlockChain {
	dbName := fmt.Sprintf(dbName, nodeID)
	//判断数据库是否存在
	if DBExists(dbName) == false {
		fmt.Println("数据库不存在!")
		os.Exit(1)
	}
	// 打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			//读取最新区块的hash
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	return &BlockChain{
		tip,
		db,
	}
}

/**
 * @Author: sundaohan
 * @Description: 查询未花费的TXOutput
 * @param addr
 * @return []*Transaction.Transaction
 */
func (bc *BlockChain) UnUTXOs(addr string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO
	spentTXOutputs := make(map[string][]int)
	blockIterator := bc.Iterator()
	for _, tx := range txs {
		if tx.IsCoinBaseTransaction() == false {
			for _, in := range tx.Vins {
				publicKeyHash := Wallet.Base58Decode([]byte(addr))
				ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4]
				//是否能够解锁
				if in.UnlockWithRipemd160Hash(ripemd160Hash) {
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}

		// Vouts

	}

	for _, tx := range txs {
		fmt.Println("**********************")
		fmt.Println(unUTXOs)
		fmt.Println(tx.Vouts)
		fmt.Println("**********************")
	WORK1:
		for index, out := range tx.Vouts {
			if out.UnlockScriptPubKeyWithAddr(addr) {
				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{
						tx.TxHash,
						index,
						out,
					}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if hash == txHashStr {

							var isUnSpentUXTO bool

							for _, outIndex := range indexArray {
								if index == outIndex {
									isUnSpentUXTO = true
									continue WORK1
								}
								if isUnSpentUXTO == false {
									utxo := &UTXO{
										tx.TxHash,
										index,
										out,
									}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {

							utxo := &UTXO{
								tx.TxHash,
								index,
								out,
							}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}
			}

		}
	}

	for {
		block := blockIterator.Next()
		fmt.Println(block)
		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]
			// txHash

			// Vins
			if tx.IsCoinBaseTransaction() == false {
				for _, in := range tx.Vins {
					//是否能够解锁
					publicKeyHash := Wallet.Base58Decode([]byte(addr))
					ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4]
					if in.UnlockWithRipemd160Hash(ripemd160Hash) {
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}

			// Vouts

		WORK2:
			for index, out := range tx.Vouts {

				if out.UnlockScriptPubKeyWithAddr(addr) {
					fmt.Println(out)
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							var f = false
							for txHash, indexArray := range spentTXOutputs {
								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {

										f = true
										continue WORK2
									}
								}
							}
							if f == false {
								utxo := &UTXO{
									tx.TxHash,
									index,
									out,
								}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else {
							utxo := &UTXO{
								tx.TxHash,
								index,
								out,
							}
							unUTXOs = append(unUTXOs, utxo)
						}
					}

					//if f == true{
					//	f = false
					//	continue
					//}
					//unUTXOs = append(unUTXOs, out)
				}
			}
		}
		fmt.Println(spentTXOutputs)
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return unUTXOs
}

/**
 * @Author: sundaohan
 * @Description: 挖掘新区块
 * @param from
 * @param to
 * @param amount
 */
func (bc *BlockChain) MineNewBlock(from, to, amount []string, nodeID string) *BlockChain {

	utxoSet := &UTXOSet{bc}

	// 建立交易数组Transaction
	var txs []*Transaction
	for index, addr := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewNormalTransaction(addr, to[index], int64(value), utxoSet, txs, nodeID)
		txs = append(txs, tx)
		//fmt.Println(tx)
	}

	// 激励
	tx := NewCoinBaseTransaction(from[0])
	txs = append(txs, tx)

	var block *Block
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("l"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})

	// 在建立新区块之前对txs进行签名验证
	_txs := []*Transaction{}

	for _, tx := range txs {
		if bc.VerifyTransaction(tx, _txs) == false {
			log.Panic("签名失败")
		}
		_txs = append(_txs, tx)
	}

	// 建立新的区块
	block = NewBlock(txs, block.Height+1, block.Hash)

	//将新区块存储到数据库
	bc.DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			b.Put(block.Hash, block.Serialize())
			b.Put([]byte("l"), block.Hash)
			bc.Tip = block.Hash
		}

		return nil
	})
	return nil
}

/**
 * @Author: sundaohan
 * @Description: 查询余额
 * @receiver bc
 * @param addr
 * @return int64
 */
func (bc *BlockChain) GetBalance(addr string) int64 {
	utxOs := bc.UnUTXOs(addr, []*Transaction{})
	var amount int64
	for _, utxo := range utxOs {
		amount = amount + utxo.TXOutput.Value
	}
	return amount
}

/**
 * @Author: sundaohan
 * @Description: 查找可用的UTXO
 * @receiver bc
 * @param from
 * @param amount
 * @return int
 * @return map[string][]int
 */
func (bc *BlockChain) FindSpendableUTXOs(from string, amount int, txs []*Transaction) (int64, map[string][]int) {
	//获取所有UTXO
	utxos := bc.UnUTXOs(from, txs)
	spendAbleUTXO := make(map[string][]int)
	//遍历数组
	var value int64
	for _, utxo := range utxos {
		value = value + utxo.TXOutput.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendAbleUTXO[hash] = append(spendAbleUTXO[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}
	}
	if value < int64(amount) {
		fmt.Printf("%s's fund is not enough\n", from)
		os.Exit(1)
	}
	return value, spendAbleUTXO
}

/**
 * @Author: sundaohan
 * @Description: 对交易进行签名
 * @receiver bc
 * @param tx
 * @param privateKey
 */
func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey, txs []*Transaction) {
	if tx.IsCoinBaseTransaction() {
		return
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, err := bc.FindTransaction(vin.TxHash, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}
	tx.Sign(privateKey, prevTXs)
}

/**
 * @Author: sundaohan
 * @Description: 找对应的交易
 * @receiver bc
 * @param ID
 * @return Transaction
 * @return error
 */
func (bc *BlockChain) FindTransaction(ID []byte, txs []*Transaction) (Transaction, error) {

	for _, tx := range txs {
		if bytes.Compare(tx.TxHash, ID) == 0 {
			return *tx, nil
		}
	}

	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Txs {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return Transaction{}, nil
}

/**
 * @Author: sundaohan
 * @Description: 验证交易合法性
 * @receiver bc
 * @param tx
 * @return bool
 */
func (bc *BlockChain) VerifyTransaction(tx *Transaction, txs []*Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, err := bc.FindTransaction(vin.TxHash, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}
	return tx.verify(prevTXs)
}

/**
 * @Author: sundaohan
 * @Description: 找到对应的字典
 * @receiver bc
 * @return map[string]*TXOutputs
 */
func (bc *BlockChain) FindUTXOMap() map[string]*TXOutputs {
	//找到未花费的输出
	bcIterator := bc.Iterator()
	//存储已花费的UTXO的信息
	spentableUTXOsMap := make(map[string][]*TXInput)

	utxoMaps := make(map[string]*TXOutputs)

	for {
		block := bcIterator.Next()
		for i := len(block.Txs) - 1; i >= 0; i-- {
			txOutputs := &TXOutputs{[]*UTXO{}}
			tx := block.Txs[i]
			if tx.IsCoinBaseTransaction() == false {
				for _, txInput := range tx.Vins {
					txHash := hex.EncodeToString(txInput.TxHash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash], txInput)
				}
			}
			txHash := hex.EncodeToString(tx.TxHash)
		WorkOutLoop:
			for index, out := range tx.Vouts {

				if tx.IsCoinBaseTransaction() {
					fmt.Println("IsCoinBaseTransaction")
				}

				txInputs := spentableUTXOsMap[txHash]
				if len(txInputs) > 0 {

					isSpent := false

					for _, in := range txInputs {
						outPublicKey := out.Ripemd160Hash
						inPublicKey := in.PublicKey

						if bytes.Compare(outPublicKey, Wallet.Ripemd160Hash(inPublicKey)) == 0 {
							if index == in.Vout {
								isSpent = true
								continue WorkOutLoop
							}
						}
					}
					if isSpent == false {
						utxo := &UTXO{tx.TxHash, index, out}
						txOutputs.UTXOS = append(txOutputs.UTXOS, utxo)
					}
				} else {
					utxo := &UTXO{tx.TxHash, index, out}
					txOutputs.UTXOS = append(txOutputs.UTXOS, utxo)
				}
			}
			//设置键值对
			utxoMaps[txHash] = txOutputs
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		//找到创世区块就退出
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}

	}
	return utxoMaps
}

/**
 * @Author: sundaohan
 * @Description: 返回区块高度
 * @receiver bc
 * @return int64
 */
func (bc *BlockChain) GetBestHeight() int64 {
	block := bc.Iterator().Next()
	return block.Height
}

/**
 * @Author: sundaohan
 * @Description: 返回当前区块链中的所有区块的哈希
 * @receiver bc
 * @return [][]byte
 */
func (bc *BlockChain) GetBlockHashes() [][]byte {
	iterator := bc.Iterator()

	var blockHashes [][]byte

	for {
		block := iterator.Next()
		blockHashes = append(blockHashes, block.Hash)

		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		//找到创世区块就退出
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return blockHashes
}

/**
 * @Author: sundaohan
 * @Description: 根据hash返回区块信息
 * @receiver bc
 * @param blockHash
 */
func (bc *BlockChain) GetBlock(blockHash []byte) (*Block, error) {

	var block *Block

	err := bc.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			blockBytes := b.Get(blockHash)

			block = DeserializeBlock(blockBytes)
		}

		return nil
	})
	return block, err
}

func (bc *BlockChain) AddBlock(block *Block) error {
	err := bc.DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			blockExist := b.Get(block.Hash)

			if blockExist != nil {
				//如果存在直接返回
				return nil
			}

			err := b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//最新区块的hash值
			blockHash := b.Get([]byte("l"))
			blockBytes := b.Get(blockHash)
			blockInDB := DeserializeBlock(blockBytes)

			if blockInDB.Height < block.Height {
				b.Put([]byte("l"), block.Hash)
				bc.Tip = block.Hash
			}
		}

		return nil
	})
	return err
}
