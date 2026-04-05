# 项目结构说明

## 更新时间
2026-04-02（功能特性 100% 完成）

## 重组目的
明确区分 Bash v1.0（生产可用）和 Go v2.0（生产可用，功能完整）版本，避免用户混淆。

## 新结构

```
kudig/
├── v1-bash/                    # v1.0 Bash 版本（生产可用）
│   ├── kudig.sh                # 主脚本
│   ├── README.md               # v1.0 完整文档
│   ├── TESTING.md              # 测试说明
│   └── reference/              # 示例诊断数据
│       └── diagnose_k8s/
│
├── v2-go/                      # v2.0 Go 版本（生产可用，功能 100%）
│   ├── cmd/kudig/              # CLI 入口（18+ 命令）
│   ├── pkg/                    # 核心包
│   │   ├── analyzer/           # 70+ 分析器（9大类）
│   │   │   ├── system/         # 系统分析器（7个）
│   │   │   ├── process/        # 进程分析器（5个）
│   │   │   ├── network/        # 网络分析器（5个）
│   │   │   ├── kernel/         # 内核分析器（5个）
│   │   │   ├── kubernetes/     # K8s 分析器（21个）
│   │   │   ├── runtime/        # 运行时分析器（4个）
│   │   │   ├── security/       # 安全分析器（11个）- CIS + RBAC
│   │   │   ├── ebpf/           # eBPF 分析器（3个）- TCP/DNS/文件I/O
│   │   │   └── servicemesh/    # 服务网格分析器（2个）- Istio/Linkerd
│   │   ├── collector/          # 数据收集层
│   │   │   ├── offline/        # 离线收集器
│   │   │   └── online/         # 在线收集器（支持并发）
│   │   ├── reporter/           # 报告生成层（Text/JSON/HTML/SARIF）
│   │   ├── rules/              # YAML 规则引擎
│   │   ├── types/              # 类型定义
│   │   ├── ai/                 # AI/LLM 辅助诊断
│   │   │   ├── provider.go     # AI Provider 接口
│   │   │   └── assistant.go    # AI 助手封装
│   │   ├── ebpf/               # eBPF 深度诊断
│   │   │   ├── probe/          # eBPF 探针管理
│   │   │   └── analyzer/       # eBPF 分析器
│   │   ├── rca/                # 根因分析 (Root Cause Analysis)
│   │   ├── autofix/            # 自动修复引擎
│   │   ├── cost/               # 成本分析
│   │   ├── scanner/            # 镜像安全扫描（Trivy 集成）
│   │   ├── tui/                # TUI 交互界面（bubbletea）
│   │   ├── legacy/             # v1.0 兼容层
│   │   ├── metrics/            # Prometheus Metrics
│   │   ├── history/            # 历史数据对比
│   │   └── notifier/           # Webhook 通知
│   ├── operator/               # Kubernetes Operator
│   │   ├── api/v1/             # CRD 类型定义
│   │   │   ├── clusterdiagnostic_types.go
│   │   │   ├── nodediagnostic_types.go
│   │   │   └── schedule_types.go
│   │   ├── controllers/        # 控制器实现
│   │   │   ├── clusterdiagnostic_controller.go
│   │   │   ├── nodediagnostic_controller.go
│   │   │   └── schedule_controller.go
│   │   ├── cmd/                # Operator 入口
│   │   ├── helm/               # Helm Chart
│   │   │   └── kudig-operator/
│   │   └── config/examples/    # 示例 CR
│   ├── build/                  # 构建输出
│   ├── charts/kudig/           # Helm Chart（CLI 部署）
│   ├── deployments/            # K8s 部署配置
│   ├── configs/                # 配置文件
│   ├── rules/                  # 规则示例
│   ├── Dockerfile              # Docker 构建
│   ├── Makefile               # 构建脚本
│   ├── go.mod
│   └── README.md               # v2.0 完整文档
│
├── docs/                       # 项目文档
│   ├── FUNCTIONAL_ROADMAP.md   # 功能路线图
│   ├── CNCF_GRADUATION_GAP_ANALYSIS.md  # CNCF 差距分析
│   ├── CODE_QUALITY_README.md
│   └── QUALITY_CHECK_SETUP.md
│
├── reference/                  # 共享参考数据（根目录）
│   └── diagnose_k8s/
│
├── README.md                   # 项目主文档（导航页）
├── TEST_REPORT.md              # 测试报告
├── TESTING.md                  # 原测试文档（保留）
├── ROADMAP.md                  # 项目路线图
├── CONTRIBUTING.md             # 贡献指南
├── CODE_OF_CONDUCT.md          # 行为准则
├── SECURITY.md                 # 安全政策
├── LICENSE                     # 许可证
└── kudig.sh                    # 根目录主脚本（保留兼容性）
```

