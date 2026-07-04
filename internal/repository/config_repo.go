package repository

import (
	"context"
	"errors"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/common/utility/context"
	"github.com/frontleaves-mc/frontleaves-yggleaf/internal/entity"
	"gorm.io/gorm"
)

// ConfigRepository 通用键值对配置仓储，负责按命名空间隔离的配置数据访问。
type ConfigRepository struct {
	log *xLog.LogNamedLogger
}

// NewConfigRepository 初始化并返回 ConfigRepository 实例。
func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{
		log: xLog.WithName(xLog.NamedREPO, "ConfigRepository"),
	}
}

// GetByKey 根据命名空间和键获取单个配置值。
//
// 当记录不存在时返回空字符串与 nil 错误（视为"未配置"而非数据库异常），
// 仅在其他数据库错误时返回 xError.DatabaseError。
func (r *ConfigRepository) GetByKey(ctx context.Context, namespace, key string) (string, *xError.Error) {
	r.log.Info(ctx, "GetByKey - 查询配置")
	var config entity.Config
	db := xCtxUtil.MustGetDB(ctx)
	result := db.WithContext(ctx).Where("namespace = ? AND key = ?", namespace, key).First(&config)
	if result.Error != nil {
		// 记录不存在视为"未配置"，返回空值而非数据库错误
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", xError.NewError(ctx, xError.DatabaseError, "查询配置失败", false, result.Error)
	}
	return config.Value, nil
}

// GetEntityByKey 根据命名空间和键获取完整的 Config 实体（含 UpdatedAt）。
//
// 当记录不存在时返回 nil + nil（视为"未配置"而非数据库异常），
// 仅在其他数据库错误时返回 xError.DatabaseError。
func (r *ConfigRepository) GetEntityByKey(ctx context.Context, namespace, key string) (*entity.Config, *xError.Error) {
	r.log.Info(ctx, "GetEntityByKey - 查询配置实体")
	var config entity.Config
	db := xCtxUtil.MustGetDB(ctx)
	result := db.WithContext(ctx).Where("namespace = ? AND key = ?", namespace, key).First(&config)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询配置失败", false, result.Error)
	}
	return &config, nil
}

// BatchGet 根据命名空间批量获取所有配置。
func (r *ConfigRepository) BatchGet(ctx context.Context, namespace string) (map[string]string, *xError.Error) {
	r.log.Info(ctx, "BatchGet - 批量查询配置")
	var configs []entity.Config
	db := xCtxUtil.MustGetDB(ctx)
	if result := db.WithContext(ctx).Where("namespace = ?", namespace).Find(&configs); result.Error != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "批量查询配置失败", false, result.Error)
	}
	result := make(map[string]string, len(configs))
	for _, c := range configs {
		result[c.Key] = c.Value
	}
	return result, nil
}

// Set 保存配置（存在则更新，不存在则创建）。
func (r *ConfigRepository) Set(ctx context.Context, namespace, key, value string) *xError.Error {
	r.log.Info(ctx, "Set - 保存配置")
	db := xCtxUtil.MustGetDB(ctx)
	var config entity.Config
	result := db.WithContext(ctx).Where("namespace = ? AND key = ?", namespace, key).First(&config)
	if result.Error == nil {
		// 存在则更新
		config.Value = value
		if result := db.WithContext(ctx).Save(&config); result.Error != nil {
			return xError.NewError(ctx, xError.DatabaseError, "更新配置失败", false, result.Error)
		}
	} else {
		// 不存在则创建（其他非 ErrRecordNotFound 的数据库错误也走创建路径，由 Create 返回错误）
		config = entity.Config{Namespace: namespace, Key: key, Value: value}
		if result := db.WithContext(ctx).Create(&config); result.Error != nil {
			return xError.NewError(ctx, xError.DatabaseError, "创建配置失败", false, result.Error)
		}
	}
	return nil
}

// BatchSet 批量保存配置。
func (r *ConfigRepository) BatchSet(ctx context.Context, namespace string, values map[string]string) *xError.Error {
	r.log.Info(ctx, "BatchSet - 批量保存配置")
	for key, value := range values {
		if xErr := r.Set(ctx, namespace, key, value); xErr != nil {
			return xErr
		}
	}
	return nil
}
