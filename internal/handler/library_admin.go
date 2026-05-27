package handler

import (
	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xUtil "github.com/bamboo-services/bamboo-base-go/common/utility"
	xResult "github.com/bamboo-services/bamboo-base-go/major/result"
	apiAdmin "github.com/frontleaves-mc/frontleaves-yggleaf/api/admin"
	apiLibrary "github.com/frontleaves-mc/frontleaves-yggleaf/api/library"
	"github.com/gin-gonic/gin"
)

// ==================== Admin Skin CRUD ====================

// AdminListSkins 管理员获取皮肤列表
//
// @Summary     [管理] 获取皮肤列表
// @Description 管理员分页查询皮肤库列表，支持按名称、公开状态、用户、系统资源筛选
// @Tags        管理员-资源库接口
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       name query string false "名称模糊搜索"
// @Param       is_public query bool false "公开状态"
// @Param       user_id query string false "创建者用户 ID"
// @Param       is_system query bool false "是否系统内置"
// @Param       page query int false "页码" default(1)
// @Param       page_size query int false "每页数量" default(20)
// @Success     200 {object} xBase.BaseResponse{data=apiLibrary.SkinListResponse} "获取成功"
// @Failure     400 {object} xBase.BaseResponse "请求参数错误"
// @Failure     401 {object} xBase.BaseResponse "未授权"
// @Failure     403 {object} xBase.BaseResponse "无权限"
// @Router      /library/admin/skins [GET]
func (h *LibraryHandler) AdminListSkins(ctx *gin.Context) {
	h.log.Info(ctx, "AdminListSkins - 管理员获取皮肤列表")

	req := xUtil.Bind(ctx, &apiAdmin.AdminSkinListRequest{}).Data()
	if req == nil {
		return
	}

	page, pageSize := h.parsePagination(ctx)
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
		if pageSize > maxPageSize {
			pageSize = maxPageSize
		}
	}

	filter := adminSkinFilterFromRequest(*req)
	skins, total, xErr := h.service.libraryLogic.AdminListSkins(ctx.Request.Context(), filter, page, pageSize)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	response := apiLibrary.SkinListResponse{
		Total: total,
		Items: skinDTOsToResponses(skins),
	}
	xResult.SuccessHasData(ctx, "获取皮肤列表成功", response)
}

// AdminCreateSkin 管理员创建皮肤
//
// @Summary     [管理] 创建皮肤
// @Description 管理员创建系统皮肤，指定名称、模型类型和纹理
// @Tags        管理员-资源库接口
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body apiAdmin.AdminCreateSkinRequest true "创建皮肤请求"
// @Success     200 {object} xBase.BaseResponse{data=apiLibrary.SkinResponse} "创建成功"
// @Failure     400 {object} xBase.BaseResponse "请求参数错误"
// @Failure     401 {object} xBase.BaseResponse "未授权"
// @Failure     403 {object} xBase.BaseResponse "无权限"
// @Failure     409 {object} xBase.BaseResponse "纹理哈希冲突"
// @Router      /library/admin/skins [POST]
func (h *LibraryHandler) AdminCreateSkin(ctx *gin.Context) {
	h.log.Info(ctx, "AdminCreateSkin - 管理员创建皮肤")

	req := xUtil.Bind(ctx, &apiAdmin.AdminCreateSkinRequest{}).Data()
	if req == nil {
		return
	}

	skin, xErr := h.service.libraryLogic.AdminCreateSkin(ctx.Request.Context(), req.Name, req.Model, req.Texture, req.IsPublic)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.SuccessHasData(ctx, "创建皮肤成功", skinDTOToResponse(skin))
}

