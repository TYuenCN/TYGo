package YSPathUtil

import (
	"os"
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
