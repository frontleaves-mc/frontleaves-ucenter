package logic

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	apiSync "github.com/frontleaves-mc/frontleaves-yggleaf/api/sync"
	bConst "github.com/frontleaves-mc/frontleaves-yggleaf/internal/constant"
	"github.com/frontleaves-mc/frontleaves-yggleaf/internal/repository"
	"github.com/frontleaves-mc/frontleaves-yggleaf/internal/repository/cache"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/common/utility/context"
)

// mavenMetadata NeoForge Maven 仓库的 maven-metadata.xml 解析结构。
type mavenMetadata struct {
	XMLName    xml.Name `xml:"metadata"`
	Versioning struct {
		Latest   string   `xml:"latest"`
		Release  string   `xml:"release"`
		Versions []string `xml:"versions>version"`
	} `xml:"versioning"`
}

// ConfigLogic 通用键值对配置业务逻辑。
//
// 通过嵌入匿名的 `logic` 结构体继承基础设施依赖，
// 通过 ConfigRepository 完成配置的读写操作，
// 通过 NeoForgeVersionCache 提供 Maven 元数据的多级缓存。
type ConfigLogic struct {
	logic
	configRepo    *repository.ConfigRepository
	neoforgeCache *cache.NeoForgeVersionCache
}

// NewConfigLogic 创建通用配置业务逻辑实例。
func NewConfigLogic(ctx context.Context) *ConfigLogic {
	db := xCtxUtil.MustGetDB(ctx)
	rdb := xCtxUtil.MustGetRDB(ctx)

	return &ConfigLogic{
		logic: logic{
			db:  db,
			rdb: rdb,
			log: xLog.WithName(xLog.NamedLOGC, "ConfigLogic"),
		},
		configRepo:    repository.NewConfigRepository(),
		neoforgeCache: &cache.NeoForgeVersionCache{RDB: rdb},
	}
}

// Get 根据命名空间和键获取配置值，未配置时返回空字符串。
func (l *ConfigLogic) Get(ctx context.Context, namespace, key string) (string, *xError.Error) {
	return l.configRepo.GetByKey(ctx, namespace, key)
}

// Set 保存配置值（存在则更新，不存在则创建）。
func (l *ConfigLogic) Set(ctx context.Context, namespace, key, value string) *xError.Error {
	return l.configRepo.Set(ctx, namespace, key, value)
}

// GetNeoForgeVersions 获取 NeoForge 版本列表（Maven XML 代理 + Redis 缓存）。
//
// 流程:
//  1. 尝试从 Redis 缓存命中 → 直接返回
//  2. 缓存未命中 → HTTP GET 拉取 Maven maven-metadata.xml
//  3. 过滤出目标 MC 版本对应的 NeoForge 版本（前缀匹配 + 排除 beta）
//  4. 按版本号降序排序，标记首个为 IsLatest / IsRecommended
//  5. 写入缓存（30min TTL）后返回
func (l *ConfigLogic) GetNeoForgeVersions(ctx context.Context) (*apiSync.NeoForgeVersionListResponse, *xError.Error) {
	// 1. 尝试缓存
	if cached, err := l.neoforgeCache.Get(ctx); err != nil {
		l.log.Warn(ctx, fmt.Sprintf("NeoForge 版本缓存读取失败，将回源: %v", err))
	} else if cached != nil {
		return cached, nil
	}

	// 2. 拉取 Maven 元数据
	client := &http.Client{Timeout: time.Duration(bConst.NeoForgeAPITimeoutSec) * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bConst.NeoForgeMavenMetadataURL, nil)
	if err != nil {
		return nil, xError.NewError(ctx, xError.ServerInternalError, "创建 Maven 元数据请求失败", false, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, xError.NewError(ctx, xError.ServerInternalError, "拉取 Maven 元数据失败", false, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, xError.NewError(ctx, xError.ServerInternalError,
			xError.ErrMessage(fmt.Sprintf("Maven 元数据返回异常状态码: %d", resp.StatusCode)), true, nil)
	}

	var metadata mavenMetadata
	if err := xml.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, xError.NewError(ctx, xError.ServerInternalError, "解析 Maven 元数据 XML 失败", false, err)
	}

	// 3. 过滤版本：仅保留目标 MC 版本对应的 NeoForge 版本，排除 beta
	filtered := make([]string, 0, len(metadata.Versioning.Versions))
	for _, v := range metadata.Versioning.Versions {
		if strings.HasPrefix(v, bConst.NeoForgeVersionPrefix) && !strings.Contains(v, "-beta") {
			filtered = append(filtered, v)
		}
	}

	// 4. 降序排序（最新优先）
	sort.Sort(sort.Reverse(sort.StringSlice(filtered)))

	// 5. 构建响应
	versions := make([]apiSync.NeoForgeVersionItem, 0, len(filtered))
	for i, v := range filtered {
		item := apiSync.NeoForgeVersionItem{Version: v}
		if i == 0 {
			item.IsLatest = true
			item.IsRecommended = true
		}
		versions = append(versions, item)
	}

	latest := ""
	recommended := ""
	if len(filtered) > 0 {
		latest = filtered[0]
		recommended = filtered[0]
	}

	result := &apiSync.NeoForgeVersionListResponse{
		MinecraftVersion: bConst.NeoForgeMinecraftVersion,
		Latest:           latest,
		Recommended:      recommended,
		Versions:         versions,
	}

	// 6. 写入缓存（best-effort，失败不影响响应）
	if err := l.neoforgeCache.Set(ctx, result); err != nil {
		l.log.Warn(ctx, fmt.Sprintf("NeoForge 版本缓存写入失败: %v", err))
	}

	return result, nil
}

// GetNeoForgeTarget 获取当前设定的 NeoForge 目标版本。
//
// 目标版本未设定时（Config 表无记录），返回零值响应（Version=""）而非错误。
func (l *ConfigLogic) GetNeoForgeTarget(ctx context.Context) (*apiSync.NeoForgeTargetResponse, *xError.Error) {
	entity, xErr := l.configRepo.GetEntityByKey(ctx, "neoforge", "target_version")
	if xErr != nil {
		return nil, xErr
	}
	if entity == nil {
		return &apiSync.NeoForgeTargetResponse{Version: "", UpdatedAt: time.Time{}}, nil
	}
	return &apiSync.NeoForgeTargetResponse{Version: entity.Value, UpdatedAt: entity.UpdatedAt}, nil
}

// SetNeoForgeTarget 设定 NeoForge 目标版本（upsert）。
func (l *ConfigLogic) SetNeoForgeTarget(ctx context.Context, version string) *xError.Error {
	return l.configRepo.Set(ctx, "neoforge", "target_version", version)
}
