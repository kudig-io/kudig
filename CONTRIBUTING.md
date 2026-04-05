# 贡献指南

感谢您对 kudig 项目的关注！我们欢迎各种形式的贡献，包括但不限于：

- 提交 bug 报告
- 提出新功能建议
- 改进文档
- 提交代码修复
- 分享使用经验

## 快速开始

1. **Fork 本仓库**
2. **克隆您的 fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/kudig.git
   cd kudig
   ```
3. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

## 开发环境设置

### v2-go (Go 版本)

```bash
cd v2-go

# 安装依赖
make deps

# 运行测试
make test

# 构建项目
make build
```

### v1-bash (Bash 版本)

```bash
cd v1-bash

# 运行测试
bash kudig --help
```

## 代码规范

### Go 代码

- 遵循 [Effective Go](https://go.dev/doc/effective_go)
- 使用 `go fmt` 格式化代码
- 通过 `golangci-lint` 检查
- 保持测试覆盖率不低于当前水平

### Bash 代码

- 遵循 [Google Shell Style Guide](https://google.github.io/styleguide/shellguide.html)
- 使用 `shellcheck` 检查脚本
- 添加适当的注释

## 提交规范

### Commit Message 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型 (type):**
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式（不影响功能的修改）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

**示例:**
```
feat(analyzer): 添加 CPU 负载分析器

实现 CPU 负载检测功能，支持离线模式分析。

Closes #123
```

## Pull Request 流程

1. **确保测试通过**
   ```bash
   make test
   make lint
   ```

2. **更新文档**
   - 如果添加了新功能，更新 README.md
   - 如果需要，更新 ROADMAP.md

3. **提交 PR**
   - 提供清晰的 PR 描述
   - 关联相关的 issue
   - 确保 CI 检查通过

4. **代码审查**
   - 维护者会进行代码审查
   - 根据反馈进行修改
   - 合并后会关闭 PR

## 报告 Bug

请使用 [GitHub Issues](https://github.com/kudig/kudig/issues) 报告 bug，并包含以下信息：

- **问题描述**: 清晰描述问题
- **复现步骤**: 详细的复现步骤
- **期望结果**: 期望的行为
- **实际结果**: 实际的行为
- **环境信息**:
  - OS 版本
  - kudig 版本
  - Go 版本（如果是 v2-go）
  - 相关日志

## 提出新功能

1. 先搜索现有 issues，避免重复
2. 使用 "feature request" 标签创建 issue
3. 描述功能的用途和预期行为
4. 讨论实现方案

## 分析器开发指南

### 目录结构

```
pkg/analyzer/
├── system/       # 系统资源分析器（7个）
├── network/      # 网络分析器（5个）
├── process/      # 进程分析器（5个）
├── kernel/       # 内核分析器（5个）
├── kubernetes/   # K8s 组件分析器（21个）
├── runtime/      # 运行时分析器（4个）
├── security/     # 安全分析器（11个）- CIS + RBAC
└── ebpf/         # eBPF 分析器（3个）- TCP/DNS/文件I/O
```

### 分析器接口

```go
type Analyzer interface {
    Name() string
    Description() string
    Category() string
    Analyze(ctx context.Context, data *types.DiagnosticData) ([]types.Issue, error)
    SupportedModes() []types.DataMode
    Dependencies() []string
}
```

### 最小实现示例

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
            "system.my_analyzer",
            "检查系统状态",
            "system",
            []types.DataMode{types.ModeOffline, types.ModeOnline},
        ),
    }
}

func (a *MyAnalyzer) Analyze(_ context.Context, data *types.DiagnosticData) ([]types.Issue, error) {
    var issues []types.Issue
    
    // 实现分析逻辑
    content, ok := data.GetRawFile("system_info")
    if !ok {
        return issues, nil
    }
    
    // 检测问题
    if strings.Contains(string(content), "error") {
        issue := types.NewIssue(
            types.SeverityWarning,
            "检测到问题",
            "MY_ISSUE_CODE",
            "问题详情描述",
            "system_info",
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

## 新增分析器类别指南

### Security 分析器

安全分析器位于 `pkg/analyzer/security/`，包含 CIS 合规和 RBAC 审计：

```go
package security

import (
    "context"
    "github.com/kudig/kudig/pkg/analyzer"
    "github.com/kudig/kudig/pkg/types"
)

type MySecurityAnalyzer struct {
    *analyzer.BaseAnalyzer
}

func NewMySecurityAnalyzer() *MySecurityAnalyzer {
    return &MySecurityAnalyzer{
        BaseAnalyzer: analyzer.NewBaseAnalyzer(
            "security.cis.my_check",  // 或 security.rbac.xxx
            "检查 CIS 规范",
            "security",
            []types.DataMode{types.ModeOnline}, // 安全分析通常需要在线模式
        ),
    }
}

