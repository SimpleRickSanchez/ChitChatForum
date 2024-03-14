package model

import (
	"math/rand"
	"time"
)

type Topic struct {
	Id        int       `json:"id"`
	Cate      string    `json:"cate"`
	Info      string    `json:"info"`
	CreatedAt time.Time `json:"created_at"`
}

func GetTopicByID(topicId int) (topic Topic) {
	db.Raw("SELECT * FROM topics WHERE id = ?", topicId).Scan(&topic)
	return
}

func GetTopicCateByID(topicId int) (cate string) {
	db.Raw("SELECT cate FROM topics WHERE id = ?", topicId).Scan(&cate)
	return
}

func IsTopicExistsByID(topicId int) bool {
	var count int
	db.Raw("SELECT COUNT(*) FROM topics WHERE id = ?", topicId).Scan(&count)
	return count != 0
}

func GetAllTopics() (topics []Topic) {
	db.Raw("SELECT * FROM topics").Scan(&topics)
	return
}

func GetAllTopicsCate() (cate []string) {
	db.Raw("SELECT cate FROM topics").Scan(&cate)
	return
}

func GetTopicsNum() int {
	var count int
	db.Raw("SELECT COUNT(*) FROM topics").Scan(&count)
	return count
}

func RandomTopicId() int {
	return rand.Intn(GetTopicsNum()) + 1
}
