/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/8/28 10:54 上午
 * Description:
 *
 */
package BLC

type UTXO struct {
	TxHash   []byte
	Index    int
	TXOutput *TXOutput
}
