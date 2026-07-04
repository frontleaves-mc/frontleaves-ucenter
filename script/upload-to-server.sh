#!/bin/bash

set -e

# ============================================================================
# 服务器环境配置（可通过环境变量覆盖）
# ============================================================================

PROD_SERVER="${PROD_SERVER:-}"
PROD_USER="${PROD_USER:-root}"
PROD_PATH="${PROD_PATH:-/root/frontleaves-yggleaf}"

TEST_SERVER="${TEST_SERVER:-172.16.100.1}"
TEST_USER="${TEST_USER:-root}"
TEST_PATH="${TEST_PATH:-/root/frontleaves-yggleaf}"

DEPLOY_PORT="${DEPLOY_PORT:-22}"
DEPLOY_SSH_KEY="${DEPLOY_SSH_KEY:-}"

readonly BRAND_PRIMARY='#2d5a27'
readonly BRAND_SECONDARY='#5c8d89'
readonly BRAND_DARK='#1a1c18'
readonly FOREGROUND='#1f2623'
readonly MUTED='#64748b'
readonly DESTRUCTIVE='#ef4444'
readonly WARNING='#f59e0b'

# ============================================================================
# 辅助函数
# ============================================================================

# 左右分栏状态行
# 参数: 左侧文本 右侧状态 状态颜色
status_row() {
    gum join --horizontal \
        "$(gum style --width 60 "$1")" \
        "$(gum style --width 15 --align right --foreground "$3" "$2")"
}

# 步骤标题（左右分栏）
# 参数: 步骤编号 步骤标题 右侧状态 状态颜色
step_header() {
    local num="$1"
    local title="$2"
    local status="${3:-进行中}"
    local color="${4:-$BRAND_SECONDARY}"

    echo ""
    gum join --horizontal \
        "$(gum style --width 60 --bold --foreground "$BRAND_SECONDARY" "步骤 $num/5  $title")" \
        "$(gum style --width 15 --align right --foreground "$color" "$status")"
    gum style --foreground "$MUTED" "────────────────────────────────────────────────────────────────────────────"
}

# 信息面板（带边框）
# 参数: 标题 内容行...
info_panel() {
    local title="$1"
    shift

    gum style \
        --border normal --border-foreground "$BRAND_PRIMARY" \
        --padding "0 2" \
        "$(gum style --bold --foreground "$BRAND_PRIMARY" "$title")" \
        "" \
        "$@"
}

# 成功面板
success_panel() {
    local title="$1"
    shift

    gum style \
        --border normal --border-foreground "$BRAND_PRIMARY" \
        --padding "0 2" \
        "$(gum style --bold --foreground "$BRAND_PRIMARY" "$title")" \
        "" \
        "$@"
}

# 错误面板
error_panel() {
    local title="$1"
    shift

    gum style \
        --border normal --border-foreground "$DESTRUCTIVE" \
        --padding "0 2" \
        "$(gum style --bold --foreground "$DESTRUCTIVE" "$title")" \
        "" \
        "$@"
}

# Banner
print_banner() {
    gum style \
        --foreground "$BRAND_PRIMARY" --bold \
        "   __  __      __                __   _____                          " \
        "  / / / /___  / /___  ____ _____/ /  / ___/___  ______   _____  _____" \
        " / / / / __ \/ / __ \/ __ \/ __  /   \__ \/ _ \/ ___/ |/ / _ \/ ___/ " \
        "/ /_/ / /_/ / / /_/ / /_/ / /_/ /   ___/ /  __/ /   | |/ /  __/ /    " \
        "\____/ .___/_/\____/\__,_/\__,_/   /____/\___/_/    |___/\___/_/     " \
        "    /_/                                                              "
}

# ============================================================================
# 强制交互式选择
# ============================================================================

print_banner
echo ""
gum style --foreground "$BRAND_SECONDARY" "$(gum style --bold '[INFO]') 请选择目标服务器..."
echo ""

