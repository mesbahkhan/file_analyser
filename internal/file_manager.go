package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	storage "storage/csv"
)

func Copy_files(move_file_list_filename string, delete_source_flag string) (int, error) {

	move_file_list_file, csv_data := storage.Open_csv_file(move_file_list_filename)
	move_file_list := storage.Read_csv_to_slice(move_file_list_file, csv_data, "")

	var error_count int = 0

	fmt.Println(len(move_file_list))
	fmt.Println(move_file_list)

	log_file := Set_log_file()

	defer log_file.Close()

	for index, row := range move_file_list {

		if index != 0 {

			log.Printf(
				"copying %s to %s \n", row[0], row[1])

			sourceFileStat, source_file_stats_error :=
				os.Stat(row[0])

			if source_file_stats_error != nil {
				log.Printf("source file error: %s\n", source_file_stats_error)
				error_count += 1
				continue
			}

			if !sourceFileStat.Mode().IsRegular() {
				log.Printf("%s", fmt.Errorf("%s is not a regular file\n", row[0]))
				error_count += 1
			}

			source, source_file_open_error := os.Open(row[0])

			if source_file_open_error != nil {
				log.Printf("source file error %s\n", source_file_stats_error)
				error_count += 1
				continue
			}
			defer source.Close()

			destination_directory := filepath.Dir(row[1])

			_, destination_directory_stats_error :=
				os.Stat(
					destination_directory)

			if os.IsNotExist(destination_directory_stats_error) {
				log.Printf("target directory %s does not exits, creating now\n", destination_directory)
				os.MkdirAll(destination_directory, os.ModePerm)

			}

			destination, destination_file_stats_error :=
				os.Create(row[1])

			if destination_file_stats_error != nil {
				log.Printf("destination file error: %s\n", source_file_stats_error)
				error_count += 1
				continue
			}

			defer destination.Close()

			if delete_source_flag == "yes" {
				file_move_error := os.Rename(source.Name(), destination.Name())

				if file_move_error != nil {
					log.Printf("cannot move file due to %v\n", file_move_error)
					error_count += 1
				} else {
					bytes_copied, file_copy_error := io.Copy(destination, source)

					if file_copy_error != nil {
						log.Printf("copied file error: %v\n", file_copy_error)
						error_count += 1
					}
					log.Printf("sucessfully copied %v bytes\n", bytes_copied)
				}

			}

		}

	}
	log.Printf("Process completed with %v errors\n", error_count)
	return error_count, nil
}
