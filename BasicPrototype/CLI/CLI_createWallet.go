/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/3 9:24 上午
 * Description:
 *
 */
package CLI

import (
	"fmt"
	"sdhChain/BasicPrototype/Wallet"
)

/**
 * @Author: sundaohan
 * @Description: 创建钱包
 * @receiver cli
 */
func (cli *Cli) createWallet(nodeID string) {
	wallets, _ := Wallet.NewWallets(nodeID)
	wallets.CreateNewWallet(nodeID)

	fmt.Println(len(wallets.WalletsMap))
}
