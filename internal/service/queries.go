package service

import (
	"github.com/vrld/ansicht/internal/db"
	"github.com/vrld/ansicht/internal/model"
)

type Queries struct {
	queries       []model.SearchQuery
	selectedIndex int
}

func NewQueries() *Queries {
	queries, err := db.GetSavedQueries()
	if err != nil || len(queries) == 0 {
		queries = []model.SearchQuery{
			{Name: "INBOX", Query: "query:INBOX"},
		}
	}

	return &Queries{
		queries:       queries,
		selectedIndex: 0,
	}
}

func (q *Queries) All() []model.SearchQuery {
	return q.queries
}

func (q *Queries) Current() (model.SearchQuery, bool) {
	if q.selectedIndex < len(q.queries) {
		return q.queries[q.selectedIndex], true
	}
	return model.SearchQuery{}, false
}

func (q *Queries) SelectedIndex() int {
	return q.selectedIndex
}

func (q *Queries) SelectNext() (ok bool) {
	return q.Select((q.selectedIndex + 1) % len(q.queries))
}

func (q *Queries) SelectPrevious() (ok bool) {
	return q.Select((q.selectedIndex - 1 + len(q.queries)) % len(q.queries))
}

func (q *Queries) SelectFirst() (ok bool) {
	return q.Select(0)
}

func (q *Queries) SelectLast() (ok bool) {
	return q.Select(len(q.queries) - 1)
}

func (q *Queries) Select(i int) (ok bool) {
	if i < 0 || i >= len(q.queries) {
		return false
	}
	q.selectedIndex = i
	return true
}

func (q *Queries) Add(query model.SearchQuery) {
	q.queries = append(q.queries, query)
}

func (q *Queries) AddQuery(name string, query string) {
	q.Add(model.SearchQuery{Name: name, Query: query})
}
