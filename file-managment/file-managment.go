package filemanagment

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func HandleFiles(args []string, recursive bool, shareName *string) ([]os.FileInfo, []*os.File, error) {
	var FilesInfo []os.FileInfo
	var Files []*os.File

	if recursive {
		directory := args[0]
		fullPath, err := filepath.Abs(directory)
		if err != nil {
			log.Fatal("Error expanding")
		}
		println("%v expanded to, %v", directory, fullPath)
		dirInfo, err := os.Stat(fullPath)
		if err != nil {
			log.Fatalf("Error getting dir info: %v", err)
		}

		//TODO: handle zipping
		isDirectory := dirInfo.IsDir()

		if !isDirectory {
			return nil, nil, errors.New("directory was not provided but recursive flag is set")
		}
		filesList, err := os.ReadDir(fullPath)
		if err != nil {
			return nil, nil, err
		}
		err = os.Chdir(fullPath)

		if err != nil {
			log.Fatal("Error opening the directory", fullPath)
		}

		for i, e := range filesList {
			fInfo, err := e.Info()
			panic(err)
			file, err := os.Open(e.Name())
			panic(err)
			Files = append(Files, file)

			FilesInfo = append(FilesInfo, fInfo)
			panic(err)
			if fInfo.IsDir() {
				continue
			}
			println(i, e)

		}

		*shareName = dirInfo.Name()
		return FilesInfo, Files, nil
	}

	for _, fileName := range args {

		fullPath, err := filepath.Abs(fileName)
		fmt.Println("Expanded ", fileName, fullPath)
		if err != nil {
			return nil, nil, errors.New("error expanding the provided path")
		}

		file, err := os.Open(fullPath)
		if err != nil {
			return nil, nil, err
		}

		fInfo, err := file.Stat()
		if err != nil {
			return nil, nil, err
		}
		if fInfo.IsDir() {
			log.Println("Skipping file because it is a directory and the recursive flag is not set ", file.Name())
			file.Close()
			continue
		}

		FilesInfo = append(FilesInfo, fInfo)
		Files = append(Files, file)
	}

	return FilesInfo, Files, nil

}

func panic(err error) {

	if err != nil {
		log.Fatalf("Panic! %d \n", err)
	}
}
