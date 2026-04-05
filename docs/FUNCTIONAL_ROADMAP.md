# Kudig 功能补齐路线图

> **文档版本**: v2.0  
> **更新日期**: 2026-04-02  
> **目标**: 基于 CNCF 毕业标准和竞品分析，系统性补齐 Kudig v2-go 的功能短板

---

## 📑 文档导航

| 文档 | 说明 | 链接 |
|------|------|------|
| **功能路线图** | 本文档 - 完整的 Phase 1-7 开发计划与完成状态 | [FUNCTIONAL_ROADMAP.md](./FUNCTIONAL_ROADMAP.md) |
| **项目主文档** | 版本选择、快速开始、功能概览 | [../README.md](../README.md) |
| **v2.0 详细文档** | 68 分析器、Operator、eBPF、AI 功能说明 | [../v2-go/README.md](../v2-go/README.md) |
| **项目结构** | 目录结构、文件说明、分析器统计 | [../STRUCTURE.md](../STRUCTURE.md) |
| **贡献指南** | 开发环境、代码规范、提交流程 | [../CONTRIBUTING.md](../CONTRIBUTING.md) |
| **CNCF 差距分析** | 与 CNCF 毕业标准的差距评估 | [CNCF_GRADUATION_GAP_ANALYSIS.md](./CNCF_GRADUATION_GAP_ANALYSIS.md) |

---

## 执行计划概览

### 已完成 (Phase 1-3)

| 阶段 | 功能项 | 优先级 | 状态 |
|------|--------|--------|------|
| Phase 1 | Prometheus Metrics 暴露 (`/metrics`) | P0 | ✅ 已完成 |
| Phase 1 | 多节点并发诊断 + 进度条 | P0 | ✅ 已完成 |
| Phase 2 | HTML Reporter | P1 | ✅ 已完成 |
| Phase 2 | Control Plane 分析器 | P1 | ✅ 已完成 |
| Phase 2 | Workload 分析器 | P1 | ✅ 已完成 |
| Phase 3 | kubectl 插件 | P2 | ✅ 已完成 |
| Phase 3 | 历史数据对比 | P2 | ✅ 已完成 |
| Phase 3 | Webhook 告警通知 | P2 | ✅ 已完成 |

### 进行中 (Phase 4-7)

| 阶段 | 功能项 | 优先级 | 预估工作量 | 状态 |
|------|--------|--------|-----------|------|
| Phase 4 | DNS 诊断分析器 | P1 | 2-3 天 | ✅ 已完成 |
| Phase 4 | 存储性能分析器 | P1 | 3-4 天 | ✅ 已完成 |
| Phase 4 | GPU/NPU 诊断 | P1 | 2-3 天 | ✅ 已完成 |
| Phase 5 | CIS 安全合规扫描 | P2 | 3-4 天 | ✅ 已完成 |
| Phase 5 | RBAC 审计分析器 | P2 | 2-3 天 | ✅ 已完成 |
| Phase 6 | Operator 模式 | P2 | 4-5 天 | ✅ 已完成 |
| Phase 7 | eBPF 深度诊断 | P3 | 5-7 天 | ✅ 已完成 |
| Phase 7 | AI/LLM 辅助诊断 | P3 | 3-4 天 | ✅ 已完成 |

---

## Phase 4：核心组件深度覆盖

### 4.1 DNS 诊断分析器

**目标**：CoreDNS 健康检查、DNS 解析延迟/失败检测。

**实现内容**：
- 新增 `pkg/analyzer/kubernetes/dns.go`
- 实现以下分析器：
  - `coredns` Pod 健康检查（kube-system 命名空间）
  - CoreDNS 日志错误分析（loop 插件错误、upstream 超时）
  - DNS 解析延迟检测（通过 dig/nslookup）
  - DNS 配置检查（resolv.conf、ndots 配置）
  - 常见 DNS 问题检测：
    - CoreDNS CrashLoopBackOff
    - DNS 请求超时
    - NXDOMAIN 频发
    - DNS 缓存问题

**验收标准**：
- 能检测 CoreDNS 服务不可用
- 能识别 DNS 解析超时问题
- 单元测试覆盖率 > 60%

---

### 4.2 存储性能分析器

**目标**：PV/PVC 挂载链分析、I/O 延迟检测。

**实现内容**：
- 新增 `pkg/analyzer/kubernetes/storage.go`
- 实现以下分析器：
  - PVC 状态检查（Pending、Lost）
  - PV 状态检查（Released、Failed）
  - StorageClass 可用性检查
  - 存储插件 Pod 健康检查（CSI driver、provisioner）
  - 挂载失败检测（MountVolume.SetUp failed）
  - 存储性能指标：
    - VolumeAttachment 状态
    - Node 存储压力条件
    - 磁盘 I/O 等待时间（通过 node exporter 指标）