// AdminGetSkinDetail 管理员获取皮肤详情
//
// @Summary     [管理] 获取皮肤详情
// @Description 管理员根据皮肤 ID 获取详情，包含装备该皮肤的档案列表
// @Tags        管理员-资源库接口
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       skin_id path string true "皮肤 ID"
// @Success     200 {object} xBase.BaseResponse{data=apiAdmin.AdminSkinDetailResponse} "获取成功"
// @Failure     400 {object} xBase.BaseResponse "请求参数错误"
// @Failure     401 {object} xBase.BaseResponse "未授权"
// @Failure     403 {object} xBase.BaseResponse "无权限"
// @Failure     404 {object} xBase.BaseResponse "资源不存在"
// @Router      /library/admin/skins/{skin_id} [GET]
func (h *LibraryHandler) AdminGetSkinDetail(ctx *gin.Context) {
	h.log.Info(ctx, "AdminGetSkinDetail - 管理员获取皮肤详情")

	skinID, err := xSnowflake.ParseSnowflakeID(ctx.Param("skin_id"))
	if err != nil {
		_ = ctx.Error(xError.NewError(ctx, xError.ParameterError, "解析皮肤 ID 失败", true, err))
		return
	}

	skinDTO, profiles, xErr := h.service.libraryLogic.AdminGetSkinDetail(ctx.Request.Context(), skinID)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.SuccessHasData(ctx, "获取皮肤详情成功", adminSkinDetailToResponse(skinDTO, profiles))
}

// AdminUpdateSkin 管理员更新皮肤
//
// @Summary     [管理] 更新皮肤
// @Description 管理员根据皮肤 ID 更新皮肤属性，支持修改名称、公开状态和纹理
// @Tags        管理员-资源库接口
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       skin_id path string true "皮肤 ID"
// @Param       request body apiAdmin.AdminUpdateSkinRequest true "更新皮肤请求"
// @Success     200 {object} xBase.BaseResponse{data=apiLibrary.SkinResponse} "更新成功"
// @Failure     400 {object} xBase.BaseResponse "请求参数错误"
// @Failure     401 {object} xBase.BaseResponse "未授权"
// @Failure     403 {object} xBase.BaseResponse "无权限"
// @Failure     404 {object} xBase.BaseResponse "资源不存在"
// @Failure     409 {object} xBase.BaseResponse "纹理哈希冲突"
// @Router      /library/admin/skins/{skin_id} [PATCH]
func (h *LibraryHandler) AdminUpdateSkin(ctx *gin.Context) {
	h.log.Info(ctx, "AdminUpdateSkin - 管理员更新皮肤")

	skinID, err := xSnowflake.ParseSnowflakeID(ctx.Param("skin_id"))
	if err != nil {
		_ = ctx.Error(xError.NewError(ctx, xError.ParameterError, "解析皮肤 ID 失败", true, err))
		return
	}

	req := xUtil.Bind(ctx, &apiAdmin.AdminUpdateSkinRequest{}).Data()
	if req == nil {
		return
	}

	updatedSkin, xErr := h.service.libraryLogic.AdminUpdateSkin(ctx.Request.Context(), skinID, req.Name, req.IsPublic, req.Texture)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.SuccessHasData(ctx, "更新皮肤成功", skinDTOToResponse(updatedSkin))
}

// AdminDeleteSkin 管理员删除皮肤
//
// @Summary     [管理] 删除皮肤
// @Description 管理员根据皮肤 ID 删除皮肤记录及关联资源
// @Tags        管理员-资源库接口
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       skin_id path string true "皮肤 ID"
// @Success     200 {object} xBase.BaseResponse "删除成功"
// @Failure     400 {object} xBase.BaseResponse "请求参数错误"
// @Failure     401 {object} xBase.BaseResponse "未授权"
// @Failure     403 {object} xBase.BaseResponse "无权限"
// @Failure     404 {object} xBase.BaseResponse "资源不存在"
// @Router      /library/admin/skins/{skin_id} [DELETE]
func (h *LibraryHandler) AdminDeleteSkin(ctx *gin.Context) {
	h.log.Info(ctx, "AdminDeleteSkin - 管理员删除皮肤")

	skinID, err := xSnowflake.ParseSnowflakeID(ctx.Param("skin_id"))
	if err != nil {
		_ = ctx.Error(xError.NewError(ctx, xError.ParameterError, "解析皮肤 ID 失败", true, err))
		return
	}

	xErr := h.service.libraryLogic.AdminDeleteSkin(ctx.Request.Context(), skinID)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.Success(ctx, "删除皮肤成功")
}

