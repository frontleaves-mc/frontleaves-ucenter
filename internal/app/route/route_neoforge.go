package route

import (
	"github.com/frontleaves-mc/frontleaves-ucenter/internal/app/middleware"
	"github.com/frontleaves-mc/frontleaves-ucenter/internal/handler"
	"github.com/gin-gonic/gin"
	bSdkMiddle "github.com/phalanx-labs/beacon-sso-sdk/middleware"
)

func (r *route) neoforgeRouter(route gin.IRouter) {
	nfHandler := handler.NewHandler[handler.NeoForgeHandler](r.context, "NeoForgeHandler")

	// 公开接口（无需认证，供同步端使用）——独立 Group
	pubGroup := route.Group("/sync/neoforge")
	pubGroup.GET("/target", nfHandler.NeoForgeTargetGet)

	// 管理员接口（需认证 + 超级管理员权限）——独立 Group
	adminGroup := route.Group("/sync/neoforge")
	adminGroup.Use(bSdkMiddle.CheckAuth(r.context))
	adminGroup.Use(middleware.User(r.context))
	adminGroup.Use(middleware.SuperAdmin(r.context))
	adminGroup.GET("/versions", nfHandler.NeoForgeVersions)
	adminGroup.POST("/target", nfHandler.NeoForgeTargetSet)
}
