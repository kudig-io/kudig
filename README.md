# kudig - Kubernetes 节点诊断工具

> **快速选择**: 
> - ✅ [v1.0 Bash 版本](v1-bash/) - **生产可用**，轻量级 Bash 脚本
> - ✅ [v2.0 Go 版本](v2-go/) - **Alpha 可用**，Go 语言重构

## 项目简介

`kudig` 是一个强大的 Kubernetes 节点诊断工具，能够自动识别各类异常情况并生成中英文对照的诊断报告。

## 版本说明

### v1.0 - Bash 版本 ✅ 生产可用

**位置**: [`v1-bash/`](v1-bash/)

- **实现**: Bash 脚本
- **状态**: 稳定可用
- **特性**: 40+ 异常检测规则，离线分析模式
- **优势**: 轻量级、无依赖、开箱即用
- **文档**: [v1-bash/README.md](v1-bash/README.md)

**快速使用**:
```bash
cd v1-bash
./kudig.sh /tmp/diagnose_1702468800
```

### v2.0 - Go 版本 ✅ Alpha 可用

**位置**: [`v2-go/`](v2-go/)

- **实现**: Go 语言
- **状态**: Alpha 阶段，功能完整
- **特性**: 
  - 双模式支持（离线 + 在线）
  - 35+ 内置分析器
  - YAML 规则引擎
  - Kubernetes 原生部署
  - Docker 镜像
- **文档**: [v2-go/README.md](v2-go/README.md)

**快速使用**:
```bash
cd v2-go
make build
./build/kudig offline /tmp/diagnose_1702468800
```

## 项目结构

```
kudig/
├── v1-bash/              # ✅ v1.0 Bash 版本（生产可用）
│   ├── kudig.sh          # 主脚本
│   ├── README.md         # v1.0 文档
│   ├── TESTING.md        # 测试说明
│   └── reference/        # 示例诊断数据
│
├── v2-go/                # ✅ v2.0 Go 版本（Alpha 可用）
│   ├── cmd/              # CLI 入口
│   ├── pkg/              # 核心包
│   │   ├── analyzer/     # 分析器（35+）
│   │   ├── collector/    # 数据收集层
│   │   ├── reporter/     # 报告生成层
│   │   └── rules/        # 规则引擎
│   ├── charts/           # Helm Chart
│   ├── Makefile          # 构建脚本
│   ├── Dockerfile        # Docker 构建
│   └── README.md         # v2.0 文档
│
├── docs/                 # 文档目录
├── scripts/              # 辅助脚本
├── tests/                # 测试文件
├── README.md             # 本文档
└── LICENSE               # 许可证
```

## 快速开始

### 使用 v1.0 Bash 版本（推荐）

1. **进入 v1-bash 目录**
```bash
cd v1-bash
```

2. **收集诊断数据**
```bash
sudo ./diagnose_k8s.sh
# 生成目录: /tmp/diagnose_1702468800
```

3. **分析诊断数据**
```bash
./kudig.sh /tmp/diagnose_1702468800
```

详细使用说明请查看 [v1-bash/README.md](v1-bash/README.md)

### 开发 v2.0 Go 版本

```bash
cd v2-go
make deps
make build
./build/kudig offline /tmp/diagnose_1702468800
```

详细开发说明请查看 [v2-go/README.md](v2-go/README.md)

## 核心特性对比

| 特性 | v1.0 Bash | v2.0 Go |
|-----|----------|--------|
| **状态** | ✅ 生产可用 | ✅ Alpha 可用 |
| **实现语言** | Bash | Go |
| **离线分析** | ✅ | ✅ |
| **在线诊断** | ❌ | ✅ |
| **检测规则** | 40+ 项 | 35+ 项 |
| **输出格式** | Text/JSON | Text/JSON |
| **排查建议** | ✅ | ✅ |
| **自定义规则** | ❌ | ✅ YAML规则 |
| **K8s部署** | ❌ | ✅ Helm Chart |
| **依赖** | 无 | Go 1.21+ |
| **跨平台** | Linux/Unix | ✅ 支持 |

## 内置检测规则（40+项）

kudig.sh 内置了多个类别的40+项检测规则，自动识别各类异常情况。

### 1. 系统资源检测（6项）
- **CPU负载检测**：检查15分钟平均负载是否超过CPU核心数2倍/4個
- **内存使用检测**：检测内存使用率是否超过85%/95%
- **磁盘空间检测**：检测挂载点使用率是否超过90%/95%
- **文件句柄检测**：检测进程文件句柄数是否过高(>50000)
- **进程/线程数检测**：检测PID泄漏，某进程线程数>5000/10000
- **Inode使用检测**：检测inode使用率是否超过90%

