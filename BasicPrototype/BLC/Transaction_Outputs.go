/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/22 8:05 下午
 * Description:
 *
 */
package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutputs struct {
	UTXOS []*UTXO
}

/**
 * @Author: sundaohan
 * @Description: 将区块序列化成字节数组
 * @receiver b
 * @return []byte
 */
func (txs *TXOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(txs)
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
func DeserializeTXOutputs(txOutputsBytes []byte) *TXOutputs {
	var txOutputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}
	return &txOutputs
}
