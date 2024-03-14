package controller

import (
	"model"
	"net/http"
	"reflect"
	"strconv"
	"util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IndexControllers struct {
	BaseController
}

type ThreadEntry struct {
	Uuid      uuid.UUID `json:"uuid"`
	Topic     string    `json:"topic"`
	Title     string    `json:"title"`
	ViewCount int       `json:"view_count"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	NumPosts  int       `json:"num_posts"`
	CreatedAt string    `json:"created_at"`
}
type PostEntry struct {
	Uuid        uuid.UUID `json:"uuid"`
	UserName    string    `json:"user_name"`
	ThreadTitle string    `json:"thread_title"`
	ThreadTopic string    `json:"thread_topic"`
	ThreadPos   int       `json:"thread_pos"`
	ThreadUuid  uuid.UUID `json:"thread_uuid"`
	Content     string    `json:"content"`
	CreatedAt   string    `json:"created_at"`
}
type PostEntryWithComments struct {
	PostEntry
	Comments []CommentEntry
}

type CommentEntry struct {
	Uuid           uuid.UUID `json:"uuid"`
	ThreadUuid     uuid.UUID `json:"thread_uuid"`
	UserName       string    `json:"user_name"`
	ReplyToName    string    `json:"reply_to_name"`
	Content        string    `json:"content"`
	ReplyToContent string    `json:"reply_to_content"`
	CreatedAt      string    `json:"created_at"`
}

func (con IndexControllers) Index(c *gin.Context) {
	var page = 1
	var err error
	var order = model.ThreadOrderByCreateTime
	var topicid = 0
	var curprivate = 0
	var userUuid string
	var contenttmpl = ""
	var entrylist = []any{}

	p := c.Query("page")

	if p != "" {
		page, err = strconv.Atoi(p)
		if err != nil || page <= 0 {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}

	o := c.Query("order")
	if o != "" {
		oo, err := strconv.Atoi(o)
		if err != nil {
			c.Redirect(http.StatusFound, "/")
			return
		}
		if util.Contains(model.GetAllThreadOrder(), model.ThreadOrder(oo)) {
			order = model.ThreadOrder(oo)
		}
	}

	t := c.Query("topic")
	if t != "" {
		tt, err := strconv.Atoi(t)
		if err != nil {
			c.Redirect(http.StatusFound, "/")
			return
		}
		topicid = tt
	}

	u := c.Query("useruuid")
	if u != "" {
		if !util.IsValidUUIDSTR(u) {
			c.Redirect(http.StatusFound, "/")
			return
		}
		userUuid = u
	}

	pr := c.Query("private")
	if pr != "" {
		prpr, err := strconv.Atoi(pr)
		if err != nil {
			c.Redirect(http.StatusFound, "/")
			return
		}
		curprivate = prpr
	}

	switch curprivate {

	case 2: // posts
		posts := model.GetPagePostsByUserUuid(page, userUuid)
		for _, post := range posts {
			t_thread := model.GetThreadById(post.ThreadId)
			entrylist = append(entrylist, PostEntry{
				Uuid:        post.Uuid,
				UserName:    model.GetUserNameById(post.UserId),
				ThreadTitle: t_thread.Title,
				ThreadTopic: model.GetTopicCateByID(t_thread.TopicId),
				ThreadPos:   post.ThreadPos,
				ThreadUuid:  t_thread.Uuid,
				Content:     post.Content,
			})
		}
		contenttmpl = "postindex"
	case 3: // comments
		comments := model.GetPageCommentByUserUuid(page, userUuid)

		for _, comment := range comments {
			postOrComment := model.GetPostOrCommentByUuid(util.BINToUUIDStr(comment.ReplyToUuid))
			entrylist = append(entrylist, CommentEntry{
				Uuid:           comment.Uuid,
				ThreadUuid:     comment.ThreadUuid,
				UserName:       model.GetUserNameById(comment.UserId),
				ReplyToName:    model.GetUserNameById(comment.ReplyToUserId),
				Content:        comment.Content,
				ReplyToContent: reflect.ValueOf(postOrComment).FieldByName("Content").Interface().(string),
			})
		}
		contenttmpl = "commentindex"
	default: // threads
		if curprivate == 0 {
			userUuid = ""
		}
		threads := model.GetPageThreads(page, order, userUuid, topicid)
		for _, thread := range threads {
			entrylist = append(entrylist, ThreadEntry{
				Uuid:      thread.Uuid,
				Topic:     model.GetTopicCateByID(thread.TopicId),
				Title:     thread.Title,
				ViewCount: thread.ViewCount,
				UserName:  model.GetUserNameById(thread.UserId),
				Content:   thread.Content,
				NumPosts:  thread.NumPosts,
			})
		}
		contenttmpl = "threadindex"
	}

	session := con.getSession(c)
	navbar, ok := session.Get("navbar").(string)
	if !ok {
		navbar = "public.navbar"
	}

	islogined := model.IsLogined(c)
	t_uuid, ok := session.Get("useruuid").(string)
	s_useruuid := t_uuid
	if !ok {
		s_useruuid = ""
	}

	tmpl := util.ParseTemplateFiles("layout", navbar, contenttmpl, "topic", "next")
	tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"topics":     model.GetAllTopics(),
		"curtopic":   topicid,
		"curorder":   order,
		"curprivate": curprivate,
		"useruuid":   s_useruuid,
		"username":   model.GetUserNameByUuid(s_useruuid),
		"list":       entrylist,
		"lenlist":    len(entrylist),
		"pagep":      page - 1,
		"page":       page,
		"pagen":      page + 1,
		"islogined":  islogined,
	})
}

func (con IndexControllers) CreateThread(c *gin.Context) {

}
func (con IndexControllers) CreatePost(c *gin.Context) {

}
func (con IndexControllers) CreateComment(c *gin.Context) {

}

func (con IndexControllers) Thread(c *gin.Context) {
	thread_uuid := c.Query("uuid")
	if thread_uuid == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	threads := model.GetThreadByUuid(thread_uuid)
	if len(threads) == 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}
	thread := threads[0]
	model.ThreadViewCountIncre(thread.Id)
	threadfull := ThreadEntry{
		Uuid:      thread.Uuid,
		Topic:     model.GetTopicCateByID(thread.TopicId),
		Title:     thread.Title,
		ViewCount: thread.ViewCount,
		UserName:  model.GetUserNameById(thread.UserId),
		Content:   thread.Content,
		NumPosts:  thread.NumPosts,
		CreatedAt: thread.CreatedAtDate(),
	}
	posts := model.GetPagePostsByThreadId(1, thread.Id)
	postEntries := []PostEntryWithComments{}
	for _, post := range posts {
		coments := model.GetAllCommentsToAPost(util.BINToUUIDStr(post.Uuid))
		commentEntries := []CommentEntry{}
		for _, comment := range coments {
			commentEntries = append(commentEntries, CommentEntry{
				UserName:    model.GetUserNameById(comment.UserId),
				ReplyToName: model.GetUserNameById(comment.ReplyToUserId),
				Content:     comment.Content,
				CreatedAt:   comment.CreatedAtDate(),
			})
		}
		postEntries = append(postEntries, PostEntryWithComments{
			PostEntry: PostEntry{
				Uuid:      post.Uuid,
				UserName:  model.GetUserNameById(post.UserId),
				ThreadPos: post.ThreadPos,
				Content:   post.Content,
				CreatedAt: post.CreatedAtDate()},
			Comments: commentEntries,
		})
	}
	session := con.getSession(c)
	navbar, ok := session.Get("navbar").(string)
	if !ok {
		navbar = "public.navbar"
	}
	threadpage, ok := session.Get("threadpage").(string)
	if !ok {
		navbar = "public.thread"
	}
	islogined := model.IsLogined(c)
	t_uuid, ok := session.Get("useruuid").(string)
	s_useruuid := t_uuid
	if !ok {
		s_useruuid = ""
	}
	tmpl := util.ParseTemplateFiles("layout", navbar, threadpage, "emptytopic", "emptynext")
	tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"useruuid":  s_useruuid,
		"username":  model.GetUserNameByUuid(s_useruuid),
		"thread":    threadfull,
		"islogined": islogined,
		"posts":     postEntries,
	})
}
