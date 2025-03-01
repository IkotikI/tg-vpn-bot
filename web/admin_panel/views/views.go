package views

import (
	"fmt"
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
var PublicDestPath string
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

func PaginationLinks(base string, p entity.Pagination, r int) *[]entity.PaginationLink {
	// If 1 page, don't make mess
	if p.TotalPages <= 1 {
		return &[]entity.PaginationLink{{Link: base, Num: 1}}
	}

	// If 2 pages, it's easy too.
	if p.TotalPages == 2 {
		return &[]entity.PaginationLink{
			{Link: base + fmt.Sprintf("?&per_page=%d", p.PerPage), Num: 1},
			{Link: base + fmt.Sprintf("?page=2&per_page=%d", p.PerPage), Num: 2},
		}
	}

	// Least case:   1,2,3...10   - Page: 1, r: 3
	// Link number:  0   r   1
	// Average capacity: 1...3,4,5,6,7...10
	// Link number:      1  -r   0  +r   1
	start := make([]entity.PaginationLink, 0, r+1)
	end := make([]entity.PaginationLink, 0, r+1)

	start = append(start, entity.PaginationLink{
		Link: base + fmt.Sprintf("?per_page=%d", p.PerPage),
		Num:  1,
	})

	n := p.Page

	// Make n-th link
	if n != 1 && n != p.TotalPages {
		end = append(end, entity.PaginationLink{
			Link: base + fmt.Sprintf("?page=%d&per_page=%d", n, p.PerPage),
			Num:  n,
		})
	}

	var l, h int64
	// Making links for left and right sides around n-th page, limited
	// by radius r.
	for i := 1; i < r; i++ {
		l = n - int64(r) + int64(i)
		h = n + int64(i)
		if l > 1 {
			start = append(start, entity.PaginationLink{
				Link: base + fmt.Sprintf("?page=%d&per_page=%d", l, p.PerPage),
				Num:  l,
			})
		}
		if h < p.TotalPages {
			end = append(end, entity.PaginationLink{
				Link: base + fmt.Sprintf("?page=%d&per_page=%d", h, p.PerPage),
				Num:  h,
			})
		}
	}

	end = append(end, entity.PaginationLink{
		Link: base + fmt.Sprintf("?page=%d&per_page=%d", p.TotalPages, p.PerPage),
		Num:  p.TotalPages,
	})

	fmt.Printf("start %+v\n", start)
	fmt.Printf("end %+v\n", end)

	links := append(start, end...)

	return &links
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
