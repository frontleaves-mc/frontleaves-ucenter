package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	bConst "github.com/frontleaves-mc/frontleaves-ucenter/internal/constant"
	apiSync "github.com/frontleaves-mc/frontleaves-ucenter/api/sync"
)

// NeoForgeVersionCache NeoForge 版本列表缓存（Redis）。
//
// 使用 String 结构存储 JSON 序列化的版本列表响应，键格式为 fyl:neoforge:versions。
// TTL 由 NeoForgeCacheTTLMinutes 控制（默认 30 分钟），防止过度请求 Maven 仓库。
type NeoForgeVersionCache xCache.Cache

// Set 缓存 NeoForge 版本列表。
func (c *NeoForgeVersionCache) Set(ctx context.Context, data *apiSync.NeoForgeVersionListResponse) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ttl := time.Duration(bConst.NeoForgeCacheTTLMinutes) * time.Minute
	return c.RDB.Set(ctx, bConst.CacheNeoForgeVersions.Get().String(), jsonBytes, ttl).Err()
}

// Get 获取缓存的 NeoForge 版本列表。
//
// 返回值:
//   - *apiSync.NeoForgeVersionListResponse: 版本列表响应，未命中缓存时返回 nil。
//   - error: 操作过程中发生的错误（redis.Nil 视为未命中，返回 nil, nil）。
func (c *NeoForgeVersionCache) Get(ctx context.Context) (*apiSync.NeoForgeVersionListResponse, error) {
	data, err := c.RDB.Get(ctx, bConst.CacheNeoForgeVersions.Get().String()).Bytes()
	if err != nil {
		if err == redis.Nil { // 缓存未命中
			return nil, nil
		}
		return nil, err
	}
	var cached apiSync.NeoForgeVersionListResponse
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}
	return &cached, nil
}
