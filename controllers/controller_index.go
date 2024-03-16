package controller

import (
	"fmt"
	"middleware"
	"model"
	"net/http"
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

type CreateInfo struct {
	Title   string `json:"title"`
	TopicId string `json:"topic_id"`
	Content string `json:"content"`
	Tuuid   string `json:"tuuid"`
	Puuid   string `json:"puuid"`
	Ruuid   string `json:"ruuid"`
	Action  string `json:"action"`
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
			entrylist = append(entrylist, CommentEntry{
				Uuid:           comment.Uuid,
				ThreadUuid:     comment.ThreadUuid,
				UserName:       model.GetUserNameById(comment.UserId),
				ReplyToName:    model.GetUserNameById(comment.ReplyToUserId),
				Content:        comment.Content,
				ReplyToContent: model.GetContentByPostOrCommentUuid(util.BINToUUIDStr(comment.ReplyToUuid)),
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
	s_useruuid, ok := session.Get("useruuid").(string)
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

func (con IndexControllers) Create(c *gin.Context) {
	var createinfo = CreateInfo{}
	err := c.ShouldBind(&createinfo)
	if err != nil {
		fmt.Println("post form not get,", err)
		c.JSON(http.StatusOK, gin.H{
			"msg": fmt.Sprintf("Error reading JSON body: %s", err),
		})
		return
	}
	user := middleware.GetUserIfLogined(c)
	if user.Email == "" {
		c.JSON(http.StatusOK, gin.H{
			"msg": "No user logined.",
		})
		return
	}

	fmt.Printf("%#v\n", createinfo)

	switch createinfo.Action {
	case "thread":
		topicid, err := strconv.Atoi(createinfo.TopicId)
		if err != nil || !model.IsTopicExistsByID(topicid) {
			c.JSON(http.StatusOK, gin.H{
				"msg":     fmt.Sprint("topic id not right", err),
				"success": false,
			})
			return
		}
		err = user.CreateThread(topicid, createinfo.Title, createinfo.Content)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":     fmt.Sprint("create thread failed", err),
				"success": false,
			})
			return
		}
		s_uuid := user.GetLatestThreadUuid()
		c.JSON(http.StatusOK, gin.H{
			"msg":     "Successfully",
			"success": true,
			"uuid":    s_uuid,
		})
		return
	case "post":
		err = user.CreatePost(model.GetThreadIdByUuid(createinfo.Tuuid), createinfo.Content)
	case "comment":
		err = user.CreateComment(createinfo.Puuid, createinfo.Ruuid, createinfo.Tuuid, createinfo.Content)
	default:
		c.JSON(http.StatusOK, gin.H{
			"msg":     "Invalid action.",
			"success": false,
		})
		return
	}
	fmt.Println("Create Error ", err)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":     msg,
		"success": false,
	})
}
func (con IndexControllers) CreateThread(c *gin.Context) {
	islogined := model.IsLogined(c)
	if !islogined {
		c.Redirect(http.StatusFound, "/")
		return
	}
	session := con.getSession(c)
	navbar, ok := session.Get("navbar").(string)
	if !ok {
		navbar = "public.navbar"
	}

	s_useruuid, ok := session.Get("useruuid").(string)
	if !ok {
		s_useruuid = ""
	}

	tmpl := util.ParseTemplateFiles("layout", navbar, "new.thread", "emptytopic", "emptynext")
	tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"islogined": islogined,
		"useruuid":  s_useruuid,
		"username":  model.GetUserNameByUuid(s_useruuid),
		"topics":    model.GetAllTopics(),
	})
}

func (con IndexControllers) Thread(c *gin.Context) {
	var page = 1
	var lenlist = 0
	var err error
	var order = 1 // time desc

	thread_uuid := c.Query("uuid")
	if thread_uuid == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}

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
			c.Redirect(http.StatusFound, "/thread?uuid="+thread_uuid)
			return
		}
		if oo == 0 {
			order = oo
		}
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
	posts := model.GetPagePostsByThreadId(page, thread.Id, order)
	lenlist = len(posts)
	postEntries := []PostEntryWithComments{}
	for _, post := range posts {
		coments := model.GetAllCommentsToAPost(util.BINToUUIDStr(post.Uuid))
		commentEntries := []CommentEntry{}
		for _, comment := range coments {
			commentEntries = append(commentEntries, CommentEntry{
				Uuid:        comment.Uuid,
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
	islogined := model.IsLogined(c)
	s_useruuid, ok := session.Get("useruuid").(string)
	if !ok {
		s_useruuid = ""
	}
	tmpl := util.ParseTemplateFiles("layout", navbar, "thread", "emptytopic", "emptynext")
	tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"useruuid":  s_useruuid,
		"username":  model.GetUserNameByUuid(s_useruuid),
		"thread":    threadfull,
		"islogined": islogined,
		"posts":     postEntries,
		"pagep":     page - 1,
		"page":      page,
		"pagen":     page + 1,
		"lenlist":   lenlist,
		"curorder":  order,
	})
}
