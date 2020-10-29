package main

import (
	"files_analyser/internal"
	"fmt"
	"github.com/c2h5oh/datasize"
	"github.com/sqweek/dialog"
	"gopkg.in/urfave/cli.v2"
	"io"
	"log"
	"os"
	"storage/csv"
)

func main() {

	var sha_algorithm string
	var analysis_flag string
	var recursivity_flag string

	app := &cli.App{
		Name:    "Go File Analysis Suite",
		Usage:   "gofiles [mode]",
		Version: "0.0.1",
		Authors: []*cli.Author{
			&cli.Author{
				"Mesbah Khan",
				"khanm@ontoledgy.io"}},
		Copyright: "copyright 2019",

		Commands: []*cli.Command{
			{
				Name:    "hash",
				Aliases: []string{"h"},
				Usage:   "use it to create a hashtable for a directory. Select hashing algorithms using -hashAlgo",
				Flags: []cli.Flag{
					&cli.StringFlag{
						"hashAlgo",
						[]string{"sha"},
						"hashing algorithm, options : sha256, sha512",
						[]string{""},
						"",
						false,
						false,
						false,
						"",
						"sha256",
						&sha_algorithm,
						false,
					},
				},
				Action: func(c *cli.Context) error {
					Start_file_hash_analysis(sha_algorithm)
					return nil

				},
			},
			{
				Name:    "copy",
				Aliases: []string{"c"},
				Usage:   "Use it to copy files using a csv loader with source and destination paths",
				Flags: []cli.Flag{
					&cli.StringFlag{
						"Analysis",
						[]string{"ana"},
						"hashing algorithm, options : sha256, sha512",
						[]string{""},
						"",
						true,
						false,
						false,
						"",
						"",
						&analysis_flag,
						false,
					},
				},
				Action: func(c *cli.Context) error {
					if analysis_flag == "yes" {
						Analyse_source_folder()
					}
					Start_file_copy()
					return nil
				},
			},
			{
				Name:    "unzip",
				Aliases: []string{"u"},
				Usage:   "use it to unzip all zips within a directory. select recursivity using -recursive",
				Flags: []cli.Flag{
					&cli.StringFlag{
						"recursivity",
						[]string{"sha"},
						"hashing algorithm, options : sha256, sha512",
						[]string{""},
						"",
						false,
						false,
						false,
						"",
						"yes",
						&recursivity_flag,
						false,
					},
				},
				Action: func(c *cli.Context) error {
					Unzip_files_in_folder(recursivity_flag)
					return nil
				},
			},
			{
				Name:    "report",
				Aliases: []string{"a"},
				Usage:   "Use get an anlysis report on folder",
				Action: func(c *cli.Context) error {
					Analyse_source_folder()
					return nil
				},
			},
		},
	}

	app.Run(os.Args)

}

func Unzip_files_in_folder(recursivity_flag string) {
	directory_name, directory_selection_error := dialog.Directory().Title("Select log file storage location").Browse()

	if directory_selection_error != nil {
		fmt.Println(directory_selection_error)
	}

	internal.Unzip_files_in_folder(directory_name, recursivity_flag)
}

func Start_file_hash_analysis(hashing_alogrithm string) {

	directory_name, directory_selection_error := dialog.Directory().Title("Select log file storage location").Browse()

	if directory_selection_error != nil {
		fmt.Println(directory_selection_error)
	}

	internal.Get_file_hashes_for_folder(directory_name, hashing_alogrithm)

}

func Start_file_copy() {

	move_file_list_filename, move_file_selection_err := dialog.File().Filter("Select mapping file", "csv").Load()
	fmt.Println(move_file_selection_err)
	//Analyse_source_folder()
	_, copy_err := Copy_files(move_file_list_filename)
	fmt.Println(copy_err)

}

func Analyse_source_folder() {

	directory_name, directory_selection_error := dialog.Directory().Title("Select location to analyse").Browse()

	if directory_selection_error != nil {
		fmt.Println(directory_selection_error)

	}

	files, file_information, _ := internal.Find_all_directory_content_recursive(directory_name)

	number_of_files := len(files)
	var total_size int64
	var largest_file_size int64
	var largest_file_name string
	var longest_path_length int
	var longest_file_name string

	longest_path_length = 0
	total_size = 0

	for _, file := range file_information {

		total_size += file.Size()
		if file.Size() > largest_file_size {
			largest_file_size = file.Size()
			largest_file_name = file.Name()
		}
	}

	log_file := internal.Set_log_file()

	for _, file := range files {

		file_length := len(file)
		if file_length > longest_path_length {
			longest_path_length = file_length
			longest_file_name = file
		}

	}

	log.Printf("Source directory Name: %v\r\n", directory_name)
	log.Printf("Log file path: %v\r\n", log_file.Name())
	log.Printf("Total Number of files: %v\r\n", number_of_files)
	log.Printf("Total size: %s\r", datasize.ByteSize(total_size).HumanReadable())
	log.Printf("Largest file name: %v\r\n", largest_file_name)
	log.Printf("Largest file size: %s\r\n", datasize.ByteSize(largest_file_size).HumanReadable())
	log.Printf("Longest file path length: %v\r\n", longest_path_length)
	log.Printf("Longest file path name: %v\r\n", longest_file_name)

	log_file.Close()

}

func Copy_files(move_file_list_filename string) (int64, error) {

	move_file_list_file, csv_data := storage.Open_csv_file(move_file_list_filename)
	move_file_list := storage.Read_csv_to_slice(move_file_list_file, csv_data, "")

	error_count := 0
	fmt.Println(len(move_file_list))
	fmt.Println(move_file_list)

	for index, row := range move_file_list {
		if index != 0 {
			log.Printf("copying %s to %s \n", row[0], row[1])
			sourceFileStat, err := os.Stat(row[0])
			if err != nil {
				log.Printf("source file error: %s\n", err)
				error_count += 1
				continue
			}

			if !sourceFileStat.Mode().IsRegular() {
				log.Printf("%s", fmt.Errorf("%s is not a regular file\n", row[0]))
			}

			source, err := os.Open(row[0])
			if err != nil {
				log.Printf("source file error %s\n", err)
				error_count += 1
				continue
			}
			defer source.Close()

			destination, err := os.Create(row[1])
			if err != nil {
				log.Printf("destination file error: %s\n", err)
				error_count += 1
				continue
			}
			defer destination.Close()
			nBytes, err := io.Copy(destination, source)
			log.Printf("copied bytes %v, error: %v\n", nBytes, err)
		}

	}
	log.Printf("Process completed with %v errors\n", error_count)
	return 0, nil
}
