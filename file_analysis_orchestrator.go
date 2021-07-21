package main

import (
	"files_analyser/internal"
	"fmt"
	"github.com/c2h5oh/datasize"
	"github.com/sqweek/dialog"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
)

func main() {

	var sha_algorithm string
	var skip_files, batch int
	var delete_source_flag string
	var recursivity_flag string

	app := &cli.App{
		Name:    "Go File Analysis Suite",
		Usage:   "gofiles [mode]",
		Version: "0.0.2",
		Authors: []*cli.Author{
			&cli.Author{
				"Mesbah Khan",
				"khanm@ontoledgy.io"}},
		Copyright: "copyright 2020",

		Commands: []*cli.Command{
			{
				Name:    "hash",
				Aliases: []string{"h"},
				Usage:   "use it to create a hashtable for a directory. Select hashing algorithms using -hashAlgo",
				Flags: []cli.Flag{
					&cli.StringFlag{
						"hashAlgo",
						[]string{"sha"},
						"hashing algorithm, options : sha256, sha512, md5",
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
					&cli.IntFlag{
						"skipFiles",
						[]string{"skip"},
						"skip a number of files",
						[]string{""},
						"",
						false,
						false,
						0,
						"",
						&skip_files,
						false,
					},
					&cli.IntFlag{
						"batchSize",
						[]string{"batch"},
						"skip a number of files",
						[]string{""},
						"",
						false,
						false,
						0,
						"",
						&batch,
						false,
					},
				},
				Action: func(c *cli.Context) error {
					Start_file_hash_analysis(sha_algorithm, skip_files, batch)
					return nil

				},
			},
			{
				Name:    "copy",
				Aliases: []string{"c"},
				Usage:   "Use it to copy files using a csv loader with source and destination paths",
				Flags: []cli.Flag{
					&cli.StringFlag{
						"deleteSource",
						[]string{"ana"},
						"enable deletion of source file : use copy -deleteSource",
						[]string{""},
						"",
						false,
						false,
						false,
						"",
						"",
						&delete_source_flag,
						false,
					},
				},
				Action: func(c *cli.Context) error {
					Start_file_copy(delete_source_flag)
					return nil
				},
			},
			{
				Name:    "unzip",
				Aliases: []string{"u"},
				Usage:   "use it to unzip all zips within a directory. Select recursivity using -recursive yes or no",
				Flags: []cli.Flag{
					&cli.StringFlag{
						"recursivity",
						[]string{"sha"},
						"",
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

	directory_name, directory_selection_error :=
		dialog.Directory().
			Title("Select log file storage location").Browse()

	if directory_selection_error != nil {
		fmt.Println(
			directory_selection_error)
	}

	internal.Unzip_files_in_folder(directory_name, recursivity_flag)
}

func Start_file_hash_analysis(hashing_alogrithm string, skip int, batch int) {

	directory_name, directory_selection_error := dialog.Directory().
		Title("Select log file storage location").Browse()

	if directory_selection_error != nil {
		fmt.Println(directory_selection_error)
	}

	internal.Get_file_hashes_for_folder(directory_name, hashing_alogrithm, skip, batch)

}

func Start_file_copy(delete_source_flag string) {

	move_file_list_filename, move_file_selection_err :=
		dialog.File().Filter("Select mapping file", "csv").Load()

	if move_file_selection_err == nil {

		_, copy_err := internal.Copy_files(
			move_file_list_filename, delete_source_flag)

		if copy_err != nil {
			fmt.Println(copy_err)
		}

	} else {
		fmt.Println(
			move_file_selection_err)

	}

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