func (a *MySecurityAnalyzer) Analyze(ctx context.Context, data *types.DiagnosticData) ([]types.Issue, error) {
    var issues []types.Issue
    
    if !data.HasK8sClient() {
        return issues, nil
    }
    
    // 实现安全检测逻辑
    // 使用 data.K8sClient 访问 K8s API
    
    return issues, nil
}

func init() {
    _ = analyzer.Register(NewMySecurityAnalyzer())
}
```

### eBPF 分析器

eBPF 分析器位于 `pkg/ebpf/analyzer/`，需要内核 4.18+：

```go
package analyzer

import (
    "context"
    "github.com/kudig/kudig/pkg/analyzer"
    "github.com/kudig/kudig/pkg/ebpf/probe"
    "github.com/kudig/kudig/pkg/types"
)

type MyEBPFAnalyzer struct {
    *analyzer.BaseAnalyzer
}

func NewMyEBPFAnalyzer() *MyEBPFAnalyzer {
    return &MyEBPFAnalyzer{
        BaseAnalyzer: analyzer.NewBaseAnalyzer(
            "ebpf.my_feature",
            "使用 eBPF 分析 XXX",
            "ebpf",
            []types.DataMode{types.ModeOnline},
        ),
    }
}

func (a *MyEBPFAnalyzer) Analyze(ctx context.Context, data *types.DiagnosticData) ([]types.Issue, error) {
    var issues []types.Issue
    
    // 检查权限
    if !a.hasPermission() {
        return issues, nil
    }
    
    // 创建 eBPF 探针
    config := &probe.Config{
        EnableTCP:     false,
        EnableDNS:     false,
        // ... 配置
    }
    
    probeMgr, err := probe.NewProbeManager(config)
    if err != nil {
        return issues, nil // 静默跳过
    }
    
    // 启动探针收集数据...
    
    return issues, nil
}
```

### Operator 开发

Operator 位于 `v2-go/operator/`，包含 CRD 和控制器：

```
operator/
├── api/v1/              # CRD 类型定义
├── controllers/         # 控制器实现
├── helm/               # Helm Chart
└── config/examples/    # 示例配置
```

添加新 CRD：
1. 在 `api/v1/` 创建类型定义文件
2. 在 `controllers/` 创建控制器
3. 在 `cmd/main.go` 注册控制器
4. 更新 Helm Chart

### AI Provider 开发

AI Provider 位于 `pkg/ai/`，支持多种 LLM 后端：

```go
package ai

// 实现 Provider 接口
type MyProvider struct {
    config *Config
}

func (p *MyProvider) Analyze(ctx context.Context, issues []types.Issue, hostname string) (*AnalysisResult, error) {
    // 调用 LLM API
}

func (p *MyProvider) Name() string {
    return "my_provider"
}
```

在 `provider.go` 的 Factory 中添加新 provider：
```go
case "my_provider":
    return NewMyProvider(config)
```

## 测试要求

- 新功能必须包含单元测试
- 测试覆盖率不应低于当前水平（目标 60%+）
- 使用表格驱动测试
- 测试函数命名: `Test<功能名>`
- eBPF 测试在 CI 环境可能跳过（需要内核支持）
- AI 测试使用 mock（避免依赖外部 API）

## 文档更新

修改代码时，请同步更新相关文档：

### 必须更新的文档

| 文件 | 何时更新 | 更新内容 |
|------|----------|----------|
| `README.md` | 新增功能 | 功能列表、使用示例 |
| `v2-go/README.md` | 新增分析器/功能 | 详细说明、配置指南 |
| `STRUCTURE.md` | 目录变更 | 项目结构说明 |
| `docs/FUNCTIONAL_ROADMAP.md` | 功能完成 | Phase 状态、完成总结 |

### 文档贡献指南

1. **保持中英文对照**：所有用户可见的文档建议中英文双语
2. **更新统计数字**：分析器数量、功能数量等保持一致
3. **添加示例代码**：新功能提供使用示例
4. **检查链接**：确保文档内链接有效
5. **代码注释**：Go 代码使用标准注释格式，导出函数必须注释

### 文档格式规范

- 使用 Markdown 格式
- 标题层级清晰（# ## ###）
- 代码块标注语言（```bash, ```go）
- 表格对齐美观
- 使用 emoji 增强可读性（✅ 🚧 ⏳）

## 行为准则

请遵守我们的 [行为准则](CODE_OF_CONDUCT.md)。

## 许可证

通过提交贡献，您同意您的贡献将在 [Apache 2.0](LICENSE) 许可证下发布。

## 联系我们

- GitHub Issues: [kudig/issues](https://github.com/kudig/kudig/issues)
- 邮件: kudig@example.com

感谢您的贡献！
