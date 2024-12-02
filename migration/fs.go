package migration

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomig/utils"
)

// NewMigrationFile generate new migration file
func NewMigrationFile(root, name, ext string, stages ...string) error {
	base := root
	if root != "" {
		base := normalizePath(root)
		if err := makeDir(base); err != nil {
			return err
		}
	}

	if name == "" || ext == "" {
		return errors.New("name and extension parameters are required")
	} else {
		name = fmt.Sprintf(
			"%s-%s.%s",
			strconv.FormatInt(time.Now().Unix(), 10),
			utils.Slugify(name), ext,
		)
	}

	if len(stages) == 0 {
		stages = []string{"main"}
	}

	// generate template
	content := make([]string, 0)
	for _, stage := range stages {
		content = append(content, fmt.Sprintf("-- { up: %s }\n\n", stage))
		content = append(content, fmt.Sprintf("-- { down: %s }\n\n", stage))
	}

	// write file
	if err := os.WriteFile(
		normalizePath(base, name),
		[]byte(strings.Join(content, "")),
		0644,
	); err != nil {
		return err
	} else {
		return nil
	}
}

// NewDirFS generate new migration filesystem from directory
func NewDirFS(root, ext string) (FS, error) {
	result := make(FS, 0)
	files := os.DirFS(root)
	err := fs.WalkDir(files, ".", func(p string, entry fs.DirEntry, err error) error {
		pth := path.Join(root, p)
		if ok, err := regexp.MatchString(`^([0-9])(.+)(\.`+ext+`)$`, entry.Name()); err != nil {
			return err
		} else if ok && !entry.IsDir() {
			if content, err := os.ReadFile(pth); err != nil {
				return err
			} else {
				result = append(result, File{name: entry.Name(), ext: ext, content: string(content)})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else {
		sort.Sort(result)
		return result, nil
	}
}

// NewEmbedFS generate new migration filesystem from embedded files
func NewEmbedFS(f embed.FS, root, ext string) (FS, error) {
	result := make(FS, 0)
	err := fs.WalkDir(f, root, func(pth string, entry fs.DirEntry, err error) error {
		if ok, err := regexp.MatchString(`^([0-9])(.+)(\.`+ext+`)$`, entry.Name()); err != nil {
			return err
		} else if ok && !entry.IsDir() {
			if content, err := os.ReadFile(pth); err != nil {
				return err
			} else {
				result = append(result, File{name: entry.Name(), ext: ext, content: string(content)})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else {
		sort.Sort(result)
		return result, nil
	}
}
