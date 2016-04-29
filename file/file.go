package file

import (
	"fmt"
	"io"
	"os"
)

//IsFileExists 目录或文件是否存在
func IsFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

//IsFile 判断是否为文件
func IsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}

//CopyFile 文件拷贝,优先使用hard link
func CopyFile(src, dest string) (err error) {
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
	err = CopyFileContents(src, dest)
	return
}

//CopyFileContents 文件内容拷贝
func CopyFileContents(src, dest string) (err error) {
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
