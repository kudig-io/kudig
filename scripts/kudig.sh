#!/usr/bin/env bash

################################################################################
# kudig - Kubernetes节点诊断日志分析工具
# 
# 功能：分析 diagnose_k8s.sh 收集的诊断日志，识别异常并输出中英文报告
# 
# 使用方法:
#   ./kudig <diagnose_dir>              # 分析指定诊断目录
#   ./kudig --json <diagnose_dir>       # 输出JSON格式
#   ./kudig --verbose <diagnose_dir>    # 详细模式
#   ./kudig --help                      # 显示帮助信息
#
# 示例:
#   ./kudig /tmp/diagnose_1702468800
#   ./kudig --json /tmp/diagnose_1702468800 > report.json
#
# 作者: kudig Team
# 版本: 1.2.0
################################################################################

set -euo pipefail

# ============================================================================
# 全局变量定义
# ============================================================================

VERSION="1.2.0"
SCRIPT_NAME=$(basename "$0")
DIAGNOSE_DIR=""
OUTPUT_FORMAT="text"  # text, json
VERBOSE=false
OUTPUT_FILE=""

# 异常数组 - 格式: "严重级别|中文名称|英文标识|详情|位置"
declare -a ANOMALIES=()

# 颜色定义（用于终端输出）
if [[ -t 1 ]]; then
    RED='\033[0;31m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    GREEN='\033[0;32m'
    NC='\033[0m' # No Color
else
    RED=''
    YELLOW=''
    BLUE=''
    GREEN=''
    NC=''
fi

# ============================================================================
# 工具函数
# ============================================================================

# 安全获取匹配计数
safe_count() {
    local file="$1"
    local pattern="$2"
    local count=$(grep -o "$pattern" "$file" 2>/dev/null | wc -l 2>/dev/null || echo "0")
    count=$(echo "$count" | tr -d '\n\r' | sed 's/[^0-9]//g' || echo "0")
    [[ -z "$count" ]] && count="0"
    echo "$count"
}

# 检查文件是否存在并可读
check_file() {
    local file="$1"
    [[ -f "$file" && -r "$file" ]]
}

# 输出信息
log_info() {
    if [[ "$VERBOSE" == true ]]; then
        echo -e "${BLUE}[INFO]${NC} $*" >&2
    fi
}

# 输出警告
log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $*" >&2
}

# 输出错误
log_error() {
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

# 输出成功
log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*" >&2
}

# 添加异常
add_anomaly() {
    local severity="$1"
    local cn_name="$2"
    local en_name="$3"
    local details="$4"
    local location="$5"
    
    # 检查是否已存在相同的异常
    local exists=false
    for anomaly in "${ANOMALIES[@]}"; do
        if [[ "$anomaly" == *"|$en_name|"* ]]; then
            exists=true
            break
        fi
    done
    
    if [[ "$exists" == false ]]; then
        ANOMALIES+=("$severity|$cn_name|$en_name|$details|$location")
    fi
}

