package filemanagment

import (
	"archive/zip"
	"errors"
	"io"
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
		logging.DebugLog("%v expanded to, %v", fileName, fullPath)
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

func ZipFiles(filesInfo []os.FileInfo, Files []*os.File) ([]os.FileInfo, []*os.File) {
	zipFile, err := os.CreateTemp("", "http-ostrich-*.zip")
	logging.DebugLog("Created zip file %s", zipFile.Name())
	if err != nil {
		logging.ErrorAndKill("Error creating the temporary file", err)
	}

	zipFileInfo, err := zipFile.Stat()
	if err != nil {
		logging.ErrorAndKill("", err)
	}

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for i, f := range Files {
		fInfo := filesInfo[i]
		fWriterInZip, err := zipWriter.Create(fInfo.Name())
		if err != nil {
			logging.ErrorAndKill("Error compressing the files", err)
		}
		logging.DebugLog("Successfully created write in the zip file for %s", fInfo.Name())

		if fInfo.IsDir() {
			logging.DebugLog("Skipping %s since it is a directory", fInfo.Name())
			f.Close()
			continue
		}

		if _, err := io.Copy(fWriterInZip, f); err != nil {
			logging.ErrorAndKill("Error trying to copy file into the zip file", err)
		}
		logging.DebugLog("Successfully copied file to the zip archive")
		f.Close()
	}
	// err = zipWriter.Close()
	// if err != nil {
	// 	logging.ErrorAndKill("Error trying to close the zip writer", err)
	// }
	logging.DebugLog("Zip file successfully populated")

	zipFile.Close()
	logging.DebugLog("Closing to reopen in read-only mode")
	zipFile, err = os.Open(zipFile.Name())
	if err != nil {
		logging.ErrorAndKill("Error re-opening the temporary zip file", err)
	}

	return []os.FileInfo{zipFileInfo}, []*os.File{zipFile}

}

func panic(err error) {

	if err != nil {
		logging.ErrorAndKill("PANIC during the file parsing", err)
	}
}
