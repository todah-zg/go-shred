package shred

import (
	"crypto/rand"
	"errors"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func getChecksum(path string) (uint32, error) {
	fp, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer fp.Close()

	finfo, err := fp.Stat()
	if err != nil {
		return 0, err
	}
	fsize := finfo.Size()
	buff := make([]byte, fsize)

	_, err = fp.Read(buff)
	if err != nil {
		return 0, err
	}

	crc32q := crc32.MakeTable(0xD5828281)
	return crc32.Checksum(buff, crc32q), nil
}

func generateRandomData(size int32) []byte {
	buff := make([]byte, size)
	rand.Read(buff)
	return buff
}

// Test the overwriteFile() by doing a quick CRC32 checksum before and after
// filling the file with random data
func TestOverwriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")
	err := ioutil.WriteFile(path, generateRandomData(10*1024), 0655)

	checksum1, err := getChecksum(path)
	if err != nil {
		t.Fatalf("Error calculating checksum for %v", path)
	}

	fp, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("Error opening file %v", path)
	}
	err = overwriteFile(fp)
	if err != nil {
		t.Fatalf("Error overwriting file:")
	}

	checksum2, err := getChecksum(path)
	if err != nil {
		t.Fatalf("Error calculating checksum for %v", path)
	}

	if checksum1 == checksum2 {
		t.Fatalf("Error: file has not changed!")
	}
}

// Test if the deleteFile() correctly removes the file
func TestDeleteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")
	err := ioutil.WriteFile(path, generateRandomData(10*1024), 0655)

	if err = deleteFile(path); err != nil {
		t.Fatalf("Error deleting the file!")
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("File persists after deletion!")
	}
}

// Test the Shred(path) function
func TestShred_PathDoesntExist(t *testing.T) {
	// TODO: implementation
}

func TestShred_PathIsDirectory(t *testing.T) {
	// TODO: implementation
}

func TestShred_FileIsDeletedOnExit(t *testing.T) {
	// TODO: implementation
}
