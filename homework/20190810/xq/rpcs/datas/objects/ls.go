package objects

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	Name string
	Type string
	Size int64
	CreateTime time.Time
	ModifyTime time.Time
	AccessTime time.Time
}

type LsRequest struct {
	Path string
}
type LsResponse struct {
	FileInfos []FileInfo
}

type Ls struct {
	BaseDir string
}


func (l *Ls) Exec(request *LsRequest, response *LsResponse) error {
	path := filepath.Join(l.BaseDir, request.Path)
	log.Printf("[debug] ls %s ", path)
	file, err := os.Open(path)

	if err != nil {
		log.Print("[error] ls file error: ", err)
		return errors.New("read dir fails")

	}
	defer file.Close()

	fileInfos, err := file.Readdir(-1)
	if err != nil {
		log.Print("[error] ls file error: ", err)
		return errors.New("read file fails")
	}
	response.FileInfos = make([]FileInfo, len(fileInfos))
	for i, fileInfo := range fileInfos {
		fileType := "F"
		if fileInfo.IsDir() {
			fileType = "D"
		}
		response.FileInfos[i] = FileInfo{
			Name: fileInfo.Name(),
			Type: fileType,
			Size: fileInfo.Size(),
			CreateTime: time.Now(),
			ModifyTime: fileInfo.ModTime(),
			AccessTime: time.Now(),

		}
	}
	return nil
}

func init() {
	rpc.Register(&Ls{BaseDir: "."})
	fmt.Println("[info] register objects.Ls")
}