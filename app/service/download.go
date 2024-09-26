package service

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"GoViewFile/library/utils"
)

// IsFileExist 判断文件是否存在，先判断重名，再根据大小判断
// 存在返回true，不存在返回false
func IsFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if filesize == info.Size() {
		return true
	}
	_ = os.Remove(filename)
	return false
}

// DownloadFile
// @title 下载文件至本地服务器
// @param url:文件地址
// @param localPath:本地存储路径
// @return string:文件路径
// @return error:错误信息
func DownloadFile(url string, localPath string) (string, error) {
	var (
		fSize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	tmpFilePath := localPath + ".download"
	client := new(http.Client)
	resp, err := client.Get(url)

	if err != nil {
		return "", err
	}
	fSize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)

	if err != nil {
		log.Println("Error: <", err, "> when get file remote size")
		return "", err
	}
	if IsFileExist(localPath, fSize) {
		return "had", nil
	}
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if resp.Body == nil {
		return "", errors.New("body is null")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := file.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	if err == nil {
		_ = file.Close()
		fileMd5 := utils.GetFileMD5(tmpFilePath)
		newPath := "cache/download/" + fileMd5 + path.Ext(localPath)
		_ = os.Rename(tmpFilePath, newPath)
		log.Printf("Download file <filename:%s, md5:%s> success\n", path.Base(localPath), fileMd5)
		return newPath, nil
	}
	return "", err
}
