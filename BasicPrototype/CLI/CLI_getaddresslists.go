/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/3 11:10 上午
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
 * @Description: 打印创建的所有钱包地址
 * @receiver cli
 */
func (cli *Cli) addressLists(nodeID string) {
	fmt.Println("打印所有的钱包地址")

	wallets, _ := Wallet.NewWallets(nodeID)
	for address, _ := range wallets.WalletsMap {
		fmt.Println(address)
	}
}
