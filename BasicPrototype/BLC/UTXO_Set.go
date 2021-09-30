/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/22 8:03 下午
 * Description:
 *
 */
package BLC

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"sdhChain/BasicPrototype/Wallet"
)

const utxoTableName = "utxoTable"

type UTXOSet struct {
	BlockChain *BlockChain
}

/**
 * @Author: sundaohan
 * @Description: 重置数据库表
 * @receiver utxoSet
 */
func (utxoSet *UTXOSet) ResetUTXOSet() {
	err := utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err != nil {
				log.Panic(err)
			}

		}
		b, _ = tx.CreateBucket([]byte(utxoTableName))
		if b != nil {
			//[string]*TXOutPuts
			txOutputsMap := utxoSet.BlockChain.FindUTXOMap()

			for keyHash, outs := range txOutputsMap {

				txHash, _ := hex.DecodeString(keyHash)

				b.Put(txHash, outs.Serialize())
			}

		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

/**
 * @Author: sundaohan
 * @Description: 根据地址到表中找到对应的utxo
 * @receiver utxoSet
 * @param addr
 * @return []*UTXO
 */
func (utxoSet *UTXOSet) findUTXOForAddress(addr string) []*UTXO {
	var utxos []*UTXO
	utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))
		//游标
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			//fmt.Printf("key = %s, value = %v\n", k, v)

			txOutputs := DeserializeTXOutputs(v)

			for _, utxo := range txOutputs.UTXOS {
				if utxo.TXOutput.UnlockScriptPubKeyWithAddr(addr) {
					utxos = append(utxos, utxo)
				}
			}

		}
		return nil
	})
	return utxos
}

/**
 * @Author: sundaohan
 * @Description: 查询余额
 * @receiver utxoSet
 * @param addr
 */
func (utxoSet *UTXOSet) GetBalance(addr string) int64 {

	UTXOS := utxoSet.findUTXOForAddress(addr)
	//
	var amount int64
	for _, utxo := range UTXOS {
		amount += utxo.TXOutput.Value
	}
	return amount
}

/**
 * @Author: sundaohan
 * @Description: 查询地址下未打包的utxo
 * @receiver utxoSet
 * @param from
 * @param txs
 * @return []*UTXO
 */
func (utxoSet *UTXOSet) FindUnpackSpendableUTXOS(from string, txs []*Transaction) []*UTXO {
	var unUTXOs []*UTXO
	spentTXOutputs := make(map[string][]int)

	for _, tx := range txs {
		if tx.IsCoinBaseTransaction() == false {
			for _, in := range tx.Vins {
				publicKeyHash := Wallet.Base58Decode([]byte(from))
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
	WORK1:
		for index, out := range tx.Vouts {
			if out.UnlockScriptPubKeyWithAddr(from) {
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
	return unUTXOs

}

/**
 * @Author: sundaohan
 * @Description: 返回金额和可用的utxo
 * @receiver utxoSet
 * @param from
 * @param amount
 * @param txs
 * @return int64
 * @return map[string][]int
 */
func (utxoSet *UTXOSet) FindSpendableUTXOS(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	unPackUTXOS := utxoSet.FindUnpackSpendableUTXOS(from, txs)

	spentableUTXO := make(map[string][]int)

	//计数
	var money int64 = 0
	//还没打包之前看钱是否够用
	for _, UTXO := range unPackUTXOS {
		money += UTXO.TXOutput.Value
		txHash := hex.EncodeToString(UTXO.TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash], UTXO.Index)
		if money >= amount {
			return money, spentableUTXO
		}
	}

	//未打包的钱还不够
	utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {
			c := b.Cursor()
		UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {

				txOutoputs := DeserializeTXOutputs(v)
				for _, utxo := range txOutoputs.UTXOS {
					money += utxo.TXOutput.Value
					txHash := hex.EncodeToString(utxo.TxHash)
					spentableUTXO[txHash] = append(spentableUTXO[txHash], utxo.Index)
					if money >= amount {
						break UTXOBREAK
					}
				}

			}
		}

		return nil
	})
	if money < amount {
		log.Panic("余额不足")
	}

	return money, spentableUTXO
}

/**
 * @Author: sundaohan
 * @Description: 更新
 * @receiver utxoSet
 */
func (utxoSet *UTXOSet) Update() {
	//最新的block
	block := utxoSet.BlockChain.Iterator().Next()
	ins := []*TXInput{}

	outsMap := make(map[string]*TXOutputs)

	//找到所有要删除的数据
	for _, tx := range block.Txs {
		for _, in := range tx.Vins {
			ins = append(ins, in)
		}
	}

	for _, tx := range block.Txs {

		utxos := []*UTXO{}
		for index, out := range tx.Vouts {
			isSpent := false
			for _, in := range ins {

				if in.Vout == index && bytes.Compare(tx.TxHash, in.TxHash) == 0 && bytes.Compare(out.Ripemd160Hash, Wallet.Ripemd160Hash(in.PublicKey)) == 0 {
					isSpent = true
					continue
				}

			}
			if isSpent == false {
				utxo := &UTXO{tx.TxHash, index, out}
				utxos = append(utxos, utxo)
			}

		}
		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.TxHash)
			outsMap[txHash] = &TXOutputs{utxos}
		}
	}

	err := utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {
			// 删除
			for _, in := range ins {
				txOutputsBytes := b.Get(in.TxHash)

				if len(txOutputsBytes) == 0 {
					continue
				}

				txOutputs := DeserializeTXOutputs(txOutputsBytes)
				UTXOS := []*UTXO{}
				isNeedDelete := false
				for _, utxo := range txOutputs.UTXOS {
					if in.Vout == utxo.Index && bytes.Compare(utxo.TXOutput.Ripemd160Hash, Wallet.Ripemd160Hash(in.PublicKey)) == 0 {
						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS, utxo)
						fmt.Println(utxo.TXOutput.Value)
					}
				}
				if isNeedDelete {
					b.Delete(in.TxHash)
					if len(UTXOS) > 0 {
						preTXOutputs := outsMap[hex.EncodeToString(in.TxHash)]
						preTXOutputs.UTXOS = append(preTXOutputs.UTXOS, UTXOS...)
						outsMap[hex.EncodeToString(in.TxHash)] = preTXOutputs
					}
				}
			}
			fmt.Println(outsMap)
			//新增
			for keyHash, outPuts := range outsMap {
				keyHashBytes, _ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes, outPuts.Serialize())
			}

		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

}