## 文件说明

### v1-bash/
- 复制 `kudig.sh` → `v1-bash/kudig.sh`
- 复制 `TESTING.md` → `v1-bash/TESTING.md`
- 复制 `reference/` → `v1-bash/reference/`
- 创建 `v1-bash/README.md`

### v2-go/
- 移动 `cmd/` → `v2-go/cmd/`
- 移动 `pkg/` → `v2-go/pkg/`
- 移动 `internal/` → `v2-go/internal/`
- 移动 `go.mod` → `v2-go/go.mod`
- 移动 `go.sum` → `v2-go/go.sum`
- 移动 `Makefile` → `v2-go/Makefile`
- 移动 `Dockerfile` → `v2-go/Dockerfile`
- 移动 `charts/` → `v2-go/charts/`
- 移动 `deployments/` → `v2-go/deployments/`
- 移动 `configs/` → `v2-go/configs/`
- 移动 `rules/` → `v2-go/rules/`
- 创建 `operator/` - Kubernetes Operator 完整实现
- 创建 `pkg/ebpf/` - eBPF 深度诊断
- 创建 `pkg/ai/` - AI/LLM 辅助诊断
- 创建 `pkg/rca/` - 根因分析
- 创建 `pkg/autofix/` - 自动修复引擎
- 创建 `pkg/cost/` - 成本分析
- 创建 `pkg/scanner/` - 镜像安全扫描
- 创建 `pkg/tui/` - TUI 交互界面
- 创建 `pkg/metrics/` - Prometheus 指标
- 创建 `pkg/history/` - 历史数据对比
- 创建 `pkg/notifier/` - Webhook 通知
- 创建 `pkg/analyzer/servicemesh/` - 服务网格分析器
- 创建 `v2-go/README.md`

### 根目录
- 更新 `README.md` - 改为导航页，明确版本区分
- 保留 `kudig.sh` - 向后兼容
- 保留 `TESTING.md` - 原有文档
- 保留 `reference/` - 共享参考数据

## 版本标识

### v1.0 Bash - 生产可用
- 目录：`v1-bash/`
- 状态：稳定，可用于生产环境
- 特性：80+ 异常检测规则，离线分析模式

### v2.0 Go - 生产可用（功能特性 100%）
- 目录：`v2-go/`
- 状态：生产阶段，功能完整
- 特性：
  - **70+ 分析器**（9 大类：系统、进程、网络、内核、K8s、运行时、安全、eBPF、服务网格）
  - 双模式（离线 + 在线）
  - 多节点并发诊断
  - YAML 规则引擎
  - **Operator 模式**（3 CRD + 控制器）
  - **eBPF 深度诊断**
  - **AI/LLM 辅助**
  - **TUI 交互模式**
  - **根因分析 (RCA)**
  - **自动修复引擎**
  - **SARIF 安全报告**
  - **Grafana Dashboard**
  - **镜像安全扫描**
  - **成本分析**
  - **多集群联邦诊断**
  - **Shell 自动补全**
  - Prometheus Metrics
  - 历史数据对比
  - Webhook 通知
  - kubectl 插件

## CLI 命令（v2.0）