### 2. 进程与服务检测（5项）
- **Kubelet服务检测**：检测kubelet服务状态（running/failed/stopped）
- **容器运行时检测**：检测docker/containerd服务状态
- **僵尸进程检测**：检测是否存在大量僵尸进程(>100/500)
- **Runc进程检测**：检测runc进程数是否异常(>1000)
- **Firewalld检测**：检测Firewalld是否在运行（可能影响K8s网络）

### 3. 网络检测（7项）
- **网卡状态检测**：检测网卡是否DOWN或发生错误
- **网卡错误检测**：检测网卡丢包、错误、冲突等问题
- **路由配置检测**：检测路由表是否有异常
- **连接追踪表检测**：检测conntrack表使用率>80%/95%
- **端口监听检测**：检测关键端口是否在监听
- **Iptables规则检测**：检测iptables规则数是否过多(>5000/10000)
- **网络延迟检测**：检测网络延迟是否异常

### 4. 内核检测（5项）
- **内核Panic检测**：检测dmesg中是否有kernel panic
- **OOM Killer检测**：检测是否发生OOM事件
- **文件系统错误检测**：检测文件系统错误（Ext4-fs error）
- **内核模块检测**：检测必要的内核模块是否加载
- **NMI Watchdog检测**：检测NMI watchdog事件

### 5. 容器运行时检测（5项）
- **Containerd状态检测**：检测containerd服务和日志
- **Docker状态检测**：检测docker服务和日志
- **容器异常退出检测**：检测容器频繁重启/异常退出
- **镜像拉取检测**：检测镜像拉取失败错误
- **运行时错误检测**：检测runc/containerd错误日志

### 6. Kubernetes组件检测（8项）
- **PLEG状态检测**：检测PLEG is not healthy错误
- **CNI插件检测**：检测CNI插件错误（network plugin not ready）
- **API Server连接检测**：检测与API Server的连接问题
- **证书过期检测**：检测Kubernetes证书是否过期
- **节点状态检测**：检测节点是否NotReady
- **Pod创建失败检测**：检测Pod创建失败错误
- **卷挂载检测**：检测卷挂载失败错误
- **Sandbox错误检测**：检测sandbox创建/删除错误

### 7. 时间同步检测（2项）
- **NTP/Chrony状态检测**：检测时间同步服务状态
- **时间偏移检测**：检测系统时间偏移是否过大

### 8. 配置检测（3项）
- **Swap配置检测**：检测Swap是否开启（K8s不建议开启）
- **SELinux检测**：检测SELinux配置是否影响K8s
- **系统参数检测**：检测关键系统参数配置



## 输出示例

### 文本格式（默认）

**无异常情况：**

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

**有异常情况：**

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

### JSON格式

```json
{
  "metadata": {
    "report_version": "1.0",
    "timestamp": "2024-12-13T13:47:00Z",
    "hostname": "k8s-node-01",
    "diagnose_dir": "/tmp/diagnose_1702468800"
  },
  "anomalies": [
    {
      "severity": "严重",
      "cn_name": "Kubelet服务未运行",
      "en_name": "KUBELET_SERVICE_DOWN",
      "details": "kubelet.service状态为failed",
      "location": "daemon_status/kubelet_status"
    },
    {
      "severity": "警告",
      "cn_name": "磁盘空间不足",
      "en_name": "DISK_SPACE_LOW",
      "details": "挂载点 / 使用率 92%",
      "location": "system_status"
    }
  ],
  "summary": {
    "critical": 1,
    "warning": 1,
    "info": 0,
    "total": 2
  }
}
```

## 文档索引

- **主文档**: [README.md](README.md) - 本文档
- **v1.0 Bash**: [v1-bash/README.md](v1-bash/README.md) - 生产可用版本
  - [v1-bash/TESTING.md](v1-bash/TESTING.md) - 测试说明
- **v2.0 Go**: [v2-go/README.md](v2-go/README.md) - 开发版本
- **质量报告**: [QUALITY_REPORT.md](QUALITY_REPORT.md) - 代码质量检查

## 版本选择指南

### 何时使用 v1.0 Bash 版本？

✅ **推荐场景**:
- 生产环境诊断
- 需要快速部署
- 无法安装额外依赖
- 离线分析 diagnose_k8s.sh 数据
- 需要详细的排查建议

### 何时使用 v2.0 Go 版本？

✅ **适用场景**:
- 需要在线实时诊断 K8s 集群
- 需要自定义 YAML 规则
- 需要 Kubernetes 原生部署（DaemonSet）
- 需要跨平台支持（Windows/macOS）
- 希望参与开源开发

## 贡献

欢迎贡献！请阅读各版本的 README 了解详情。

## 许可证

本项目采用 Apache License 2.0 许可证。

---

**版本说明**: v1.0 Bash 版本为生产可用的稳定版本；v2.0 Go 版本功能完整，处于 Alpha 阶段，欢迎测试反馈。
