package model

import (
	"math/rand"
	"time"
	"util"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Post struct {
	Id        int       `json:"id"`
	Uuid      uuid.UUID `json:"uuid"`
	Content   string    `json:"content"`
	UserId    int       `json:"user_id"`
	ThreadId  int       `json:"thread_id"`
	ThreadPos int       `json:"thread_pos"`
	CreatedAt time.Time `json:"created_at"`
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func GetPagePostsByThreadId(page int, threadid int, order int) (posts []Post) {
	orderstr := " DESC"
	if order == 1 {
		orderstr = " ASC"
	}
	sqlstr1 := "SELECT * FROM posts WHERE thread_id = ? ORDER BY thread_pos"
	sqlstr2 := " LIMIT ?,?"
	db.Raw(sqlstr1+orderstr+sqlstr2,
		threadid,
		PostsPerPage*(page-1),
		PostsPerPage).Scan(&posts)
	return
}
func GetPagePostsByUserUuid(page int, useruuid string) (posts []Post) {
	if !IsUserExistsByUuid(useruuid) {
		return
	}
	db.Raw("SELECT * FROM posts WHERE user_id = ? ORDER BY created_at DESC LIMIT ?,?",
		GetUserIdByUuid(useruuid),
		PostsPerPage*(page-1),
		PostsPerPage).Scan(&posts)
	return
}
func GetPostById(postId int) (posts []Post) {
	db.Raw("SELECT * FROM posts WHERE id=?", postId).Scan(&posts)
	return
}
func GetPostAuthorId(postId int) (author int) {
	db.Raw("SELECT user_id FROM posts WHERE id=?", postId).Scan(&author)
	return
}
func GetPostByUuid(uuid string) (posts []Post) {
	if !util.IsValidUUIDSTR(uuid) {
		return
	}
	db.Raw("SELECT * FROM posts WHERE uuid=UUID_TO_BIN(?)", uuid).Scan(&posts)
	return
}
func GetPostIdByUuid(uuid string) (id int) {
	if !util.IsValidUUIDSTR(uuid) {
		return
	}
	db.Raw("SELECT id FROM posts WHERE uuid=UUID_TO_BIN(?)", uuid).Scan(&id)
	return
}
func IsPostExists(postId int) bool {
	var count int
	db.Raw("SELECT COUNT(*) FROM posts WHERE id=?", postId).Scan(&count)
	return count != 0
}

func GetPostsNum() int {
	var count int
	db.Raw("SELECT COUNT(*) FROM posts").Scan(&count)
	return count
}
func RandomPostId() int {
	return rand.Intn(GetPostsNum()) + 1
}

// TODO
// Update a post

// Delete a post

// Delete all posts
