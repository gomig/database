package migration

type Migrated struct {
	Name  string `db:"name"`
	Stage string `db:"stage"`
}

type Summary []Migrated

// IsEmpty check if summary is empty
func (summary Summary) IsEmpty() bool {
	return len(summary) == 0
}

// Names get migrated files
func (summary Summary) Names() []string {
	result := make([]string, 0)
	for _, migration := range summary {
		result = append(result, migration.Name)
	}
	return result
}

// GroupByStage group migrated files by stage
func (summary Summary) GroupByStage() map[string][]string {
	result := make(map[string][]string)
	for _, file := range summary {
		result[file.Stage] = append(result[file.Stage], file.Name)
	}
	return result
}

// GroupByFile group migrated files by file name
func (summary Summary) GroupByFile() map[string][]string {
	result := make(map[string][]string)
	for _, file := range summary {
		result[file.Name] = append(result[file.Name], file.Stage)
	}
	return result
}
