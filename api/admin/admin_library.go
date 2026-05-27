package admin

import (
	library "github.com/frontleaves-mc/frontleaves-yggleaf/api/library"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
)

// AdminCreateSkinRequest 创建皮肤请求
type AdminCreateSkinRequest struct {
	Name     string  `json:"name" binding:"required"`
	Model    uint8   `json:"model" binding:"required,oneof=1 2"`
	Texture  string  `json:"texture" binding:"required"`
	IsPublic *bool   `json:"is_public,omitempty"`
}

// AdminUpdateSkinRequest 更新皮肤请求
type AdminUpdateSkinRequest struct {
	Name     *string `json:"name,omitempty"`
	IsPublic *bool   `json:"is_public,omitempty"`
	Texture  *string `json:"texture,omitempty"`
}

// AdminSkinListRequest 皮肤列表请求
type AdminSkinListRequest struct {
	Name     string `form:"name" binding:"omitempty,max=64"`
	IsPublic *bool  `form:"is_public"`
	UserID   string `form:"user_id" binding:"omitempty,max=64"`
	IsSystem *bool  `form:"is_system"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// AdminSkinFilter 皮肤过滤器（内部类型）
type AdminSkinFilter struct {
	Name     *string
	IsPublic *bool
	UserID   *xSnowflake.SnowflakeID
	IsSystem *bool
}

// AdminSkinDetailResponse 皮肤详情响应
type AdminSkinDetailResponse struct {
	library.SkinResponse
	EquippedProfiles []AdminEquippedProfile `json:"equipped_profiles"`
}

// AdminEquippedProfile 装备档案信息
type AdminEquippedProfile struct {
	ID       string `json:"id"`
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

// AdminCreateCapeRequest 创建披风请求
type AdminCreateCapeRequest struct {
	Name     string `json:"name" binding:"required"`
	Texture  string `json:"texture" binding:"required"`
	IsPublic *bool  `json:"is_public,omitempty"`
}

// AdminUpdateCapeRequest 更新披风请求
type AdminUpdateCapeRequest struct {
	Name     *string `json:"name,omitempty"`
	IsPublic *bool   `json:"is_public,omitempty"`
	Texture  *string `json:"texture,omitempty"`
}

// AdminCapeListRequest 披风列表请求
type AdminCapeListRequest struct {
	Name     string `form:"name" binding:"omitempty,max=64"`
	IsPublic *bool  `form:"is_public"`
	UserID   string `form:"user_id" binding:"omitempty,max=64"`
	IsSystem *bool  `form:"is_system"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// AdminCapeFilter 披风过滤器（内部类型）
type AdminCapeFilter struct {
	Name     *string
	IsPublic *bool
	UserID   *xSnowflake.SnowflakeID
	IsSystem *bool
}

// AdminCapeDetailResponse 披风详情响应
type AdminCapeDetailResponse struct {
	library.CapeResponse
	EquippedProfiles []AdminEquippedProfile `json:"equipped_profiles"`
}