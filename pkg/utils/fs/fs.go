// fs package provides utility functions to work with files
package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ExecDir function returns the directory of the executable file.
// it is used to run all the commands in the directory of the executable file.
// this maintaning the integrity of the data.
func ExecDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", errors.New("failed to get executable path: " + err.Error())
	}
	return filepath.Dir(execPath), nil
}

// Path function returns the full path of a file.
// it is use to retrieve the path from the executable file.
func Path(path string) (string, error) {
	execDir, err := ExecDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(execDir, path), nil
}

// DBfile function returns the full path of the database file.
// it is used to retrieve the path of the database file.
func DBfile() (string, error) {
	dbfile, err := Path("data.db")
	if err != nil {
		return "", err
	}
	return dbfile, nil
}

// Backup function creates a backup of the database file.
// it is used to create a backup of the database file before making changes to the database.
func Backup() error {
	dbfile, err := DBfile() // get the database file path
	if err != nil {
		return err
	}

	// create a backup file with a timestamp
	backupfile, err := Path(fmt.Sprintf("data_backup_%s.db", time.Now().Format("20060102150405")))
	if err != nil {
		return err
	}

	// open the database file
	input, err := os.Open(dbfile)
	if err != nil {
		return errors.New("failed to open database file: " + err.Error())
	}

	defer input.Close()

	// create the backup file
	output, err := os.Create(backupfile)
	if err != nil {
		return errors.New("failed to create backup file: " + err.Error())
	}

	defer output.Close()

	// copy the database file to the backup file
	_, err = io.Copy(output, input)
	if err != nil {
		return errors.New("failed to copy database file: " + err.Error())
	}

	return nil
}
