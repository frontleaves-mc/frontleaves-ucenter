package repository

import (
	"context"
	"errors"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	"github.com/frontleaves-mc/frontleaves-yggleaf/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GameProfileRepo 游戏档案仓储，负责游戏档案数据访问。
type GameProfileRepo struct {
	db  *gorm.DB
	log *xLog.LogNamedLogger
}

// NewGameProfileRepo 初始化并返回 GameProfileRepo 实例。
func NewGameProfileRepo(db *gorm.DB) *GameProfileRepo {
	return &GameProfileRepo{
		db:  db,
		log: xLog.WithName(xLog.NamedREPO, "GameProfileRepo"),
	}
}

// Create 创建游戏档案记录。
func (r *GameProfileRepo) Create(ctx context.Context, tx *gorm.DB, profile *entity.GameProfile) (*entity.GameProfile, *xError.Error) {
	r.log.Info(ctx, "Create - 创建游戏档案")

	if err := r.pickDB(ctx, tx).Create(profile).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建游戏档案失败", true, err)
	}
	return profile, nil
}

// ExistsByUUID 检查指定 UUID 是否已存在。
func (r *GameProfileRepo) ExistsByUUID(ctx context.Context, tx *gorm.DB, uuid string) (bool, *xError.Error) {
	r.log.Info(ctx, "ExistsByUUID - 检查 UUID 是否存在")

	var count int64
	if err := r.pickDB(ctx, tx).Model(&entity.GameProfile{}).Where("uuid = ?", uuid).Count(&count).Error; err != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "查询游戏档案 UUID 失败", true, err)
	}
	return count > 0, nil
}

// ExistsByNameExceptID 检查用户名是否存在（排除指定档案 ID）。
func (r *GameProfileRepo) ExistsByNameExceptID(ctx context.Context, tx *gorm.DB, name string, profileID xSnowflake.SnowflakeID) (bool, *xError.Error) {
	r.log.Info(ctx, "ExistsByNameExceptID - 检查用户名是否存在")

	var count int64
	query := r.pickDB(ctx, tx).Model(&entity.GameProfile{}).Where("name = ?", name)
	if !profileID.IsZero() {
		query = query.Where("id <> ?", profileID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "查询游戏档案名称失败", true, err)
	}
	return count > 0, nil
}

// GetByIDAndUserID 根据档案 ID 与用户 ID 查询游戏档案。
func (r *GameProfileRepo) GetByIDAndUserID(ctx context.Context, tx *gorm.DB, profileID xSnowflake.SnowflakeID, userID xSnowflake.SnowflakeID, forUpdate bool) (*entity.GameProfile, bool, *xError.Error) {
	r.log.Info(ctx, "GetByIDAndUserID - 根据档案 ID 与用户 ID 获取游戏档案")

	query := r.pickDB(ctx, tx).Model(&entity.GameProfile{}).Where("id = ? AND user_id = ?", profileID, userID)
	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var profile entity.GameProfile
	err := query.First(&profile).Error
	if err == nil {
		return &profile, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询游戏档案失败", true, err)
}

// GetByID 根据档案 ID 查询游戏档案。
func (r *GameProfileRepo) GetByID(ctx context.Context, tx *gorm.DB, profileID xSnowflake.SnowflakeID) (*entity.GameProfile, bool, *xError.Error) {
	r.log.Info(ctx, "GetByID - 根据档案 ID 获取游戏档案")

	var profile entity.GameProfile
	err := r.pickDB(ctx, tx).Model(&entity.GameProfile{}).Where("id = ?", profileID).First(&profile).Error
	if err == nil {
		return &profile, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询游戏档案失败", true, err)
}

// GetDetailByID 根据档案 ID 与用户 ID 查询游戏档案详情（含关联皮肤和披风）。
func (r *GameProfileRepo) GetDetailByID(ctx context.Context, tx *gorm.DB, profileID xSnowflake.SnowflakeID, userID xSnowflake.SnowflakeID) (*entity.GameProfile, bool, *xError.Error) {
	r.log.Info(ctx, "GetDetailByID - 获取游戏档案详情")

	var profile entity.GameProfile
	err := r.pickDB(ctx, tx).
		Model(&entity.GameProfile{}).
		Preload("SkinLibrary").
		Preload("CapeLibrary").
		Where("id = ? AND user_id = ?", profileID, userID).
		First(&profile).Error
	if err == nil {
		return &profile, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询游戏档案详情失败", true, err)
}

// ListByUserID 根据用户 ID 查询其所有游戏档案。
func (r *GameProfileRepo) ListByUserID(ctx context.Context, tx *gorm.DB, userID xSnowflake.SnowflakeID) ([]entity.GameProfile, *xError.Error) {
	r.log.Info(ctx, "ListByUserID - 根据用户 ID 获取游戏档案列表")

	var profiles []entity.GameProfile
	if err := r.pickDB(ctx, tx).Model(&entity.GameProfile{}).Where("user_id = ?", userID).Order("created_at ASC").Find(&profiles).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询游戏档案列表失败", true, err)
	}
	return profiles, nil
}

// UpdateName 更新指定档案的用户名。
func (r *GameProfileRepo) UpdateName(ctx context.Context, tx *gorm.DB, profileID xSnowflake.SnowflakeID, name string) (*entity.GameProfile, *xError.Error) {
	r.log.Info(ctx, "UpdateName - 更新游戏档案用户名")

	if err := r.pickDB(ctx, tx).Model(&entity.GameProfile{}).Where("id = ?", profileID).Update("name", name).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新游戏档案名称失败", true, err)
	}

	updatedProfile, found, xErr := r.GetByID(ctx, tx, profileID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.ResourceNotFound, "游戏档案不存在", true)
	}
	return updatedProfile, nil
}

