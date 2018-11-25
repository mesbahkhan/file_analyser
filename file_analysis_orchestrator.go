package main

import (
	"files_analyser/internal"
	"fmt"
	"os"
	"strconv"
)

func main() {

	var file_hash_sha256 string

	var file_information_row []string
	var file_information_table [][]string

	directory_name := "C:\\S\\Go"
	output_report_file_name := "output.csv"
	files_list, directory_contents, file_walk_errors := internal.Find_all_directory_content_recursive(directory_name)

	if file_walk_errors != nil {

		fmt.Println("error reading directory : %s", file_walk_errors)

		return
	}

	fmt.Println(files_list)

	for file_index, file := range files_list {

		file_information_row = nil
		file_information_row = append(file_information_row, directory_contents[file_index].Name())
		file_information_row = append(file_information_row, file)
		file_information_row = append(file_information_row, strconv.FormatInt(directory_contents[file_index].Size(), 10))
		file_information_row = append(file_information_row, directory_contents[file_index].ModTime().String())
		file_information_row = append(file_information_row, directory_contents[file_index].Mode().String())
		file_information_row = append(file_information_row, strconv.FormatBool(directory_contents[file_index].IsDir()))

		if directory_contents[file_index].IsDir() != true {
			file_hash_sha256 = internal.Calculate_file_hash_sha256(file)
			file_information_row = append(file_information_row, file_hash_sha256)
		}

		file_information_table = append(file_information_table, file_information_row)

	}

	output_file, err :=
		os.OpenFile(
			output_report_file_name,
			os.O_CREATE|os.O_WRONLY,
			0777)

	if err != nil {
		os.Exit(1)
	}

	internal.Write_2d_slice_set_to_csv(file_information_table, output_file)

}
