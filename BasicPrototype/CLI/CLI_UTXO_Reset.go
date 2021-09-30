/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/23 9:28 下午
 * Description:
 *
 */
package CLI

import (
	"fmt"
	"sdhChain/BasicPrototype/BLC"
)

func (cli *Cli) resetUTXOSet(nodeID string) {

	blockchain := BLC.GetBlockChainObject(nodeID)
	defer blockchain.DB.Close()
	utxoSet := &BLC.UTXOSet{blockchain}
	utxoSet.ResetUTXOSet()
	fmt.Println(blockchain.FindUTXOMap())
}
