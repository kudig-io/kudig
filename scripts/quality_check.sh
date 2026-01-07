#!/usr/bin/env bash

################################################################################
# Quality Check Script for kudig.sh
# 自定义代码质量检查规则
#
# Usage:
#   ./quality_check.sh [options] <file1.sh> [file2.sh ...]
#
# Options:
#   --verbose    输出详细信息
#   --help       显示帮助信息
################################################################################

set -euo pipefail

# 颜色定义
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

VERBOSE=false
CHECK_FAILED=0
TOTAL_ISSUES=0

# 显示帮助
show_help() {
    cat << EOF
Usage: $0 [options] <file1.sh> [file2.sh ...]

Custom quality checks for Shell scripts.

Options:
    --verbose    Show detailed output
    --help       Show this help message

Examples:
    $0 kudig.sh
    $0 --verbose *.sh
EOF
    exit 0
}

# 日志函数
log_info() {
    if [[ "$VERBOSE" == true ]]; then
        echo -e "${BLUE}[INFO]${NC} $*"
    fi
}

log_pass() {
    echo -e "${GREEN}[PASS]${NC} $*"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
    ((TOTAL_ISSUES++)) || true
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $*"
    CHECK_FAILED=1
    ((TOTAL_ISSUES++)) || true
}

# 检查1：文件头部注释完整性
check_file_header() {
    local file="$1"
    log_info "Checking file header: $file"
    
    # 检查是否有 shebang
    if ! head -1 "$file" | grep -q '^#!/'; then
        log_fail "$file: Missing shebang line"
        return
    fi
    
    # 检查是否有功能描述
    if ! head -20 "$file" | grep -qi '# .*功能\|# .*用途\|# .*description'; then
        log_warn "$file: Missing function description in header"
    else
        log_pass "$file: File header complete"
    fi
}

