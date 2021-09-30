/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/24 11:48 下午
 * Description:
 *
 */
package BLC

import (
	"bytes"
	"sdhChain/BasicPrototype/Wallet"
)

type TXInput struct {
	// 交易Hash
	TxHash []byte
	// 存储TXOutput在Vout里面的索引
	Vout int
	//数字签名
	Signature []byte
	//公钥
	PublicKey []byte
}

/**
 * @Author: sundaohan
 * @Description: 判断当前消费是否为传入地址的
 * @receiver txInput
 * @param addr
 * @return bool
 */
func (txInput *TXInput) UnlockWithRipemd160Hash(ripemd160Hash []byte) bool {
	publicKey := Wallet.Ripemd160Hash(txInput.PublicKey)
	return bytes.Compare(publicKey, ripemd160Hash) == 0
}
