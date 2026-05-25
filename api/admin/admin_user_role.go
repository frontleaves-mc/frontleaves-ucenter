package admin

import "time"

// AdminUpdateRoleRequest 管理员更新用户角色请求。
type AdminUpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=PLAYER ADMIN"`
}

// AdminUpdateRoleResponse 管理员更新用户角色响应。
type AdminUpdateRoleResponse struct {
	UserID    string    `json:"user_id"`
	RoleName  string    `json:"role_name"`
	UpdatedAt time.Time `json:"updated_at"`
}
