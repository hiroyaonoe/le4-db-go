package search

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/domain"
	"github.com/hiroyaonoe/le4-db-go/lib/builder"
	"github.com/hiroyaonoe/le4-db-go/pkg/auth"
	"github.com/jmoiron/sqlx"
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

	searchTag := c.Query("tag")
	tagID, err := strconv.Atoi(searchTag)
	if err != nil {
		tagID = -1 // all
	}

	searchStartDateStr := c.Query("start_date")
	searchEndDateStr := c.Query("end_date")
	jst, _ := time.LoadLocation("Asia/Tokyo")
	layout := "2006-01-02"

	query := "SELECT thread_id, title, created_at, user_id, user_name, category_id, category_name FROM threads_with_user_category"

	args := make([]interface{}, len(words), len(words)+4)
	for i, v := range words {
		args[i] = "%" + v + "%"
	}
	var categoryBuilder builder.Builder
	if categoryID >= 0 {
		args = append(args, categoryID)
		categoryBuilder = builder.Word("category_id = ?")
	} else { // searchCategory == -1(all) の場合
		categoryBuilder = builder.Null()
	}

	var tagBuilder builder.Builder
	if tagID >= 0 {
		query += " NATURAL JOIN add_tags"
		args = append(args, tagID)
		tagBuilder = builder.Word("tag_id = ?")
	} else { // searchTag == -1(all) の場合
		tagBuilder = builder.Null()
	}

	var startDateBuilder builder.Builder
	if searchStartDate, err := time.ParseInLocation(layout, searchStartDateStr, jst); err == nil {
		args = append(args, searchStartDate)
		startDateBuilder = builder.Word("created_at >= ?")
	} else {
		startDateBuilder = builder.Null()
	}

	var endDateBuilder builder.Builder
	if searchEndDate, err := time.ParseInLocation(layout, searchEndDateStr, jst); err == nil {
		args = append(args, searchEndDate)
		endDateBuilder = builder.Word("created_at <= ?")
	} else {
		endDateBuilder = builder.Null()
	}

	likeBuilders := make([]builder.Builder, len(words))
	for i := 0; i < len(words); i++ {
		likeBuilders[i] = builder.Word("title LIKE ?")
	}
	queryBuilder := builder.Or(likeBuilders...)
	queryBuilder = builder.And(queryBuilder, categoryBuilder, tagBuilder, startDateBuilder, endDateBuilder)
	queryBuilder = builder.Where(builder.Word(query), queryBuilder)
	query = queryBuilder.Build()
	query = db.Rebind(query)
	log.Println(query)
	log.Println(args)

	threads := []domain.Thread{}
	err = db.Select(&threads, query, args...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	AddTags := []domain.Tag{}
	query = "SELECT tag_id, name, thread_id FROM tag_with_thread_id WHERE thread_id IN (:thread_id)"
	query, argsT, err := sqlx.Named(query, threads)
	query, argsT, err = sqlx.In(query, argsT)
	query = db.Rebind(query)
	err = db.Select(&AddTags, query, argsT...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	numComments := []domain.Thread{}
	query = "SELECT thread_id, num_comment " +
		"FROM num_comments " +
		"WHERE thread_id IN (:thread_id)"
	query, argsT, err = sqlx.Named(query, threads)
	query, argsT, err = sqlx.In(query, argsT)
	query = db.Rebind(query)
	err = db.Select(&numComments, query, argsT...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	indexThread := map[int]*domain.Thread{}
	for i := 0; i < len(threads); i++ {
		indexThread[threads[i].ThreadID] = &threads[i]
	}
	for _, v := range AddTags {
		t := indexThread[v.ThreadID]
		t.Tags = append(t.Tags, v)
	}
	for _, v := range numComments {
		t := indexThread[v.ThreadID]
		t.NumComment = v.NumComment
	}

	query = "SELECT content, comment_id, thread_id, thread_title, created_at, user_id, user_name " +
		"FROM comments_with_user_thread " +
		"NATURAL JOIN link_categories " +
		"NATURAL JOIN categories"

	if tagID >= 0 {
		query += " NATURAL JOIN add_tags"
	}
	for i := 0; i < len(words); i++ {
		likeBuilders[i] = builder.Word("content LIKE ?")
	}
	queryBuilder = builder.Or(likeBuilders...)
	queryBuilder = builder.And(queryBuilder, categoryBuilder, tagBuilder, startDateBuilder, endDateBuilder)
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

	tags := []domain.Tag{}
	err = db.Select(&tags, "SELECT tag_id, name FROM tags")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	loginUserID := auth.GetUserID(c)
	loginUserName := auth.GetUserName(c)

	c.HTML(http.StatusOK, "search.html", gin.H{
		"threads":         threads,
		"comments":        comments,
		"categoryID":      categoryID,
		"categories":      categories,
		"query":           searchQuery,
		"start_date":      searchStartDateStr,
		"end_date":        searchEndDateStr,
		"tagID":           tagID,
		"tags":            tags,
		"login_user_id":   loginUserID,
		"login_user_name": loginUserName,
	})
}
