/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/28 11:44 上午
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
 * @Description: 转账方法
 * @param from
 * @param to
 * @param amount
 */
func (cli *Cli) send(from, to, amount []string, nodeID string, mineNow bool) {
	blockChain := BLC.GetBlockChainObject(nodeID)
	defer blockChain.DB.Close()
	if mineNow {
		blockChain.MineNewBlock(from, to, amount, nodeID)
		utxoSet := &BLC.UTXOSet{blockChain}
		//转账成功后，进行更新
		utxoSet.Update()
	} else {
		//把交易发送给矿工节点验证
		fmt.Println("由矿工节点处理")
	}
}
