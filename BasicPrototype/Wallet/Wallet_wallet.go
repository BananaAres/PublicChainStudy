/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/31 4:49 下午
 * Description:
 *
 */
package Wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const version = byte(0x00)
const addressCheckSumLen = 4

/**
 * @Author: sundaohan
 * @Description: 钱包结构体
 */
type Wallet struct {
	//私钥
	PrivateKey ecdsa.PrivateKey
	//公钥
	PublicKey []byte
}

/**
 * @Author: sundaohan
 * @Description: 钱包构造方法
 * @return *Wallet
 */
func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

/**
 * @Author: sundaohan
 * @Description: 构造公私钥对
 * @return ecdsa.PrivateKey
 * @return []byte
 */
func newKeyPair() (ecdsa.PrivateKey, []byte) {

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

/**
 * @Author: sundaohan
 * @Description: 为一个公钥生成checksum
 * @param payload
 * @return []byte
 */
func CheckSum(payload []byte) []byte {

	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressCheckSumLen]
}

/**
 * @Author: sundaohan
 * @Description: 根据公钥生成地址
 * @receiver w
 * @return string
 */
func (w *Wallet) GetAddress() []byte {
	// hash160
	ripemd160Hash := Ripemd160Hash(w.PublicKey)

	version_ripemd160Hash := append([]byte{version}, ripemd160Hash...)

	checkSumBytes := CheckSum(version_ripemd160Hash)
	// 25字节
	bytes := append(version_ripemd160Hash, checkSumBytes...)

	return Base58Encode(bytes)
}

/**
 * @Author: sundaohan
 * @Description: 生成ripemd160Hash
 * @receiver w
 * @param publicKey
 * @return []byte
 */
func Ripemd160Hash(publicKey []byte) []byte {
	// 256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)
	//160
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)
	return ripemd160.Sum(nil)
}

func IsValidForAddress(addr []byte) bool {
	version_public_checksumBytes := Base58Decode(addr)
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes)-addressCheckSumLen:]
	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes)-addressCheckSumLen]
	checkBytes := CheckSum(version_ripemd160)
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}
	return false
}
