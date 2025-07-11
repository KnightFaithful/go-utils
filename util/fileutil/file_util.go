package fileutil

import (
	"bufio"
	"example.com/m/util/utilerror"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

func ReadFile(filePath string) ([]byte, *utilerror.UtilError) {
	fp, err := os.Open(filePath)
	if err != nil {
		return nil, utilerror.NewError(err.Error())
	}
	defer fp.Close()

	buff := make([]byte, 1024)
	var res []byte
	for {
		length, err := fp.Read(buff)
		if err == io.EOF || length < 0 {
			break
		}
		res = append(res, buff[:length]...)
	}
	return res, nil
}

// ReadFolder 获取该路径所有文件名，注意结果不是文件名，而是相对路径
func ReadFolder(filePath string) ([]string, *utilerror.UtilError) {
	var files []string

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, utilerror.NewError(err.Error())
	}
	return files, nil
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CreateFolderIfNotExist(filePath string) *utilerror.UtilError {
	parts := strings.Split(filePath, "/")
	if len(parts) <= 1 {
		return nil
	}
	//如果最后一个是文件名，则去掉文件名
	if strings.Contains(parts[len(parts)-1], ".") {
		parts = parts[:len(parts)-1]
	}
	folder := strings.Join(parts, "/")
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return utilerror.NewError(err.Error())
	}
	return nil
}

func WriteByteList(value []byte, filePath string) *utilerror.UtilError {
	if err := CreateFolderIfNotExist(filePath); err != nil {
		return err.Mark()
	}

	// 打开文件，如果不存在则创建
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	defer file.Close()

	// 创建一个带缓冲的写入器
	writer := bufio.NewWriter(file)

	// 写入数据
	if _, err := writer.Write(value); err != nil {
		return utilerror.NewError(err.Error())
	}

	// 刷新缓冲
	if err := writer.Flush(); err != nil {
		return utilerror.NewError(err.Error())
	}
	return nil
}

func WriteString(str string, filePath string) *utilerror.UtilError {
	return WriteByteList([]byte(str), filePath)
}

func GetStringMemorySize(s string) int64 {
	// 获取字符串长度
	strLen := uintptr(len(s))

	// 获取单个字符的大小
	charSize := unsafe.Sizeof(s[0])

	// 计算固定开销
	fixedOverhead := 2

	// 计算字符串内存大小
	totalSize := strLen*charSize + uintptr(fixedOverhead)

	return int64(totalSize)
}

func RemoveFile(filePathName string) *utilerror.UtilError {
	if !FileExists(filePathName) {
		return nil
	}
	err := os.Remove(filePathName)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	return nil
}

func DeletePath(path string) *utilerror.UtilError {
	err := os.RemoveAll(path)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	return nil
}

func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	paths := strings.Split(dstName, "/")
	dstPath := strings.Join(paths[:len(paths)-1], "/")
	if CreateFolderIfNotExist(dstPath) != nil {
		return
	}
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}
