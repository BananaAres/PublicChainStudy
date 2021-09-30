/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/29 3:00 下午
 * Description:
 *
 */
package BLC

type Version struct {
	Version    int    //版本
	BestHeight int64  //当前节点区块高度
	AddrFrom   string // 当前节点地址
}
