package logger

import (
	"fmt"
	"os"
)

// check if file or directory exists
func checkNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// checkPermission checks if the permission is denied
func checkPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// isNotExistMkDir checks if the directory exists, if not, create it
func isNotExistMkDir(src string) error {
	if notExist := checkNotExist(src); notExist == true {
		if err := mkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// mkDir creates a directory
func mkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// mustOpen opens a file, if the file does not exist, it will be created
func mustOpen(fileName, dir string) (*os.File, error) {
	perm := checkPermission(dir)
	if perm == true {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}

	err := isNotExistMkDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error during make dir %s, err: %s", dir, err)
	}

	f, err := os.OpenFile(dir+string(os.PathSeparator)+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err: %s", err)
	}

	return f, nil
}