**验收标准**：
- 能检测 PVC 绑定失败
- 能识别存储挂载问题
- 单元测试覆盖率 > 60%

---

### 4.3 GPU/NPU 诊断

**目标**：nvidia-smi、DCGM 集成（AI 训练场景必需）。

**实现内容**：
- 新增 `pkg/analyzer/kubernetes/gpu.go`
- 实现以下分析器：
  - GPU 节点标签检查（nvidia.com/gpu.present）
  - NVIDIA 驱动检查（通过 node annotation）
  - GPU Pod 资源请求检查
  - GPU 共享模式检查（MIG、time-slicing）
  - GPU 显存压力检测
  - GPU 温度/功耗检查（通过 DCGM 指标，如可用）
  - 常见 GPU 问题：
    - GPU 设备插件未运行
    - GPU 资源分配失败
    - CUDA 版本不匹配

**验收标准**：
- 能识别 GPU 节点问题
- 支持 NVIDIA GPU 检测
- 单元测试覆盖率 > 60%

---

### Phase 4 完成总结

| 分析器 | 文件 | Issue 类型 |
|--------|------|-----------|
| CoreDNS | `kubernetes.coredns` | 未部署、CrashLoopBackOff、镜像拉取失败、重启过多、Endpoints 为空、缺少上游配置 |
| DNS Pod 配置 | `kubernetes.dns.pod_config` | ndots:1 配置、大量 Pod 使用 Default DNS 策略 |
| PVC | `kubernetes.pvc` | Pending 超时、Lost、未调度 |
| PV | `kubernetes.pv` | Failed、Released 超时、容量为 0 |
| StorageClass | `kubernetes.storageclass` | 无默认 SC、无回收策略、Immediate 绑定模式建议 |
| VolumeAttachment | `kubernetes.volumeattachment` | 未成功附加 |
| CSI Driver | `kubernetes.csi` | CSI Driver Pod 缺失 |
| Storage Pod | `kubernetes.storage.pod` | 存储调度失败 |
| GPU Node | `kubernetes.gpu.node` | 设备插件未运行、GPU 不可分配、节点状态异常 |
| GPU Pod | `kubernetes.gpu.pod` | 调度失败、资源限制不匹配、分数 GPU 请求 |
| GPU Share | `kubernetes.gpu.share` | 单例模式使用 |
| NPU | `kubernetes.npu` | Ascend 设备插件未运行 |

**新增分析器总数**: 12  
**当前总分析器数**: 54

---

## Phase 5：安全合规

### 5.1 CIS 安全合规扫描

**目标**：CIS Kubernetes Benchmark 检测。

**实现内容**：
- 新增 `pkg/analyzer/security/cis.go`
- 实现关键检查项：
  - API Server 安全配置（匿名认证、不安全端口）
  - etcd 安全配置（客户端证书认证）
  - Kubelet 安全配置（认证、授权、只读端口）
  - Pod Security 检查（特权容器、root 用户、危险挂载）
  - Network Policy 检查
  - Secret 管理检查

**验收标准**：
- 覆盖 CIS Benchmark Level 1 关键项
- 提供修复建议
- 单元测试覆盖率 > 60%

---

### 5.2 RBAC 审计分析器

**目标**：RBAC 配置审计、权限过度授予检测。

**实现内容**：
- 新增 `pkg/analyzer/security/rbac.go`
- 实现以下分析器：
  - 集群管理员角色绑定检测（cluster-admin）
  - 特权 ServiceAccount 检测
  - 默认 ServiceAccount 自动挂载检查
  - 危险权限检测（create pods、escalate）
  - 过期 Token 检测
  - 未使用 Role/ClusterRole 检测

**验收标准**：
- 能识别过度权限授予
- 提供最小权限建议
- 单元测试覆盖率 > 60%

---

### Phase 5 完成总结

| 分析器 | 文件 | Issue 类型 |
|--------|------|-----------|
| CIS API Server | `security.cis.apiserver` | 匿名认证、不安全端口、不安全绑定地址 |
| CIS etcd | `security.cis.etcd` | 未启用客户端证书认证、自动 TLS |
| CIS Kubelet | `security.cis.kubelet` | Kubelet 配置检查 |
| CIS Pod Security | `security.cis.pod` | 特权容器、HostPID、HostNetwork、root 用户运行 |
| CIS Network Policy | `security.cis.network` | 缺少 Network Policy 的命名空间 |
| CIS Secret | `security.cis.secret` | default 命名空间存放 Secret |
| RBAC Admin | `security.rbac.admin` | 过多 cluster-admin 绑定 |
| RBAC ServiceAccount | `security.rbac.serviceaccount` | default SA 自动挂载 Token |
| RBAC Dangerous | `security.rbac.dangerous` | 危险权限 (escalate, create pods, secrets, nodes proxy) |
| RBAC Unused | `security.rbac.unused` | 未使用的 ClusterRole |
| RBAC Token | `security.rbac.token` | 多个 Token 的 ServiceAccount |

