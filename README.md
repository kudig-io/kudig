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
- **特性**: 43项异常检测规则，离线分析模式
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
| **检测规则** | 43 项 | 35+ 项 |
| **输出格式** | Text/JSON | Text/JSON |
| **排查建议** | ✅ | ✅ |
| **自定义规则** | ❌ | ✅ YAML规则 |
| **K8s部署** | ❌ | ✅ Helm Chart |
| **依赖** | 无 | Go 1.21+ |
| **跨平台** | Linux/Unix | ✅ 支持 |

## 内置检测规则（43项）

kudig.sh 内置了多个类别的43项检测规则，自动识别各类异常情况。

### 1. 系统资源检测（6项）
- **CPU负载检测**：检查15分钟平均负载是否超过CPU核心数2倍/4倍
- **内存使用检测**：检测内存使用率是否超过85%/95%
- **磁盘空间检测**：检测挂载点使用率是否超过90%/95%
- **文件句柄检测**：检测进程文件句柄数是否过高(>50000)
- **进程/线程数检测**：检测PID泄漏，某进程线程数>5000/10000
- **磁盘使用检测**：检测磁盘使用率是否超过90%

### 2. 进程与服务检测（6项）
- **Kubelet服务检测**：检测kubelet服务状态（running/failed/stopped）
- **容器运行时检测**：检测docker/containerd服务状态
- **ps命令检测**：检测ps命令是否挂起
- **D状态进程检测**：检测是否存在D状态（不可中断睡眠）进程
- **Runc进程检测**：检测runc进程数是否异常(>1000)
- **Firewalld检测**：检测Firewalld是否在运行（可能影响K8s网络）

### 3. 网络检测（5项）
- **连接追踪表检测**：检测conntrack表使用率>80%/95%
- **网卡状态检测**：检测网卡是否DOWN（排除lo/veth）
- **默认路由检测**：检测是否配置默认路由
- **端口监听检测**：检测kubelet端口(10250)是否在监听
- **Iptables规则检测**：检测iptables规则数是否过多(>50000)

### 4. 内核检测（7项）
- **内核Panic检测**：检测dmesg中是否有kernel panic
- **OOM Killer检测**：检测是否发生OOM事件
- **Messages日志OOM检测**：检测/var/log/messages中的OOM事件
- **文件系统错误检测**：检测文件系统是否变为只读
- **磁盘IO错误检测**：检测磁盘IO错误次数
- **内核模块检测**：检测是否有内核模块加载失败
- **NMI Watchdog检测**：检测NMI watchdog事件

### 5. 容器运行时检测（4项）
- **Docker启动检测**：检测docker启动错误
- **Docker存储驱动检测**：检测存储驱动错误
- **Containerd容器创建检测**：检测容器创建失败次数
- **镜像拉取检测**：检测镜像拉取失败错误

### 6. Kubernetes组件检测（9项）
- **PLEG状态检测**：检测PLEG is not healthy错误
- **CNI插件检测**：检测CNI插件错误（network plugin not ready）
- **Kubelet证书检测**：检测Kubernetes证书是否过期
- **API Server连接检测**：检测与API Server的连接问题
- **Kubelet认证检测**：检测认证失败错误
- **Pod驱逐检测**：检测Pod驱逐事件
- **节点状态检测**：检测节点是否Ready
- **磁盘压力检测**：检测节点磁盘压力
- **内存压力检测**：检测节点内存压力

### 7. 时间同步检测（1项）
- **NTP/Chrony状态检测**：检测时间同步服务状态（ntpd/chronyd）

### 8. 配置检测（5项）
- **Swap配置检测**：检测Swap是否开启（K8s不建议开启）
- **IP转发检测**：检测net.ipv4.ip_forward是否启用
- **bridge-nf-call-iptables检测**：检测桥接流量是否经过iptables
- **ulimit检测**：检测open files限制是否过低
- **SELinux检测**：检测SELinux是否为Enforcing



## 输出示例

### 文本格式（默认）

**正常情况（有少量警告）：**

