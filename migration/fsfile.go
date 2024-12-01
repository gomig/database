package migration

import (
	"regexp"
	"slices"
	"strconv"

	"github.com/gomig/utils"
)

// File migration file
type File struct {
	name    string
	content string
	ext     string
}

// Name get file full name
func (file File) Name() string {
	return file.name
}

// Content get file content
func (file File) Content() string {
	return file.content
}

// Get timestamp part of filename
func (file File) Timestamp() int64 {
	if res, err := strconv.ParseInt(utils.ExtractNumbers(file.name), 10, 64); err == nil {
		return res
	} else {
		return 0
	}
}

// Extension get file extension
func (file File) Extension() string {
	return file.ext
}

// HumanizedName get file readable name withut timestamp and extension
func (file File) HumanizedName() string {
	return regexp.
		MustCompile(`^(\d+-)|(\.`+file.ext+`)$`).
		ReplaceAllString(file.name, "")
}

// Is Compare file name in humanize format
func (file File) Is(name string) bool {
	return file.HumanizedName() == regexp.
		MustCompile(`^(\d+-)|(\.`+file.ext+`)$`).
		ReplaceAllString(name, "")
}

// UpScripts get up section scripts
func (file File) UpScripts(stage string) ([]string, error) {
	if scripts, err := upScripts(file.content, stage); err != nil {
		return nil, err
	} else {
		return scripts, nil
	}
}

// DownScripts get down section scripts
func (file File) DownScripts(stage string) ([]string, error) {
	if scripts, err := downScripts(file.content, stage); err != nil {
		return nil, err
	} else {
		return scripts, nil
	}
}

// isMigrated check if file name is in migrated names and should skipped
func (file File) isMigrated(migrated ...string) bool {
	return slices.Contains(migrated, file.name)
}
