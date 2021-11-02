package search

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/le4-db-go/db"
	"github.com/hiroyaonoe/le4-db-go/pkg/comment"
	"github.com/hiroyaonoe/le4-db-go/pkg/thread"
)

func Get(c *gin.Context) {
	db, err := db.NewDB()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	searchQuery := c.Query("query")
	words := strings.Fields(searchQuery)
	wordsI := make([]interface{}, len(words))
	for i, v := range words {
		wordsI[i] = "%"+v+"%"
	}

	query := "SELECT thread_id, title, created_at, user_id, users.name AS user_name " + 
		"FROM threads " + 
		"NATURAL JOIN post_threads " + 
		"NATURAL JOIN users"
	if len(words) > 0 {
		likeQuery := make([]string, len(words))
		for i := 0; i < len(words); i++ {
			likeQuery[i] = fmt.Sprintf("title LIKE $%d", i+1)
		}
		query = query + " WHERE " + strings.Join(likeQuery, " OR ")
	}

	threads := []thread.Thread{}
	err = db.Select(&threads, query, wordsI...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	query = "SELECT comments.content, comments.comment_id, comments.thread_id, threads.title AS thread_title, post_comments.created_at, post_comments.user_id, users.name AS user_name " + 
		"FROM post_comments " + 
		"JOIN comments ON post_comments.thread_id = comments.thread_id AND post_comments.comment_id = comments.comment_id " + 
		"JOIN threads ON comments.thread_id = threads.thread_id " + 
		"JOIN users ON post_comments.user_id = users.user_id"
	if len(words) > 0 {
		likeQuery := make([]string, len(words))
		for i := 0; i < len(words); i++ {
			likeQuery[i] = fmt.Sprintf("comments.content LIKE $%d", i+1)
		}
		query = query + " WHERE " + strings.Join(likeQuery, " OR ")
	}

	comments := []comment.Comment{}
	err = db.Select(&comments, query, wordsI...)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "search.html", gin.H{
		"threads": threads,
		"comments": comments,
		"query": searchQuery,
	})
}
