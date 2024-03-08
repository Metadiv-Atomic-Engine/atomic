package base

import (
	"github.com/Metadiv-Atomic-Engine/sql"
)

type RequestListing struct {
	sql.Pagination
	sql.Sorting
	Keyword string `form:"keyword"`
}

func (r *RequestListing) BuildSimilarClause(fields ...string) *sql.Clause {
	var clauses = make([]*sql.Clause, 0)
	for _, field := range fields {
		clauses = append(clauses, sql.Similar(field, r.Keyword))
	}
	if len(clauses) == 0 {
		return nil
	}
	return sql.Or(clauses...)
}
