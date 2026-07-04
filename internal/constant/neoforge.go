package bConst

// NeoForge 版本同步相关常量配置。

const (
	// NeoForgeMavenMetadataURL Maven 元数据 XML 地址（含全部已发布版本）。
	NeoForgeMavenMetadataURL = "https://maven.neoforged.net/releases/net/neoforged/neoforge/maven-metadata.xml"

	// NeoForgeAPITimeoutSec 单次 HTTP 请求超时（秒）。
	NeoForgeAPITimeoutSec = 15

	// NeoForgeMinecraftVersion 当前服务的 Minecraft 基础版本。
	NeoForgeMinecraftVersion = "1.21.1"

	// NeoForgeVersionPrefix 目标版本前缀（仅保留该 MC 版本对应的 NeoForge 版本）。
	NeoForgeVersionPrefix = "21.1."

	// NeoForgeCacheTTLMinutes Redis 缓存存活时长（分钟）。
	NeoForgeCacheTTLMinutes = 30

	// NeoForge 配置命名空间与键（存储在 config 表中）
	NeoForgeConfigNamespace  = "neoforge"
	NeoForgeConfigKeyVersion = "target_version"
)
