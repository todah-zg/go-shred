package shred

import (
	"crypto/rand"
	"os"
)

func Shred(path string) error {

	// Open the file and do basic sanity checking

	fp, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	// defer file deletion and closing in reverse (stack) order
	defer deleteFile(path)
	defer fp.Close()

	finfo, err := fp.Stat()
	if err != nil || finfo.IsDir() {
		return err
	}

	// Overwrite the file 3 times with random data

	for i := 0; i < 3; i++ {
		err := overwriteFile(fp)
		if err != nil {
			return err
		}
	}

	// defered Close and Delete here

	return nil
}

func overwriteFile(fp *os.File) error {

	finfo, _ := fp.Stat()
	size := finfo.Size()
	buff := make([]byte, size)

	_, err := rand.Read(buff)
	if err != nil {
		return err
	}

	_, err = fp.WriteAt(buff, 0)
	if err != nil {
		return err
	}

	return nil
}

func deleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
