package YSZipUtil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"yuensoft.com/YSPathUtil"
)

// 创建Zip文件包
// zipFileName string：为Zip文件包的全路径＋名字
// includeFiles []string：为Zip文件包内要包含的文件
func CreateZip(zipFileName string, includeFiles []string) error {

	var (
		zipFile *os.File
		err     error
	)

	//检查路径，没有路径即创建
	if isExist, _ := YSPathUtil.IsPathExist(zipFileName); !isExist {
		dir := path.Dir(zipFileName)
		os.MkdirAll(dir, 0755)

	}
	//创建压缩文件并打开
	zipFile, err = os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	//zip.Writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	//逐个文件，写入zip包
	for _, includeFileName := range includeFiles {
		fmt.Printf("Zip <<-- %s", includeFileName)
		if err := writeFileToZip(zipWriter, includeFileName); err != nil {
			return err
		}
	}

	return nil
}

// 逐个文件，写入zip包
// zipper *zip.Writer：要写入的zip包的zip.Writer
// includeFileName string：本次写入的文件的全路径＋名字
func writeFileToZip(zipper *zip.Writer, includeFileName string) error {
	//打开待写入文件
	includeFile, err := os.Open(includeFileName)
	if err != nil {
		return err
	}
	defer includeFile.Close()

	//获取文件描述
	includeFileInfo, err := includeFile.Stat()
	if err != nil {
		return err
	}

	//zip.FileInfoHeader
	zipFileHeader, err := zip.FileInfoHeader(includeFileInfo)
	if err != nil {
		return err
	}

	//修改文件描述的Header，截断路径，只保留文件名
	//否则，解压的时候，可能按压缩进来时候的文件路径来解压，解压回到原来的位置，而不是当前目录
	zipFileHeader.Name = path.Base(includeFileName)

	//用zip.FileInfoHeader，创建zip包内的一个项，并获得io.Writer，准备写入文件
	zipFileWriter, err := zipper.CreateHeader(zipFileHeader)
	if err != nil {
		return err
	}

	//写入本次的文件
	_, err = io.Copy(zipFileWriter, includeFile)
	return err
}