// ==================== Admin Cape CRUD ====================

// AdminListCapes 管理员获取披风列表
//
// @Summary     [管理] 获取披风列表
// @Description 管理员分页查询披风库列表，支持按名称、公开状态、用户、系统资源筛选
// @Tags        管理员-资源库接口
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       name query string false "名称模糊搜索"
// @Param       is_public query bool false "公开状态"
// @Param       user_id query string false "创建者用户 ID"
// @Param       is_system query bool false "是否系统内置"
// @Param       page query int false "页码" default(1)
// @Param       page_size query int false "每页数量" default(20)
// @Success     200 {object} xBase.BaseResponse{data=apiLibrary.CapeListResponse} "获取成功"
// @Failure     400 {object} xBase.BaseResponse "请求参数错误"
// @Failure     401 {object} xBase.BaseResponse "未授权"
// @Failure     403 {object} xBase.BaseResponse "无权限"
// @Router      /library/admin/capes [GET]
func (h *LibraryHandler) AdminListCapes(ctx *gin.Context) {
	h.log.Info(ctx, "AdminListCapes - 管理员获取披风列表")

	req := xUtil.Bind(ctx, &apiAdmin.AdminCapeListRequest{}).Data()
	if req == nil {
		return
	}

	page, pageSize := h.parsePagination(ctx)
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
		if pageSize > maxPageSize {
			pageSize = maxPageSize
		}
	}

	filter := adminCapeFilterFromRequest(*req)
	capes, total, xErr := h.service.libraryLogic.AdminListCapes(ctx.Request.Context(), filter, page, pageSize)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	response := apiLibrary.CapeListResponse{
		Total: total,
		Items: capeDTOsToResponses(capes),
	}
	xResult.SuccessHasData(ctx, "获取披风列表成功", response)
}

// AdminCreateCape 管理员创建披风
//
// @Summary     [管理] 创建披风
// @Description 管理员创建系统披风，指定名称和纹理
// @Tags        管理员-资源库接口
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body apiAdmin.AdminCreateCapeRequest true "创建披风请求"
// @Success     200 {object} xBase.BaseResponse{data=apiLibrary.CapeResponse} "创建成功"
// @Failure     400 {object} xBase.BaseResponse "请求参数错误"
// @Failure     401 {object} xBase.BaseResponse "未授权"
// @Failure     403 {object} xBase.BaseResponse "无权限"
// @Failure     409 {object} xBase.BaseResponse "纹理哈希冲突"
// @Router      /library/admin/capes [POST]
func (h *LibraryHandler) AdminCreateCape(ctx *gin.Context) {
	h.log.Info(ctx, "AdminCreateCape - 管理员创建披风")

	req := xUtil.Bind(ctx, &apiAdmin.AdminCreateCapeRequest{}).Data()
	if req == nil {
		return
	}

	cape, xErr := h.service.libraryLogic.AdminCreateCape(ctx.Request.Context(), req.Name, req.Texture, req.IsPublic)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.SuccessHasData(ctx, "创建披风成功", capeDTOToResponse(cape))
}