| 命令 | 功能 | 状态 |
|------|------|------|
| `offline` | 离线分析诊断数据 | ✅ |
| `online` | 在线实时诊断 | ✅ |
| `analyze` | 离线分析（alias） | ✅ |
| `legacy` | v1.0 兼容模式 | ✅ |
| `rules` | YAML 规则引擎 | ✅ |
| `list-analyzers` | 列出分析器 | ✅ |
| `history` | 历史数据管理 | ✅ |
| `tui` | TUI 交互模式 | ✅ 新增 |
| `rca` | 根因分析 | ✅ 新增 |
| `fix` | 自动修复 | ✅ 新增 |
| `cost` | 成本分析 | ✅ 新增 |
| `scan` | 镜像安全扫描 | ✅ 新增 |
| `grafana` | Grafana Dashboard 导出 | ✅ 新增 |
| `pprof` | 性能剖析 | ✅ 新增 |
| `trace` | 分布式追踪 | ✅ 新增 |
| `multicluster` | 多集群诊断 | ✅ 新增 |
| `completion` | Shell 补全 | ✅ 新增 |

## 使用指南

### 使用 v1.0 Bash（生产推荐）
```bash
cd v1-bash
./kudig.sh /tmp/diagnose_1702468800
```

### 使用 v2.0 Go（功能完整）
```bash
cd v2-go
make build
./build/kudig offline /tmp/diagnose_1702468800
./build/kudig online --all-nodes
./build/kudig rules --file rules/custom.yaml /tmp/diagnose_1702468800
./build/kudig tui                    # TUI 模式
./build/kudig rca                    # 根因分析
./build/kudig cost                   # 成本分析
./build/kudig scan nginx:latest      # 镜像扫描
./build/kudig grafana > dashboard.json  # Grafana Dashboard
source <(./build/kudig completion bash) # Shell 补全
```

### 使用 Operator 模式
```bash
cd v2-go/operator
helm install kudig-operator ./helm/kudig-operator -n kudig-system --create-namespace
kubectl apply -f config/examples/cluster-diagnostic.yaml
```

## 文档更新

- ✅ 根目录 `README.md` - 项目导航页
- ✅ `v1-bash/README.md` - v1.0 完整文档
- ✅ `v2-go/README.md` - v2.0 开发文档（包含 Phase 4-8 全部功能）
- ✅ `STRUCTURE.md` - 项目结构说明（本文档）
- ✅ `docs/CNCF_GRADUATION_GAP_ANALYSIS.md` - CNCF 毕业差距分析
- ✅ 各 README 相互链接

## 向后兼容

根目录保留了 `kudig.sh` 和 `TESTING.md`，确保：
- 现有脚本路径仍可用
- 现有文档链接不失效
- 用户可以平滑过渡

## 优势

1. **清晰的版本区分**：目录名明确表示版本和状态
2. **独立的文档**：每个版本都有完整的 README
3. **降低混淆**：用户一眼就能看出哪个版本可用
4. **独立开发**：v2.0 可以独立演进，不影响 v1.0
5. **易于维护**：各版本文件完全分离
6. **功能完整**：v2.0 包含所有计划功能（100%）

## 分析器统计（v2.0）

| 类别 | 数量 | 说明 |
|------|------|------|
| system | 7 | CPU、内存、磁盘、Swap、文件句柄等 |
| process | 5 | Kubelet、容器运行时、PID 泄漏等 |
| network | 5 | 连接追踪、网卡、路由、iptables 等 |
| kernel | 5 | Panic、OOM、文件系统、模块等 |
| kubernetes | 21 | Pod、Node、Deployment、DNS、存储、GPU 等 |
| runtime | 4 | Docker、Containerd 状态检查 |
| security | 11 | CIS 合规、RBAC 审计 |
| ebpf | 3 | TCP、DNS、文件 I/O 实时追踪 |
| servicemesh | 2 | Istio、Linkerd 服务网格诊断 |
| **总计** | **70+** | 覆盖 K8s 节点诊断全场景 |

## 注意事项

- 根目录的 `kudig.sh` 是 v1.0 的副本，保持同步更新
- `scripts/kudig.sh` 可能是备份，需要确认是否保留
- Git 历史保持完整，可以追溯文件移动
- v2.0 Operator 需要 Kubernetes 1.21+ 集群
- eBPF 功能需要 Linux 内核 4.18+ 和 CAP_BPF 权限
- AI 功能需要配置相应的 API 密钥
- 镜像扫描功能需要安装 Trivy
