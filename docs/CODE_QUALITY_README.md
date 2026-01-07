# 代码质量检查体系

本文档概述 kudig.sh 项目的代码质量检查体系的实施情况。

## 📋 目录

- [概览](#概览)
- [已完成的工作](#已完成的工作)
- [文件清单](#文件清单)
- [使用指南](#使用指南)
- [CI/CD 集成](#cicd-集成)
- [下一步计划](#下一步计划)

## 概览

根据代码质量检查设计文档，我们已经完成了**第一阶段（基础建设）**和**第二阶段（CI 集成）**的主要任务。

### 实施进度

- ✅ **第一阶段：基础建设** (100%)
  - ShellCheck 和 shfmt 配置
  - Pre-commit hook 脚本
  - 自定义质量检查脚本
  - EditorConfig 配置

- ✅ **第二阶段：CI 集成** (100%)
  - GitHub Actions 工作流配置
  - bats-core 测试框架集成
  - 单元测试用例编写
  - 测试辅助函数库

- ⏳ **第三阶段：完善与优化** (计划中)
  - 补充更多测试用例
  - 多环境兼容性测试
  - 性能优化

- ⏳ **第四阶段：持续改进** (长期)
  - 规则优化与演进
  - 团队培训与分享

## 已完成的工作

### 1. 工具链配置

#### ShellCheck 配置
- 文件：`.shellcheckrc`
- 功能：配置 Shell 脚本静态分析规则
- 特点：
  - 启用所有检查
  - 排除不适用的规则（SC2034, SC1090, SC1091）
  - 设置适当的严重级别

#### shfmt 配置
- 文件：`.shfmtrc` (文档说明)
- 格式规范：
  - 缩进 4 个空格
  - 二元操作符换行
  - Case 语句缩进
  - 重定向符号后加空格

#### EditorConfig
- 文件：`.editorconfig`
- 功能：统一不同编辑器的代码风格
- 覆盖：Shell 脚本、Markdown、YAML、JSON 文件

### 2. Git Hook 集成

#### Pre-commit Hook
- 文件：`scripts/pre-commit`
- 功能：在代码提交前自动执行质量检查
- 检查项：
  1. ShellCheck 语法检查
  2. shfmt 格式检查
  3. 自定义规则检查
- 特点：
  - 仅检查变更的 .sh 文件
  - 友好的错误提示
  - 支持 `--no-verify` 跳过

### 3. 自定义质量检查

#### 质量检查脚本
- 文件：`scripts/quality_check.sh`
- 检查模块：
  1. 文件头部注释完整性
  2. 全局变量命名规范（大写字母+下划线）
  3. 函数命名规范（小写字母+下划线）
  4. 函数注释完整性
  5. 退出码使用规范（0/1/2）
  6. 算术运算安全性（set -e 模式）
  7. 硬编码路径检测
  8. TODO/FIXME 标记统计

#### 使用方式
```bash
# 检查单个文件
./scripts/quality_check.sh kudig.sh

# 检查多个文件
./scripts/quality_check.sh *.sh

# 详细模式
./scripts/quality_check.sh --verbose kudig.sh
```

### 4. CI/CD 集成

#### GitHub Actions 工作流
- 文件：`.github/workflows/code-quality.yml`
- 触发条件：
  - Push 到 main/develop 分支
  - Pull Request 创建或更新
  - 每日定时执行（凌晨 2 点）
  - 手动触发

#### CI 任务矩阵
1. **ShellCheck 静态分析**
   - 运行 ShellCheck 扫描
   - 上传检查报告

2. **代码格式检查**
   - shfmt 格式验证
   - 格式问题提示

3. **自定义规则检查**
   - 执行 quality_check.sh
   - 验证项目特定规范

4. **安全扫描**
   - ShellCheck 安全规则
   - 危险命令模式检测

5. **单元测试**
   - bats-core 测试执行
   - 测试覆盖率统计（规划中）

6. **集成测试**
   - 基本功能验证
   - JSON 输出验证
   - 使用真实数据测试

7. **兼容性测试**
   - Ubuntu 20.04/22.04
   - Debian 11
   - 多环境并行测试

8. **质量报告汇总**
   - 生成质量摘要
   - PR 自动评论
   - 上传报告制品

### 5. 测试框架

#### Bats-core 集成
- 文件：`tests/test_helper.bash`
- 功能：提供测试辅助函数库
- 主要函数：
  - `create_mock_diagnose_dir()` - 创建模拟诊断目录
  - `create_system_status_with_load()` - 生成负载数据
  - `create_memory_info()` - 生成内存数据
  - `create_kubelet_status()` - 生成服务状态
  - `validate_json_output()` - 验证 JSON 格式
  - `assert_anomaly_exists()` - 断言异常存在

#### 单元测试用例

**基础功能测试** (`tests/unit/test_basic_functions.bats`):
- 脚本存在性和可执行性
- 帮助信息显示
- 版本信息显示
- 参数验证
- JSON 输出格式
- Verbose 模式
- 真实数据处理

**系统资源检查测试** (`tests/unit/test_system_check.bats`):
- CPU 负载检测（正常/偏高/严重）
- 内存使用率检测（正常/偏高/严重）
- 磁盘空间检测（正常/不足/严重不足）

## 文件清单

### 配置文件
```
.shellcheckrc                           # ShellCheck 配置
.shfmtrc                                # shfmt 配置说明
.editorconfig                           # 编辑器配置
```

### 脚本文件
```
scripts/
├── pre-commit                          # Git pre-commit hook
└── quality_check.sh                    # 自定义质量检查脚本
```

### CI/CD 配置
```
.github/workflows/
└── code-quality.yml                    # GitHub Actions 工作流
```

### 测试文件
```
tests/
├── test_helper.bash                    # 测试辅助函数
└── unit/
    ├── test_basic_functions.bats       # 基础功能测试
    └── test_system_check.bats          # 系统检查测试
```

### 文档
```
docs/
├── QUALITY_CHECK_SETUP.md              # 安装配置指南
└── CODE_QUALITY_README.md              # 本文档
```

## 使用指南

### 本地开发

#### 1. 安装工具

参考 [安装配置指南](./QUALITY_CHECK_SETUP.md) 安装必要工具：
- ShellCheck
- shfmt
- bats-core（可选）

#### 2. 配置 Git Hook

```bash
# 复制 pre-commit hook
cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

#### 3. 运行检查

```bash
# 手动运行 ShellCheck
shellcheck kudig.sh

# 手动运行格式检查
shfmt -d -i 4 -bn -ci -sr kudig.sh

# 运行自定义检查
bash scripts/quality_check.sh kudig.sh

# 运行单元测试
bats tests/unit/
```

### 编辑器集成

#### VSCode

安装推荐扩展：
- ShellCheck (timonwong.shellcheck)
- shell-format (foxundermoon.shell-format)

配置会自动读取项目的 `.shellcheckrc` 和 `.editorconfig`。

#### Vim/Neovim

使用 ALE 插件进行实时检查和自动格式化。

## CI/CD 集成

### 自动触发

CI 工作流会在以下情况自动运行：
- 推送到 main 或 develop 分支
- 创建或更新 Pull Request
- 每天凌晨 2 点定时执行

### 质量门禁

以下情况 CI 会失败：
- ShellCheck Error 级别问题
- 格式检查失败
- 自定义规则检查失败
- 安全扫描发现问题
- 单元测试失败
- 集成测试失败

### PR 集成

CI 会在 PR 中自动评论质量报告，包括：
- 各项检查的通过状态
- 整体通过率
- 详细问题列表（如有）

## 下一步计划

### 短期（1-2周）

1. **补充测试用例**
   - 网络检查测试
   - 内核检查测试
   - Kubernetes 组件检查测试
   - 完整的集成测试

2. **优化检查规则**
   - 根据实际运行情况调整规则
   - 减少误报
   - 提高检查效率

3. **完善文档**
   - 添加更多使用示例
   - 编写故障排除指南
   - 记录常见问题

### 中期（3-4周）

1. **多环境测试增强**
   - 添加 CentOS 7/8 测试
   - 添加 Aliyun Linux 测试
   - 添加 Kylin 测试

2. **代码覆盖率**
   - 集成 kcov
   - 上传到 Codecov
   - 设置覆盖率目标（80%）

3. **性能优化**
   - 优化 CI 执行时间
   - 实现增量检查
   - 添加缓存机制

### 长期（持续）

1. **规则演进**
   - 季度规则评审
   - 引入新的最佳实践
   - 社区规则同步

2. **团队培训**
   - 编写培训文档
   - 开展内部分享
   - Code Review 反馈

3. **质量度量**
   - 建立质量指标看板
   - 趋势分析
   - 持续改进

## 参考资源

- [设计文档](../.qoder/quests/code-quality-check.md) - 完整的设计文档
- [安装指南](./QUALITY_CHECK_SETUP.md) - 工具安装和配置
- [ShellCheck 文档](https://github.com/koalaman/shellcheck/wiki)
- [shfmt 文档](https://github.com/mvdan/sh)
- [bats-core 文档](https://bats-core.readthedocs.io/)

## 问题反馈

如有问题或建议，请：
1. 查看 [安装指南](./QUALITY_CHECK_SETUP.md) 的常见问题部分
2. 提交 GitHub Issue
3. 联系项目维护者

---

**最后更新**: 2026-01-07  
**文档版本**: 1.0  
**实施状态**: 第一、二阶段完成
