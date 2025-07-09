package zip

import (
	"archive/tar"
	"compress/gzip"
	"github.com/bySouffle/bysutil/file"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func TarXVF(archive multipart.File, unzipPath string) (string, error) {
	_, err := archive.Seek(0, 0)
	if err != nil {
		return "", err
	}

	gzipReader, err := gzip.NewReader(archive)
	if err != nil {
		return "", err
	}
	defer func(gzipReader *gzip.Reader) {
		err := gzipReader.Close()
		if err != nil {
			log.Error(err)
		}
	}(gzipReader)

	tarReader := tar.NewReader(gzipReader)
	// 遍历 tar 文件中的每个文件
	var targetPath string
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// 读取完所有文件，退出循环
			break
		}
		if err != nil {
			return "", err
		}

		// 获取文件的路径和名称
		targetPath = filepath.Join(unzipPath, header.Name)
		if header.FileInfo().IsDir() {
			//	创建目录
			err := os.Mkdir(targetPath, os.FileMode(header.Mode))
			if err != nil {
				log.Error(err)
			}
		} else {
			// 创建文件
			file, err := file.CreateFile(targetPath, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Error(err)
				}
			}(file)

			// 将文件内容复制到目标文件
			_, err = io.Copy(file, tarReader)
			if err != nil {
				log.Error(err)
			}

			// 设置文件权限
			err = os.Chmod(targetPath, 0775)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return targetPath, nil
}
