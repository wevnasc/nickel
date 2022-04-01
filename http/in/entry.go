package in

import "nickel/core/domain"

type Entry struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Amount      float64  `json:"amount"`
	Tags        []string `json:"tags"`
	Type        string   `json:"type"`
}

func (e *Entry) Domain() *domain.Entry {

	tags := make([]domain.Tag, len(e.Tags))

	for idx, t := range e.Tags {
		tags[idx] = domain.Tag{
			Name: t,
		}
	}

	return &domain.Entry{
		ID:          e.ID,
		Description: e.Description,
		Amount:      e.Amount,
		Tags:        tags,
		Type:        domain.Type(e.Type),
	}
}