SERVER_TYPE=$(gum choose "默认" "测试服务器")

if [ "$SERVER_TYPE" = "默认" ]; then
    DEPLOY_SERVER="${DEPLOY_FRONTLEAVES_SERVER:-comprehensive.frontleaves.com}"
    DEPLOY_USER="${DEPLOY_USER:-root}"
    DEPLOY_PATH="${DEPLOY_PATH:-/root/frontleaves-yggleaf}"
    ENV_FILE=".env.prod"
    gum style --foreground "$BRAND_PRIMARY" "$(gum style --bold '[SUCCESS]') 已选择默认服务器: $DEPLOY_SERVER"
elif [ "$SERVER_TYPE" = "测试服务器" ]; then
    ENV_FILE=".env.test"
    if [ -n "$TEST_SERVER" ]; then
        DEPLOY_SERVER="$TEST_SERVER"
        DEPLOY_USER="$TEST_USER"
        DEPLOY_PATH="$TEST_PATH"
        gum style --foreground "$BRAND_PRIMARY" "$(gum style --bold '[SUCCESS]') 已选择测试服务器: $DEPLOY_SERVER"
    else
        DEPLOY_SERVER=$(gum input --placeholder "请输入服务器地址" --prompt "测试服务器地址: ")
        DEPLOY_USER=$(gum input --placeholder "root" --value "root" --prompt "SSH 用户: ")
        DEPLOY_PATH=$(gum input --placeholder "/root/frontleaves-yggleaf" --value "/root/frontleaves-yggleaf" --prompt "部署路径: ")
    fi
fi

# ============================================================================
# 步骤 1: 部署配置
# ============================================================================

print_banner
step_header "1" "部署配置" "📝 配置中" "$WARNING"

info_panel "配置详情" \
    "$(status_row "服务器地址" "$DEPLOY_SERVER" "✅ 已配置" "$BRAND_PRIMARY")" \
    "$(status_row "SSH 用户" "$DEPLOY_USER" "✅ 已配置" "$BRAND_PRIMARY")" \
    "$(status_row "SSH 端口" "$DEPLOY_PORT" "✅ 已配置" "$BRAND_PRIMARY")" \
    "$(status_row "部署路径" "$DEPLOY_PATH" "✅ 已配置" "$BRAND_PRIMARY")" \
    "$(status_row "环境文件" "$ENV_FILE" "✅ 已选择" "$BRAND_PRIMARY")"

# ============================================================================
# 步骤 2: 环境检查
# ============================================================================

step_header "2" "环境检查" "🔍 检查中" "$WARNING"

# 检查 rsync
if ! command -v rsync &> /dev/null; then
    error_panel "检查结果" \
        "$(status_row "rsync" "未安装" "❌ 失败" "$DESTRUCTIVE")"
    echo ""
    gum style --foreground "$BRAND_SECONDARY" "$(gum style --bold '[INFO]') 请先安装 rsync:"
    gum style --foreground "$MUTED" "  Ubuntu/Debian: apt-get install rsync"
    gum style --foreground "$MUTED" "  CentOS/RHEL:   yum install rsync"
    gum style --foreground "$MUTED" "  macOS:         brew install rsync"
    exit 1
fi

# 检查 ssh
if ! command -v ssh &> /dev/null; then
    error_panel "检查结果" \
        "$(status_row "ssh" "未安装" "❌ 失败" "$DESTRUCTIVE")"
    exit 1
fi

# 检查文件
FILES_TO_UPLOAD=("docker-compose.yml" "Makefile.run" "$ENV_FILE")
DIRS_TO_UPLOAD=("template")

file_status_rows=()
for file in "${FILES_TO_UPLOAD[@]}"; do
    if [ ! -f "$file" ]; then
        error_panel "检查结果" \
            "$(status_row "文件检查" "$file 不存在" "❌ 失败" "$DESTRUCTIVE")"
        exit 1
    fi
    file_status_rows+=("$(status_row "检查文件 $file" "找到" "✅ 通过" "$BRAND_PRIMARY")")