**新增分析器总数**: 11  
**当前总分析器数**: 65

---

## Phase 6：Operator 模式

### 6.1 kudig-operator

**目标**：Kubernetes Operator 原生体验。

**实现内容**：
- 创建 `v2-go/operator/` 目录结构
- 实现以下 CRD：
  - `ClusterDiagnostic`：集群级诊断任务
  - `NodeDiagnostic`：节点级诊断任务
  - `DiagnosticSchedule`：定时诊断调度
- 实现 Controller 逻辑：
  - 监听 CR 创建/更新
  - 创建 Job/DaemonSet 执行诊断
  - 收集结果并写入 CR 状态
  - 触发告警通知
- 支持 Helm Chart 部署 Operator

**验收标准**：
- 可通过 kubectl apply 创建诊断任务
- 诊断结果可通过 kubectl describe 查看
- Operator 自身高可用部署

---

### Phase 6 完成总结

| 组件 | 文件 | 功能 |
|------|------|------|
| CRD Types | `operator/api/v1/` | ClusterDiagnostic, NodeDiagnostic, DiagnosticSchedule 类型定义 |
| Cluster Controller | `controllers/clusterdiagnostic_controller.go` | 监听 CR 创建，启动 Job 执行诊断，更新状态 |
| Node Controller | `controllers/nodediagnostic_controller.go` | 支持节点选择器，在多个节点上执行诊断 |
| Schedule Controller | `controllers/schedule_controller.go` | 定时触发诊断任务 (@hourly, @daily, @weekly) |
| Operator Main | `cmd/main.go` | Operator 启动入口，初始化 Manager 和 Controllers |
| Helm Chart | `helm/kudig-operator/` | 完整的 Helm Chart 包含 CRD, RBAC, Deployment |
| 示例 | `config/examples/` | 3 个示例 YAML 文件展示如何使用 CRD |
| 文档 | `README.md` | Operator 安装和使用文档 |

**目录结构**:
```
operator/
├── api/v1/                    # CRD 类型定义
│   ├── clusterdiagnostic_types.go
│   ├── nodediagnostic_types.go
│   ├── schedule_types.go
│   └── groupversion_info.go
├── controllers/               # 控制器实现
│   ├── clusterdiagnostic_controller.go
│   ├── nodediagnostic_controller.go
│   └── schedule_controller.go
├── cmd/
│   └── main.go               # Operator 入口
├── helm/
│   └── kudig-operator/       # Helm Chart
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
│           ├── deployment.yaml
│           ├── rbac.yaml
│           ├── crds.yaml
│           └── _helpers.tpl
├── config/examples/          # 示例配置
│   ├── cluster-diagnostic.yaml
│   ├── node-diagnostic.yaml
│   └── schedule.yaml
└── README.md                 # 使用文档
```

---

## Phase 7：创新方向（可选）

### 7.1 eBPF 深度诊断

**目标**：系统调用追踪、网络包级分析。

**实现内容**：
- 调研 eBPF 方案（cilium/ebpf-go、bpftrace）
- 实现以下诊断能力：
  - TCP 连接追踪（重传、延迟）
  - DNS 查询追踪
  - 系统调用延迟分析
  - 文件 I/O 追踪
  - 网络丢包检测

**验收标准**：
- 无需修改应用代码
- 支持内核 4.18+（BPF Type Format）
- 提供可视化追踪结果

---

### Phase 7 (eBPF) 完成总结

| 分析器 | 文件 | 功能 |
|--------|------|------|
| TCP 分析器 | `ebpf.tcp` | TCP 连接追踪、重传检测、延迟分析 |
| DNS 分析器 | `ebpf.dns` | DNS 查询追踪、失败率分析、延迟分析 |
| 文件 I/O 分析器 | `ebpf.fileio` | 文件操作追踪、I/O 延迟分析 |
| Probe 管理器 | `pkg/ebpf/probe/` | eBPF 程序管理、事件收集、统计计算 |

