/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/22 6:28 下午
 * Description:
 *
 */
package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"sdhChain/BasicPrototype/Utils"
)

type ProofOfWork struct {
	Block  *Block   //当前要验证的区块
	target *big.Int //大数存储
}

//256位hash里面至少有16个零
const targetBit = 16

/**
 * @Author: sundaohan
 * @Description:创建新的工作量证明对象
 * @param block
 * @return *ProofOfWork
 */
func NewProofOfWork(block *Block) *ProofOfWork {
	//创建一个初始值为1的target
	target := big.NewInt(1)
	//左移256-targetBit
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{
		block,
		target,
	}
}

/**
 * @Author: sundaohan
 * @Description: 拼接数据，返回字符数组
 * @receiver pow
 * @param nonce
 * @return []byte
 */
func (proofOfWork *ProofOfWork) prepareData(nonce int) []byte {
	join := bytes.Join(
		[][]byte{
			proofOfWork.Block.PreBlockHash,
			proofOfWork.Block.HashTransactions(),
			Utils.IntToHex(proofOfWork.Block.TimeStamp),
			Utils.IntToHex(int64(targetBit)),
			Utils.IntToHex(int64(nonce)),
			Utils.IntToHex(int64(proofOfWork.Block.Height)),
		},
		[]byte{},
	)
	return join
}

/**
 * @Author: sundaohan
 * @Description: 检查hash是否合法
 * @receiver proofOfWork
 * @return bool
 */
func (proofOfWork *ProofOfWork) IsValid() bool {

	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.Block.Hash)
	if proofOfWork.target.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}

/**
 * @Author: sundaohan
 * @Description: 挖矿操作
 * @receiver ProofOfWork
 * @return []byte
 * @return int64
 */
func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {

	nonce := 0
	var hashInt big.Int //存储我们新生成的hash
	var hash [32]byte
	for {
		// 将BLOCK属性拼接成字节数组
		dataBytes := proofOfWork.prepareData(nonce)
		// 生成hash, sum256返回32位需要转换为64位
		hash = sha256.Sum256(dataBytes)
		// 将hash存储到hashInt,采取hash[:]将切片转换为64位
		hashInt.SetBytes(hash[:])
		fmt.Printf("\r%x", hash)
		// 判断hashInt是否小于Block里面的target
		// x < y -1
		// x == y 0
		// x > y 1
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			//判断有效性，如果满足条件，跳出循环
			break
		}
		nonce = nonce + 1
	}
	return hash[:], int64(nonce)
}
