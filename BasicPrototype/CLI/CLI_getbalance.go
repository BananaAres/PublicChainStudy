/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/28 11:45 上午
 * Description:
 *
 */
package CLI

import (
	"fmt"
	"sdhChain/BasicPrototype/BLC"
)

/**
 * @Author: sundaohan
 * @Description: 查询余额
 * @receiver cli
 * @param addr
 * @param nodeID
 */
func (cli *Cli) getBalance(addr string, nodeID string) {
	//fmt.Println(addr)
	bc := BLC.GetBlockChainObject(nodeID)
	defer bc.DB.Close()
	utxoSet := &BLC.UTXOSet{BlockChain: bc}
	amount := utxoSet.GetBalance(addr)
	fmt.Printf("%s一共有%d个Token\n", addr, amount)
}