// AdminGetCapeDetail 管理员获取披风详情
//
// @Summary     [管理] 获取披风详情
// @Description 管理员根据披风 ID 获取详情，包含装备该披风的档案列表
// @Tags        管理员-资源库接口
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       cape_id path string true "披风 ID"
// @Success     200 {object} xBase.BaseResponse{data=apiAdmin.AdminCapeDetailResponse} "获取成功"
// @Failure     400 {object} xBase.BaseResponse "请求参数错误"
// @Failure     401 {object} xBase.BaseResponse "未授权"
// @Failure     403 {object} xBase.BaseResponse "无权限"
// @Failure     404 {object} xBase.BaseResponse "资源不存在"
// @Router      /library/admin/capes/{cape_id} [GET]
func (h *LibraryHandler) AdminGetCapeDetail(ctx *gin.Context) {
	h.log.Info(ctx, "AdminGetCapeDetail - 管理员获取披风详情")

	capeID, err := xSnowflake.ParseSnowflakeID(ctx.Param("cape_id"))
	if err != nil {
		_ = ctx.Error(xError.NewError(ctx, xError.ParameterError, "解析披风 ID 失败", true, err))
		return
	}

	capeDTO, profiles, xErr := h.service.libraryLogic.AdminGetCapeDetail(ctx.Request.Context(), capeID)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.SuccessHasData(ctx, "获取披风详情成功", adminCapeDetailToResponse(capeDTO, profiles))
}

// AdminUpdateCape 管理员更新披风
//
// @Summary     [管理] 更新披风
// @Description 管理员根据披风 ID 更新披风属性，支持修改名称、公开状态和纹理
// @Tags        管理员-资源库接口
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       cape_id path string true "披风 ID"
// @Param       request body apiAdmin.AdminUpdateCapeRequest true "更新披风请求"
// @Success     200 {object} xBase.BaseResponse{data=apiLibrary.CapeResponse} "更新成功"
// @Failure     400 {object} xBase.BaseResponse "请求参数错误"
// @Failure     401 {object} xBase.BaseResponse "未授权"
// @Failure     403 {object} xBase.BaseResponse "无权限"
// @Failure     404 {object} xBase.BaseResponse "资源不存在"
// @Failure     409 {object} xBase.BaseResponse "纹理哈希冲突"
// @Router      /library/admin/capes/{cape_id} [PATCH]
func (h *LibraryHandler) AdminUpdateCape(ctx *gin.Context) {
	h.log.Info(ctx, "AdminUpdateCape - 管理员更新披风")

	capeID, err := xSnowflake.ParseSnowflakeID(ctx.Param("cape_id"))
	if err != nil {
		_ = ctx.Error(xError.NewError(ctx, xError.ParameterError, "解析披风 ID 失败", true, err))
		return
	}

	req := xUtil.Bind(ctx, &apiAdmin.AdminUpdateCapeRequest{}).Data()
	if req == nil {
		return
	}

	updatedCape, xErr := h.service.libraryLogic.AdminUpdateCape(ctx.Request.Context(), capeID, req.Name, req.IsPublic, req.Texture)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.SuccessHasData(ctx, "更新披风成功", capeDTOToResponse(updatedCape))
}

// AdminDeleteCape 管理员删除披风
//
// @Summary     [管理] 删除披风
// @Description 管理员根据披风 ID 删除披风记录及关联资源
// @Tags        管理员-资源库接口
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       cape_id path string true "披风 ID"
// @Success     200 {object} xBase.BaseResponse "删除成功"
// @Failure     400 {object} xBase.BaseResponse "请求参数错误"
// @Failure     401 {object} xBase.BaseResponse "未授权"
// @Failure     403 {object} xBase.BaseResponse "无权限"
// @Failure     404 {object} xBase.BaseResponse "资源不存在"
// @Router      /library/admin/capes/{cape_id} [DELETE]
func (h *LibraryHandler) AdminDeleteCape(ctx *gin.Context) {
	h.log.Info(ctx, "AdminDeleteCape - 管理员删除披风")

	capeID, err := xSnowflake.ParseSnowflakeID(ctx.Param("cape_id"))
	if err != nil {
		_ = ctx.Error(xError.NewError(ctx, xError.ParameterError, "解析披风 ID 失败", true, err))
		return
	}

	xErr := h.service.libraryLogic.AdminDeleteCape(ctx.Request.Context(), capeID)
	if xErr != nil {
		_ = ctx.Error(xErr)
		return
	}

	xResult.Success(ctx, "删除披风成功")
}
