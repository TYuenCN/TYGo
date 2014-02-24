package path

import (
	//"fmt"
	"github.com/nu7hatch/uuid"
	"os"
	"path/filepath"
	"strings"
)

// 判断路径是否存在（文件或目录）
// 		import (
//			"yuensoft.com/path/YSPathUtil"
//		)
func IsPathExist(aPath string) (isExist bool, isDir bool) {
	isExist, isDir = false, false

	fileInfo, err := os.Stat(aPath)

	if err != nil {
		return false, false
	} else {
		isExist = true
	}

	if fileInfo.IsDir() {
		isDir = true
	}

	return isExist, isDir
}

// 创建并返回文件指针，各级目录一齐创建
func CreateFileBesidesAllDir(flPath string) (*os.File, error) {
	basePath := filepath.Dir(flPath)
	err := os.MkdirAll(basePath, 0666)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(flPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	} else {
		return f, nil
	}
}

// 创建UUID为名，指定后缀的文件
//		ext string 后缀名，不含"."
//		basePath string 文件路径名，"/mnt/hhcehua/images"
func CreateFileWithNewUUIDNameUseExtendBesideMkAllDir(ext string, basePath string) (*os.File, string, error) {
	u5, _ := uuid.NewV4()
	var nmExt string
	if strings.HasPrefix(ext, `.`) {
		nmExt = u5.String() + ext
	} else {
		nmExt = u5.String() + `.` + ext
	}

	flPath := basePath + nmExt
	fl, err := CreateFileBesidesAllDir(flPath)
	return fl, flPath, err
}
