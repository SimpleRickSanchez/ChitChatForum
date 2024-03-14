package model

import (
	"fmt"
	"math/rand"
	"time"
	"util"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Comment struct {
	Id            int       `json:"id"`
	Uuid          uuid.UUID `json:"uuid"`
	Content       string    `json:"content"`
	UserId        int       `json:"user_id"`
	ReplyToUserId int       `json:"reply_to_user_id"`
	ReplyToUuid   uuid.UUID `json:"reply_to_uuid"`
	PostUuid      uuid.UUID `json:"post_uuid"`
	ThreadUuid    uuid.UUID `json:"thread_uuid"`
	CreatedAt     time.Time `json:"created_at"`
}

func (comment *Comment) CreatedAtDate() string {
	return comment.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func GetPageCommentByPostUuid(page int, postUuid string) (comments []Comment) {
	db.Raw("SELECT * FROM comments WHERE post_uuid = UUID_TO_BIN(?) ORDER BY created_at DESC LIMIT ?,?",
		postUuid,
		CommentsPerPage*(page-1),
		CommentsPerPage).Scan(&comments)
	return
}
func GetPageCommentByUserUuid(page int, useruuid string) (comments []Comment) {
	if !IsUserExistsByUuid(useruuid) {
		return
	}
	db.Raw("SELECT * FROM comments WHERE user_id = ? ORDER BY created_at DESC LIMIT ?,?",
		GetUserIdByUuid(useruuid),
		CommentsPerPage*(page-1),
		CommentsPerPage).Scan(&comments)
	return
}
func GetAllCommentsToAPost(postUuid string) (comments []Comment) {
	db.Raw("SELECT * FROM comments WHERE post_uuid=UUID_TO_BIN(?)", postUuid).Scan(&comments)
	return
}
func GetAllReplierIdsToAPost(postUuid string) (ids []int) {
	db.Raw("SELECT user_id FROM comments WHERE post_uuid=UUID_TO_BIN(?)", postUuid).Scan(&ids)
	return
}

func IsCommentExists(commentid int) bool {
	var count int
	db.Raw("SELECT COUNT(*) FROM comments WHERE id=?", commentid).Scan(&count)
	return count != 0
}
func GetCommentsNum() int {
	var count int
	db.Raw("SELECT COUNT(*) FROM comments").Scan(&count)
	return count
}
func RandomCommentsId() int {
	return rand.Intn(GetCommentsNum()) + 1
}
func GetCommentByUuid(uuid string) (comments []Comment) {
	if !util.IsValidUUIDSTR(uuid) {
		return
	}
	db.Raw("SELECT * FROM comments WHERE uuid=UUID_TO_BIN(?)", uuid).Scan(&comments)
	return
}
func GetPostOrCommentByUuid(uuid string) interface{} {
	if !util.IsValidUUIDSTR(uuid) {
		fmt.Println(uuid)
		return nil
	}
	post := GetPostByUuid(uuid)
	if len(post) != 0 {
		return post[0]
	}
	comment := GetCommentByUuid(uuid)
	if len(comment) != 0 {
		return comment[0]
	}
	return nil
}
