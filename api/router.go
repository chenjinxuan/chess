// @APIVersion 1.0
// @Title API
// @Description
// @Schemes http
// @Host texasholdem.com
package main

import (
	"chess/api/components/middleware"
	"chess/api/controllers/debug"
	"chess/common/config"
<<<<<<< HEAD
        "chess/api/controllers/auth"
=======
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"time"
>>>>>>> 87708c63be6678e8be7b79406073e91603ca9d2d
)

func InitRouter() {
	if config.Api.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, If-Modified-Since, x-requested-with",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	//router.Use(system_maintenance.CheckSystemMaintenance(*configServerPath))
	//router.Use(colletion.Filter(config.C.Filter))
	router.Use(middleware.SetContext())

	if config.Api.Debug {
		debugRouter := router.Group("/debug")
		{
			// pprof
			debugRouter.GET("/pprof/", c_debug.PprofIndex())
			debugRouter.GET("/pprof/heap", c_debug.PprofHeap())
			debugRouter.GET("/pprof/goroutine", c_debug.PprofGoroutine())
			debugRouter.GET("/pprof/block", c_debug.PprofBlock())
			debugRouter.GET("/pprof/threadcreate", c_debug.PprofThreadCreate())
			debugRouter.GET("/pprof/cmdline", c_debug.PprofCmdline())
			debugRouter.GET("/pprof/profile", c_debug.PprofProfile())
			debugRouter.GET("/pprof/symbol", c_debug.PprofSymbol())

			debugRouter.GET("/ip", c_debug.IP)
			//debugRouter.GET("/config", c_debug.Config)
		}
	}

	// @SubApi /auth - 授权相关 [/auth/]
	authRouter := router.Group("/auth")
	{
<<<<<<< HEAD
		authRouter.POST("/login", c_auth.Login) // 账号密码登录
		authRouter.POST("/login/quick", c_auth.LoginMobile) // 手机号快速登录
		authRouter.POST("/login/tp", c_auth.TpLogin) // 第三方登录
		authRouter.POST("/logout", nil) // 登出，销毁token
		authRouter.POST("/token/info", c_auth.TokenInfo) //获取token信息
		authRouter.POST("/token/refresh", c_auth.TokenRefrash) // 刷新token
		authRouter.POST("/register/mobile", c_auth.RegisterMobile)
		authRouter.POST("/password/reset", c_auth.PasswordReset)
=======
		authRouter.POST("/login", nil)         // 账号密码登录
		authRouter.POST("/login/quick", nil)   // 手机号快速登录
		authRouter.POST("/login/tp", nil)      // 第三方登录
		authRouter.POST("/logout", nil)        // 登出，销毁token
		authRouter.POST("/token/info", nil)    //获取token信息
		authRouter.POST("/token/refresh", nil) // 刷新token
>>>>>>> 87708c63be6678e8be7b79406073e91603ca9d2d
	}

	// @SubApi /verify - 验证码相关 [/verify/]
	verifyRouter := router.Group("/verify")
	{
		verifyRouter.POST("/send", nil)
		verifyRouter.POST("/check", nil)
	}

	// @SubApi /client - 客户端配置相关 [/client/]
	clientRouter := router.Group("/client")
	{
		clientRouter.GET("/upgrade", nil)
	}

	router.Run(config.Api.Port)
}