done

for dir in "${DIRS_TO_UPLOAD[@]}"; do
    if [ ! -d "$dir" ]; then
        error_panel "检查结果" \
            "$(status_row "目录检查" "$dir 不存在" "❌ 失败" "$DESTRUCTIVE")"
        exit 1
    fi
    file_status_rows+=("$(status_row "检查目录 $dir" "找到" "✅ 通过" "$BRAND_PRIMARY")")
done

success_panel "检查结果" \
    "$(status_row "rsync" "已安装" "✅ 通过" "$BRAND_PRIMARY")" \
    "$(status_row "ssh" "已安装" "✅ 通过" "$BRAND_PRIMARY")" \
    "${file_status_rows[@]}"

# ============================================================================
# 步骤 3: 测试服务器连接
# ============================================================================

step_header "3" "测试服务器连接" "🌐 连接中" "$WARNING"

SSH_CMD="ssh -p ${DEPLOY_PORT} -o ConnectTimeout=10 -o StrictHostKeyChecking=no"
if [ -n "$DEPLOY_SSH_KEY" ]; then
    SSH_CMD="$SSH_CMD -i $DEPLOY_SSH_KEY"
fi

info_panel "连接信息" \
    "$(status_row "目标服务器" "$DEPLOY_SERVER" "🎯" "$BRAND_SECONDARY")" \
    "$(status_row "SSH 端口" "$DEPLOY_PORT" "🔌" "$BRAND_SECONDARY")" \
    "$(status_row "连接用户" "$DEPLOY_USER" "👤" "$BRAND_SECONDARY")"

if gum spin --spinner dot --title "正在测试 SSH 连接..." -- \
    $SSH_CMD ${DEPLOY_USER}@${DEPLOY_SERVER} "echo > /dev/null 2>&1"; then
    success_panel "连接结果" \
        "$(status_row "SSH 连接测试" "成功" "✅ 通过" "$BRAND_PRIMARY")"
else
    error_panel "连接结果" \
        "$(status_row "SSH 连接测试" "失败" "❌ 失败" "$DESTRUCTIVE")"
    echo ""
    gum style --foreground "$BRAND_SECONDARY" "$(gum style --bold '[INFO]') 请检查:"
    gum style --foreground "$MUTED" "  1. 服务器地址是否正确"
    gum style --foreground "$MUTED" "  2. SSH 服务是否运行"
    gum style --foreground "$MUTED" "  3. 网络连接是否正常"
    gum style --foreground "$MUTED" "  4. SSH 密钥是否正确 (如使用密钥认证)"
    exit 1
fi

# ============================================================================
# 步骤 4: 上传文件到服务器
# ============================================================================

step_header "4" "上传文件到服务器" "📤 上传中" "$WARNING"

info_panel "上传信息" \
    "$(status_row "目标服务器" "$DEPLOY_SERVER" "🎯" "$BRAND_SECONDARY")" \
    "$(status_row "部署路径" "$DEPLOY_PATH" "📁" "$BRAND_SECONDARY")" \
    "$(status_row "文件数量" "4 个文件 + 1 个目录" "📦" "$BRAND_SECONDARY")"

echo ""

RSYNC_CMD="rsync -avz --progress"
RSYNC_CMD="$RSYNC_CMD --exclude='.env'"
RSYNC_CMD="$RSYNC_CMD --exclude='.git'"
RSYNC_CMD="$RSYNC_CMD --exclude='node_modules'"
RSYNC_CMD="$RSYNC_CMD --exclude='*.log'"
RSYNC_CMD="$RSYNC_CMD -e \"ssh -p ${DEPLOY_PORT} -o StrictHostKeyChecking=no\""

