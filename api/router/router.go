package router

import (
    "github.com/gin-gonic/gin"
    . "github.com/liyinda/google-authenticator/api/apis"
    "net/http"
    "fmt"
    "github.com/liyinda/google-authenticator/middleware/jwt"
    "github.com/gin-gonic/contrib/sessions"
    "github.com/gin-contrib/cors"
)

func InitRouter() *gin.Engine {
    router := gin.Default()

    router.Use(cors.New(cors.Config{
        //AllowOrigins:     []string{"*"},
        AllowOrigins:     []string{"http://192.168.30.18", "http://localhost"},
        AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS", "DELETE"},
        //AllowHeaders:     []string{"Content-Type,Authorization,X-Token,Access-Control-Allow-Origin"},
        AllowHeaders:     []string{"Content-Type","Authorization","X-Token","Access-Control-Allow-Origin"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        AllowOriginFunc: func(origin string) bool {
            return origin == "*"
        },
    }))

    //引用静态资源
    router.LoadHTMLGlob("dist/*.html")
    router.LoadHTMLFiles("static/*/*")
    router.Static("/static", "./dist/static")
    router.StaticFile("/vue/", "dist/index.html")

    //设置sessions
    store := sessions.NewCookieStore([]byte("secret"))
    router.Use(sessions.Sessions("mysession", store))

    //登录入口
    passport := router.Group("/passport")
    {
        passport.POST("/login", Login) 
        passport.POST("/logout", Logout) 
    }

    //用户管理入口
    home := router.Group("/home")
    home.Use(jwt.JWT())
    {
        home.GET("/userinfo", Userinfo) 
        home.POST("/useradd", Useradd) 
        //home.POST("/useredit", Useredit) 
        home.PUT("/useredit", Useredit) 
        home.GET("/userlist", Userlist) 
        home.DELETE("/userdel", Userdel) 
    }
    home.Use(AuthRequired())

    //定义默认路由
    router.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, gin.H{
            "status": 404,
            "error":  "404, page not exists!",
        })
    })

    //接口管理入口
    api := router.Group("/api")
    {
        api.GET("/apiqrcode", ApiQrcode)
    }


    //设置cookie
    router.GET("/cookie", func(c *gin.Context) {

        cookie, err := c.Cookie("gin_cookie")

        if err != nil {
            cookie = "NotSet"
            c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
        }

        fmt.Printf("Cookie value: %s \n", cookie)
    })

    return router
}