// UpdateSkinLibraryID 更新指定档案的关联皮肤库 ID。
func (r *GameProfileRepo) UpdateSkinLibraryID(ctx context.Context, tx *gorm.DB, profileID xSnowflake.SnowflakeID, skinLibraryID *xSnowflake.SnowflakeID) (*entity.GameProfile, *xError.Error) {
	r.log.Info(ctx, "UpdateSkinLibraryID - 更新游戏档案关联皮肤")

	if err := r.pickDB(ctx, tx).Model(&entity.GameProfile{}).Where("id = ?", profileID).Update("skin_library_id", skinLibraryID).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新游戏档案皮肤失败", true, err)
	}

	updatedProfile, found, xErr := r.GetByID(ctx, tx, profileID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.ResourceNotFound, "游戏档案不存在", true)
	}
	return updatedProfile, nil
}

// UpdateCapeLibraryID 更新指定档案的关联披风库 ID。
func (r *GameProfileRepo) UpdateCapeLibraryID(ctx context.Context, tx *gorm.DB, profileID xSnowflake.SnowflakeID, capeLibraryID *xSnowflake.SnowflakeID) (*entity.GameProfile, *xError.Error) {
	r.log.Info(ctx, "UpdateCapeLibraryID - 更新游戏档案关联披风")

	if err := r.pickDB(ctx, tx).Model(&entity.GameProfile{}).Where("id = ?", profileID).Update("cape_library_id", capeLibraryID).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新游戏档案披风失败", true, err)
	}

	updatedProfile, found, xErr := r.GetByID(ctx, tx, profileID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.ResourceNotFound, "游戏档案不存在", true)
	}
	return updatedProfile, nil
}

// AdminGameProfileFilter 管理员游戏档案列表筛选条件。
type AdminGameProfileFilter struct {
	Name     *string                    // 档案名模糊匹配
	UserIDs  []xSnowflake.SnowflakeID   // 用户 ID 列表（由 logic 层按用户名查询后传入）
	UUID     *string                    // UUID 精确匹配
}

// ListAdmin 分页查询游戏档案列表（管理员用）。
//
// 支持按档案名模糊匹配、用户 ID 列表过滤和 UUID 精确筛选，
// 默认按更新时间降序排列，预加载关联用户信息。
// 用户名筛选由 logic 层先查 UserRepo 获取 user IDs 后传入，避免跨表 JOIN。
func (r *GameProfileRepo) ListAdmin(ctx context.Context, tx *gorm.DB, page, pageSize int, filter AdminGameProfileFilter) ([]entity.GameProfile, int64, *xError.Error) {
	r.log.Info(ctx, "ListAdmin - 管理员分页查询游戏档案列表")

	query := r.pickDB(ctx, tx).Model(&entity.GameProfile{})

	if filter.Name != nil && *filter.Name != "" && len(filter.UserIDs) > 0 {
		query = query.Where("name ILIKE ? OR user_id IN ?", "%"+*filter.Name+"%", filter.UserIDs)
	} else if filter.Name != nil && *filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
	} else if len(filter.UserIDs) > 0 {
		query = query.Where("user_id IN ?", filter.UserIDs)
	}
	if filter.UUID != nil && *filter.UUID != "" {
		query = query.Where("uuid = ?", *filter.UUID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "查询游戏档案总数失败", true, err)
	}

	var profiles []entity.GameProfile
	offset := (page - 1) * pageSize
	if err := query.Preload("User").
		Order("updated_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&profiles).Error; err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "查询游戏档案列表失败", true, err)
	}

	return profiles, total, nil
}

// GetDetailByIDAdmin 根据档案 ID 查询游戏档案详情（管理员用，无用户归属限制）。
//
// 预加载关联皮肤库、披风库和用户信息。
func (r *GameProfileRepo) GetDetailByIDAdmin(ctx context.Context, tx *gorm.DB, profileID xSnowflake.SnowflakeID) (*entity.GameProfile, bool, *xError.Error) {
	r.log.Info(ctx, "GetDetailByIDAdmin - 管理员获取游戏档案详情")

	var profile entity.GameProfile
	err := r.pickDB(ctx, tx).
		Model(&entity.GameProfile{}).
		Preload("SkinLibrary").
		Preload("CapeLibrary").
		Preload("User").
		Where("id = ?", profileID).
		First(&profile).Error
	if err == nil {
		return &profile, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询游戏档案详情失败", true, err)
}

func (r *GameProfileRepo) pickDB(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}
