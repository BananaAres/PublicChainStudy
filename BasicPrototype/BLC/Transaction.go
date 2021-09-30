/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/24 4:55 下午
 * Description:
 *
 */
package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
	"sdhChain/BasicPrototype/Utils"
	"sdhChain/BasicPrototype/Wallet"
	"time"
)

/**
 * @Author: sundaohan
 * @Description: UTXO
 */
type Transaction struct {
	//交易Hash
	TxHash []byte
	//输入
	Vins []*TXInput
	//输出
	Vouts []*TXOutput
}

/**
 * @Author: sundaohan
 * @Description: 判断当前交易是否是CoinBase交易
 * @receiver tx
 * @return bool
 */
func (tx *Transaction) IsCoinBaseTransaction() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

/**
 * @Author: sundaohan
 * @Description: 创建创世区块时的Transaction
 * @param address
 * @return *Transaction
 */
func NewCoinBaseTransaction(addr string) *Transaction {

	txInput := &TXInput{[]byte{}, -1, nil, []byte{}}
	txOutput := NewTXOutPut(10, addr)
	txCoinBase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	//设置Hash值
	txCoinBase.HashTransaction()

	return txCoinBase
}

/**
 * @Author: sundaohan
 * @Description: 生成交易hash
 * @receiver tx
 * @return []byte
 */
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	time := Utils.IntToHex(time.Now().Unix())
	resultBytes := bytes.Join([][]byte{time, result.Bytes()}, []byte{})
	hash := sha256.Sum256(resultBytes)
	tx.TxHash = hash[:]
}

/**
 * @Author: sundaohan
 * @Description: 转账时产生的Transaction
 * @param from
 * @param to
 * @param amount
 * @return *Transaction
 */
func NewNormalTransaction(from, to string, amount int64, utxoSet *UTXOSet, txs []*Transaction, nodeID string) *Transaction {

	wallets, _ := Wallet.NewWallets(nodeID)
	wallet := wallets.WalletsMap[from]

	money, spendableUTXODic := utxoSet.FindSpendableUTXOS(from, amount, txs)

	var txInputs []*TXInput
	var txOutputs []*TXOutput
	//消费
	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txInput := &TXInput{txHashBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}
	}
	//转账
	txOutput := NewTXOutPut(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput = NewTXOutPut(int64(money)-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)
	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	//设置Hash值
	tx.HashTransaction()
	//进行签名
	utxoSet.BlockChain.SignTransaction(tx, wallet.PrivateKey, txs)
	return tx
}

/**
 * @Author: sundaohan
 * @Description: 序列化
 * @receiver tx
 * @return []byte
 */
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

/**
 * @Author: sundaohan
 * @Description: 对Transaction中每个input签名
 * @receiver tx
 * @param privKey
 * @param prevTXs
 */
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinBaseTransaction() {
		return
	}

	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vins {
		prevTX := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTX.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		// 签名
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[inID].Signature = signature
	}

}

/**
 * @Author: sundaohan
 * @Description:
 */
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vins {
		inputs = append(inputs, &TXInput{vin.TxHash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TXOutput{vout.Value, vout.Ripemd160Hash})
	}

	txCopy := Transaction{tx.TxHash, inputs, outputs}

	return txCopy

}

/**
 * @Author: sundaohan
 * @Description: 验证数字签名
 * @receiver tx
 * @return bool
 */
func (tx *Transaction) verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinBaseTransaction() {
		return true
	}

	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vins {
		prevTX := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTX.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		//私钥ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			return false
		}
	}
	return true
}
