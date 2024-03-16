package controller

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"model"
	"net/http"
	"os"
	"sync"
	"util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	BaseController
}

func (con UserController) Login(c *gin.Context) {
	t := util.ParseTemplateFiles("layout", "public.navbar", "login", "emptytopic", "emptynext")
	t.ExecuteTemplate(c.Writer, "layout", gin.H{})
}

func (con UserController) Check(c *gin.Context) {
	user := model.User{}
	err := c.ShouldBind(&user)
	if err == nil {
		users := model.GetUserByEmail(user.Email)
		if len(users) != 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": fmt.Sprintf("Welcome back, %v", users[0].Name),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Error reading JSON body: %s", err),
	})
}
func (con UserController) CheckExists(c *gin.Context) {
	user := model.User{}
	err := c.ShouldBind(&user)
	if err == nil {
		users := model.GetUserByEmail(user.Email)
		if len(users) != 0 {
			c.JSON(http.StatusOK, gin.H{
				"emailValid": false,
				"msg":        "Email already registered.",
			})
			return
		}
		if !util.IsValidEmail(user.Email) {
			c.JSON(http.StatusOK, gin.H{
				"emailValid": false,
				"msg":        "Email invalid.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"emailValid": true,
			"msg":        "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"emailValid": false,
		"msg":        fmt.Sprintf("Error reading JSON body: %s", err),
	})
}
func (con UserController) Salt(c *gin.Context) {
	user := model.User{}
	err := c.ShouldBind(&user)
	if err == nil {
		users := model.GetUserByEmail(user.Email)
		if len(users) != 0 {
			c.JSON(http.StatusOK, gin.H{
				"salt": users[0].Salt,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"salt": util.CreateSalt(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Error reading JSON body: %s", err),
	})

}

func (con UserController) Auth(c *gin.Context) {
	user := model.User{}
	err := c.ShouldBind(&user)
	if err == nil {
		if user.Auth() {
			full_user := model.GetUserByEmail(user.Email)[0]
			con.setSessionLogined(c, full_user)
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"msg":     "Login successfully.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "Wrong Password or Email.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"msg":     fmt.Sprintf("Error reading JSON body: %s", err),
	})
}

func (con UserController) setSessionLogined(c *gin.Context, user model.User) {
	session := con.getSession(c)
	session.Set("auth", "#*$%"+user.Email+"#*$%")
	session.Set("useruuid", util.BINToUUIDStr(user.Uuid))
	session.Set("userid", user.Id)
	session.Set("useremail", user.Email)
	session.Save()
}

func (con UserController) Logout(c *gin.Context) {
	if model.IsLogined(c) {
		session := con.getSession(c)
		session.Clear()
		session.Options(sessions.Options{Path: "/", MaxAge: -1})
		session.Save()
	}
	c.Redirect(http.StatusFound, "/")
}

func (con UserController) SignUp(c *gin.Context) {
	salt := util.CreateSalt()
	session := con.getSession(c)
	session.Set("salt", salt)
	session.Save()
	t := util.ParseTemplateFiles("layout", "public.navbar", "signup", "emptytopic", "emptynext")
	t.ExecuteTemplate(c.Writer, "layout", gin.H{
		"salt": salt,
	})
}

func (con UserController) DoSignUp(c *gin.Context) {
	user := model.User{}
	err := c.ShouldBind(&user)
	session := con.getSession(c)
	strAny := session.Get("salt")
	salt := strAny.(string)
	user.Salt = salt
	if err == nil {
		if !model.IsUserExists(user.Email) {
			err := user.CreateUser()
			if err == nil {
				full_user := model.GetUserByEmail(user.Email)[0]
				con.setSessionLogined(c, full_user)
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"msg":     "Sign up sucessfully!",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"success": false,
					"msg":     fmt.Sprintf("Create user falied: %s", err),
				})
			}
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"msg":     "Email already registered.",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"msg":     fmt.Sprintf("Error reading JSON body: %s", err),
	})

}

func (con UserController) ForgeUser(c *gin.Context) {
	randomUsers(100)
	c.JSON(http.StatusOK, gin.H{
		"msg": "forge completed",
	})
}

func randomUsers(n int) {
	sem := make(chan struct{}, 80)
	errlog := make(chan string)
	var wg sync.WaitGroup

	file, err := os.Create("errlog.txt")
	if err != nil {

		fmt.Println("Error creating file:", err)
	}
	defer file.Close()

	pwd := "123456"

	for i := range n {

		wg.Add(1)
		go func(wg *sync.WaitGroup, sem chan struct{}, errlog chan<- string) {
			sem <- struct{}{}
			defer wg.Done()
			defer func() { <-sem }()

			hasher := md5.New()
			salt := util.CreateSalt()
			hasher.Write([]byte(pwd + salt))
			t_user := model.User{
				Uuid:   util.CreateUUIDBin(),
				Name:   fmt.Sprintf("%v_%v", i, util.RandomName()),
				Email:  util.RandomEmail(),
				Pwdmd5: hex.EncodeToString(hasher.Sum(nil)),
				Salt:   salt,
			}
			err := t_user.CreateUser()
			if err != nil {
				errlog <- err.Error()
			}
			user := model.GetUserByEmail(t_user.Email)[0]

			for range 3 {
				rtitle, rcontent := util.RandomThread()
				err := user.CreateThread(model.RandomTopicId(), rtitle, rcontent)
				if err != nil {
					errlog <- err.Error()
				}
			}
			for range 12 {
				rthreadid := model.RandomThreadId()
				err := user.CreatePost(rthreadid, util.RandomContent())
				if err != nil {
					errlog <- err.Error()
				}
				model.ThreadViewCountIncre(rthreadid)

			}
			for range 25 {
				post := model.GetPostById(model.RandomPostId())[0]
				model.ThreadViewCountIncre(post.ThreadId)
				t_thread := model.GetThreadById(post.ThreadId)

				var possibleReplyMap = make(map[int]uuid.UUID)
				possibleReplyMap[post.UserId] = post.Uuid // post author

				for _, v := range model.GetAllCommentsToAPost(util.BINToUUIDStr(post.Uuid)) {
					possibleReplyMap[v.UserId] = v.Uuid
				}
				possibleIds := util.MapKeys(possibleReplyMap)
				r_reply_to_id := possibleIds[rand.Intn(len(possibleIds))]

				err := user.CreateComment(
					util.BINToUUIDStr(post.Uuid),
					util.BINToUUIDStr(possibleReplyMap[r_reply_to_id]),
					util.BINToUUIDStr(t_thread.Uuid),
					util.RandomContent())
				if err != nil {
					errlog <- err.Error()
				}
			}
		}(&wg, sem, errlog)
	}
	//closer
	go func() {
		wg.Wait()
		close(sem)
		close(errlog)
	}()

	for errstr := range errlog {
		file.WriteString(errstr + "\n")
	}

}

// func (con UserController) Killsleeper(c *gin.Context) {
// 	SleepConnsKiller(model.GetDB())
// }
// func SleepConnsKiller(db *gorm.DB) {
// 	var sleepIds []int
// 	db.Raw("SELECT id FROM INFORMATION_SCHEMA.processlist").Scan(&sleepIds)
// 	fmt.Printf("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx  %v", len(sleepIds))
// 	for id := range sleepIds {
// 		db.Exec("KILL ?", id)
// 	}
// }
// test only
// func (con UserController) Test(c *gin.Context) {
// 	users := model.QueryUser("fin7514@163.com")

// 	c.JSON(http.StatusOK, gin.H{
// 		"result": users,
// 	})
// }
// func (con UserController) PTest(c *gin.Context) {
// 	var user model.User
// 	err := c.ShouldBind(&user)

// 	// c.JSON(http.StatusOK, gin.H{
// 	// 	"result": user,
// 	// 	"err":    err,
// 	// })
// 	fmt.Printf("%#v, %v", user, err)
// }

/*--------------------------------------------------------*/
