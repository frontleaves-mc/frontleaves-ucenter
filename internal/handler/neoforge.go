package handler

import (
	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	xUtil "github.com/bamboo-services/bamboo-base-go/common/utility"
	apiSync "github.com/frontleaves-mc/frontleaves-yggleaf/api/sync"
	"github.com/gin-gonic/gin"
)

// ==================== NeoForge Version List ====================

// NeoForgeVersions 获取 NeoForge 版本列表
//
// @Summary     [管理员] NeoForge 版本列表
// @Description 代理 Maven 元数据 XML，过滤 MC 1.21.1 对应的 NeoForge 版本，Redis 缓存 30 分钟
// @Tags        NeoForge
// @Produce     json
// @Success     200 {object} xBase.BaseResponse{data=apiSync.NeoForgeVersionListResponse} "获取成功"
// @Failure     500 {object} xBase.BaseResponse "服务器内部错误"
// @Router      /sync/neoforge/versions [GET]
func (h *NeoForgeHandler) NeoForgeVersions(ctx *gin.Context) {
	h.log.Info(ctx, "NeoForgeVersions - 获取版本列表")

	data, xErr := h.service.configLogic.GetNeoForgeVersions(ctx)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.SuccessHasData(ctx, "获取成功", data)
}

// ==================== NeoForge Target Version ====================

// NeoForgeTargetGet 获取目标 NeoForge 版本
//
// @Summary     [公开] 获取目标版本
// @Description 读取当前配置的目标 NeoForge 版本号
// @Tags        NeoForge
// @Produce     json
// @Success     200 {object} xBase.BaseResponse{data=apiSync.NeoForgeTargetResponse} "获取成功"
// @Router      /sync/neoforge/target [GET]
func (h *NeoForgeHandler) NeoForgeTargetGet(ctx *gin.Context) {
	h.log.Info(ctx, "NeoForgeTargetGet - 获取目标版本")

	data, xErr := h.service.configLogic.GetNeoForgeTarget(ctx)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.SuccessHasData(ctx, "获取成功", data)
}

// NeoForgeTargetSet 设置目标 NeoForge 版本
//
// @Summary     [管理员] 设置目标版本
// @Description 设置当前推荐的目标 NeoForge 版本号
// @Tags        NeoForge
// @Accept      json
// @Produce     json
// @Param       request body apiSync.SetNeoForgeTargetRequest true "设置目标版本请求"
// @Success     200 {object} xBase.BaseResponse "设置成功"
// @Failure     400 {object} xBase.BaseResponse "参数错误"
// @Router      /sync/neoforge/target [POST]
func (h *NeoForgeHandler) NeoForgeTargetSet(ctx *gin.Context) {
	h.log.Info(ctx, "NeoForgeTargetSet - 设置目标版本")

	req := xUtil.Bind(ctx, &apiSync.SetNeoForgeTargetRequest{}).Data()
	if req == nil {
		return
	}

	if xErr := h.service.configLogic.SetNeoForgeTarget(ctx, req.Version); xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.Success(ctx, "设置成功")
}
