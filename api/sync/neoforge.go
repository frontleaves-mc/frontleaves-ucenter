package sync

import "time"

// NeoForgeVersionItem NeoForge 单个版本条目。
type NeoForgeVersionItem struct {
	Version       string `json:"version"`       // 版本号，如 "21.1.235"
	IsLatest      bool   `json:"isLatest"`      // 是否为最新版本
	IsRecommended bool   `json:"isRecommended"` // 是否为推荐版本（稳定版最新）
}

// NeoForgeVersionListResponse NeoForge 版本列表响应。
type NeoForgeVersionListResponse struct {
	MinecraftVersion string                `json:"minecraftVersion"` // Minecraft 基础版本
	Latest           string                `json:"latest"`           // 最新版本号
	Recommended      string                `json:"recommended"`      // 推荐版本号
	Versions         []NeoForgeVersionItem `json:"versions"`         // 版本列表（按版本号降序排列）
}

// NeoForgeTargetResponse NeoForge 目标版本响应。
type NeoForgeTargetResponse struct {
	Version   string    `json:"version"`   // 当前设定的目标版本号（未设定时为空字符串）
	UpdatedAt time.Time `json:"updatedAt"` // 目标版本最后更新时间（未设定时为零值）
}

// SetNeoForgeTargetRequest 设置 NeoForge 目标版本请求。
type SetNeoForgeTargetRequest struct {
	Version string `json:"version" binding:"required"` // 目标版本号，如 "21.1.235"
}
