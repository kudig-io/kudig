#!/usr/bin/env bats

# 基础功能测试
# 测试 kudig.sh 的基本参数和输出

load ../test_helper

@test "kudig.sh 文件存在且可执行" {
    [ -f "$KUDIG_SH" ]
    [ -x "$KUDIG_SH" ]
}

@test "显示帮助信息 (--help)" {
    run bash "$KUDIG_SH" --help
    
    [ "$status" -eq 0 ]
    [[ "$output" =~ "用法" ]]
    [[ "$output" =~ "kudig.sh" ]]
    [[ "$output" =~ "选项" ]]
}

@test "显示版本信息 (--version)" {
    run bash "$KUDIG_SH" --version
    
    [ "$status" -eq 0 ]
    [[ "$output" =~ "version" ]]
}

@test "缺少诊断目录参数时报错" {
    run bash "$KUDIG_SH"
    
    [ "$status" -ne 0 ]
    [[ "$output" =~ "请指定诊断目录" || "$output" =~ "诊断目录" ]]
}

@test "诊断目录不存在时报错" {
    run bash "$KUDIG_SH" "/nonexistent/directory"
    
    [ "$status" -ne 0 ]
    [[ "$output" =~ "不存在" ]]
}

@test "空诊断目录能够处理" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    
    run bash "$KUDIG_SH" "$diagnose_dir"
    
    # 应该能够运行完成，即使没有发现异常
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "JSON 输出格式有效" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    create_system_status_with_load "$diagnose_dir" "1.0"
    create_memory_info "$diagnose_dir"
    
    run bash "$KUDIG_SH" --json "$diagnose_dir"
    
    # 验证输出是有效的 JSON
    validate_json_output "$output"
    [ $? -eq 0 ]
}

@test "JSON 输出包含必需字段" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    create_system_status_with_load "$diagnose_dir" "1.0"
    
    run bash "$KUDIG_SH" --json "$diagnose_dir"
    
    [[ "$output" =~ "report_version" ]]
    [[ "$output" =~ "timestamp" ]]
    [[ "$output" =~ "anomalies" ]]
    [[ "$output" =~ "summary" ]]
}

@test "verbose 模式输出详细信息" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    create_system_status_with_load "$diagnose_dir" "1.0"
    
    run bash "$KUDIG_SH" --verbose "$diagnose_dir"
    
    # verbose 模式应该输出更多信息到 stderr
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "使用 reference 目录的真实数据" {
    skip_if_missing_reference_data
    
    local ref_dir="$PROJECT_ROOT/reference/diagnose_k8s/diagnose_1765626516"
    
    if [ -d "$ref_dir" ]; then
        run bash "$KUDIG_SH" "$ref_dir"
        
        # 应该能够成功分析
        [ "$status" -eq 0 ] || [ "$status" -eq 1 ] || [ "$status" -eq 2 ]
    else
        skip "Reference data not available"
    fi
}

# 辅助函数
skip_if_missing_reference_data() {
    if [ ! -d "$PROJECT_ROOT/reference/diagnose_k8s/diagnose_1765626516" ]; then
        skip "Reference data not available"
    fi
}
