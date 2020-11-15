package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip_files_in_folder(directory_name string, recursivity_flag string) {

	var zip_files []string
	//add directory information

	files_list, _, file_walk_errors :=
		Find_all_directory_content_recursive(
			directory_name)

	if file_walk_errors != nil {
		print(file_walk_errors)
	}

	zip_files = get_zips(files_list)

	unzip_file_set(zip_files)

}

func get_zips(files_list []string) []string {

	var zip_files []string
	var file_is_zip bool
	var file_extension string //TODO - use the file infromation to process this

	for _, file := range files_list {

		file_extension = file[len(file)-3:]

		file_is_zip = file_extension == "zip"

		if file_is_zip {
			zip_files = append(zip_files, file)
		}

	}

	return zip_files
}

func unzip_file_set(zip_files []string) {
	//TODO add logging
	var zip_files_within_zip []string

	for _, zip_file := range zip_files {

		unzipped_file_names, zip_errors := unzip_file(zip_file)

		if zip_errors != nil {
			print(zip_errors)
		}

		zip_files_within_zip = get_zips(unzipped_file_names)

		if zip_files != nil {
			unzip_file_set(zip_files_within_zip)
		}

	}

}

func unzip_file(src string) ([]string, error) {

	var dest string
	var filenames []string
	dest = src + ".unzipped"

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil

}
