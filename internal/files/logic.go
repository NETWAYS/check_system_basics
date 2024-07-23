package files

import (
	//"fmt"
	"io/fs"
	"os"
)

func GetFileList(config FilesConfig) ([]os.FileInfo, error) {

	// TODO
	return []os.FileInfo{}, nil
}

func FilterFiles(config FilesConfig, rawFileList []os.FileInfo) ([]os.FileInfo, error) {

	// TODO
	return []os.FileInfo{}, nil
}

var fileList = make([]os.FileInfo, 0)

func GenFileListOfDir(directory string) ([]os.FileInfo, error) {
	//fmt.Printf("GenFileListOfDir: %s\n", directory)

	fs.WalkDir(os.DirFS(directory), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == "." {
			return nil
		}

		fi, err := os.Lstat(path)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			subList, err := GenFileListOfDir(path + string(os.PathSeparator) + fi.Name())
			if err != nil {
				return err
			}

			fileList = append(fileList, subList...)
		}

		fileList = append(fileList, fi)

		return nil
	})

	return fileList, nil
}
