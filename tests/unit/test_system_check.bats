#!/usr/bin/env bats

# 系统资源检查测试
# 测试 check_system_resources 函数的各种场景

load ../test_helper

@test "检测正常 CPU 负载" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    
    # 创建 system_info 声明 4 核 CPU
    echo "processor       : 0" > "$diagnose_dir/system_info"
    echo "processor       : 1" >> "$diagnose_dir/system_info"
    echo "processor       : 2" >> "$diagnose_dir/system_info"
    echo "processor       : 3" >> "$diagnose_dir/system_info"
    
    # 创建正常负载（15分钟负载 = 2.0，低于 4*2=8）
    create_system_status_with_load "$diagnose_dir" "2.0"
    create_memory_info "$diagnose_dir"
    
    run bash "$KUDIG_SH" "$diagnose_dir"
    
    # 不应该报告 CPU 负载过高
    ! assert_anomaly_exists "$output" "HIGH_SYSTEM_LOAD"
}

@test "检测偏高 CPU 负载" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    
    # 4 核 CPU
    echo "processor       : 0" > "$diagnose_dir/system_info"
    echo "processor       : 1" >> "$diagnose_dir/system_info"
    echo "processor       : 2" >> "$diagnose_dir/system_info"
    echo "processor       : 3" >> "$diagnose_dir/system_info"
    
    # 偏高负载（15分钟负载 = 10.0，介于 4*2=8 和 4*4=16 之间）
    create_system_status_with_load "$diagnose_dir" "10.0"
    create_memory_info "$diagnose_dir"
    
    run bash "$KUDIG_SH" "$diagnose_dir"
    
    # 应该报告负载偏高（警告级别）
    assert_anomaly_exists "$output" "ELEVATED_SYSTEM_LOAD"
    [ "$status" -eq 1 ] || [ "$status" -eq 2 ]
}

@test "检测严重 CPU 负载" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    
    # 4 核 CPU
    echo "processor       : 0" > "$diagnose_dir/system_info"
    echo "processor       : 1" >> "$diagnose_dir/system_info"
    echo "processor       : 2" >> "$diagnose_dir/system_info"
    echo "processor       : 3" >> "$diagnose_dir/system_info"
    
    # 严重负载（15分钟负载 = 20.0，超过 4*4=16）
    create_system_status_with_load "$diagnose_dir" "20.0"
    create_memory_info "$diagnose_dir"
    
    run bash "$KUDIG_SH" "$diagnose_dir"
    
    # 应该报告负载过高（严重级别）
    assert_anomaly_exists "$output" "HIGH_SYSTEM_LOAD"
    [ "$status" -eq 2 ]
}

@test "检测正常内存使用" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    create_system_status_with_load "$diagnose_dir" "1.0"
    
    # 总内存 8GB，可用 4GB（使用率 50%）
    create_memory_info "$diagnose_dir" "8000000" "4000000"
    
    run bash "$KUDIG_SH" "$diagnose_dir"
    
    # 不应该报告内存问题
    ! assert_anomaly_exists "$output" "HIGH_MEMORY_USAGE"
    ! assert_anomaly_exists "$output" "ELEVATED_MEMORY_USAGE"
}

@test "检测内存使用率偏高" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    create_system_status_with_load "$diagnose_dir" "1.0"
    
    # 总内存 8GB，可用 0.8GB（使用率 90%）
    create_memory_info "$diagnose_dir" "8000000" "800000"
    
    run bash "$KUDIG_SH" "$diagnose_dir"
    
    # 应该报告内存使用率偏高
    assert_anomaly_exists "$output" "ELEVATED_MEMORY_USAGE"
    [ "$status" -eq 1 ] || [ "$status" -eq 2 ]
}

@test "检测内存使用率严重" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    create_system_status_with_load "$diagnose_dir" "1.0"
    
    # 总内存 8GB，可用 0.2GB（使用率 97.5%）
    create_memory_info "$diagnose_dir" "8000000" "200000"
    
    run bash "$KUDIG_SH" "$diagnose_dir"
    
    # 应该报告内存使用率过高（严重级别）
    assert_anomaly_exists "$output" "HIGH_MEMORY_USAGE"
    [ "$status" -eq 2 ]
}

@test "检测正常磁盘空间" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    create_system_status_with_load "$diagnose_dir" "1.0"
    create_memory_info "$diagnose_dir"
    
    # 磁盘使用率 50%
    create_system_status_with_disk "$diagnose_dir" "50"
    
    run bash "$KUDIG_SH" "$diagnose_dir"
    
    # 不应该报告磁盘空间问题
    ! assert_anomaly_exists "$output" "DISK_SPACE_LOW"
    ! assert_anomaly_exists "$output" "DISK_SPACE_CRITICAL"
}

@test "检测磁盘空间不足" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    create_system_status_with_load "$diagnose_dir" "1.0"
    create_memory_info "$diagnose_dir"
    
    # 磁盘使用率 92%
    cat > "$diagnose_dir/system_status" <<EOF
----- run df -h -----
Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1        50G   46G    4G  92% /
----- End of df -----
EOF
    
    run bash "$KUDIG_SH" "$diagnose_dir"
    
    # 应该报告磁盘空间不足
    assert_anomaly_exists "$output" "DISK_SPACE_LOW"
    [ "$status" -eq 1 ] || [ "$status" -eq 2 ]
}

@test "检测磁盘空间严重不足" {
    local diagnose_dir=$(create_mock_diagnose_dir)
    create_system_status_with_load "$diagnose_dir" "1.0"
    create_memory_info "$diagnose_dir"
    
    # 磁盘使用率 97%
    cat > "$diagnose_dir/system_status" <<EOF
----- run df -h -----
Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1        50G   48G    2G  97% /
----- End of df -----
EOF
    
    run bash "$KUDIG_SH" "$diagnose_dir"
    
    # 应该报告磁盘空间严重不足
    assert_anomaly_exists "$output" "DISK_SPACE_CRITICAL"
    [ "$status" -eq 2 ]
}
