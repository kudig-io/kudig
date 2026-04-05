# kudig v2.0 - Go 版本

[![codecov](https://codecov.io/gh/kudig/kudig/branch/main/graph/badge.svg)](https://codecov.io/gh/kudig/kudig)
[![CI](https://github.com/kudig/kudig/actions/workflows/code-quality.yml/badge.svg)](https://github.com/kudig/kudig/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kudig/kudig)](https://goreportcard.com/report/github.com/kudig/kudig)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

> ✅ **Production 可用** - Go 语言重构版本，功能完整，生产就绪

## 简介

`kudig` v2.0 是使用 Go 语言重构的下一代 Kubernetes 节点诊断工具，在 v1.0 Bash 版本基础上进行了全面升级。

## 核心特性

- **Go 语言实现**：性能提升，跨平台支持
- **双模式支持**：
  - 离线模式：分析 diagnose_k8s.sh 数据
  - 在线模式：实时诊断 K8s 集群（通过 K8s API）
- **68 内置分析器**：涵盖系统、进程、网络、内核、Kubernetes、运行时、安全、eBPF 等维度
- **YAML 规则引擎**：支持自定义诊断规则
- **Kubernetes 原生部署**：提供 Helm Chart 和 DaemonSet 支持
- **Docker 镜像**：开箱即用的容器化部署

## 开发状态

当前版本已完成所有主要功能实现：

- ✅ 项目结构框架
- ✅ 基础类型定义
- ✅ 68 分析器实现
- ✅ 离线数据收集层
- ✅ 在线数据收集层（K8s API）
- ✅ 报告生成层（Text/JSON）
- ✅ YAML 规则引擎
- ✅ Helm Chart
- ✅ Dockerfile
- ✅ 完整测试
- ✅ 代码质量检查

## 项目结构

```
v2-go/
├── cmd/kudig/          # CLI 入口
├── pkg/
│   ├── analyzer/       # 分析器框架
│   │   ├── system/     # 系统分析器
│   │   ├── process/    # 进程分析器
│   │   ├── network/    # 网络分析器
│   │   ├── kernel/     # 内核分析器
│   │   ├── kubernetes/ # K8s 分析器
│   │   └── runtime/    # 运行时分析器
│   ├── collector/      # 数据收集层
│   │   ├── offline/    # 离线收集器
│   │   └── online/     # 在线收集器
│   ├── reporter/       # 报告生成层
│   ├── rules/          # 规则引擎
│   ├── types/          # 公共类型
│   └── legacy/         # v1.0 兼容层
├── charts/             # Helm Chart
├── rules/              # 示例规则
├── Dockerfile          # Docker 构建文件
├── Makefile           # 构建脚本
├── go.mod
└── go.sum
```

## 构建

### 使用 Make（Linux/macOS）

```bash
# 下载依赖
make deps

# 构建
make build

# 运行测试
make test

# 构建所有平台
make build-all
```

### 直接使用 Go 命令（跨平台）

```bash
# 下载依赖
go mod tidy

# 构建
go build -buildvcs=false -o kudig ./cmd/kudig

# 构建 Windows 版本
go build -buildvcs=false -o kudig.exe ./cmd/kudig

# 运行测试
go test ./pkg/...
```

## 使用方法

### 离线模式

```bash
# 分析诊断数据
kudig offline /tmp/diagnose_1702468800

# 详细模式
kudig offline -v /tmp/diagnose_1702468800

# JSON 格式输出
kudig offline --format json /tmp/diagnose_1702468800

# 保存到文件
kudig offline -o report.txt /tmp/diagnose_1702468800
```

### 在线模式

```bash
# 使用默认 kubeconfig
kudig online

# 指定节点
kudig online --node worker-1

# 检查所有节点
kudig online --all-nodes

# 指定 kubeconfig
kudig online --kubeconfig ~/.kube/config
```

### 规则模式

```bash
# 使用自定义规则文件
kudig rules --file rules/custom.yaml /tmp/diagnose_1702468800

# 使用规则目录
kudig rules --dir rules/ /tmp/diagnose_1702468800

# 列出所有规则
kudig rules --list
```

### 兼容模式

```bash
# 使用原版 kudig.sh 脚本（需要 bash）
kudig legacy /tmp/diagnose_1702468800
```

### 列出分析器

```bash
kudig list-analyzers
```

## 开发

### 环境要求

- Go 1.25+
- Make（可选，用于简化构建）
- golangci-lint（代码质量检查）

### 快速开始

```bash
# 克隆项目
git clone https://github.com/kudig/kudig.git
cd kudig/v2-go

# 下载依赖
make deps

# 构建项目
make build

# 运行测试
make test

# 生成覆盖率报告
make test-coverage
```

### 开发工具

```bash
# 代码格式化
make fmt

# 代码检查
make vet

# 运行 linter
make lint

# 自动修复 lint 问题
make lint-fix

# 依赖管理
make tidy
```

### 添加新的分析器

1. 在 `pkg/analyzer/<category>/` 目录下创建新的分析器文件
2. 实现 `analyzer.Analyzer` 接口
3. 在 `init()` 函数中注册分析器

示例：

```go
package system

import (
    "context"
    "github.com/kudig/kudig/pkg/analyzer"
    "github.com/kudig/kudig/pkg/types"
)

type MyAnalyzer struct {
    *analyzer.BaseAnalyzer
}

func NewMyAnalyzer() *MyAnalyzer {
    return &MyAnalyzer{
        BaseAnalyzer: analyzer.NewBaseAnalyzer(
            "system.my_analyzer",  // 唯一标识
            "检查系统状态",         // 中文描述
            "system",              // 分类
            []types.DataMode{types.ModeOffline, types.ModeOnline}, // 支持的模式
        ),
    }
}

func (a *MyAnalyzer) Analyze(ctx context.Context, data *types.DiagnosticData) ([]types.Issue, error) {
    var issues []types.Issue
    
    // 实现分析逻辑
    if /* 检测到问题 */ {
        issue := types.NewIssue(
            types.SeverityWarning,
            "问题中文名",
            "ISSUE_CODE",
            "问题详情",
            "数据来源",
        ).WithRemediation("修复建议")
        issue.AnalyzerName = a.Name()
        issues = append(issues, *issue)
    }
    
    return issues, nil
}

func init() {
    _ = analyzer.Register(NewMyAnalyzer())
}
```

## 新增功能（Phase 4-7）

### Phase 4: 核心组件深度覆盖 ✅

#### DNS 诊断
- **CoreDNS 健康检查**：Pod 状态、CrashLoopBackOff、镜像拉取失败
- **DNS 配置分析**：Pod DNS 策略、ndots 配置、Default DNS 使用统计

#### 存储分析
- **PVC/PV 状态检查**：Pending、Lost、Released 状态检测
- **StorageClass 配置**：默认 SC 检查、回收策略、绑定模式
- **CSI Driver 监控**：CSI 插件 Pod 健康检查
- **VolumeAttachment**：存储挂载失败检测

#### GPU/NPU 诊断
- **NVIDIA GPU**：设备插件检查、资源分配、MIG 模式
- **华为 Ascend NPU**：Ascend 设备插件状态监控

### Phase 5: 安全合规 ✅

#### CIS 安全扫描
- **API Server**：匿名认证、不安全端口、绑定地址检查
- **etcd**：客户端证书认证、自动 TLS 检测
- **Kubelet**：匿名认证、只读端口配置
- **Pod Security**：特权容器、HostPID、HostNetwork、root 用户运行
- **Network Policy**：缺少 Network Policy 的命名空间检测
- **Secret 管理**：default 命名空间 Secret 存放检查

#### RBAC 审计
- **权限分析**：cluster-admin 绑定数量监控
- **危险权限**：escalate、create pods、secrets 访问检测
- **ServiceAccount**：default SA 自动挂载 Token 检查
- **未使用角色**：未使用的 ClusterRole 清理建议

### Phase 6: Operator 模式 ✅

#### Kubernetes Operator
- **ClusterDiagnostic CRD**：集群级诊断任务
- **NodeDiagnostic CRD**：节点级诊断任务（支持节点选择器）
- **DiagnosticSchedule CRD**：定时诊断调度（@hourly/@daily/@weekly）
- **Helm Chart**：完整的 Operator 部署包

```bash
# 安装 Operator
helm install kudig-operator ./operator/helm/kudig-operator -n kudig-system --create-namespace

# 创建集群诊断
kubectl apply -f operator/config/examples/cluster-diagnostic.yaml

# 创建定时任务
kubectl apply -f operator/config/examples/schedule.yaml
```

### Phase 7: 创新功能 ✅

#### eBPF 深度诊断（内核 4.18+）
- **TCP 分析**：连接追踪、重传检测、延迟分析
- **DNS 分析**：查询追踪、失败率分析、延迟分析
- **文件 I/O**：文件操作追踪、I/O 延迟分析

```bash
# eBPF 诊断自动启用（需要 CAP_BPF 权限）
./build/kudig online --all-nodes
```

#### AI/LLM 辅助诊断
- **智能分析**：诊断结果摘要、根因分析
- **修复建议**：自动推荐修复命令和步骤
- **多语言支持**：中文/英文自动切换
- **多提供商**：支持 OpenAI、阿里云通义千问、私有化部署（Ollama）

```bash
# 配置 AI（环境变量）
export KUDIG_AI_PROVIDER="openai"  # 或 qwen, ollama
export KUDIG_AI_API_KEY="sk-xxx"
export KUDIG_AI_LANGUAGE="zh"      # 或 en

# 使用 AI 分析
./build/kudig online --ai-analysis
```

### Phase 8: 功能完善 ✅ (功能特性 100%)

#### TUI 交互模式
- **bubbletea 终端 UI**：直观的菜单驱动界面
- **实时诊断进度**：诊断进度可视化
- **交互式结果浏览**：支持键盘导航和详情查看

```bash
# 启动 TUI 模式
kudig tui
```

#### 服务网格诊断
- **Istio 诊断**：istiod、ingress/egress gateway、sidecar 状态
- **Linkerd 诊断**：控制平面、proxy 状态、性能指标
- **mTLS 检查**：证书状态、配置一致性

```bash
# 服务网格诊断（自动检测）
kudig online
```

#### 根因分析 (RCA)
- **智能关联**：多症状关联推导根因
- **置信度评估**：每个根因附带置信度评分
- **修复建议**：针对根因提供系统级修复方案

```bash
# 执行根因分析
kudig rca

# 离线模式 RCA
kudig rca /tmp/diagnose_1702468800
```

#### 自动修复引擎
- **安全修复**：低风险修复自动执行
- **分级确认**：高风险操作需要用户确认
- **修复回滚**：支持修复操作记录和回滚

```bash
# 查看可修复问题（干跑模式）
kudig fix --dry-run

# 执行修复
kudig fix --confirm
```

#### SARIF 安全报告
- **GitHub/CodeQL 兼容**：标准 SARIF 2.1.0 格式
- **安全扫描集成**：与 CI/CD 安全扫描流程集成
- **漏洞追踪**：支持自动修复建议

```bash
# 生成 SARIF 报告
kudig online --format sarif --output report.sarif
```

#### Grafana Dashboard
- **官方 Dashboard**：预配置的 Grafana JSON
- **多维可视化**：问题分布、趋势分析、资源监控
- **Prometheus 集成**：与 metrics 服务无缝集成

```bash
# 导出 Grafana Dashboard
kudig grafana > kudig-dashboard.json
```

#### 镜像安全扫描
- **Trivy 集成**：自动检测镜像 CVE
- **多严重级别**：CRITICAL/HIGH/MEDIUM/LOW 分级
- **修复建议**：提供漏洞修复版本建议

```bash
# 扫描指定镜像
kudig scan nginx:latest

# 扫描集群所有镜像
kudig scan --all-images
```

#### 成本分析
- **资源成本估算**：基于 AWS/GCP/Azure 定价
- **优化建议**：识别资源浪费和优化机会
- **多维度分析**：按命名空间、工作负载分类

```bash
# 成本分析
kudig cost
```

#### 性能剖析 (pprof)
- **CPU Profile**：诊断性能热点
- **内存分析**：检测内存泄漏
- **Goroutine 追踪**：并发问题分析

```bash
# 启动 pprof 服务器
kudig pprof --port 6060
```

#### 分布式追踪
- **OpenTelemetry 支持**：追踪诊断操作
- **Jaeger 集成**：导出追踪数据到 Jaeger
- **性能瓶颈定位**：分析诊断流程耗时

```bash
# 启用追踪
kudig trace --jaeger http://localhost:14268
```

#### 多集群联邦诊断
- **跨集群诊断**：同时诊断多个 K8s 集群
- **统一视图**：汇总跨集群问题
- **上下文管理**：支持 kubeconfig 多上下文

```bash
# 诊断所有上下文
kudig multicluster --all-contexts

# 诊断指定上下文
kudig multicluster --contexts prod-cluster,dr-cluster
```

#### Shell 自动补全
- **Bash/Zsh/Fish 支持**：主流 Shell 全兼容
- **命令补全**：子命令和标志自动补全
- **动态提示**：根据上下文提供智能提示

```bash
# Bash
source <(kudig completion bash)

# Zsh
source <(kudig completion zsh)

# Fish
kudig completion fish | source
```

## 高级功能配置

### 环境变量

| 变量 | 说明 | 示例 |
|------|------|------|
| `KUDIG_AI_PROVIDER` | AI 提供商 | `openai`, `qwen`, `ollama` |
| `KUDIG_AI_API_KEY` | API 密钥 | `sk-...` |
| `KUDIG_AI_MODEL` | 模型名称 | `gpt-4` |
| `KUDIG_AI_LANGUAGE` | 输出语言 | `zh`, `en` |
| `KUDIG_SLACK_WEBHOOK_URL` | Slack 通知 | `https://hooks.slack.com/...` |
| `KUDIG_DINGTALK_WEBHOOK_URL` | 钉钉通知 | `https://oapi.dingtalk.com/...` |

### Prometheus Metrics

```bash
# 启动 metrics 服务
./build/kudig online --serve --metrics-port 9090

# 访问指标
curl http://localhost:9090/metrics
```

## 开发路线图

### v2.0 ✅ 已完成 (功能特性 100%)

#### 核心功能
- [x] 完成离线分析模式
- [x] 实现 **70+ 分析器**（9 大类：系统、进程、网络、内核、K8s、运行时、安全、eBPF、服务网格）
- [x] 生成文本/JSON/HTML/**SARIF** 报告
- [x] 兼容 v1.0 数据格式
- [x] 添加在线诊断模式
- [x] 实现 YAML 规则引擎

#### 部署与集成
- [x] 添加 Helm Chart
- [x] Dockerfile 构建
- [x] Prometheus Metrics 支持
- [x] **Grafana Dashboard** 官方支持
- [x] kubectl 插件

#### 诊断能力
- [x] 多节点并发诊断
- [x] 历史数据对比
- [x] Webhook 通知（Slack/钉钉/企业微信）
- [x] **根因分析 (RCA)** 智能关联
- [x] **自动修复引擎** 安全修复

#### 高级功能
- [x] **Operator 模式**（3 CRD + 控制器）
- [x] **DNS/存储/GPU 诊断**
- [x] **CIS 安全/RBAC 审计**
- [x] **eBPF 深度诊断**
- [x] **AI/LLM 辅助诊断**
- [x] **服务网格诊断**（Istio/Linkerd）
- [x] **镜像安全扫描**（Trivy 集成）
- [x] **成本分析** 资源成本估算

#### 用户体验
- [x] **TUI 交互模式**（bubbletea）
- [x] **Shell 自动补全**（Bash/Zsh/Fish）
- [x] **性能剖析**（pprof 支持）
- [x] **分布式追踪**（OpenTelemetry）
- [x] **多集群联邦诊断**

#### 质量保障
- [x] 完善错误处理
- [x] 性能优化
- [x] 完整文档
- [x] 生产环境测试
- [x] 正式发布

## 贡献

v2.0 版本欢迎社区贡献。如果您想参与开发，请查看以下文档：

- [贡献指南](../CONTRIBUTING.md) - 详细的贡献流程和规范
- [行为准则](../CODE_OF_CONDUCT.md) - 社区行为准则
- [安全政策](../SECURITY.md) - 安全漏洞报告流程

快速开始：

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

## 社区

- [GitHub Issues](https://github.com/kudig/kudig/issues) - Bug 报告和功能建议
- [GitHub Discussions](https://github.com/kudig/kudig/discussions) - 一般性讨论和问答

## 许可证

Apache License 2.0

## 相关链接

- [v1.0 Bash 版本](../v1-bash/) - ✅ 生产可用
- [项目主页](../)

---

**注意**: 此版本已达到生产就绪状态，所有主要功能都已实现并经过测试。
