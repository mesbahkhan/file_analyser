package internal

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"sync"
)

func Calculate_file_hash_sha256(file_path_and_name string) string {

	wait_group := new(sync.WaitGroup)
	file, err := os.Open(file_path_and_name)

	if err != nil {
		panic(err.Error())
	}

	file_hash := sha256.New()

	// 2 channels: used to give green light for reading into buffer b1 or b2
	read_data_channel, read_status_channel := make(chan int, 1), make(chan int, 1)

	// 2 channels: used to give green light for hashing the content of b1 or b2
	hash_data_channel, hash_status_channel := make(chan int, 1), make(chan int, 1)

	// Start signal: Allow b1 to be read and hashed
	read_data_channel <- 1
	hash_data_channel <- 1

	wait_group.Add(1)

	go hashHelper(file, file_hash, read_data_channel, read_status_channel, hash_data_channel, hash_status_channel, wait_group)

	wait_group.Add(1)

	hashHelper(file, file_hash, read_status_channel, read_data_channel, hash_status_channel, hash_data_channel, wait_group)
	wait_group.Wait()

	file_hash_code := fmt.Sprintf("%x", file_hash.Sum(nil))

	return file_hash_code
}

func hashHelper(f *os.File, h hash.Hash, mayRead <-chan int, readDone chan<- int, mayHash <-chan int, hashDone chan<- int, wait_group *sync.WaitGroup) {

	for b, hasMore := make([]byte, 8192<<10), true; hasMore; {
		<-mayRead
		n, err := f.Read(b)
		if err != nil {
			if err == io.EOF {
				hasMore = false
			} else {
				panic(err)
			}
		}
		readDone <- 1

		<-mayHash
		_, err = h.Write(b[:n])
		if err != nil {
			panic(err)
		}
		hashDone <- 1

	}
	wait_group.Done()

}