if [ -n "$DEPLOY_SSH_KEY" ]; then
    RSYNC_CMD="$RSYNC_CMD -i $DEPLOY_SSH_KEY"
fi

eval "$RSYNC_CMD docker-compose.yml Makefile.run $ENV_FILE template ${DEPLOY_USER}@${DEPLOY_SERVER}:${DEPLOY_PATH}/"

echo ""
success_panel "上传结果" \
    "$(status_row "文件上传" "完成" "✅ 成功" "$BRAND_PRIMARY")"

# ============================================================================
# 步骤 5: 配置服务器端文件
# ============================================================================

step_header "5" "配置服务器端文件" "⚙️  配置中" "$WARNING"

info_panel "配置操作" \
    "$(status_row "重命名 Makefile.run" "→ Makefile" "📝" "$BRAND_SECONDARY")" \
    "$(status_row "重命名 $ENV_FILE" "→ .env" "📝" "$BRAND_SECONDARY")" \
    "$(status_row "设置权限" "chmod +x Makefile" "🔐" "$BRAND_SECONDARY")"

if $SSH_CMD ${DEPLOY_USER}@${DEPLOY_SERVER} \
    "cd ${DEPLOY_PATH} && mv -f Makefile.run Makefile && chmod +x Makefile && mv -f $ENV_FILE .env && chmod 600 .env"; then
    success_panel "配置结果" \
        "$(status_row "远程文件配置" "成功" "✅ 完成" "$BRAND_PRIMARY")"
else
    error_panel "配置结果" \
        "$(status_row "远程文件配置" "失败" "❌ 失败" "$DESTRUCTIVE")"
    echo ""
    gum style --foreground "$WARNING" "$(gum style --bold '[WARN]') 文件重命名失败，请手动执行:"
    echo ""
    gum style --foreground "$MUTED" "  ssh ${DEPLOY_USER}@${DEPLOY_SERVER}"
    gum style --foreground "$MUTED" "  cd ${DEPLOY_PATH}"
    gum style --foreground "$MUTED" "  mv Makefile.run Makefile"
    gum style --foreground "$MUTED" "  chmod +x Makefile"
    gum style --foreground "$MUTED" "  mv $ENV_FILE .env"
    gum style --foreground "$MUTED" "  chmod 600 .env"
    echo ""
fi

# ============================================================================
# 完成面板
# ============================================================================

echo ""
gum style \
    --border double --border-foreground "$BRAND_PRIMARY" \
    --padding "1 3" --margin "1 0" \
    "$(gum style --bold --foreground "$BRAND_PRIMARY" '🎉 上传完成')" \
    "" \
    "$(gum style --foreground "$BRAND_SECONDARY" "文件已成功上传到服务器")" \
    "" \
    "$(gum style --foreground "$MUTED" "服务器: ${DEPLOY_SERVER}")" \
    "$(gum style --foreground "$MUTED" "路径:   ${DEPLOY_PATH}")" \
    "" \
    "$(gum style --bold --foreground "$BRAND_PRIMARY" "后续步骤:")" \
    "$(gum style --foreground "$MUTED" "1. SSH 登录服务器:")" \
    "$(gum style --foreground 240 "   ssh ${DEPLOY_USER}@${DEPLOY_SERVER}")" \
    "" \
    "$(gum style --foreground "$MUTED" "2. 进入部署目录:")" \
    "$(gum style --foreground 240 "   cd ${DEPLOY_PATH}")" \
    "" \
    "$(gum style --foreground "$MUTED" "3. 启动服务:")" \
    "$(gum style --foreground 240 "   make run")" \
    "" \
    "$(gum style --foreground "$MUTED" "4. 查看日志:")" \
    "$(gum style --foreground 240 "   make logs")"

echo ""
gum style --foreground "$BRAND_PRIMARY" "$(gum style --bold '[SUCCESS]') 所有操作已完成！"
echo ""
