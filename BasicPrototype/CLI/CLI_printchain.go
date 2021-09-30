/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/28 11:43 上午
 * Description:
 *
 */
package CLI

import (
	"sdhChain/BasicPrototype/BLC"
)

/**
 * @Author: sundaohan
 * @Description: 遍历区块链
 * @receiver cli
 */
func (cli *Cli) printChain(nodeID string) {

	bc := BLC.GetBlockChainObject(nodeID)
	defer bc.DB.Close()
	bc.PrintChain()
}
