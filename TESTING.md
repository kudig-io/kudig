# kudig.sh 测试说明

## 测试环境要求

`kudig.sh` 是一个Bash脚本，需要在Linux环境或具有Bash环境的系统上运行。

### 在Linux/Unix系统上测试

1. **准备诊断数据**：
```bash
# 在Kubernetes节点上收集诊断数据
sudo ./diagnose_k8s.sh

# 会生成类似 /tmp/diagnose_1702468800 的目录
```

2. **运行kudig.sh**：
```bash
# 添加执行权限
chmod +x kudig.sh

# 基本测试
./kudig.sh /tmp/diagnose_1702468800

# 详细模式测试（显示每一步检查项）
./kudig.sh --verbose /tmp/diagnose_1702468800

# JSON格式测试
./kudig.sh --json /tmp/diagnose_1702468800

# 保存到文件测试
./kudig.sh -o report.txt /tmp/diagnose_1702468800

# 组合测试
./kudig.sh --verbose --json -o report.json /tmp/diagnose_1702468800
```

3. **验证退出码**：
```bash
./kudig.sh /tmp/diagnose_1702468800
echo "Exit code: $?"
# 0 = 无异常
# 1 = 有警告/提示
# 2 = 有严重异常
```

### 在Windows WSL上测试

如果使用Windows系统，可以通过WSL (Windows Subsystem for Linux) 运行：

```bash
# 在WSL中
cd /mnt/c/Users/Allen/Documents/GitHub/kudig.sh
chmod +x kudig.sh

# 创建测试目录（模拟诊断数据）
./create_test_data.sh

# 运行测试
./kudig.sh ./test_diagnose_dir
```

### 使用Git Bash测试

如果安装了Git for Windows，可以使用Git Bash：

```bash
cd /c/Users/Allen/Documents/GitHub/kudig.sh
bash kudig.sh --help
```

## 功能验证清单

- [x] 帮助信息显示正常 (`--help`)
- [x] 版本信息显示正常 (`--version`)
- [x] 能够正确解析诊断目录
- [x] 系统资源检测功能正常（CPU、内存、磁盘、文件句柄、进程线程、Inode）
- [x] 进程服务检测功能正常（Kubelet、容器运行时、僵尸进程、Runc、Firewalld）
- [x] 网络检测功能正常（网卡、路由、连接追踪、端口、Iptables）
- [x] 内核检测功能正常（Panic、OOM、文件系统、模块）
- [x] 容器运行时检测功能正常（Containerd、Docker、容器异常、镜像拉取）
- [x] Kubernetes组件检测功能正常（PLEG、CNI、API Server、证书、节点状态、Pod、卷挂载、Sandbox）
- [x] 时间同步检测功能正常（NTP/Chrony）
- [x] 配置检测功能正常（Swap、SELinux、系统参数）
- [x] 异常去重功能正常
- [x] 异常排序功能正常
- [x] 文本格式输出正常
- [x] JSON格式输出正常
- [x] 文件保存功能正常
- [x] 退出码正确
- [x] 详细模式显示每个检查项状态
- [x] 异常时显示排查建议

## 预期输出示例

### 无异常情况

```
================================================================
  kudig.sh v1.0.0 - Kubernetes节点诊断分析工具
================================================================

诊断目录: /tmp/diagnose_1702468800
分析时间: 2024-12-13 21:47:00

开始诊断检查...

========== 系统资源检查 ==========
  [✓] CPU负载: 正常 (15min负载: 0.34, CPU核心: 4)
  [✓] 内存使用: 正常 (使用率: 31%)
  [✓] 磁盘空间: 正常 (所有挂载点使用率<90%)
  ...

========== 诊断结果汇总 ==========

✓ 未检测到异常
节点状态良好！
```

退出码: 0

### 有异常情况

```
================================================================
  kudig.sh v1.0.0 - Kubernetes节点诊断分析工具
================================================================

诊断目录: /tmp/diagnose_1702468800
分析时间: 2024-12-13 21:47:00

开始诊断检查...

========== 系统资源检查 ==========
  [✓] CPU负载: 正常 (15min负载: 0.34, CPU核心: 4)
  [✓] 内存使用: 正常 (使用率: 31%)
  [✗] 磁盘空间 [/]: 不足 (使用率: 92%)
    → 建议: 检查占用空间大的目录: du -sh /* | sort -rh | head
  ...

========== 进程与服务检查 ==========
  [✗] Kubelet服务: failed
    → 建议: 检查kubelet日志: journalctl -u kubelet -n 100; systemctl restart kubelet
  ...

================================================================
  诊断结果汇总
================================================================

-------------------------------------------
【严重级别】异常项
-------------------------------------------
[严重] Kubelet服务未运行 | KUBELET_SERVICE_DOWN
  详情: kubelet.service状态为failed
  位置: daemon_status/kubelet_status

-------------------------------------------
【警告级别】异常项
-------------------------------------------
[警告] 磁盘空间不足 | DISK_SPACE_LOW
  详情: 挂载点 / 使用率 92%
  位置: system_status

-------------------------------------------
异常统计
-------------------------------------------
严重: 1项
警告: 1项
提示: 0项
总计: 2项
```

退出码: 2

## 自动化测试脚本

可以创建一个自动化测试脚本来验证各个功能：

```bash
#!/bin/bash

echo "开始kudig.sh功能测试..."

# 测试1: 帮助信息
echo "测试1: 帮助信息"
./kudig.sh --help
if [ $? -eq 0 ]; then
    echo "✓ 帮助信息测试通过"
else
    echo "✗ 帮助信息测试失败"
fi

# 测试2: 版本信息
echo "测试2: 版本信息"
./kudig.sh --version
if [ $? -eq 0 ]; then
    echo "✓ 版本信息测试通过"
else
    echo "✗ 版本信息测试失败"
fi

# 测试3: 空目录测试
echo "测试3: 空目录测试"
mkdir -p /tmp/test_empty_diagnose
./kudig.sh /tmp/test_empty_diagnose
rm -rf /tmp/test_empty_diagnose

# 测试4: JSON输出格式
echo "测试4: JSON输出格式"
./kudig.sh --json /tmp/diagnose_* | python -m json.tool > /dev/null
if [ $? -eq 0 ]; then
    echo "✓ JSON格式测试通过"
else
    echo "✗ JSON格式测试失败"
fi

echo "测试完成！"
```

## 注意事项

1. 脚本需要在Linux环境或支持Bash的环境下运行
2. 诊断目录必须是由 `diagnose_k8s.sh` 生成的完整目录
3. 脚本不需要root权限，普通用户即可运行
4. 脚本只读取诊断数据，不会修改任何文件

## 常见问题

**Q: 在Windows上如何测试？**
A: 推荐使用WSL (Windows Subsystem for Linux) 或Git Bash。

**Q: 如果诊断目录不完整会怎样？**
A: 脚本会显示警告但继续分析可用的文件，不会中断执行。

**Q: 为什么某些检测项没有结果？**
A: 可能是对应的日志文件在诊断数据中缺失，这是正常的。

**Q: 如何验证JSON输出格式正确？**
A: 可以使用 `jq` 或 `python -m json.tool` 验证JSON格式。
