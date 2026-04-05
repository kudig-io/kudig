# Kudig v2.0 功能特性清单

> **版本**: v2.0  
> **更新日期**: 2026-04-02  
> **状态**: 功能特性 100% 完成

---

## 概览

Kudig v2.0 是一个功能完整的 Kubernetes 节点诊断工具，包含 **70+ 分析器**、**18+ CLI 命令**、**9 大功能类别**，覆盖从离线分析到在线诊断、从单节点到多集群的全场景需求。

---

## 核心功能统计

| 指标 | 数量 | 说明 |
|------|------|------|
| 分析器 | 70+ | 9 大类：系统、进程、网络、内核、K8s、运行时、安全、eBPF、服务网格 |
| CLI 命令 | 18+ | 覆盖诊断、分析、修复、扫描、导出等全场景 |
| 报告格式 | 4 | Text、JSON、HTML、SARIF |
| 部署方式 | 4 | Binary、Docker、Helm、Operator |

---

## CLI 命令清单

### 诊断命令

| 命令 | 功能 | 示例 |
|------|------|------|
| `offline` | 离线分析诊断数据 | `kudig offline /tmp/diagnose_1702468800` |
| `online` | 在线实时诊断 | `kudig online --all-nodes` |
| `analyze` | 离线分析（alias） | `kudig analyze /tmp/diagnose_1702468800` |
| `legacy` | v1.0 兼容模式 | `kudig legacy /tmp/diagnose_1702468800` |
| `rules` | YAML 规则引擎 | `kudig rules --file rules/custom.yaml` |
| `multicluster` | 多集群诊断 | `kudig multicluster --all-contexts` |

### 高级功能命令

| 命令 | 功能 | 示例 |
|------|------|------|
| `tui` | TUI 交互模式 | `kudig tui` |
| `rca` | 根因分析 | `kudig rca` |
| `fix` | 自动修复 | `kudig fix --dry-run` |
| `cost` | 成本分析 | `kudig cost` |
| `scan` | 镜像安全扫描 | `kudig scan nginx:latest` |
| `grafana` | Grafana Dashboard | `kudig grafana > dashboard.json` |
| `pprof` | 性能剖析 | `kudig pprof --port 6060` |
| `trace` | 分布式追踪 | `kudig trace --jaeger http://localhost:14268` |

### 工具命令

| 命令 | 功能 | 示例 |
|------|------|------|
| `list-analyzers` | 列出分析器 | `kudig list-analyzers` |
| `history` | 历史数据管理 | `kudig history list` |
| `completion` | Shell 补全 | `kudig completion bash` |

---

## 分析器分类（70+）

### 1. System（系统）- 7 个

检测系统资源使用情况：

- CPU 负载检测
- 内存使用检测
- 磁盘空间检测
- 文件句柄检测
- 进程/线程数检测
- Swap 配置检测
- 系统时间检测

### 2. Process（进程）- 5 个

检测进程和服务状态：

- Kubelet 服务检测
- 容器运行时检测
- PID 泄漏检测
- D 状态进程检测
- Zombie 进程检测

### 3. Network（网络）- 5 个

检测网络配置和状态：

- 连接追踪表检测
- 网卡状态检测
- 默认路由检测
- 端口监听检测
- iptables 规则检测

### 4. Kernel（内核）- 5 个

检测内核状态和日志：

- Kernel Panic 检测
- OOM Killer 检测
- 文件系统错误检测
- 磁盘 IO 错误检测
- 内核模块检测

### 5. Kubernetes（K8s）- 21 个

检测 Kubernetes 组件和资源配置：

**控制平面**
- API Server 连接检测
- etcd 健康检测
- Scheduler 状态检测
- Controller Manager 检测
- 证书过期检测

**工作负载**
- Pod 状态检测
- Deployment 健康检测
- StatefulSet 检测
- DaemonSet 检测
- Job/CronJob 检测

**存储**
- PVC/PV 状态检测
- StorageClass 检测
- CSI Driver 检测

**网络**
- Service 状态检测
- Endpoint 检测
- DNS (CoreDNS) 检测
- NetworkPolicy 检测

**资源**
- GPU 资源检测
- 节点资源压力检测
- 资源配额检测

### 6. Runtime（运行时）- 4 个

检测容器运行时状态：

- Docker 状态检测
- Containerd 状态检测
- CRI-O 状态检测
- 镜像拉取检测

### 7. Security（安全）- 11 个

检测安全合规性：

**CIS Benchmark**
- API Server 安全配置
- etcd 安全配置
- Kubelet 安全配置
- Pod 安全策略
- NetworkPolicy 检测

**RBAC**
- ClusterRole 绑定检测
- 危险权限检测
- ServiceAccount 检测
- 未使用角色检测

**其他**
- 特权容器检测
- Secret 管理检测

