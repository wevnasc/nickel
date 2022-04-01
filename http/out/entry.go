package out

import "nickel/core/domain"

type Entry struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Amount      float64  `json:"amount"`
	Tags        []string `json:"tags"`
	Type        string   `json:"type"`
}

func TagListFrom(tags []domain.Tag) []string {
	out := make([]string, len(tags))

	for idx, t := range tags {
		out[idx] = t.Name
	}

	return out
}

func EntryFrom(entry *domain.Entry) *Entry {
	return &Entry{
		ID:          entry.ID,
		Description: entry.Description,
		Amount:      entry.Amount,
		Tags:        TagListFrom(entry.Tags),
		Type:        string(entry.Type),
	}
}

func EntriesListFrom(entries []domain.Entry) []Entry {
	out := make([]Entry, len(entries))

	for idx, e := range entries {
		out[idx] = *EntryFrom(&e)
	}

	return out
}
