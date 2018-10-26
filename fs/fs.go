package fs

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// IsExists checks the given path if exists.
func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsFile checks the given path is a file.
func IsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}

// Copy copies the src to dest with hard link for preference.
func Copy(dest, src string) (err error) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return
	}
	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", srcInfo.Name(), srcInfo.Mode().String())
	}
	destInfo, err := os.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(destInfo.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", destInfo.Name(), destInfo.Mode().String())
		}
		if os.SameFile(srcInfo, destInfo) {
			return
		}
	}
	if err = os.Link(src, dest); err == nil {
		return
	}
	err = CopyContents(src, dest)
	return
}

// CopyContents copies the src to dest with contents.
func CopyContents(dest, src string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// ReadLines reads the file lines to slice.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// WriteLines writes the lines to the specific file.
func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
