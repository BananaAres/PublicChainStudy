/**
 * Author:  sundaohan
 * Version: 1.0.0
 * Date:    2021/9/28 4:37 下午
 * Description:
 *
 */
package BLC

import "crypto/sha256"

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

type MerkleTree struct {
	RootNode *MerkleNode
}

/**
 * @Author: sundaohan
 * @Description: 创建默克尔树
 * @param data
 * @return *MerkleNode
 */
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	//创建叶子节点
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}
	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
			if len(newLevel)%2 != 0 {
				newLevel = append(newLevel, newLevel[len(newLevel)-1])
			}
		}
		nodes = newLevel
	}
	mTree := MerkleTree{&nodes[0]}
	return &mTree
}

/**
 * @Author: sundaohan
 * @Description: 创建默克尔树的节点
 * @param left
 * @param right
 * @param data
 * @return *MerkleNode
 */
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}
	//创建叶子节点
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		//非叶子节点
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}
	mNode.Left = left
	mNode.Right = right
	return &mNode

}
