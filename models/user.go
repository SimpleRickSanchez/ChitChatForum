package model

import (
	"database/sql"
	"fmt"
	"time"
	"util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UuidStruct struct {
	Uuid uuid.UUID
}

type User struct {
	Id        int       `json:"id"`
	Uuid      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Pwdmd5    string    `json:"pwdmd5"`
	Salt      string    `json:"salt"`
	CreatedAt time.Time `json:"created_at"`
}

var Txop = &sql.TxOptions{
	Isolation: sql.LevelRepeatableRead,
}

// 用来自动匹配的结构体必须一一对应数据库表中的colname及其类型
//自动匹配数据库表中的colname，首字母自动转换为小写，如果数据库中使用下划线则，article_news对应go文件中ArticleNews

func (user User) CreateUser() error {
	if user.isUserValidToWrite() {
		err := db.Transaction(func(db *gorm.DB) error {
			var count int
			db.Raw("SELECT COUNT(*) FROM users WHERE email = ? ", user.Email).Scan(&count)
			if count == 0 {
				err := db.Exec(`INSERT INTO users (uuid,name,email,pwdmd5,salt) 
					VALUES(UUID_TO_BIN(?),?,?,?,?)`,
					util.CreateUUIDStr(),
					user.Name,
					user.Email,
					user.Pwdmd5,
					user.Salt).Error
				return err
			}
			return fmt.Errorf("user :%#v not exists", user.Email)
		}, Txop)

		return err
	}
	return fmt.Errorf("user :%#v not valid", user.Email)
}

func (user User) Auth() bool {
	users_in_db := GetUserByEmail(user.Email)
	if len(users_in_db) != 0 &&
		user.Pwdmd5 == users_in_db[0].Pwdmd5 &&
		user.Email == users_in_db[0].Email {
		return true
	}
	return false
}

func (user User) isUserValidToWrite() bool {
	return user.Name != "" &&
		len(user.Name) <= CommonStringMax &&
		util.IsValidEmail(user.Email) &&
		len(user.Email) <= CommonStringMax &&
		util.IsValidMD5(user.Pwdmd5) &&
		util.IsValidSalt(user.Salt)
}

func GetUserByEmail(email string) (users []User) {
	db.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&users)
	return
}

func GetUserNameById(userid int) (name string) {
	db.Raw("SELECT name FROM users WHERE id = ?", userid).Scan(&name)
	return
}

func GetUserNameByUuid(uuid string) (name string) {
	if !util.IsValidUUIDSTR(uuid) {
		return
	}
	db.Raw("SELECT name FROM users WHERE uuid = UUID_TO_BIN(?)", uuid).Scan(&name)
	return
}

func IsUserExists(email string) bool {
	var count int
	db.Raw("SELECT COUNT(*) FROM users WHERE email = ? ", email).Scan(&count)
	return count != 0
}

func IsUserExistsById(id int) bool {
	var count int
	db.Raw("SELECT COUNT(*) FROM users WHERE id = ? ", id).Scan(&count)
	return count != 0
}
func IsUserExistsByUuid(uuid string) bool {
	if !util.IsValidUUIDSTR(uuid) {
		return false
	}
	var count int
	db.Raw("SELECT COUNT(*) FROM users WHERE uuid = UUID_TO_BIN(?) ", uuid).Scan(&count)
	return count != 0
}

func IsLogined(c *gin.Context) bool {
	session := GetSession(c)
	strAny := session.Get("auth")
	islogined, _ := strAny.(string)
	return len(islogined) > 5 // TODO
}
func GetUserByUuid(uuid string) (users []User) {
	if !util.IsValidUUIDSTR(uuid) {
		return
	}
	db.Raw("SELECT * FROM users WHERE uuid = UUID_TO_BIN(?)", uuid).Scan(&users)
	return
}
func GetUserIdByUuid(uuid string) (id int) {
	if !util.IsValidUUIDSTR(uuid) {
		return -1
	}
	db.Raw("SELECT id FROM users WHERE uuid = UUID_TO_BIN(?)", uuid).Scan(&id)
	return
}
func GetUesrUuidByEmail(email string) (uuid string) {
	if !util.IsValidEmail(email) {
		return
	}
	uuidStruct := UuidStruct{}
	db.Raw("SELECT uuid FROM users WHERE email = ?", email).Scan(&uuidStruct)
	return util.BINToUUIDStr(uuidStruct.Uuid)
}
func (user User) GetLatestThreadUuid() (uuid string) {
	uuidStruct := UuidStruct{}
	db.Raw("SELECT t.uuid FROM threads t WHERE t.user_id = ? ORDER BY created_at DESC LIMIT 0,1", user.Id).Scan(&uuidStruct)
	return util.BINToUUIDStr(uuidStruct.Uuid)
}