```
================================================================
  kudig.sh v1.0.0 - Kubernetes节点诊断分析工具
================================================================

诊断目录: /tmp/diagnose_1702468800
分析时间: 2026-01-16 10:50:01

开始诊断检查...

========== 系统资源检查 ==========
  [✓] CPU负载: 正常 (15min负载: 0.34, CPU核心: 8)
  [✓] 内存使用: 正常 (使用率: 32%)
  [✓] 磁盘空间: 正常 (所有挂载点使用率<90%)
  [✓] 文件句柄: 正常 (最大: 273)
  [✓] 进程/线程数: 正常 (最大线程数: 33)
  [✓] 磁盘使用: 正常 (所有挂载点<90%)

========== 进程与服务检查 ==========
  [✓] Kubelet服务: running
  [✓] 容器运行时: docker=unknown, containerd=running
  [✓] ps命令: 正常
  [✓] D状态进程: 未发现
  [✓] runc进程: 正常
  [✓] Firewalld: 已关闭

========== 网络状态检查 ==========
  [✓] 连接跟踪表: 正常 (517/262144, 0%)
  [!] 网卡状态: 部分down (kube-ipvs0,nodelocaldns)
  [✓] 默认路由: 已配置
  [✓] Kubelet端口(10250): 正常监听
  [✓] iptables规则: 正常 (44 条)

========== 内核状态检查 ==========
  [✓] 内核Panic: 未发现
  [✓] OOM Killer: 未触发
  [✓] messages日志OOM: 未发现
  [✓] 文件系统: 正常
  [✓] 磁盘IO: 正常 (0 次错误)
  [!] 内核模块: 存在加载失败
  [!] NMI Watchdog: 被触发

========== 容器运行时检查 ==========
  [✓] Docker启动: 正常
  [✓] Docker存储驱动: 正常
  [✓] Containerd容器创建: 正常 (0 次失败)
  [✓] 镜像拉取: 正常 (0 次失败)

========== Kubernetes组件检查 ==========
  [✓] PLEG状态: 健康
  [✓] CNI网络插件: 正常
  [✓] Kubelet证书: 正常
  [✓] API Server连接: 正常 (0 次失败)
  [✓] Kubelet认证: 正常
  [✓] Pod驱逐: 未发现
  [✓] 节点状态: Ready
  [✓] 磁盘压力: 无
  [✓] 内存压力: 无

========== 时间同步检查 ==========
  [!] 时间同步: ntpd=unknown, chronyd=unknown (建议启用)

========== 系统配置检查 ==========
  [✓] Swap配置: 已禁用
  [✓] IP转发: 已启用
  [✓] bridge-nf-call-iptables: 已启用
  [✓] ulimit open files: 正常
  [✓] SELinux: 非Enforcing

================================================================
  诊断结果汇总
================================================================
=== Kubernetes节点诊断异常报告 ===
诊断时间: 2026-01-16 10:50:11
节点信息: k8s-node-01
分析目录: /tmp/diagnose_1702468800

-------------------------------------------
【警告级别】异常项
-------------------------------------------
[警告] 网卡接口down | NETWORK_INTERFACE_DOWN
  详情: 以下网卡处于down状态: kube-ipvs0,nodelocaldns
  位置: network_info

[警告] 内核模块加载失败 | KERNEL_MODULE_LOAD_FAILED
  详情: 存在内核模块加载失败
  位置: logs/dmesg.log

-------------------------------------------
【提示级别】异常项
-------------------------------------------
[提示] 时间同步服务未运行 | TIME_SYNC_SERVICE_DOWN
  详情: ntpd和chronyd服务均未运行
  位置: service_status

-------------------------------------------
异常统计
-------------------------------------------
严重: 0 项
警告: 2 项
提示: 1 项
总计: 3 项
```

**有严重异常情况：**

```
========== 进程与服务检查 ==========
  [✗] Kubelet服务: failed
    → 建议: 检查kubelet日志: journalctl -u kubelet -n 100; systemctl restart kubelet

========== 系统资源检查 ==========
  [✗] 磁盘空间 [/]: 严重不足 (使用率: 96%)
    → 建议: 检查占用空间大的目录: du -sh /* | sort -rh | head

-------------------------------------------
【严重级别】异常项
-------------------------------------------
[严重] Kubelet服务未运行 | KUBELET_SERVICE_DOWN
  详情: kubelet.service状态为failed
  位置: daemon_status/kubelet_status

[严重] 磁盘空间严重不足 | DISK_SPACE_CRITICAL
  详情: 挂载点 / 使用率 96%
  位置: system_status
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
