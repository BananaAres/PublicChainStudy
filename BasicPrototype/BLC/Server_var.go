/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/29 8:34 下午
 * Description: 存储节点全局变量
 *
 */
package BLC

//localhost:3000 主节点地址
var knowNodes = []string{"localhost:3000"}

//全局变量，节点地址
var nodeAddr string

//存储hash值
var transactionArray [][]byte
