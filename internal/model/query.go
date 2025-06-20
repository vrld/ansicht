package model

type SearchQuery struct {
	Query string
	Name  string
}

type SearchResult struct {
	Query   *SearchQuery
	Threads []Thread
}