### 8. eBPF - 3 个

基于 eBPF 的深度诊断（内核 4.18+）：

- TCP 连接分析（重传、延迟）
- DNS 查询分析（失败率、延迟）
- 文件 I/O 分析（延迟、频率）

### 9. Service Mesh（服务网格）- 2 个

检测服务网格健康状态：

- Istio 诊断（istiod、sidecar）
- Linkerd 诊断（控制平面、proxy）

---

## 特色功能

### TUI 交互模式

基于 bubbletea 的终端用户界面：
- 菜单驱动的操作界面
- 实时诊断进度显示
- 键盘导航支持
- 问题详情交互查看

### 根因分析 (RCA)

智能根因分析引擎：
- 多症状关联分析
- 置信度评分
- 系统级修复建议
- 10+ 内置 RCA 规则

### 自动修复引擎

安全的自动修复功能：
- 低风险修复自动执行
- 高风险操作确认机制
- 修复操作记录
- 支持修复回滚

### 成本分析

Kubernetes 资源成本估算：
- AWS/GCP/Azure 定价模型
- 多维度成本分析
- 资源优化建议
- 成本趋势预测

### 镜像安全扫描

集成 Trivy 的镜像扫描：
- CVE 漏洞检测
- 多严重级别分类
- 修复版本建议
- CI/CD 集成支持

### 多集群联邦诊断

跨集群诊断能力：
- 多 kubeconfig 上下文支持
- 并行集群诊断
- 统一问题视图
- 跨集群问题关联

### SARIF 安全报告

标准安全报告格式：
- SARIF 2.1.0 格式
- GitHub/CodeQL 兼容
- CI/CD 安全门禁集成
- 漏洞追踪和修复建议

### Grafana Dashboard

官方监控仪表板：
- 预配置 Dashboard JSON
- 问题分布可视化
- 趋势分析图表
- Prometheus 集成

### 性能剖析 (pprof)

Go 性能分析支持：
- CPU Profile
- Heap Profile
- Goroutine 分析
- 内存泄漏检测

### 分布式追踪

OpenTelemetry 支持：
- 诊断操作追踪
- Jaeger 集成
- 性能瓶颈定位
- 调用链分析

### Shell 自动补全

主流 Shell 支持：
- Bash 补全
- Zsh 补全
- Fish 补全
- PowerShell 补全

---

## 报告格式

### Text（文本）

默认文本格式，适合终端查看：
- 彩色输出支持
-  severity 分级显示
- 修复建议提示

### JSON

机器可读格式，适合集成：
- 结构化数据
- 元数据完整
- 易于解析处理

### HTML

可视化报告，适合分享：
- 图表可视化
- 交互式浏览
- 响应式设计

### SARIF

安全扫描标准格式：
- GitHub 安全面板集成
- CodeQL 兼容
- CI/CD 门禁支持

---

## 部署方式

### Binary

直接下载可执行文件：
```bash
# Linux/macOS/Windows 支持
./kudig offline /tmp/diagnose
```

### Docker

容器化部署：
```bash
docker run -v /tmp/diagnose:/data kudig/kudig offline /data
```

### Helm

Kubernetes 部署：
```bash
helm install kudig ./charts/kudig
```

### Operator

Kubernetes Operator 模式：
```bash
helm install kudig-operator ./operator/helm/kudig-operator
kubectl apply -f clusterdiagnostic.yaml
```

---

## 环境变量配置

| 变量 | 说明 | 示例 |
|------|------|------|
| KUDIG_AI_PROVIDER | AI 提供商 | openai, qwen, ollama |
| KUDIG_AI_API_KEY | API 密钥 | sk-xxx |
| KUDIG_AI_MODEL | 模型名称 | gpt-4 |
| KUDIG_AI_LANGUAGE | 输出语言 | zh, en |
| KUDIG_SLACK_WEBHOOK_URL | Slack 通知 | https://hooks.slack.com/... |
| KUDIG_DINGTALK_WEBHOOK_URL | 钉钉通知 | https://oapi.dingtalk.com/... |

---

## 未来规划

虽然功能特性已达到 100%，但我们仍在持续改进：

### 近期目标（2026 Q2）
- 测试覆盖率提升至 80%+
- SLSA Level 2 合规
- OpenSSF Best Practices Badge

### 中期目标（2026 Q3）
- 插件系统支持
- 更多分析器
- 性能优化

### 长期目标（2026 Q4+）
- CNCF Sandbox 申请
- 社区生态建设
- 企业级支持

---

## 参考文档

- [v2-go README](../v2-go/README.md) - v2.0 详细文档
- [ROADMAP](../ROADMAP.md) - 项目路线图
- [CNCF 差距分析](CNCF_GRADUATION_GAP_ANALYSIS.md) - CNCF 毕业差距分析
