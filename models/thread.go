package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"util"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreadErrorType int
type ThreadOrder int

const (
	ThreadErrorNone         ThreadErrorType = 0
	ThreadErrorUser         ThreadErrorType = 1
	ThreadErrorTopic        ThreadErrorType = 2
	ThreadErrorTitle        ThreadErrorType = 3
	ThreadErrorContent      ThreadErrorType = 4
	ThreadOrderByCreateTime ThreadOrder     = 1
	ThreadOrderByViewCount  ThreadOrder     = 2
	ThreadOrderByNumPosts   ThreadOrder     = 3
)

type Thread struct {
	Id        int       `json:"id"`
	Uuid      uuid.UUID `json:"uuid"`
	TopicId   int       `json:"topic_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserId    int       `json:"user_id"`
	LastPos   int       `json:"last_pos"`
	ViewCount int       `json:"view_count"`
	NumPosts  int       `json:"num_posts"`
	CreatedAt time.Time `json:"created_at"`
}

func GetAllThreadOrder() []ThreadOrder {
	return []ThreadOrder{
		ThreadOrderByCreateTime,
		ThreadOrderByViewCount,
		ThreadOrderByNumPosts}
}

// format the CreatedAt date to display nicely on the screen
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

// verify the number of posts
func (thread *Thread) NumPostsVerify() (count int) {
	db.Raw("SELECT count(*) FROM posts WHERE thread_id = ?", thread.Id).Scan(&count)
	return
}

func isThreadValidToWrite(user User, thread Thread) (t ThreadErrorType, err error) {
	if !IsUserExists(user.Email) {
		return ThreadErrorUser, fmt.Errorf("user %v not exists", user.Email)
	}
	if !IsTopicExistsByID(thread.TopicId) {
		return ThreadErrorTopic, fmt.Errorf("topic not exists, %v", thread.TopicId)
	}
	if len(thread.Title) == 0 {
		return ThreadErrorTitle, fmt.Errorf("title empty")
	}
	if len(thread.Title) > CommonStringMax {
		return ThreadErrorTitle, fmt.Errorf("title too long %v", thread.Title)
	}
	if len(thread.Content) == 0 {
		return ThreadErrorContent, fmt.Errorf("content empty")
	}
	if len(thread.Content) > CommonTextMax {
		return ThreadErrorContent, fmt.Errorf("content too long, %v", thread.Content)
	}
	return ThreadErrorNone, nil
}

func IsThreadExists(threadId int) bool {
	var count int
	db.Raw("SELECT COUNT(*) FROM threads WHERE id = ?", threadId).Scan(&count)
	return count != 0
}
func GetThreadsNum() int {
	var count int
	db.Raw("SELECT COUNT(*) FROM threads").Scan(&count)
	return count
}
func RandomThreadId() int {
	return rand.Intn(GetThreadsNum()) + 1
}

func GetThreadById(threadId int) (thread Thread) {
	db.Raw("SELECT * FROM threads WHERE id = ?", threadId).Scan(&thread)
	return
}

func GetThreadByUuid(uuid string) (threads []Thread) {
	if !util.IsValidUUIDSTR(uuid) {
		return
	}
	db.Raw("SELECT * FROM threads WHERE uuid = UUID_TO_BIN(?)", uuid).Scan(&threads)
	return
}
func GetThreadUuidById(threadid int) (uuid uuid.UUID) {
	uuidStruct := UuidStruct{}
	db.Raw("SELECT uuid FROM threads WHERE id = ?", threadid).Scan(&uuidStruct)
	return uuidStruct.Uuid
}
func GetThreadIdByUuid(uuid string) (id int) {
	db.Raw("SELECT id FROM threads WHERE uuid = UUID_TO_BIN(?)", uuid).Scan(&id)
	return
}
func GetThreadTitleById(threadid int) (title string) {
	db.Raw("SELECT title FROM threads WHERE id = ?", threadid).Scan(&title)
	return
}

func GetPageThreads(page int, order ThreadOrder, userUuid string, curTopicID int) (threads []Thread) {
	var sqlStr = []string{
		"SELECT uuid, topic_id, title, content, user_id, view_count, num_posts FROM threads",
		"WHERE",
		"ORDER BY",
		"LIMIT ?,?"}

	switch order {
	case ThreadOrderByCreateTime:
		sqlStr[2] = fmt.Sprint(sqlStr[2], " created_at DESC")
	case ThreadOrderByViewCount:
		sqlStr[2] = fmt.Sprint(sqlStr[2], " view_count DESC")
	case ThreadOrderByNumPosts:
		sqlStr[2] = fmt.Sprint(sqlStr[2], " num_posts DESC")
	}

	if IsUserExistsByUuid(userUuid) { // private
		if IsTopicExistsByID(curTopicID) {
			sqlStr[1] = fmt.Sprint(sqlStr[1], " user_id = ", GetUserIdByUuid(userUuid), " AND topic_id = ", curTopicID)
		} else {
			sqlStr[1] = fmt.Sprint(sqlStr[1], " user_id = ", GetUserIdByUuid(userUuid))
		}
	} else { //public
		if IsTopicExistsByID(curTopicID) {
			sqlStr[1] = fmt.Sprint(sqlStr[1], " topic_id = ", curTopicID)
		} else {
			sqlStr[1] = ""
		}
	}

	db.Raw(strings.Join(sqlStr, " "), ThreadsPerPage*(page-1), ThreadsPerPage).Scan(&threads)
	return threads
}

func ThreadViewCountIncre(threadid int) error {
	err := db.Transaction(func(db *gorm.DB) error {
		var viewcount int
		db.Raw("SELECT view_count FROM threads WHERE id=? FOR UPDATE", threadid).Scan(&viewcount)
		db.Exec("UPDATE threads SET view_count=? WHERE id=?", viewcount+1, threadid)
		return nil
	}, Txop)
	return err
}

// TODO
// update a thread
// delete a thread
