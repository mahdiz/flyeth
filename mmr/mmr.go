package main

import (
	"fmt"
	"math/rand"
	"math"
	"crypto/sha256"
	"encoding/hex"
)

type Hash []byte

type Node struct {
	Left      *Node
	Hash      Hash
	Right     *Node
	LeafCount int
}

var MainRootNode *Node

var AllHashes []Hash

func IsPowerTwo(x int) bool {
	return (x & (x - 1)) == 0
}

// Highest power of two less than n
func HighestPowerTwo(n int) int {
	if IsPowerTwo(n) {
		n--
	}
	exp := 0
	for {
		if n <= 1 {
			break
		}
		if (n % 2 == 0) {
			n = n / 2
		} else {
			n = (n-1)/2
		}
		exp++
	}
	return int(math.Pow(2, float64(exp)))
}

func GetHash(b []byte) Hash {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}

func AddNode(root *Node, hash Hash) *Node {

	if root == nil {
		root = &Node{Hash: hash, LeafCount: 1}
	} else if IsPowerTwo(root.LeafCount) {
		node := &Node{Hash: hash, LeafCount: 1}
		newRoot := &Node{
			Hash: GetHash(append(root.Hash, hash...)),
			Left: root,
			Right: node,
			LeafCount: root.LeafCount + 1,
		}
		root = newRoot
	} else {
		root.Right = AddNode(root.Right, hash)
		root.Hash = GetHash(append(root.Right.Hash, root.Left.Hash...))
		root.LeafCount = root.LeafCount + 1
	}
	return root
}

func PrintMmr(root *Node, indent string, right bool) {
	if root == nil {
		return
	}
	fmt.Print(indent)
	if right {
		fmt.Print("\\-")
		indent += "\t"
	} else {
		fmt.Print("|-")
		indent += "|\t"
	}
	str := hex.EncodeToString(root.Hash)
	fmt.Println("0x" + str)

	PrintMmr(root.Left, indent, false)
	PrintMmr(root.Right, indent, true)
}

func GetMmrProof(root *Node, targetNumber int, m int) []Hash {
	if root == nil {//|| m + targetNumber > root.LeafCount {
		return []Hash{}
	}
	var proof []Hash

	if root.LeafCount == 1 {
		return []Hash { root.Hash }
	}
	if (targetNumber <= m + HighestPowerTwo(root.LeafCount)) {
		childProof := GetMmrProof(root.Left, targetNumber, m)
		proof = append([]Hash{root.Hash, childProof[0], root.Right.Hash}, childProof[1:]...)
	} else {
		m += HighestPowerTwo(root.LeafCount)
		childProof := GetMmrProof(root.Right, targetNumber, m)
		proof = append([]Hash{root.Hash, root.Left.Hash}, childProof...)
	}
	return proof
}

func main() {
	var mmr *Node
	randomBytes := make([]byte, 2)

	for i := 0; i < 500; i++ {
		rand.Read(randomBytes)
		blockHash := GetHash(randomBytes)
		mmr = AddNode(mmr, blockHash)
	}

	PrintMmr(mmr, "", true)

	GetMmrProof(mmr, 3, 0)
}
