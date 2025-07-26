package service

import (
	"github.com/vrld/ansicht/internal/db"
	"github.com/vrld/ansicht/internal/model"
)

type queries struct {
	queries       []model.SearchQuery
	selectedIndex int
}

var queriesInstance *queries

func Queries() *queries {
	if queriesInstance != nil {
		return queriesInstance
	}

	savedQueries, err := db.GetSavedQueries()
	if err != nil || len(savedQueries) == 0 {
		Logger().Warning(err.Error())
		savedQueries = []model.SearchQuery{
			{Name: "INBOX", Query: "query:INBOX"},
		}
	}

	queriesInstance = &queries{
		queries:       savedQueries,
		selectedIndex: 0,
	}
	return queriesInstance
}

func (q *queries) All() []model.SearchQuery {
	return q.queries
}

func (q *queries) Current() (model.SearchQuery, bool) {
	if q.selectedIndex < len(q.queries) {
		return q.queries[q.selectedIndex], true
	}
	return model.SearchQuery{}, false
}

func (q *queries) SelectedIndex() int {
	return q.selectedIndex
}

func (q *queries) SelectNext() (ok bool) {
	return q.Select((q.selectedIndex + 1) % len(q.queries))
}

func (q *queries) SelectPrevious() (ok bool) {
	return q.Select((q.selectedIndex - 1 + len(q.queries)) % len(q.queries))
}

func (q *queries) SelectFirst() (ok bool) {
	return q.Select(0)
}

func (q *queries) SelectLast() (ok bool) {
	return q.Select(len(q.queries) - 1)
}

func (q *queries) Select(i int) (ok bool) {
	if i < 0 || i >= len(q.queries) {
		return false
	}
	q.selectedIndex = i
	return true
}

func (q *queries) Add(query model.SearchQuery) {
	q.queries = append(q.queries, query)
}

func (q *queries) AddQuery(name string, query string) {
	q.Add(model.SearchQuery{Name: name, Query: query})
}
