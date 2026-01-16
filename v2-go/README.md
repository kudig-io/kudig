# kudig v2.0 - Go 版本

> ✅ **Alpha 可用** - Go 语言重构版本，功能完整

## 简介

`kudig` v2.0 是使用 Go 语言重构的下一代 Kubernetes 节点诊断工具，在 v1.0 Bash 版本基础上进行了全面升级。

## 核心特性

- **Go 语言实现**：性能提升，跨平台支持
- **双模式支持**：
  - 离线模式：分析 diagnose_k8s.sh 数据
  - 在线模式：实时诊断 K8s 集群（通过 K8s API）
- **35+ 内置分析器**：涵盖系统、进程、网络、内核、Kubernetes、运行时等维度
- **YAML 规则引擎**：支持自定义诊断规则
- **Kubernetes 原生部署**：提供 Helm Chart 和 DaemonSet 支持
- **Docker 镜像**：开箱即用的容器化部署

## 开发状态

当前版本已完成主要功能实现：

- ✅ 项目结构框架
- ✅ 基础类型定义
- ✅ 35+ 分析器实现
- ✅ 离线数据收集层
- ✅ 在线数据收集层（K8s API）
- ✅ 报告生成层（Text/JSON）
- ✅ YAML 规则引擎
- ✅ Helm Chart
- ✅ Dockerfile

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
├── internal/           # 内部包
├── charts/             # Helm Chart
├── deployments/        # K8s 部署配置
├── configs/            # 配置文件
├── rules/              # 示例规则
├── Dockerfile          # Docker 构建文件
├── Makefile           # 构建脚本
├── go.mod
└── go.sum
```

## 构建

```bash
# 下载依赖
make deps

# 构建
make build

# 运行测试
make test

# 构建所有平台
make build-all

# 构建 Docker 镜像
make docker
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
- Make

### 构建命令

```bash
make build          # 构建单平台二进制
make build-all      # 构建多平台二进制
make test           # 运行测试
make lint           # 代码检查
make docker         # 构建 Docker 镜像
```

## 开发路线图

### v2.0-alpha ✅ 已完成
- [x] 完成离线分析模式
- [x] 实现 35+ 分析器
- [x] 生成文本/JSON 报告
- [x] 兼容 v1.0 数据格式

### v2.0-beta ✅ 已完成
- [x] 添加在线诊断模式
- [x] 实现 YAML 规则引擎
- [x] 添加 Helm Chart
- [x] Dockerfile 构建

### v2.0-stable（进行中）
- [ ] 完善错误处理
- [ ] 性能优化
- [ ] 完整文档
- [ ] 生产环境测试
- [ ] 正式发布

## 贡献

v2.0 版本欢迎社区贡献。如果您想参与开发：

1. Fork 本项目
2. 创建功能分支
3. 提交 Pull Request

## 许可证

Apache License 2.0

## 相关链接

- [v1.0 Bash 版本](../v1-bash/) - ✅ 生产可用
- [项目主页](../)

---

**注意**: 此版本为 Alpha 阶段，主要功能已完成但可能存在边缘情况问题。如需高度稳定的版本，请使用 [v1.0 Bash 版本](../v1-bash/)。
