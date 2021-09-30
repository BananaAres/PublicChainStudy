/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/29 4:42 下午
 * Description:
 *
 */
package BLC

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sdhChain/BasicPrototype/Utils"
)

/**
 * @Author: sundaohan
 * @Description:
 * @param toAddr
 * @param bc
 */
func sendVersion(toAddr string, bc *BlockChain) {
	bestHeight := bc.GetBestHeight()
	payload := Utils.GobEncode(Version{Utils.NODE_VERSION, bestHeight, nodeAddr})
	request := append(Utils.CommandToBytes(Utils.COMMAND_VERSION), payload...)

	sendData(toAddr, request)
}

/**
 * @Author: sundaohan
 * @Description: 发送数据
 * @param to
 * @param data
 */
func sendData(to string, data []byte) {
	fmt.Println("客户端向服务器发送数据")
	conn, err := net.Dial(Utils.PROTOCOL, to)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

/**
 * @Author: sundaohan
 * @Description: 处理 COMMAND_GETBLOCKS
 * @param toAddr
 */
func sendGetBlocks(toAddr string) {
	payload := Utils.GobEncode(GetBlocks{nodeAddr})
	request := append(Utils.CommandToBytes(Utils.COMMAND_GETBLOCKS), payload...)
	sendData(toAddr, request)
}

/**
 * @Author: sundaohan
 * @Description: 主节点将自己的所有区块hash发送给钱包节点
 * @param toAddr
 * @param command
 * @param hashes
 */
func sendInv(toAddr string, kind string, hashes [][]byte) {
	payload := Utils.GobEncode(Inv{nodeAddr, kind, hashes})
	request := append(Utils.CommandToBytes(Utils.COMMAND_INV), payload...)
	sendData(toAddr, request)
}

/**
 * @Author: sundaohan
 * @Description: 钱包节点向主节点请求区块数据
 * @param toAddr
 * @param kind
 * @param blockHash
 */
func sendGetData(toAddr string, kind string, blockHash []byte) {
	payload := Utils.GobEncode(GetData{nodeAddr, kind, blockHash})
	request := append(Utils.CommandToBytes(Utils.COMMAND_GETDATA), payload...)
	sendData(toAddr, request)
}

/**
 * @Author: sundaohan
 * @Description: 主节点将数据发送给钱包节点
 * @param toAddr
 * @param block
 */
func sendBlock(toAddr string, block *Block) {
	payload := Utils.GobEncode(BlockData{nodeAddr, block})
	request := append(Utils.CommandToBytes(Utils.COMMAND_BLOCK), payload...)
	sendData(toAddr, request)
}
