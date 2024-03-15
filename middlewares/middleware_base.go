package middleware

import (
	"fmt"
	"model"
	"net/http"
	"time"
	"util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Timer(c *gin.Context) {
	start := time.Now().UnixMilli()
	c.Set("midw", "^^^^share data from midware^^^")
	fmt.Println("middleware-----------timer start")
	c.Next() // 暂停此中间件，执行完router注册的后续function再继续执行此中间件
	// c.Abort() //终止执行router注册的后续function，继续执行此中间件
	fmt.Println("middleware timer-----------total:", float64(time.Now().UnixMilli()-start)/1000.0, " s")
	// fmt.Println("middleware-----------timer end")

}

func SetCookie(c *gin.Context) {

	cookie, err := c.Cookie("gin_cookie")

	if err != nil {
		cookie = "NotSet"
		c.SetSameSite(http.SameSiteLaxMode)
		// SameSiteDefaultMode //设置为null
		// SameSiteLaxMode //默认值，阻止发送cookie，但对超链接放行，链接、预加载、GET表单会发送，POST表单、iframe、AJAX、Image不会发送
		// SameSiteStrictMode //跨站则阻止发送所有cookie
		// SameSiteNoneMode //无论是否跨站都会发送cookie，但浏览器会限制只有在secure为true时SameSiteNoneMode才生效
		c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		// maxage设置为0或-1会使浏览器删除该cookie
		// secure true则在https中有效
		// httpOnly true则无法通过JS等程序读取cookie，防止xss攻击
		// 一个cookie只能存一个键值对
	}

	fmt.Printf("Cookie value: %s \n", cookie)

}

func CheckSessionAuth(c *gin.Context) {
	navbar := "public.navbar"
	if IsLogined(c) {
		navbar = "private.navbar"

	}
	session := sessions.DefaultMany(c, model.UserSession)
	session.Set("navbar", navbar)
	session.Save()
}
func IsLogined(c *gin.Context) bool {
	session := sessions.DefaultMany(c, model.UserSession)
	strAny := session.Get("auth")
	islogined, _ := strAny.(string)
	return len(islogined) > 5 // TODO
}

func GetUserIfLogined(c *gin.Context) (user model.User) {
	if IsLogined(c) {
		session := sessions.DefaultMany(c, model.UserSession)
		return model.User{
			Uuid:  util.UUIDStrToBin(session.Get("useruuid").(string)),
			Id:    session.Get("userid").(int),
			Email: session.Get("useremail").(string),
		}
	}
	return model.User{}
}
