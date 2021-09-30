/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/28 11:45 上午
 * Description:
 *
 */
package CLI

import "sdhChain/BasicPrototype/BLC"

/**
 * @Author: sundaohan
 * @Description: 创建创世区块
 * @receiver cli
 * @param address
 * @param nodeID
 */
func (cli *Cli) createGenesisBlockChain(address string, nodeID string) {
	blockChain := BLC.CreateBlockChainWithGenesisBlock(address, nodeID)
	defer blockChain.DB.Close()
	utxoSet := &BLC.UTXOSet{BlockChain: blockChain}
	utxoSet.ResetUTXOSet()
}
