package filemanagment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/v1ctorio/http-ostrich/logging"
)

func HandleFiles(args []string, recursive bool, shareName *string) ([]os.FileInfo, []*os.File, error) {
	var FilesInfo []os.FileInfo
	var Files []*os.File

	if recursive {
		directory := args[0]

		fullPath, err := filepath.Abs(directory)
		if err != nil {
			logging.ErrorAndKill("Error expanding the provided path", err)
		}
		logging.DebugLog("%v expanded to, %v", directory, fullPath)
		dirInfo, err := os.Stat(fullPath)
		if err != nil {
			logging.ErrorAndKill("Error getting dir info", err)
		}

		//TODO: handle zipping
		isDirectory := dirInfo.IsDir()

		if !isDirectory {
			return nil, nil, errors.New("directory was not provided but recursive flag is set")
		}
		filesList, err := os.ReadDir(fullPath)
		panic(err)
		err = os.Chdir(fullPath)

		panic(err)
		for _, e := range filesList {
			fInfo, err := e.Info()
			panic(err)
			file, err := os.Open(e.Name())
			panic(err)
			Files = append(Files, file)

			FilesInfo = append(FilesInfo, fInfo)
			panic(err)
			if fInfo.IsDir() {
				logging.DebugLog("Skipping directory %s", fInfo.Name())
				continue
			}
			logging.DebugLog("File discovered: %s", e.Name())

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
			println("Skipping file because it is a directory and the recursive flag is not set ", file.Name())
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
		logging.ErrorAndKill("PANIC during the file parsing", err)
	}
}
