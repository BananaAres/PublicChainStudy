/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/24 11:49 下午
 * Description:
 *
 */
package BLC

import (
	"bytes"
	"sdhChain/BasicPrototype/Wallet"
)

type TXOutput struct {
	Value         int64
	Ripemd160Hash []byte //公钥
}

/**
 * @Author: sundaohan
 * @Description: 上锁
 * @receiver txOutput
 */
func (txOutput *TXOutput) Lock(addr string) {
	publicKeyHash := Wallet.Base58Decode([]byte(addr))
	txOutput.Ripemd160Hash = publicKeyHash[1 : len(publicKeyHash)-4]
}

/**
 * @Author: sundaohan
 * @Description: 解锁
 * @receiver txOutput
 * @param addr
 * @return bool
 */
func (txOutput *TXOutput) UnlockScriptPubKeyWithAddr(addr string) bool {

	publicKeyHash := Wallet.Base58Decode([]byte(addr))
	hash160 := publicKeyHash[1 : len(publicKeyHash)-4]
	return bytes.Compare(txOutput.Ripemd160Hash, hash160) == 0
}

/**
 * @Author: sundaohan
 * @Description: 创建新交易
 * @param value
 * @param addr
 * @return *TXOutput
 */
func NewTXOutPut(value int64, addr string) *TXOutput {
	tXOutPut := &TXOutput{value, nil}
	//设置Ripemd160Hash
	tXOutPut.Lock(addr)
	return tXOutPut
}
