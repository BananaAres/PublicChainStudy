/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/29 4:41 下午
 * Description:
 *
 */
package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sdhChain/BasicPrototype/Utils"
)

func handleConnection(conn net.Conn, bc *BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Receive a Message:%s\n", request[:Utils.COMMANDLENGTH])
	command := Utils.BytesToCommand(request[:Utils.COMMANDLENGTH])
	fmt.Println(command)
	switch command {
	case Utils.COMMAND_VERSION:
		handleVersion(request, bc)
	case Utils.COMMAND_ADDR:
		handleAddr(request, bc)
	case Utils.COMMAND_BLOCK:
		handleBlock(request, bc)
	case Utils.COMMAND_GETBLOCKS:
		handleGetBlocks(request, bc)
	case Utils.COMMAND_GETDATA:
		handleGetData(request, bc)
	case Utils.COMMAND_INV:
		handleInv(request, bc)
	case Utils.COMMAND_TX:
		handleTX(request, bc)
	default:
		fmt.Println("Unknown command!")

	}
	defer conn.Close()
}

/**
 * @Author: sundaohan
 * @Description: 钱包节点检查自己的区块信息是否完整
 * @param request
 * @param bc
 */
func handleVersion(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload Version
	dataBytes := request[Utils.COMMANDLENGTH:]
	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	//version
	// 1. version
	// 2. BestHeight
	// 3.节点地址

	bestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if bestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	} else if bestHeight < foreignerBestHeight {
		// 向主节点要信息
		sendGetBlocks(payload.AddrFrom)
	}
}
func handleAddr(request []byte, bc *BlockChain) {

}
func handleBlock(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload BlockData

	dataBytes := request[Utils.COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	block := payload.Block
	bc.AddBlock(block)

	if len(transactionArray) > 0 {
		sendGetData(payload.AddrFrom, Utils.BLOCK_TYPE, transactionArray[0])
		transactionArray = transactionArray[1:]
	} else {
		utxoSet := &UTXOSet{bc}
		utxoSet.ResetUTXOSet()
	}

}

/**
 * @Author: sundaohan
 * @Description: 主节点将全部区块hash发送
 * @param request
 * @param bc
 */
func handleGetBlocks(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload GetBlocks
	dataBytes := request[Utils.COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, Utils.BLOCK_TYPE, blocks)
}

/**
 * @Author: sundaohan
 * @Description: 主节点处理请求区块信息的请求
 * @param request
 * @param bc
 */
func handleGetData(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload GetData
	dataBytes := request[Utils.COMMANDLENGTH:]
	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}
	if payload.Type == Utils.BLOCK_TYPE {
		block, err := bc.GetBlock(payload.Hash)
		if err != nil {
			return
		}
		sendBlock(payload.AddrFrom, block)
	}

	if payload.Type == Utils.TX_TYPE {

	}

}

/**
 * @Author: sundaohan
 * @Description: 钱包节点接收到主节点发送的区块hash
 * @param request
 * @param bc
 */
func handleInv(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload Inv
	dataBytes := request[Utils.COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == Utils.BLOCK_TYPE {
		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, Utils.BLOCK_TYPE, blockHash)

		if len(payload.Items) >= 1 {
			transactionArray = payload.Items[1:]
		}

	}

	if payload.Type == Utils.TX_TYPE {

	}

}
func handleTX(request []byte, bc *BlockChain) {

}
