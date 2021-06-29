package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	dataPath = "./"
	//fileName = "test.txt"
)

type FileInfo struct {
	Name string
	MD5  string
}

func main() {
	if len(os.Args) > 1 {
		dataPath = os.Args[1]
	}

	files, _ := grabDirectory(dataPath)
	searchDuplicates(files)

}

func grabDirectory(dataPath string) ([]string, error) {
	//fmt.Printf("Scan from dir - %s\n", dataPath)

	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		fmt.Printf("error occurred with a ReadDir: %v", err.Error())
	}

	var filesArray []string
	for _, file := range files {
		filePath := filepath.Join(dataPath, file.Name())
		if file.IsDir() {
			files, err := grabDirectory(filePath)
			if err != nil {
				return nil, err
			}
			filesArray = append(filesArray, files...)
		}
		//if filepath.Ext(strings.TrimSpace(filePath)) == ".txt" {
			filesArray = append(filesArray, filePath)
		//}
	}

	return filesArray, nil
}

func getMD5SumString(file *FileInfo) (string, error) {
	f, err := os.Open(file.Name)
	if err != nil{
		return "", err
	}
	defer f.Close()

	md5 := md5.New()
	_, err = io.Copy(md5, f)
	if err != nil{
		return "", err
	}
	file.MD5 = fmt.Sprintf("%X", md5.Sum(nil))
	return file.MD5, nil
}

func searchDuplicates(files []string) {
	data := map[int64][]*FileInfo{}

	for _, file := range files{
		fileinfo, _ := os.Stat(file)
		data[fileinfo.Size()] = append(data[fileinfo.Size()], &FileInfo{Name: file})
	}

	hashes := map[string][]*FileInfo{}
	for _, list := range data{
		if len(list) > 1 {
			for _, file := range list{
				if hash, err := getMD5SumString(file); hash != "" && err == nil{
					hashes[hash] = append(hashes[hash], file)
				}
			}
			for _, list := range hashes{
				if len(list) > 1{
					fmt.Print("Files that have the same content: ")
					var str []string
					for _, file := range list{
						str = append(str, file.Name)
					}
					fmt.Println(str)
				}
			}
		}
	}
}