**目录结构**:
```
pkg/ebpf/
├── probe/
│   ├── probe.go              # eBPF Probe 管理器
│   └── probe_test.go         # 单元测试
└── analyzer/
    ├── tcp_analyzer.go       # TCP 连接分析
    ├── dns_analyzer.go       # DNS 查询分析
    ├── fileio_analyzer.go    # 文件 I/O 分析
    ├── init.go               # 分析器注册
    └── tcp_analyzer_test.go  # 单元测试
```

**新增分析器**: 3 个 (ebpf.tcp, ebpf.dns, ebpf.fileio)  
**当前总分析器数**: 68

---

### 7.2 AI/LLM 辅助诊断

**目标**：集成 LLM 解释诊断结果、给出修复建议。

**实现内容**：
- 新增 `pkg/ai/` 包
- 实现功能：
  - 诊断结果摘要生成
  - 根因分析推理
  - 修复建议智能推荐
  - 多语言支持（中英文自动切换）
- 支持多种 LLM 后端：
  - OpenAI GPT API
  - 阿里云通义千问
  - 私有化部署（Ollama、vLLM）

**验收标准**：
- 修复建议准确率 > 80%
- 响应延迟 < 5s
- 支持离线/在线模式

---

### Phase 7 (AI/LLM) 完成总结

| 组件 | 文件 | 功能 |
|------|------|------|
| Provider | `ai/provider.go` | AI Provider 接口、OpenAI 实现、配置管理 |
| Assistant | `ai/assistant.go` | AI 助手封装、问题分析、修复建议 |
| 配置 | `KUDIG_AI_*` 环境变量 | 支持多提供商 (OpenAI/Qwen/Ollama) |
| 功能 | - | 诊断摘要生成、根因分析、修复建议、中英文支持 |

**环境变量配置**:
```bash
export KUDIG_AI_PROVIDER="openai"      # 或 qwen, ollama
export KUDIG_AI_API_KEY="sk-xxx"
export KUDIG_AI_MODEL="gpt-4"
export KUDIG_AI_LANGUAGE="zh"          # 或 en
```

**目录结构**:
```
pkg/ai/
├── provider.go       # AI Provider 接口与 OpenAI 实现
├── assistant.go      # AI 助手封装
└── provider_test.go  # 单元测试
```

---

## 依赖项跟踪

| 功能 | 新增依赖 | 备注 |
|------|---------|------|
| Prometheus Metrics | `github.com/prometheus/client_golang` | ✅ 已添加 |
| 进度条 | `github.com/schollz/progressbar/v3` | ✅ 已添加 |
| DNS 诊断 | `k8s.io/dns` (可选) | 标准库为主 |
| GPU 诊断 | `github.com/NVIDIA/go-nvml` | NVIDIA 官方库 |
| Operator | `sigs.k8s.io/controller-runtime` | Kubernetes Operator 框架 |
| eBPF | `github.com/cilium/ebpf` | 内核 4.18+ |
| LLM | `github.com/sashabaranov/go-openai` | OpenAI SDK |

---

## 验收总清单

### 已完成 (Phase 1-7) ✅ 全部完成！
- [x] Prometheus Metrics 正常暴露
- [x] 多节点并发诊断完成
- [x] HTML Reporter 可用
- [x] Control Plane 分析器补充 (etcd, scheduler, controller-manager, apiserver-latency)
- [x] Workload 分析器补充 (deployment, statefulset, daemonset, pdb)
- [x] kubectl 插件可用
- [x] 历史数据对比可用
- [x] Webhook 通知可用
- [x] DNS 诊断分析器 (coredns, dns.pod_config)
- [x] 存储性能分析器 (pvc, pv, storageclass, volumeattachment, csi, storage.pod)
- [x] GPU/NPU 诊断 (gpu.node, gpu.pod, gpu.share, npu)
- [x] CIS 安全合规扫描 (cis.apiserver, cis.etcd, cis.kubelet, cis.pod, cis.network, cis.secret)
- [x] RBAC 审计分析器 (rbac.admin, rbac.serviceaccount, rbac.dangerous, rbac.unused, rbac.token)
- [x] Operator 模式 (ClusterDiagnostic, NodeDiagnostic, DiagnosticSchedule CRD + Controllers + Helm Chart)
- [x] eBPF 深度诊断 (ebpf.tcp, ebpf.dns, ebpf.fileio)
- [x] AI/LLM 辅助诊断 (OpenAI/Qwen/Ollama 支持)

**最终统计**:
- **分析器总数**: 68 个
- **CRD 数量**: 3 个
- **Operator Controllers**: 3 个
- **Helm Charts**: 1 套
- **AI Provider 支持**: 3 种 (OpenAI, Qwen, Ollama)

### 🎉 开发任务全部完成！
