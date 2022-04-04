package domain

type Type string

const (
	Income  Type = "INCOME"
	Expense Type = "EXPENSE"
)

type Tag struct {
	ID   string
	Name string
}

type Entry struct {
	ID          string
	Description string
	Amount      float64
	Tags        []Tag
	Type        Type
}

func Names(tags []Tag) []string {
	names := make([]string, len(tags))

	for _, t := range tags {
		names = append(names, t.Name)
	}

	return names
}

func Contains(tags []Tag, name string) bool {
	for _, v := range tags {
		if v.Name == name {
			return true
		}
	}
	return false
}

//Return all tags not found in the current tags
func MissingTags(currentTags []Tag, tags []Tag) []Tag {
	var missing []Tag

	for _, t := range tags {
		if !Contains(currentTags, t.Name) {
			missing = append(missing, t)
		}
	}

	return missing
}
