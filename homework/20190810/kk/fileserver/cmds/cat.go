package cmds

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

type CatRequest struct {
	Path string
}

type CatResponse struct {
	Content []byte
}

type Cat struct {
	BaseDir string
}

func (c *Cat) Exec(request *CatRequest, response *CatResponse) error {
	path := filepath.Join(c.BaseDir, request.Path)
	log.Printf("[debug] cat %s", path)
	file, err := os.Open(path)
	if err != nil {
		log.Print("[error] cat file error: ", err)
		return errors.New("读取文件错误")
	}
	defer file.Close()

	cxt := make([]byte, 1024)
	n, err := file.Read(cxt)
	if err != nil && err != io.EOF {
		log.Print("[error] cat file error: ", err)
		return errors.New("读取文件错误")
	}
	response.Content = cxt[:n]
	return nil
}