// Create a new thread
func (user User) CreateThread(topicId int, title, content string) (err error) {

	thread := Thread{
		TopicId: topicId,
		Title:   title,
		Content: content,
		UserId:  user.Id,
	}
	t, err := isThreadValidToWrite(user, thread)
	if t != ThreadErrorNone {

		return err
	}

	err = db.Transaction(func(db *gorm.DB) error {
		err := db.Exec(`INSERT INTO threads (uuid, topic_id, title, content, user_id)
			VALUES (UUID_TO_BIN(?),?,?,?,?)`,
			util.CreateUUIDStr(),
			topicId,
			title,
			content,
			user.Id).Error
		return err
	}, Txop)

	return
}

func (user User) CreatePost(threadid int, content string) (err error) {
	if !IsThreadExists(threadid) {
		return fmt.Errorf("thread not exists: %v", threadid)
	}
	if !IsUserExists(user.Email) {
		return fmt.Errorf("user not exists")
	}
	if len(content) > CommonTextMax {
		return fmt.Errorf("content too long")
	}
	if len(content) == 0 {
		return fmt.Errorf("content empty")
	}
	err = db.Transaction(func(db *gorm.DB) error {

		var lastpos, numposts int
		db.Raw("SELECT last_pos FROM threads WHERE id=? FOR UPDATE", threadid).Scan(&lastpos)
		db.Exec("UPDATE threads SET last_pos=? WHERE id=?", lastpos+1, threadid)
		db.Raw("SELECT num_posts FROM threads WHERE id=? FOR UPDATE", threadid).Scan(&numposts)
		db.Exec("UPDATE threads SET num_posts=? WHERE id=?", numposts+1, threadid)
		db.Exec(`INSERT INTO posts (uuid, content, user_id, thread_id, thread_pos)
		VALUES (UUID_TO_BIN(?),?,?,?,?)`,
			util.CreateUUIDStr(),
			content,
			user.Id,
			threadid,
			lastpos+1)
		return nil
	}, Txop)
	return
}

func (user User) CreateComment(postUuid, replyToUuid, threadUuid, content string) (err error) {
	postid := GetPostIdByUuid(postUuid)
	if !IsUserExists(user.Email) {
		return fmt.Errorf("user not exists")
	}
	replyTo := GetUserIdByPostOrCommentUuid(replyToUuid)
	if !IsUserExistsById(replyTo) {
		return fmt.Errorf("the user to be replied not exists, %v %v", replyTo, postid)
	}
	if !util.Contains(GetAllReplierIdsToAPost(postUuid), replyTo) && replyTo != GetPostAuthorId(postid) {
		return fmt.Errorf("the user to be replied not valid")
	}
	if len(content) > CommonTextMax {
		return fmt.Errorf("content too long")
	}
	if len(content) == 0 {
		return fmt.Errorf("content empty")
	}
	err = db.Transaction(func(db *gorm.DB) error {
		db.Exec(`INSERT INTO comments (uuid, content, user_id, reply_to_user_id, reply_to_uuid, post_uuid, thread_uuid)
			VALUES (UUID_TO_BIN(?),?,?,?,UUID_TO_BIN(?),UUID_TO_BIN(?),UUID_TO_BIN(?))`,
			util.CreateUUIDStr(),
			content,
			user.Id,
			replyTo,
			replyToUuid,
			postUuid,
			threadUuid)
		return nil
	}, Txop)
	return
}

// TODO
// update a user info
// update a thread
// update a post
// update a comment
// delete a user
// delete a thread
// delete a post
// delete a comment
