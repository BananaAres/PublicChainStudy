/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/29 1:46 下午
 * Description:
 *
 */
package BLC

import (
	"fmt"
	"log"
	"net"
	"sdhChain/BasicPrototype/Utils"
)

/**
 * @Author: sundaohan
 * @Description: 启动服务
 * @param nodeID
 * @param minerAdd
 */
func StartServer(nodeID string, minerAdd string) {
	//当前节点的Ip地址
	nodeAddr = fmt.Sprintf("localhost:%s", nodeID)

	ln, err := net.Listen(Utils.PROTOCOL, nodeAddr)

	if err != nil {
		log.Panic(err)
	}

	defer ln.Close()
	bc := GetBlockChainObject(nodeID)
	if nodeAddr != knowNodes[0] {
		//此节点是钱包节点或者矿工节点，需要向主节点发送请求同步数据
		sendVersion(knowNodes[0], bc)
	}

	for {
		//接收客户端发送来的数据
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleConnection(conn, bc)
	}

}