# 检查2：全局变量命名规范
check_global_variables() {
    local file="$1"
    log_info "Checking global variable naming: $file"
    
    # 提取全局变量声明（排除函数内部）
    local bad_vars=$(awk '
        /^[a-z][a-zA-Z0-9_]*=/ && !/^    / && !/^\t/ {
            var = $0
            gsub(/=.*/, "", var)
            if (var !~ /^(local|readonly)/ && var != "set") {
                print NR ": " var
            }
        }
    ' "$file")
    
    if [[ -n "$bad_vars" ]]; then
        log_warn "$file: Global variables should use UPPER_CASE:"
        while IFS= read -r line; do
            echo "    $line"
        done <<< "$bad_vars"
    else
        log_pass "$file: Global variable naming OK"
    fi
}

# 检查3：函数命名规范
check_function_names() {
    local file="$1"
    log_info "Checking function naming: $file"
    
    # 查找不符合规范的函数名（应使用小写+下划线）
    local bad_funcs=$(grep -n '^[A-Z][a-zA-Z0-9_]*()' "$file" || true)
    
    if [[ -n "$bad_funcs" ]]; then
        log_warn "$file: Function names should use lowercase_with_underscores:"
        echo "$bad_funcs" | while IFS= read -r line; do
            echo "    $line"
        done
    else
        log_pass "$file: Function naming OK"
    fi
}

# 检查4：关键函数注释完整性
check_function_comments() {
    local file="$1"
    log_info "Checking function comments: $file"
    
    local missing_comments=0
    
    # 提取所有函数定义
    local functions=$(grep -n '^[a-z_][a-z0-9_]*()' "$file" | awk -F: '{print $1":"$2}' || true)
    
    if [[ -z "$functions" ]]; then
        log_pass "$file: No functions to check"
        return
    fi
    
    while IFS=: read -r line_num func_name; do
        func_name=$(echo "$func_name" | sed 's/().*//')
        
        # 检查函数前3行是否有注释
        local start_line=$((line_num - 3))
        [[ $start_line -lt 1 ]] && start_line=1
        
        local has_comment=$(sed -n "${start_line},${line_num}p" "$file" | grep -c '^#' || true)
        
        if [[ $has_comment -eq 0 ]] && [[ ! "$func_name" =~ ^log_ ]]; then
            if [[ "$VERBOSE" == true ]]; then
                log_warn "$file:$line_num: Function '$func_name' missing comment"
            fi
            ((missing_comments++)) || true
        fi
    done <<< "$functions"
    
    if [[ $missing_comments -gt 0 ]]; then
        log_warn "$file: $missing_comments functions missing comments"
    else
        log_pass "$file: Function comments OK"
    fi
}

# 检查5：退出码使用规范
check_exit_codes() {
    local file="$1"
    log_info "Checking exit code usage: $file"
    
    # 查找 exit 语句，检查是否使用了非标准退出码
    local bad_exits=$(grep -n 'exit [^012]' "$file" | grep -v '#' || true)
    
    if [[ -n "$bad_exits" ]]; then
        log_warn "$file: Consider using standard exit codes (0, 1, 2):"
        echo "$bad_exits" | while IFS= read -r line; do
            echo "    $line"
        done
    else
        log_pass "$file: Exit code usage OK"
    fi
}

# 检查6：set -e 模式下的算术运算
check_arithmetic_safety() {
    local file="$1"
    log_info "Checking arithmetic safety: $file"
    
    # 检查是否使用了 set -e
    if ! grep -q 'set -e\|set -euo' "$file"; then
        log_pass "$file: Not using 'set -e', arithmetic check skipped"
        return
    fi
    
    # 查找递增/递减操作，检查是否有 || true 保护
    local unsafe_arithmetic=$(grep -n '(([^)]*++\|--[^)]*))\s*$' "$file" | grep -v '|| true' || true)
    
    if [[ -n "$unsafe_arithmetic" ]]; then
        log_warn "$file: Arithmetic operations in 'set -e' mode should use '|| true':"
        echo "$unsafe_arithmetic" | while IFS= read -r line; do
            echo "    $line"
        done
    else
        log_pass "$file: Arithmetic safety OK"
    fi
}

# 检查7：硬编码路径
check_hardcoded_paths() {
    local file="$1"
    log_info "Checking for hardcoded paths: $file"
    
    # 查找可疑的硬编码路径（排除注释和常见系统路径）
    local hardcoded=$(grep -n '="/[^"]*"' "$file" | 
                      grep -v '#' | 
                      grep -v '/dev/null\|/etc/\|/var/\|/tmp/\|/proc/\|/sys/' | 
                      grep -v 'PATH=' || true)
    
    if [[ -n "$hardcoded" ]]; then
        log_warn "$file: Potential hardcoded paths found:"
        echo "$hardcoded" | head -5 | while IFS= read -r line; do
            echo "    $line"
        done
    else
        log_pass "$file: No hardcoded paths detected"
    fi
}

# 检查8：TODO/FIXME 标记
check_todo_fixme() {
    local file="$1"
    log_info "Checking for TODO/FIXME: $file"
    
    local todos=$(grep -ni 'TODO\|FIXME' "$file" || true)
    
    if [[ -n "$todos" ]]; then
        log_info "$file: Found TODO/FIXME markers:"
        echo "$todos" | while IFS= read -r line; do
            echo "    $line"
        done
    fi
}

# 主检查函数
check_file() {
    local file="$1"
    
    if [[ ! -f "$file" ]]; then
        log_fail "File not found: $file"
        return
    fi
    
    echo ""
    echo -e "${BLUE}=== Checking: $file ===${NC}"
    
    check_file_header "$file"
    check_global_variables "$file"
    check_function_names "$file"
    check_function_comments "$file"
    check_exit_codes "$file"
    check_arithmetic_safety "$file"
    check_hardcoded_paths "$file"
    check_todo_fixme "$file"
}

# 主函数
main() {
    local files=()
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --verbose)
                VERBOSE=true
                shift
                ;;
            --help)
                show_help
                ;;
            -*)
                echo "Unknown option: $1"
                show_help
                ;;
            *)
                files+=("$1")
                shift
                ;;
        esac
    done
    
    # 检查是否有文件参数
    if [[ ${#files[@]} -eq 0 ]]; then
        echo "Error: No files specified"
        show_help
    fi
    
    echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  Custom Quality Check for Shell       ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
    
    # 检查每个文件
    for file in "${files[@]}"; do
        check_file "$file"
    done
    
    # 输出总结
    echo ""
    echo -e "${BLUE}=== Summary ===${NC}"
    if [[ $CHECK_FAILED -eq 1 ]]; then
        echo -e "${RED}✗ Quality check FAILED${NC}"
        echo -e "Total issues found: $TOTAL_ISSUES"
        exit 1
    else
        if [[ $TOTAL_ISSUES -gt 0 ]]; then
            echo -e "${YELLOW}⚠ Quality check PASSED with warnings${NC}"
            echo -e "Total warnings: $TOTAL_ISSUES"
        else
            echo -e "${GREEN}✓ Quality check PASSED${NC}"
        fi
        exit 0
    fi
}

# 执行主函数
main "$@"
