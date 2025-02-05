package debug

import (
	"fmt"
	"io/fs"
)

func DisplayFileSystem(filesystem fs.FS) error {
	err := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
	return err
}