# 解析命令行参数
parse_arguments() {
    local args=()
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                echo "kudig version $VERSION"
                exit 0
                ;;
            --verbose)
                VERBOSE=true
                shift
                ;;
            --json)
                OUTPUT_FORMAT="json"
                shift
                ;;
            -o|--output)
                if [[ $# -gt 1 ]]; then
                    OUTPUT_FILE="$2"
                    shift 2
                else
                    log_error "缺少输出文件参数"
                    show_help
                    exit 1
                fi
                ;;
            -*)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
            *)
                args+=("$1")
                shift
                ;;
        esac
    done
    
    if [[ ${#args[@]} -ne 1 ]]; then
        log_error "请指定诊断目录"
        show_help
        exit 1
    fi
    
    DIAGNOSE_DIR="${args[0]}"
    
    if [[ ! -d "$DIAGNOSE_DIR" ]]; then
        log_error "诊断目录不存在: $DIAGNOSE_DIR"
        exit 1
    fi
}

# 显示帮助信息
show_help() {
    cat << EOF
用法: $SCRIPT_NAME [选项] <诊断目录>

Kubernetes节点诊断日志分析工具
分析 diagnose_k8s.sh 收集的诊断数据，识别异常并生成报告

选项:
    -h, --help              显示此帮助信息
    -v, --version           显示版本信息
    --verbose               详细输出模式
    --json                  输出JSON格式
    -o, --output <文件>     保存报告到指定文件

参数:
    <诊断目录>              diagnose_k8s.sh 生成的诊断目录路径

示例:
    $SCRIPT_NAME /tmp/diagnose_1702468800
    $SCRIPT_NAME --json /tmp/diagnose_1702468800 > report.json
    $SCRIPT_NAME --verbose -o report.txt /tmp/diagnose_1702468800
EOF
}

# ============================================================================
# 检测规则实现
# ============================================================================

# 检测系统资源
check_system_resources() {
    log_info "检查系统资源..."
    
    local system_status="$DIAGNOSE_DIR/system_status"
    local memory_info="$DIAGNOSE_DIR/memory_info"
    
    # 检查CPU负载
    if check_file "$system_status"; then
        local load=$(grep -E "load average:" "$system_status" | awk '{print $12}' | sed 's/,//' 2>/dev/null || echo "0")
        local cpu_cores=$(grep -c "processor" /proc/cpuinfo 2>/dev/null || echo "4")
        local load_threshold=$((cpu_cores * 2))
        local load_critical=$((cpu_cores * 4))
        
        if [[ $(echo "$load > $load_critical" | bc -l 2>/dev/null || echo "0") -eq 1 ]]; then
            add_anomaly "严重" "CPU负载过高" "CPU_LOAD_HIGH" "15分钟平均负载: $load, CPU核心: $cpu_cores" "$system_status"
            echo -e "  [✗] CPU负载: 严重过高 (15min负载: $load, CPU核心: $cpu_cores)" >&2
        elif [[ $(echo "$load > $load_threshold" | bc -l 2>/dev/null || echo "0") -eq 1 ]]; then
            add_anomaly "警告" "CPU负载过高" "CPU_LOAD_HIGH" "15分钟平均负载: $load, CPU核心: $cpu_cores" "$system_status"
            echo -e "  [!] CPU负载: 过高 (15min负载: $load, CPU核心: $cpu_cores)" >&2
        else
            echo -e "  [✓] CPU负载: 正常 (15min负载: $load, CPU核心: $cpu_cores)" >&2
        fi
    else
        echo -e "  [-] CPU负载: 文件不存在 ($system_status)" >&2
    fi
    
    # 检查内存使用
    if check_file "$memory_info"; then
        local mem_total=$(grep "MemTotal:" "$memory_info" | awk '{print $2}' 2>/dev/null || echo "0")
        local mem_free=$(grep "MemFree:" "$memory_info" | awk '{print $2}' 2>/dev/null || echo "0")
        local mem_available=$(grep "MemAvailable:" "$memory_info" | awk '{print $2}' 2>/dev/null || echo "0")
        
        if [[ "$mem_total" != "0" && "$mem_available" != "0" ]]; then
            local mem_used=$((mem_total - mem_available))
            local mem_usage=$((mem_used * 100 / mem_total))
            
            if [[ $mem_usage -gt 95 ]]; then
                add_anomaly "严重" "内存使用过高" "MEMORY_USAGE_HIGH" "使用率: ${mem_usage}%" "$memory_info"
                echo -e "  [✗] 内存使用: 严重过高 (使用率: ${mem_usage}%)" >&2
            elif [[ $mem_usage -gt 85 ]]; then
                add_anomaly "警告" "内存使用过高" "MEMORY_USAGE_HIGH" "使用率: ${mem_usage}%" "$memory_info"
                echo -e "  [!] 内存使用: 过高 (使用率: ${mem_usage}%)" >&2
            else
                echo -e "  [✓] 内存使用: 正常 (使用率: ${mem_usage}%)" >&2
            fi
        else
            echo -e "  [-] 内存使用: 无法获取信息" >&2
        fi
    else
        echo -e "  [-] 内存使用: 文件不存在 ($memory_info)" >&2
    fi
    
    # 检查磁盘空间
    if check_file "$system_status"; then
        local disk_issues=false
        while IFS= read -r line; do
            if [[ "$line" =~ ^/dev/ ]]; then
                local usage=$(echo "$line" | awk '{print $5}' | sed 's/%//' 2>/dev/null || echo "0")
                local mount=$(echo "$line" | awk '{print $6}' 2>/dev/null || echo "/")
                
                if [[ $usage -gt 95 ]]; then
                    add_anomaly "严重" "磁盘空间严重不足" "DISK_SPACE_CRITICAL" "挂载点 $mount 使用率 ${usage}%" "$system_status"
                    echo -e "  [✗] 磁盘空间 [$mount]: 严重不足 (使用率: ${usage}%)" >&2
                    disk_issues=true
                elif [[ $usage -gt 90 ]]; then
                    add_anomaly "警告" "磁盘空间不足" "DISK_SPACE_LOW" "挂载点 $mount 使用率 ${usage}%" "$system_status"
                    echo -e "  [!] 磁盘空间 [$mount]: 不足 (使用率: ${usage}%)" >&2
                    disk_issues=true
                fi
            fi
        done < <(grep -A 10 "----- run df -h -----" "$system_status" | grep -v "-----" | grep -v "Filesystem")
        
        if [[ "$disk_issues" == false ]]; then
            echo -e "  [✓] 磁盘空间: 正常 (所有挂载点使用率<90%)" >&2
        fi
    else
        echo -e "  [-] 磁盘空间: 文件不存在 ($system_status)" >&2
    fi
    
    # 检查文件句柄
    if check_file "$system_status"; then
        local fh_count=$(grep "open files" "$system_status" | awk '{print $3}' 2>/dev/null || echo "0")
        if [[ $fh_count -gt 50000 ]]; then
            add_anomaly "警告" "文件句柄数过高" "FILE_HANDLES_HIGH" "当前文件句柄数: $fh_count" "$system_status"
            echo -e "  [!] 文件句柄: 过高 (当前: $fh_count)" >&2
        else
            echo -e "  [✓] 文件句柄: 正常 (当前: $fh_count)" >&2
        fi
    else
        echo -e "  [-] 文件句柄: 文件不存在 ($system_status)" >&2
    fi
}

# 检测进程与服务
check_processes_services() {
    log_info "检查进程与服务..."
    
    local daemon_status="$DIAGNOSE_DIR/daemon_status"
    local kubelet_status="$daemon_status/kubelet_status"
    local docker_status="$daemon_status/docker_status"
    local containerd_status="$daemon_status/containerd_status"
    
    # 检查Kubelet服务
    if check_file "$kubelet_status"; then
        if grep -q "active (running)" "$kubelet_status"; then
            echo -e "  [✓] Kubelet服务: running" >&2
        elif grep -q "failed" "$kubelet_status"; then
            add_anomaly "严重" "Kubelet服务未运行" "KUBELET_SERVICE_DOWN" "kubelet.service状态为failed" "$kubelet_status"
            echo -e "  [✗] Kubelet服务: failed" >&2
        else
            add_anomaly "警告" "Kubelet服务未运行" "KUBELET_SERVICE_DOWN" "kubelet.service状态异常" "$kubelet_status"
            echo -e "  [!] Kubelet服务: 未运行" >&2
        fi
    else
        echo -e "  [-] Kubelet服务: 状态未知" >&2
    fi
    
    # 检查容器运行时
    local docker_state="unknown"
    local containerd_state="unknown"
    
    if check_file "$docker_status"; then
        if grep -q "active (running)" "$docker_status"; then
            docker_state="running"
        elif grep -q "failed" "$docker_status"; then
            docker_state="failed"
        fi
    fi
    
    if check_file "$containerd_status"; then
        if grep -q "active (running)" "$containerd_status"; then
            containerd_state="running"
        elif grep -q "failed" "$containerd_status"; then
            containerd_state="failed"
        fi
    fi
    
    if [[ "$docker_state" == "failed" || "$containerd_state" == "failed" ]]; then
        add_anomaly "严重" "容器运行时异常" "CONTAINER_RUNTIME_ERROR" "docker=$docker_state, containerd=$containerd_state" "$daemon_status"
        echo -e "  [✗] 容器运行时: 异常 (docker=$docker_state, containerd=$containerd_state)" >&2
    elif [[ "$docker_state" == "running" || "$containerd_state" == "running" ]]; then
        echo -e "  [✓] 容器运行时: docker=$docker_state, containerd=$containerd_state" >&2
    else
        echo -e "  [-] 容器运行时: 状态未知" >&2
    fi
}

# 检测网络状态
check_network() {
    log_info "检查网络状态..."
    
    local network_info="$DIAGNOSE_DIR/network_info"
    local system_status="$DIAGNOSE_DIR/system_status"
    
    # 检查连接跟踪表
    if check_file "$system_status"; then
        local conntrack=$(grep "nf_conntrack_count" "$system_status" 2>/dev/null || echo "0")
        local conntrack_max=$(grep "nf_conntrack_max" "$system_status" 2>/dev/null || echo "262144")
        
        if [[ -n "$conntrack" && -n "$conntrack_max" ]]; then
            local conntrack_usage=$(echo "$conntrack * 100 / $conntrack_max" | bc 2>/dev/null || echo "0")
            
            if [[ $conntrack_usage -gt 95 ]]; then
                add_anomaly "严重" "连接跟踪表使用率过高" "CONNTRACK_FULL" "使用率: ${conntrack_usage}%" "$system_status"
                echo -e "  [✗] 连接跟踪表: 严重过高 (使用率: ${conntrack_usage}%)" >&2
            elif [[ $conntrack_usage -gt 80 ]]; then
                add_anomaly "警告" "连接跟踪表使用率过高" "CONNTRACK_HIGH" "使用率: ${conntrack_usage}%" "$system_status"
                echo -e "  [!] 连接跟踪表: 过高 (使用率: ${conntrack_usage}%)" >&2
            else
                echo -e "  [✓] 连接跟踪表: 正常 (使用率: ${conntrack_usage}%)" >&2
            fi
        else
            echo -e "  [-] 连接跟踪表: 无法获取信息" >&2
        fi
    else
        echo -e "  [-] 连接跟踪表: 文件不存在 ($system_status)" >&2
    fi
    
    # 检查网卡状态
    if check_file "$network_info"; then
        local down_interfaces=$(grep -B2 "state DOWN" "$network_info" | grep -E "^[0-9]+: " | awk -F: '{print $2}' | tr -d ' ' | grep -v lo | grep -v veth 2>/dev/null || echo "")
        
        if [[ -n "$down_interfaces" ]]; then
            add_anomaly "警告" "网卡接口down" "NETWORK_INTERFACE_DOWN" "以下网卡处于down状态: $down_interfaces" "$network_info"
            echo -e "  [!] 网卡状态: 部分down ($down_interfaces)" >&2
        else
            echo -e "  [✓] 网卡状态: 正常" >&2
        fi
    else
        echo -e "  [-] 网卡状态: 文件不存在 ($network_info)" >&2
    fi
    
    # 检查默认路由
    if check_file "$network_info"; then
        if grep -q "default via" "$network_info"; then
            echo -e "  [✓] 默认路由: 已配置" >&2
        else
            add_anomaly "警告" "默认路由未配置" "DEFAULT_ROUTE_MISSING" "未找到默认路由" "$network_info"
            echo -e "  [!] 默认路由: 未配置" >&2
        fi
    else
        echo -e "  [-] 默认路由: 文件不存在 ($network_info)" >&2
    fi
}

# 检测内核状态
check_kernel() {
    log_info "检查内核状态..."
    
    local dmesg_log="$DIAGNOSE_DIR/logs/dmesg.log"
    local messages_log="$DIAGNOSE_DIR/logs/messages"
    
    # 检查内核Panic
    if check_file "$dmesg_log"; then
        local panic_count=$(safe_count "$dmesg_log" "kernel panic")
        if [[ $panic_count -gt 0 ]]; then
            add_anomaly "严重" "内核Panic" "KERNEL_PANIC" "检测到 $panic_count 次kernel panic" "$dmesg_log"
            echo -e "  [✗] 内核Panic: 检测到 $panic_count 次" >&2
        else
            echo -e "  [✓] 内核Panic: 未发现" >&2
        fi
    else
        echo -e "  [-] 内核Panic: dmesg.log不存在" >&2
    fi
    
    # 检查OOM Killer
    if check_file "$dmesg_log"; then
        local oom_count=$(safe_count "$dmesg_log" "Out of memory: Kill process")
        if [[ $oom_count -gt 0 ]]; then
            add_anomaly "警告" "OOM Killer触发" "OOM_KILLER_TRIGGERED" "检测到 $oom_count 次OOM事件" "$dmesg_log"
            echo -e "  [!] OOM Killer: 触发 $oom_count 次" >&2
        else
            echo -e "  [✓] OOM Killer: 未触发" >&2
        fi
    else
        echo -e "  [-] OOM Killer: dmesg.log不存在" >&2
    fi
    
    # 检查messages日志中的OOM
    if check_file "$messages_log"; then
        local messages_oom_count=$(safe_count "$messages_log" "Out of memory")
        if [[ $messages_oom_count -gt 0 ]]; then
            add_anomaly "警告" "系统日志OOM" "SYSTEM_LOG_OOM" "检测到 $messages_oom_count 次OOM事件" "$messages_log"
            echo -e "  [!] messages日志OOM: 检测到 $messages_oom_count 次" >&2
        else
            echo -e "  [✓] messages日志OOM: 未发现" >&2
        fi
    else
        echo -e "  [-] messages日志: 文件不存在" >&2
    fi
}

# 检测容器运行时
check_container_runtime() {
    log_info "检查容器运行时..."
    
    local docker_log="$DIAGNOSE_DIR/logs/docker.log"
    local containerd_log="$DIAGNOSE_DIR/logs/containerd.log"
    
    # 检查Docker日志
    if check_file "$docker_log"; then
        local docker_error_count=$(safe_count "$docker_log" "error")
        if [[ $docker_error_count -gt 100 ]]; then
            add_anomaly "警告" "Docker错误过多" "DOCKER_ERRORS" "Docker日志中错误数量: $docker_error_count" "$docker_log"
            echo -e "  [!] Docker日志: 错误过多 ($docker_error_count 条)" >&2
        else
            echo -e "  [✓] Docker日志: 正常" >&2
        fi
    else
        echo -e "  [-] Docker日志: 文件不存在" >&2
    fi
    
    # 检查Containerd日志
    if check_file "$containerd_log"; then
        local containerd_error_count=$(safe_count "$containerd_log" "error")
        if [[ $containerd_error_count -gt 100 ]]; then
            add_anomaly "警告" "Containerd错误过多" "CONTAINERD_ERRORS" "Containerd日志中错误数量: $containerd_error_count" "$containerd_log"
            echo -e "  [!] Containerd日志: 错误过多 ($containerd_error_count 条)" >&2
        else
            echo -e "  [✓] Containerd日志: 正常" >&2
        fi
    else
        echo -e "  [-] Containerd日志: 文件不存在" >&2
    fi
}

# 检测Kubernetes组件
check_kubernetes() {
    log_info "检查Kubernetes组件..."
    
    local kubelet_log="$DIAGNOSE_DIR/logs/kubelet.log"
    
    # 检查PLEG状态
    if check_file "$kubelet_log"; then
        local pleg_count=$(safe_count "$kubelet_log" "PLEG is not healthy")
        if [[ $pleg_count -gt 0 ]]; then
            add_anomaly "警告" "PLEG状态异常" "PLEG_UNHEALTHY" "检测到 $pleg_count 次PLEG异常" "$kubelet_log"
            echo -e "  [!] PLEG状态: 异常" >&2
        else
            echo -e "  [✓] PLEG状态: 健康" >&2
        fi
    else
        echo -e "  [-] PLEG状态: kubelet.log不存在" >&2
    fi
    
    # 检查CNI网络插件
    if check_file "$kubelet_log"; then
        local cni_count=$(safe_count "$kubelet_log" "network plugin not ready")
        if [[ $cni_count -gt 0 ]]; then
            add_anomaly "警告" "CNI网络插件异常" "CNI_PLUGIN_ERROR" "检测到 $cni_count 次CNI错误" "$kubelet_log"
            echo -e "  [!] CNI网络插件: 异常" >&2
        else
            echo -e "  [✓] CNI网络插件: 正常" >&2
        fi
    else
        echo -e "  [-] CNI网络插件: kubelet.log不存在" >&2
    fi
    
    # 检查API Server连接
    if check_file "$kubelet_log"; then
        local api_count=$(safe_count "$kubelet_log" "Failed to connect to API server")
        if [[ $api_count -gt 5 ]]; then
            add_anomaly "警告" "API Server连接失败" "API_SERVER_CONNECTION_ERROR" "检测到 $api_count 次连接失败" "$kubelet_log"
            echo -e "  [!] API Server连接: 失败 $api_count 次" >&2
        else
            echo -e "  [✓] API Server连接: 正常" >&2
        fi
    else
        echo -e "  [-] API Server连接: kubelet.log不存在" >&2
    fi
}

# 检测时间同步
check_time_sync() {
    log_info "检查时间同步..."
    
    local service_status="$DIAGNOSE_DIR/service_status"
    
    if check_file "$service_status"; then
        local ntp_status="unknown"
        local chrony_status="unknown"
        
        if grep -q "ntpd" "$service_status"; then
            if grep -A 5 "ntpd" "$service_status" | grep -q "active (running)"; then
                ntp_status="running"
            else
                ntp_status="stopped"
            fi
        fi
        
        if grep -q "chronyd" "$service_status"; then
            if grep -A 5 "chronyd" "$service_status" | grep -q "active (running)"; then
                chrony_status="running"
            else
                chrony_status="stopped"
            fi
        fi
        
        if [[ "$ntp_status" == "running" || "$chrony_status" == "running" ]]; then
            echo -e "  [✓] 时间同步服务: ntpd=$ntp_status, chronyd=$chrony_status" >&2
        else
            add_anomaly "提示" "时间同步服务未运行" "TIME_SYNC_SERVICE_DOWN" "ntpd和chronyd服务均未运行" "$service_status"
            echo -e "  [!] 时间同步服务: ntpd=$ntp_status, chronyd=$chrony_status (建议启用)" >&2
        fi
    else
        echo -e "  [-] 时间同步服务: 文件不存在" >&2
    fi
}

# 检测系统配置
check_system_config() {
    log_info "检查系统配置..."
    
    local system_status="$DIAGNOSE_DIR/system_status"
    
    # 检查Swap配置
    if check_file "$system_status"; then
        if grep -q "Swap:.*0 kB" "$system_status"; then
            echo -e "  [✓] Swap配置: 已禁用" >&2
        else
            add_anomaly "提示" "Swap已启用" "SWAP_ENABLED" "Kubernetes不建议启用Swap" "$system_status"
            echo -e "  [!] Swap配置: 已启用 (Kubernetes不建议启用)" >&2
        fi
    else
        echo -e "  [-] Swap配置: 文件不存在" >&2
    fi
    
    # 检查IP转发
    if check_file "$system_status"; then
        if grep -q "net.ipv4.ip_forward = 1" "$system_status"; then
            echo -e "  [✓] IP转发: 已启用" >&2
        else
            add_anomaly "警告" "IP转发未启用" "IP_FORWARD_DISABLED" "Kubernetes需要启用IP转发" "$system_status"
            echo -e "  [✗] IP转发: 未启用 (Kubernetes需要启用)" >&2
        fi
    else
        echo -e "  [-] IP转发: 文件不存在" >&2
    fi
    
    # 检查bridge-nf-call-iptables
    if check_file "$system_status"; then
        if grep -q "net.bridge.bridge-nf-call-iptables = 1" "$system_status"; then
            echo -e "  [✓] bridge-nf-call-iptables: 已启用" >&2
        else
            add_anomaly "警告" "bridge-nf-call-iptables未启用" "BRIDGE_NF_CALL_IPTABLES_DISABLED" "Kubernetes需要启用此选项" "$system_status"
            echo -e "  [✗] bridge-nf-call-iptables: 未启用 (Kubernetes需要启用)" >&2
        fi
    else
        echo -e "  [-] bridge-nf-call-iptables: 文件不存在" >&2
    fi
}

# ============================================================================
# 主函数
# ============================================================================

main() {
    parse_arguments "$@"
    
    echo -e "${BLUE}================================================================${NC}" >&2
    echo -e "${BLUE}  kudig v$VERSION - Kubernetes节点诊断分析工具${NC}" >&2
    echo -e "${BLUE}================================================================${NC}" >&2
    echo "" >&2
    echo -e "诊断目录: $DIAGNOSE_DIR" >&2
    echo -e "分析时间: $(date '+%Y-%m-%d %H:%M:%S')" >&2
    echo "" >&2
    echo -e "开始诊断检查..." >&2
    echo "" >&2
    
    # 执行各项检查
    echo -e "========== 系统资源检查 ==========" >&2
    check_system_resources
    echo "" >&2
    
    echo -e "========== 进程与服务检查 ==========" >&2
    check_processes_services
    echo "" >&2
    
    echo -e "========== 网络状态检查 ==========" >&2
    check_network
    echo "" >&2
    
    echo -e "========== 内核状态检查 ==========" >&2
    check_kernel
    echo "" >&2
    
    echo -e "========== 容器运行时检查 ==========" >&2
    check_container_runtime
    echo "" >&2
    
    echo -e "========== Kubernetes组件检查 ==========" >&2
    check_kubernetes
    echo "" >&2
    
    echo -e "========== 时间同步检查 ==========" >&2
    check_time_sync
    echo "" >&2
    
    echo -e "========== 系统配置检查 ==========" >&2
    check_system_config
    echo "" >&2
    
    # 生成报告
    echo -e "${BLUE}================================================================${NC}" >&2
    echo -e "${BLUE}  诊断结果汇总${NC}" >&2
    echo -e "${BLUE}================================================================${NC}" >&2
    echo -e "=== Kubernetes节点诊断异常报告 ===" >&2
    echo -e "诊断时间: $(date '+%Y-%m-%d %H:%M:%S')" >&2
    echo -e "节点信息: $(hostname 2>/dev/null || echo "未知")" >&2
    echo -e "分析目录: $DIAGNOSE_DIR" >&2
    echo "" >&2
    
    # 分类异常
    local critical=()
    local warning=()
    local info=()
    
    for anomaly in "${ANOMALIES[@]}"; do
        IFS='|' read -r severity cn_name en_name details location <<< "$anomaly"
        case "$severity" in
            "严重")
                critical+=("$anomaly")
                ;;
            "警告")
                warning+=("$anomaly")
                ;;
            "提示")
                info+=("$anomaly")
                ;;
        esac
    done
    
    # 输出严重异常
    if [[ ${#critical[@]} -gt 0 ]]; then
        echo -e "-------------------------------------------" >&2
        echo -e "【严重级别】异常项" >&2
        echo -e "-------------------------------------------" >&2
        for anomaly in "${critical[@]}"; do
            IFS='|' read -r severity cn_name en_name details location <<< "$anomaly"
            echo -e "[${RED}严重${NC}] $cn_name | $en_name" >&2
            echo -e "  详情: $details" >&2
            echo -e "  位置: $location" >&2
            echo "" >&2
        done
    fi
    
    # 输出警告异常
    if [[ ${#warning[@]} -gt 0 ]]; then
        echo -e "-------------------------------------------" >&2
        echo -e "【警告级别】异常项" >&2
        echo -e "-------------------------------------------" >&2
        for anomaly in "${warning[@]}"; do
            IFS='|' read -r severity cn_name en_name details location <<< "$anomaly"
            echo -e "[${YELLOW}警告${NC}] $cn_name | $en_name" >&2
            echo -e "  详情: $details" >&2
            echo -e "  位置: $location" >&2
            echo "" >&2
        done
    fi
    
    # 输出提示异常
    if [[ ${#info[@]} -gt 0 ]]; then
        echo -e "-------------------------------------------" >&2
        echo -e "【提示级别】异常项" >&2
        echo -e "-------------------------------------------" >&2
        for anomaly in "${info[@]}"; do
            IFS='|' read -r severity cn_name en_name details location <<< "$anomaly"
            echo -e "[${BLUE}提示${NC}] $cn_name | $en_name" >&2
            echo -e "  详情: $details" >&2
            echo -e "  位置: $location" >&2
            echo "" >&2
        done
    fi
    
    # 输出异常统计
    echo -e "-------------------------------------------" >&2
    echo -e "异常统计" >&2
    echo -e "-------------------------------------------" >&2
    echo -e "严重: ${#critical[@]} 项" >&2
    echo -e "警告: ${#warning[@]} 项" >&2
    echo -e "提示: ${#info[@]} 项" >&2
    echo -e "总计: ${#ANOMALIES[@]} 项" >&2
    echo "" >&2
    
    # 输出结论
    if [[ ${#critical[@]} -eq 0 && ${#warning[@]} -eq 0 && ${#info[@]} -eq 0 ]]; then
        echo -e "${GREEN}✓ 未检测到异常${NC}" >&2
        echo -e "节点状态良好！" >&2
        exit 0
    elif [[ ${#critical[@]} -gt 0 ]]; then
        exit 2
    else
        exit 1
    fi
}

# 执行主函数
main "$@"
