package domain

type Tag struct {
	ID   string
	Name string
}

type TagList []Tag

func (tt TagList) Names() []string {
	names := make([]string, len(tt))

	for _, t := range tt {
		names = append(names, t.Name)
	}

	return names
}

func (tt TagList) Contains(name string) bool {
	for _, v := range tt {
		if v.Name == name {
			return true
		}
	}
	return false
}

//Return all tags not found in the current tags
func (tt TagList) Diff(tags TagList) []Tag {
	var missing []Tag

	for _, t := range tt {
		if !tags.Contains(t.Name) {
			missing = append(missing, t)
		}
	}

	return missing
}
