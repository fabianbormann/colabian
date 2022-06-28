package persistence

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"colabian/core"
)

func Load (path string) []core.Block {
	blockchain := []core.Block{}

	file, error := os.Open(path)

	if os.IsNotExist(error) {
		fmt.Println("The provided file does not exist yet so a new blockchain was initialized.")
		return blockchain
	}

	if error != nil {
        fmt.Println(error)
		fmt.Println("Failed to read the blockchain file so a new blockchain was initialized.")
		return blockchain
	}

	reader, error := gzip.NewReader(file)

	if error != nil {
        fmt.Println(error)
		fmt.Println("Failed to read the blockchain file so a new blockchain was initialized.")
		return blockchain
	}

	defer file.Close()
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	blockCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "version") {
			fmt.Printf("Read blockchain from file %s (Version %s)\n", path, strings.Split(line, " ")[1])
		} else if strings.HasPrefix(line, "block") {
			extract := strings.Split(line, "|")
			_, id, hash, predecessor, created, data := extract[0], extract[1], extract[2], extract[3], extract[4], extract[5]
			const layout = "2006-01-02 15:04:05.999999999 -0700 MST"
			
			blockId, _ := strconv.ParseUint(id, 10, 32)
			blockCreated, error := time.Parse(layout, strings.Split(created, " m=")[0])

			if error != nil {
				fmt.Println(error)
				fmt.Println("Failed to read the blockchain file so a new blockchain was initialized.")
				return []core.Block{}
			}

			block := core.Block{
				Id: blockId,
				Hash: hash,
				Predecessor: predecessor,
				Created: blockCreated,
				Data: data,
			}
	
			blockchain = append(blockchain, block)
			blockCount += 1
		}
	}

	fmt.Printf("Successfully read %d blocks\n", blockCount)
	return blockchain
}

func Save (blockchain []core.Block, path string, version string) {
	data := "version " + version + "\n"
	data += fmt.Sprintf("length %d\n", len(blockchain))
	
	for _, block := range blockchain {
		data +=	fmt.Sprintf("block|%d|%s|%s|%s|%s\n", block.Id, block.Hash, block.Predecessor, block.Created.String(), block.Data)
	}

	if path != "" && !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

    var compressed bytes.Buffer
    gz := gzip.NewWriter(&compressed)
    if _, error := gz.Write([]byte(data)); error != nil {
        fmt.Println(error)
		return
    }
    if error := gz.Close(); error != nil {
        fmt.Println(error)
		return
    }

	error := ioutil.WriteFile(path + "blockchain.dat", compressed.Bytes(), 0777)
    if error != nil {
        fmt.Println(error)
    }
}