package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Id uint64
	Hash string
	Predecessor string
	Created time.Time
	Data string
	Nonce int
}

func computeHash(predecessorHash string, data string) (string, int) {
	matchCondition := false
	nonce := 0
	minedHash := ""

	for !matchCondition {
		hash := sha256.New()
		hash.Write([]byte(data))
		hash.Write([]byte(predecessorHash))
		hash.Write([]byte(fmt.Sprintf("%d", nonce)))
		minedHash = string(hash.Sum(nil))

		if strings.HasPrefix(minedHash, "0a") {
			matchCondition = true
		}
		nonce += 1
	}

	fmt.Println(minedHash, nonce)

	return minedHash, nonce
}

func Mine(data string, blockchain []Block) []Block {
	if len(blockchain) == 0 {
		hash := sha256.New()
		hash.Write([]byte(data))

		block := Block{
			0,
			string(hex.EncodeToString(hash.Sum(nil))),
			"aeebad4a796fcc2e15dc4c6061b45ed9b373f26adfc798ca7d2d8cc58182718e",
			time.Now(),
			data,
			0,
		}

		return append(blockchain, block)
	} else {
		predecessor := blockchain[len(blockchain) - 1]
		minedHash, nonce := computeHash(predecessor.Hash, data)

		block := Block{
			(predecessor.Id + 1),
			string(hex.EncodeToString([]byte(minedHash))),
			predecessor.Hash,
			time.Now(),
			data,
			nonce,
		}

		return append(blockchain, block)
	}
}

func Summarize(blockchain []Block) { 
	if len(blockchain) == 0 {
		fmt.Printf("blockchain length is %d\n", len(blockchain))
		return
	}

	block := blockchain[len(blockchain) - 1] 

	timestamp := block.Created
	created := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
	timestamp.Year(), timestamp.Month(), timestamp.Day(),
	timestamp.Hour(), timestamp.Minute(), timestamp.Second())

	fmt.Printf("blockchain length is %d\n", len(blockchain))
	fmt.Printf("The last block mined has the content:\n")
	fmt.Printf("Id: %d\nHash: %s\nPredecessor: %s\nCreated at:%s\nData: %s\nNonce: %d\n", block.Id + 1, block.Hash, block.Predecessor, created, block.Data, block.Nonce)
}