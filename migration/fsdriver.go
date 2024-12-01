package migration

import (
	"slices"

	"github.com/gomig/utils"
)

// FS migration file system
type FS []File

// Len get files length
func (files FS) Len() int {
	return len(files)
}

// Swap swap item i and j
func (files FS) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

// Less check if name timestamp is smaller
func (files FS) Less(i, j int) bool {
	return files[i].Timestamp() < files[j].Timestamp()
}

// Copy create a new fresh copy from file system
func (files FS) Copy() FS {
	if files.Len() > 1 {
		result := make([]File, len(files))
		copy(result, files)
		return result
	} else {
		return files
	}
}

// Reverse reverse array order
func (files FS) Reverse() FS {
	result := files.Copy()
	slices.Reverse(result)
	return result
}

// Filter filter files by name
func (files FS) Filter(names ...string) FS {
	if len(names) > 0 {
		result := make(FS, 0)
		for _, name := range names {
			for _, file := range files {
				if file.Is(utils.Slugify(name)) {
					result = append(result, file)
				}
			}
		}
		return result
	}
	return files
}

// FilterMigrated exclude files from file system
func (files FS) FilterMigrated(names ...string) FS {
	result := make(FS, 0)
	for _, file := range files {
		if file.isMigrated(names...) {
			result = append(result, file)
		}
	}
	return result
}

// ExcludeMigrated exclude files from file system
func (files FS) ExcludeMigrated(names ...string) FS {
	if len(names) > 0 {
		result := make(FS, 0)
		for _, file := range files {
			if !file.isMigrated(names...) {
				result = append(result, file)
			}
		}
		return result
	}
	return files
}
