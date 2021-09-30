/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/29 1:39 下午
 * Description:
 *
 */
package CLI

import (
	"fmt"
	"os"
	"sdhChain/BasicPrototype/BLC"
	"sdhChain/BasicPrototype/Wallet"
)

/**
 * @Author: sundaohan
 * @Description: 指定挖矿节点
 * @receiver cli
 * @param nodeID
 * @param minerAdd
 */
func (cli *Cli) startNode(nodeID string, minerAdd string) {
	//启动服务器

	if minerAdd == "" || Wallet.IsValidForAddress([]byte(minerAdd)) {
		//启动服务器
		fmt.Printf("启动服务器:localhost:%s", nodeID)
		BLC.StartServer(nodeID, minerAdd)
	} else {
		fmt.Println("指定的地址无效")
		os.Exit(0)
	}

}
