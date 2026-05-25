package admin

import "time"

// AdminGameProfileListRequest 管理员游戏档案列表查询参数。
type AdminGameProfileListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword" binding:"omitempty,max=64"`
	UUID     string `form:"uuid" binding:"omitempty,max=36"`
}

// AdminGameProfileListResponse 管理员游戏档案列表响应。
type AdminGameProfileListResponse struct {
	List  []AdminGameProfileListItem `json:"list"`
	Total int64                      `json:"total"`
	Page  int                        `json:"page"`
	Size  int                        `json:"size"`
}

// AdminGameProfileListItem 游戏档案列表项（精简字段）。
type AdminGameProfileListItem struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Username       string    `json:"username"`
	UUID           string    `json:"uuid"`
	Name           string    `json:"name"`
	SkinLibraryID  *string   `json:"skin_library_id,omitempty"`
	CapeLibraryID  *string   `json:"cape_library_id,omitempty"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AdminGameProfileDetailResponse 管理员游戏档案详情响应。
type AdminGameProfileDetailResponse struct {
	ID             string                  `json:"id"`
	UserID         string                  `json:"user_id"`
	UUID           string                  `json:"uuid"`
	Name           string                  `json:"name"`
	SkinLibraryID  *string                 `json:"skin_library_id,omitempty"`
	CapeLibraryID  *string                 `json:"cape_library_id,omitempty"`
	UpdatedAt      time.Time               `json:"updated_at"`
	User           AdminGameProfileUser    `json:"user"`
}

// AdminGameProfileUser 游戏档案关联用户信息。
type AdminGameProfileUser struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
}

// AdminForceSetSkinRequest 管理员强制设置皮肤请求。
type AdminForceSetSkinRequest struct {
	SkinLibraryID *string `json:"skin_library_id,omitempty"` // 皮肤库 ID（null 或空字符串表示卸下）
}

// AdminForceSetCapeRequest 管理员强制设置披风请求。
type AdminForceSetCapeRequest struct {
	CapeLibraryID *string `json:"cape_library_id,omitempty"` // 披风库 ID（null 或空字符串表示卸下）
}
