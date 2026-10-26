package main

import (
	"net/http"
	"strconv"
)

// paginate: parse ?page=&limit= dan return offset, limit
func paginate(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return (page - 1) * limit, limit
}

func pageMeta(r *http.Request, total int) map[string]interface{} {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	return map[string]interface{}{"page": page, "limit": limit, "total": total, "total_pages": totalPages}
}
