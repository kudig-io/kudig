# 测试辅助函数
# 供 bats 测试用例使用的公共函数

# 加载 bats 支持库（如果可用）
if [ -f "/usr/lib/bats-support/load.bash" ]; then
    load '/usr/lib/bats-support/load.bash'
fi

if [ -f "/usr/lib/bats-assert/load.bash" ]; then
    load '/usr/lib/bats-assert/load.bash'
fi

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "$BATS_TEST_FILENAME")/.." && pwd)"

# 被测脚本路径
KUDIG_SH="$PROJECT_ROOT/kudig.sh"

# 测试数据目录
FIXTURES_DIR="$PROJECT_ROOT/tests/fixtures"

# 临时目录
TEST_TEMP_DIR=""

# Setup 函数 - 在每个测试前运行
setup() {
    # 创建临时目录
    TEST_TEMP_DIR="$(mktemp -d)"
    
    # 确保脚本存在且可执行
    if [ ! -f "$KUDIG_SH" ]; then
        skip "kudig.sh not found at $KUDIG_SH"
    fi
    
    chmod +x "$KUDIG_SH"
}

# Teardown 函数 - 在每个测试后运行
teardown() {
    # 清理临时目录
    if [ -n "$TEST_TEMP_DIR" ] && [ -d "$TEST_TEMP_DIR" ]; then
        rm -rf "$TEST_TEMP_DIR"
    fi
}

# 创建模拟诊断目录
create_mock_diagnose_dir() {
    local dir_name="${1:-diagnose_mock}"
    local diagnose_dir="$TEST_TEMP_DIR/$dir_name"
    
    mkdir -p "$diagnose_dir"
    mkdir -p "$diagnose_dir/daemon_status"
    mkdir -p "$diagnose_dir/logs"
    
    # 创建基本文件
    touch "$diagnose_dir/system_info"
    touch "$diagnose_dir/system_status"
    touch "$diagnose_dir/service_status"
    touch "$diagnose_dir/memory_info"
    touch "$diagnose_dir/network_info"
    
    echo "$diagnose_dir"
}

# 创建包含 CPU 负载数据的模拟文件
create_system_status_with_load() {
    local diagnose_dir="$1"
    local load_15min="${2:-1.0}"
    
    cat > "$diagnose_dir/system_status" <<EOF
----- run uptime -----
 12:00:00 up 10 days,  1:23,  2 users,  load average: 0.5, 1.0, $load_15min
----- End of uptime -----
EOF
}

# 创建包含内存信息的模拟文件
create_memory_info() {
    local diagnose_dir="$1"
    local total_kb="${2:-8000000}"
    local avail_kb="${3:-4000000}"
    
    cat > "$diagnose_dir/memory_info" <<EOF
MemTotal:        $total_kb kB
MemFree:         2000000 kB
MemAvailable:    $avail_kb kB
Buffers:         100000 kB
Cached:          1000000 kB
EOF
}

# 创建包含磁盘信息的模拟文件
create_system_status_with_disk() {
    local diagnose_dir="$1"
    local usage_percent="${2:-50}"
    
    cat > "$diagnose_dir/system_status" <<EOF
----- run df -h -----
Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1        50G   ${usage_percent}G   20G  ${usage_percent}% /
/dev/sda2       100G   30G   70G  30% /data
----- End of df -----
EOF
}

# 创建 kubelet 服务状态文件
create_kubelet_status() {
    local diagnose_dir="$1"
    local status="${2:-running}"
    
    mkdir -p "$diagnose_dir/daemon_status"
    
    if [ "$status" = "running" ]; then
        cat > "$diagnose_dir/daemon_status/kubelet_status" <<EOF
● kubelet.service - Kubernetes Kubelet
   Loaded: loaded (/usr/lib/systemd/system/kubelet.service; enabled)
   Active: active (running) since Mon 2024-01-01 00:00:00 UTC; 1 day ago
EOF
    elif [ "$status" = "failed" ]; then
        cat > "$diagnose_dir/daemon_status/kubelet_status" <<EOF
● kubelet.service - Kubernetes Kubelet
   Loaded: loaded (/usr/lib/systemd/system/kubelet.service; enabled)
   Active: failed (Result: exit-code) since Mon 2024-01-01 00:00:00 UTC; 1 hour ago
EOF
    else
        cat > "$diagnose_dir/daemon_status/kubelet_status" <<EOF
● kubelet.service - Kubernetes Kubelet
   Loaded: loaded (/usr/lib/systemd/system/kubelet.service; enabled)
   Active: inactive (dead)
EOF
    fi
}

# 运行 kudig.sh 并捕获输出
run_kudig() {
    local diagnose_dir="$1"
    shift
    
    run bash "$KUDIG_SH" "$@" "$diagnose_dir"
}

# 运行 kudig.sh JSON 模式
run_kudig_json() {
    local diagnose_dir="$1"
    
    run bash "$KUDIG_SH" --json "$diagnose_dir"
}

# 验证 JSON 输出格式
validate_json_output() {
    local json_output="$1"
    
    # 检查是否为有效 JSON
    echo "$json_output" | python3 -m json.tool > /dev/null 2>&1
    return $?
}

# 从 JSON 输出中提取字段
extract_json_field() {
    local json_output="$1"
    local field="$2"
    
    echo "$json_output" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('$field', ''))"
}

# 计算 JSON 中异常数量
count_anomalies_in_json() {
    local json_output="$1"
    
    echo "$json_output" | python3 -c "import sys, json; data=json.load(sys.stdin); print(len(data.get('anomalies', [])))"
}

# 检查输出中是否包含特定异常
assert_anomaly_exists() {
    local output="$1"
    local anomaly_code="$2"
    
    echo "$output" | grep -q "$anomaly_code"
}

# 断言退出码
assert_exit_code() {
    local expected="$1"
    local actual="$status"
    
    if [ "$actual" != "$expected" ]; then
        echo "Expected exit code $expected, got $actual"
        return 1
    fi
}
