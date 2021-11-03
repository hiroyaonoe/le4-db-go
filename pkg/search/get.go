package search

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/pkg/builder"
)

func Get(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	searchQuery := c.Query("query")
	words := strings.Fields(searchQuery)

	searchCategory := c.Query("category")
	categoryID, err := strconv.Atoi(searchCategory)
	if err != nil {
		categoryID = -1 // all
	}

	query := "SELECT thread_id, title, created_at, user_id, users.name AS user_name, categories.category_id, categories.name AS category_name " +
		"FROM threads " +
		"NATURAL JOIN post_threads " +
		"NATURAL JOIN users " +
		"NATURAL JOIN link_categories " +
		"JOIN categories ON categories.category_id = link_categories.category_id "

	var categoryBuilder builder.Builder
	args := make([]interface{}, len(words), len(words)+1)
	for i, v := range words {
		args[i] = "%" + v + "%"
	}
	if categoryID >= 0 {
		args = append(args, categoryID)
		categoryBuilder = builder.Word("categories.category_id = ?")
	} else { // searchCategory == -1(all) の場合
		categoryBuilder = builder.Null()
	}

	likeBuilders := make([]builder.Builder, len(words))
	for i := 0; i < len(words); i++ {
		likeBuilders[i] = builder.Word("title LIKE ?")
	}
	queryBuilder := builder.Or(likeBuilders...)
	queryBuilder = builder.And(queryBuilder, categoryBuilder)
	queryBuilder = builder.Where(builder.Word(query), queryBuilder)
	query = queryBuilder.Build()
	query = db.Rebind(query)

	threads := []domain.Thread{}
	err = db.Select(&threads, query, args...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	query = "SELECT comments.content, comments.comment_id, comments.thread_id, threads.title AS thread_title, post_comments.created_at, post_comments.user_id, users.name AS user_name " +
		"FROM post_comments " +
		"JOIN comments ON post_comments.thread_id = comments.thread_id AND post_comments.comment_id = comments.comment_id " +
		"JOIN threads ON comments.thread_id = threads.thread_id " +
		"JOIN users ON post_comments.user_id = users.user_id " +
		"JOIN link_categories ON link_categories.thread_id = threads.thread_id " +
		"JOIN categories ON categories.category_id = link_categories.category_id "
	
	for i := 0; i < len(words); i++ {
		likeBuilders[i] = builder.Word("comments.content LIKE ?")
	}
	queryBuilder = builder.Or(likeBuilders...)
	queryBuilder = builder.And(queryBuilder, categoryBuilder)
	queryBuilder = builder.Where(builder.Word(query), queryBuilder)
	query = queryBuilder.Build()
	query = db.Rebind(query)

	comments := []domain.Comment{}
	err = db.Select(&comments, query, args...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	categories := []domain.Category{}
	err = db.Select(&categories, "SELECT category_id, name FROM categories")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "search.html", gin.H{
		"threads":    threads,
		"comments":   comments,
		"categoryID": categoryID,
		"categories": categories,
		"query": searchQuery,
	})
}
