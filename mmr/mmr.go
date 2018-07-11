package mmr


import (
	"fmt"
	"math/rand"
	"github.com/ethereum/go-ethereum/common"
)

type Node struct {
	Left      *Node
	Hash      common.Hash
	Right     *Node
	LeafCount int
}

var MainRootNode *Node

var AllHashes []common.Hash

func IsPowerTwo(x int) bool {
	return (x & (x - 1)) == 0
}

func AddNode(root *Node, hash common.Hash) *Node {

	if root == nil {
		node := &Node{ Hash: hash }
		root = &Node{
			Left: node,
			Hash: hash,
			LeafCount: 1,
		}
	} else if IsPowerTwo(root.LeafCount) {
		node := &Node{Hash: hash}
		newRoot := &Node{
			Hash: common.BytesToHash(append(root.Hash.Bytes(), hash.Bytes()...)),
			Left: root,
			Right: node,
			LeafCount: root.LeafCount + 1,
		}
		root = newRoot
	} else {
		root.Right = AddNode(root.Right, hash)
		root.Hash = common.BytesToHash(append(root.Right.Hash.Bytes(), root.Left.Hash.Bytes()...))
		root.LeafCount = root.LeafCount + 1
	}
	return root
}


func main() {
	var mmr *Node
	randomHash := make([]byte, 4)

	for i := 0; i < 100; i++ {
		rand.Read(randomHash)
		blockHash := common.BytesToHash(randomHash)
		mmr = AddNode(mmr, blockHash)
		fmt.Println(mmr.LeafCount)
	}
}
