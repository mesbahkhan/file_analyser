package internal

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"storage/csv"
	"strconv"
	"sync"
)

func Get_file_hashes_for_folder(directory_name string, hash_algorithm string) {
	var file_hash_sha256 string

	var file_information_row []string
	var file_information_table [][]string

	output_report_file_name := "output.csv"
	files_list, directory_contents, file_walk_errors := Find_all_directory_content_recursive(directory_name)

	if file_walk_errors != nil {

		fmt.Println("error reading directory : %s", file_walk_errors)

		return
	}
	log_file := Set_log_file()

	defer log_file.Close()

	for file_index, file := range files_list {
		log.Printf("calculting has for %s\r\n", file)
		file_information_row = nil
		file_information_row = append(
			file_information_row,
			strconv.Itoa(file_index),
			directory_contents[file_index].Name(),
			file,
			strconv.FormatInt(directory_contents[file_index].Size(), 10),
			directory_contents[file_index].ModTime().String(),
			directory_contents[file_index].Mode().String(),
			strconv.FormatBool(directory_contents[file_index].IsDir()))

		if directory_contents[file_index].IsDir() != true {
			file_hash_sha256 = Calculate_file_hash(file, hash_algorithm)
			file_information_row = append(file_information_row, file_hash_sha256)
		}

		file_information_table = append(file_information_table, file_information_row)

	}

	output_file, _ :=
		storage.Open_csv_file(
			output_report_file_name)

	storage.Write_2d_slice_set_to_csv(file_information_table, output_file)

}

func Calculate_file_hash(file_path_and_name string, hash_algorithm string) string {

	wait_group := new(sync.WaitGroup)
	file, err := os.Open(file_path_and_name)

	if err != nil {
		panic(err.Error())
	}
	var file_hash hash.Hash

	switch hash_algorithm {
	case "sha256":
		file_hash = sha256.New()
	case "sha512":
		file_hash = sha512.New()
	}

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
