package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-auth-otp-service/src/database/scopes"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func QueryParametersBuilderMiddleware(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		builder := &scopes.BuilderModel{}

		page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(scopes.DefaultPage)))
		if err != nil {
			page = scopes.DefaultPage
		}
		builder.Page = uint(page)

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(scopes.DefaultPageSize)))
		if err != nil {
			pageSize = scopes.DefaultPageSize
		}
		builder.PageSize = uint(pageSize)

		sortOrder := c.Query("sort_order")
		if sortOrder != "desc" && sortOrder != "asc" {
			sortOrder = "desc"
		}
		builder.SortOrder = sortOrder

		if createdAfterStr := c.Query("created_after"); createdAfterStr != "" {
			createdAfter, err := time.Parse(time.RFC3339, createdAfterStr)
			if err == nil {
				builder.CreatedAfter = &createdAfter
			}
		}

		if createdBeforeStr := c.Query("created_before"); createdBeforeStr != "" {
			createdBefore, err := time.Parse(time.RFC3339, createdBeforeStr)
			if err == nil {
				builder.CreatedBefore = &createdBefore
			}
		}

		filters := make(map[string]interface{})
		likes := make(map[string]interface{})
		modelType := reflect.TypeOf(model)

		sortBy := c.DefaultQuery("sort_by", "created_at")
		filtersFromQuery := c.QueryMap("filters")
		likesFromQuery := c.QueryMap("likes")

		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			if jsonTag, ok := field.Tag.Lookup("json"); ok && jsonTag != "" {
				jsonTag = strings.Split(jsonTag, ",")[0]
				if filterable, ok := field.Tag.Lookup("filter"); ok && filterable == "true" {
					if value, exists := filtersFromQuery[jsonTag]; exists {
						filters[jsonTag] = value
					}
				}

				if likeable, ok := field.Tag.Lookup("like"); ok && likeable == "true" {
					if value, exists := likesFromQuery[jsonTag]; exists {
						likes[jsonTag] = value
					}
				}

				if sortable, ok := field.Tag.Lookup("sort"); ok && sortable == "true" {
					if sortBy == jsonTag {
						builder.SortBy = sortBy
					}
				}
			}
		}
		builder.Filters = filters
		builder.Likes = likes

		c.Set("query_parameters_builder", builder)
		c.Next()
	}
}
