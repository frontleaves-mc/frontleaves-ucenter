package route

import (
	"github.com/frontleaves-mc/frontleaves-ucenter/internal/app/middleware"
	"github.com/frontleaves-mc/frontleaves-ucenter/internal/handler"
	"github.com/gin-gonic/gin"
	bSdkMiddle "github.com/phalanx-labs/beacon-sso-sdk/middleware"
)

func (r *route) userRouter(route gin.IRouter) {
	userHandler := handler.NewHandler[handler.UserHandler](r.context, "UserHandler")

	userGroup := route.Group("/user")
	userGroup.Use(bSdkMiddle.CheckAuth(r.context))
	userGroup.Use(middleware.User(r.context))

	{
		userGroup.GET("/info", userHandler.UserCurrent)
		userGroup.PUT("/game-password", userHandler.UpdateGamePassword)
	}
}
