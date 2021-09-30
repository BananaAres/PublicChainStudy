/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/22 3:29 下午
 * Description:
 *
 */
package Utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
)

/**
 * @Author: sundaohan
 * @Description: int64转字节数组
 * @param num
 * @return []byte
 */
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

/**
 * @Author: sundaohan
 * @Description: 标准的JSON字符串转数组
 * @param json
 * @return []string
 */
func JSONtoArray(jsonStr string) []string {
	var sArr []string
	if err := json.Unmarshal([]byte(jsonStr), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}

/**
 * @Author: sundaohan
 * @Description: 字节数组反转
 * @param data
 */
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

/**
 * @Author: sundaohan
 * @Description: version转字节数组
 * @param command
 * @return []byte
 */
func CommandToBytes(command string) []byte {
	var bytes [COMMANDLENGTH]byte
	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

/**
 * @Author: sundaohan
 * @Description: 字节数组转version
 * @param bytes
 * @return string
 */
func BytesToCommand(bytes []byte) string {
	var command []byte
	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}

/**
 * @Author: sundaohan
 * @Description: 将结构体序列化成字节数组
 * @param data
 * @return []byte
 */
func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
