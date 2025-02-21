package views

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"vpn-tg-bot/pkg/e"
	"vpn-tg-bot/web/admin_panel/entity"
)

var BasePath string
var PublicPath string
var PublicAssetsMeta map[string]entity.FileMeta

var executableTimestampS int64 = 0

func Versioned(path string) string {
	// // This variable need to handle incompatible "/" in path.
	// var filePath string
	// if strings.HasPrefix(path, "/") {
	// 	filePath = path[1:]
	// } else {
	// 	filePath = path
	// }

	// var s int64

	// meta, ok := PublicAssetsMeta[filePath]
	// if ok {
	// 	s = meta.ModTime.Unix()
	// } else {
	// 	log.Printf("[WARNING] FileMeta is not specified for path \"%s\"", filePath)
	// 	s = time.Now().Unix()
	// }

	// if s <= 0 {
	// 	log.Printf("[WARNING] ModTime is negative for path \"%s\", value %d", filePath, s)
	// 	s = time.Now().Unix()
	// }

	// fmt.Printf("PublicAssetsMeta:\n%+v\n", PublicAssetsMeta)

	var s int64
	if executableTimestampS == 0 {
		t, err := getExecutableTimestamp()
		if err == nil {
			s = t.Unix()
			executableTimestampS = s
			log.Printf("[INFO] admin_panel/views: Got executable timestamp: %v, memorized as %d", t, s)

		} else {
			log.Printf("[WARNING] admin_panel/views: Can't get executable timestamp: %v. For versioning will use current time.", err)
			s = time.Now().Unix()
		}
	} else {
		s = executableTimestampS
	}

	return path + "?v=" + strconv.FormatInt(s, 10)
}

func getExecutableTimestamp() (t time.Time, err error) {
	execPath, err := os.Executable()
	if err != nil {
		return time.Time{}, e.Wrap("Can't get executable path: %v", err)
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return time.Time{}, e.Wrap("Can't evaluate symlinks for executable path: %v", err)
	}

	info, err := os.Stat(execPath)
	if err != nil {
		return time.Time{}, e.Wrap("Can't get info for executable path: %v", err)
	}

	return info.ModTime(), nil
}
