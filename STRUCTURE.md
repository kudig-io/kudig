# 项目结构重组说明

## 重组时间
2024-12-13（更新于 2026-01-15）

## 重组目的
明确区分 Bash v1.0（生产可用）和 Go v2.0（Alpha 可用）版本，避免用户混淆。

## 新结构

```
kudig/
├── v1-bash/                    # ✅ v1.0 Bash 版本（生产可用）
│   ├── kudig.sh                # 主脚本
│   ├── README.md               # v1.0 完整文档
│   ├── TESTING.md              # 测试说明
│   └── reference/              # 示例诊断数据
│       └── diagnose_k8s/
│
├── v2-go/                      # ✅ v2.0 Go 版本（Alpha 可用）
│   ├── cmd/kudig/              # CLI 入口
│   ├── pkg/                    # 核心包
│   │   ├── analyzer/           # 35+ 分析器
│   │   ├── collector/          # 数据收集（offline/online）
│   │   ├── reporter/           # 报告生成（text/json）
│   │   ├── rules/              # YAML 规则引擎
│   │   ├── types/              # 类型定义
│   │   └── legacy/             # v1 兼容层
│   ├── build/                  # 构建输出（kudig.exe）
│   ├── charts/kudig/           # Helm Chart
│   ├── deployments/            # K8s 部署配置
│   ├── configs/                # 配置文件
│   ├── rules/                  # 规则示例
│   ├── Dockerfile              # Docker 构建
│   ├── Makefile                # 构建脚本
│   ├── go.mod                  # Go 模块
│   ├── go.sum                  # Go 依赖
│   └── README.md               # v2.0 完整文档
│
├── docs/                       # 项目文档
│   ├── CODE_QUALITY_README.md
│   └── QUALITY_CHECK_SETUP.md
│
├── reference/                  # 共享参考数据（根目录）
│   └── diagnose_k8s/
│
├── README.md                   # 项目主文档（导航页）
├── TEST_REPORT.md              # 测试报告
├── TESTING.md                  # 原测试文档（保留）
├── LICENSE                     # 许可证
└── kudig.sh                    # 根目录主脚本（保留兼容性）
```

## 文件移动记录

### v1-bash/
- ✅ 复制 `kudig.sh` → `v1-bash/kudig.sh`
- ✅ 复制 `TESTING.md` → `v1-bash/TESTING.md`
- ✅ 复制 `reference/` → `v1-bash/reference/`
- ✅ 创建 `v1-bash/README.md`

### v2-go/
- ✅ 移动 `cmd/` → `v2-go/cmd/`
- ✅ 移动 `pkg/` → `v2-go/pkg/`
- ✅ 移动 `internal/` → `v2-go/internal/`
- ✅ 移动 `go.mod` → `v2-go/go.mod`
- ✅ 移动 `go.sum` → `v2-go/go.sum`
- ✅ 移动 `Makefile` → `v2-go/Makefile`
- ✅ 移动 `Dockerfile` → `v2-go/Dockerfile`
- ✅ 移动 `charts/` → `v2-go/charts/`
- ✅ 移动 `deployments/` → `v2-go/deployments/`
- ✅ 移动 `configs/` → `v2-go/configs/`
- ✅ 移动 `rules/` → `v2-go/rules/`
- ✅ 创建 `v2-go/README.md`

### 根目录
- ✅ 更新 `README.md` - 改为导航页，明确版本区分
- ✅ 保留 `kudig.sh` - 向后兼容
- ✅ 保留 `TESTING.md` - 原有文档
- ✅ 保留 `reference/` - 共享参考数据

## 版本标识

### v1.0 Bash - ✅ 生产可用
- 目录：`v1-bash/`
- 状态：稳定，可用于生产环境
- 特性：40+ 异常检测规则，离线分析模式

### v2.0 Go - ✅ Alpha 可用
- 目录：`v2-go/`
- 状态：Alpha 阶段，功能完整，可用于测试
- 特性：双模式（离线+在线）、35+ 分析器、YAML规则引擎、Helm Chart、Docker

## 使用指南

### 使用 v1.0 Bash（生产推荐）
```bash
cd v1-bash
./kudig.sh /tmp/diagnose_1702468800
```

### 使用 v2.0 Go（测试可用）
```bash
cd v2-go
make build
./build/kudig offline /tmp/diagnose_1702468800
./build/kudig online --node worker-1
./build/kudig rules --file rules/custom.yaml /tmp/diagnose_1702468800
```

## 文档更新

- ✅ 根目录 `README.md` - 项目导航页
- ✅ `v1-bash/README.md` - v1.0 完整文档
- ✅ `v2-go/README.md` - v2.0 开发文档
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

## 注意事项

- 根目录的 `kudig.sh` 是 v1.0 的副本，保持同步更新
- `scripts/kudig.sh` 可能是备份，需要确认是否保留
- Git 历史保持完整，可以追溯文件移动
