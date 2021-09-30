/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/24 2:32 下午
 * Description:
 *
 */
package CLI

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sdhChain/BasicPrototype/Utils"
	"sdhChain/BasicPrototype/Wallet"
)

type Cli struct {
	//BlockChain *BLC.BlockChain
}

/**
 * @Author: sundaohan
 * @Description: 输出帮助信息
 */
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\taddresslists --输出所有钱包地址")
	fmt.Println("\tcreateblockchain -address --创建创世区块")
	fmt.Println("\tcreatewallet --创建钱包")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细")
	fmt.Println("\tprintchain --输出区块信息")
	fmt.Println("\tgetbalance -address --获取账户信息")
	fmt.Println("\tresetUTXO --重置")
	fmt.Println("\tstartnode -miner ADDRESS --启动节点服务器，并且指定挖矿奖励的地址")
}

/**
 * @Author: sudaohan
 * @Description: 检查输入是否合规
 */
func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

/**
 * @Author: sundaohan
 * @Description: 实现解析命令行操作
 * @receiver cli
 */
func (cli *Cli) Run() {
	isValidArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	fmt.Printf("NODE_ID:%s\n", nodeID)

	resetCmd := flag.NewFlagSet("resetUTXO", flag.ExitOnError)
	addressListsCmd := flag.NewFlagSet("addresslists", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账源地址")
	flagTo := sendBlockCmd.String("to", "", "转账目的地址")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额")
	flagCreateBlockChainAddress := createBlockChainCmd.String("address", "", "创世区块的地址")
	getBalanceWithAddress := getBalanceCmd.String("address", "", "要查询某一个账号的余额")
	flagMine := sendBlockCmd.Bool("mine", false, "是否在当前节点中立即验证")
	flagMiner := startNodeCmd.String("miner", "", "定义挖矿奖励地址")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "resetUTXO":
		err := resetCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addresslists":
		err := addressListsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}
		from := Utils.JSONtoArray(*flagFrom)
		to := Utils.JSONtoArray(*flagTo)

		for index, fromAddress := range from {
			if Wallet.IsValidForAddress([]byte(fromAddress)) == false || Wallet.IsValidForAddress([]byte(to[index])) == false {
				fmt.Printf("地址不合法.......")
				printUsage()
				os.Exit(1)
			}
		}

		amount := Utils.JSONtoArray(*flagAmount)
		cli.send(from, to, amount, nodeID, *flagMine)
	}

	if printChainCmd.Parsed() {
		cli.printChain(nodeID)
	}

	if resetCmd.Parsed() {
		cli.resetUTXOSet(nodeID)
	}

	if addressListsCmd.Parsed() {
		cli.addressLists(nodeID)
	}
	if createWalletCmd.Parsed() {
		// 创建钱包
		cli.createWallet(nodeID)
	}
	if createBlockChainCmd.Parsed() {
		if Wallet.IsValidForAddress([]byte(*flagCreateBlockChainAddress)) == false {
			fmt.Println("地址不合法......")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockChain(*flagCreateBlockChainAddress, nodeID)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceWithAddress == "" {
			fmt.Println("地址不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceWithAddress, nodeID)
	}
	if startNodeCmd.Parsed() {
		cli.startNode(nodeID, *flagMiner)
	}
}
