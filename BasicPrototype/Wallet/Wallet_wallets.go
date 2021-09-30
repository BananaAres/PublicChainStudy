/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/1 3:12 下午
 * Description:
 *
 */
package Wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "Wallets_%s.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

/**
 * @Author: sundaohan
 * @Description: 钱包集合构造方法
 * @return *Wallets
 */
func NewWallets(nodeID string) (*Wallets, error) {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string]*Wallet)
		return wallets, err
	}
	filecontent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(filecontent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	return &wallets, nil
}

/**
 * @Author: sundaohan
 * @Description: 创建新钱包
 * @receiver w
 */
func (w *Wallets) CreateNewWallet(nodeID string) {
	wallet := NewWallet()
	fmt.Printf("Address : %s\n", wallet.GetAddress())
	w.WalletsMap[string(wallet.GetAddress())] = wallet
	//把所有数据保存起来
	w.SaveWallets(nodeID)
}

/**
 * @Author: sundaohan
 * @Description: 将钱包信息写入文件
 * @receiver w
 */
func (w *Wallets) SaveWallets(nodeID string) {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	var content bytes.Buffer
	// 注册的目的，是为了可以序列化任何类型
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}
	// 将序列化以后的数据写入到文件，原来文件的数据会被覆盖
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}

}
