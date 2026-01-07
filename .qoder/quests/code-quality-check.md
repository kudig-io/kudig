# 代码质量检查设计文档

## 1. 设计目标

为 kudig.sh 项目建立完整的代码质量检查体系，确保 Bash 脚本代码的可靠性、可维护性和规范性，在代码提交、合并和发布前自动发现潜在问题。

## 2. 业务价值

### 2.1 核心价值
- **提升代码质量**：通过自动化检查发现语法错误、安全隐患和不规范代码
- **降低维护成本**：统一编码规范，减少代码审查负担，提高可读性
- **增强可靠性**：在开发阶段发现潜在缺陷，避免生产环境故障
- **保障兼容性**：确保脚本在不同 Linux 发行版和 Bash 版本中正常运行

### 2.2 适用场景
- 开发人员本地提交前检查
- CI/CD 流程中的自动化门禁
- 代码合并前的质量保证
- 版本发布前的最终验证

## 3. 质量检查维度

### 3.1 语法正确性检查

#### 检查目标
确保 Bash 脚本语法完全符合规范，能够在目标环境中正确解析执行。

#### 检查项
- Bash 语法错误检测
- 条件表达式正确性
- 变量引用完整性
- 函数定义规范性
- 循环结构正确性
- 重定向语法合法性

#### 工具选择
- **ShellCheck**：业界标准的 Shell 脚本静态分析工具
  - 检测语法错误
  - 识别常见陷阱
  - 提供最佳实践建议
  - 支持多种警告级别

#### 质量标准
- 零语法错误（Error级别）
- 警告（Warning级别）不超过 5 个
- 信息提示（Info级别）酌情处理

### 3.2 代码规范检查

#### 检查目标
统一代码风格，提升可读性和团队协作效率。

#### 检查项

##### 命名规范
- 变量命名约定
  - 全局常量：大写字母加下划线（如 `SEVERITY_CRITICAL`）
  - 局部变量：小写字母加下划线（如 `diagnose_dir`）
  - 函数名：小写字母加下划线（如 `check_system_resources`）
- 避免使用保留字
- 避免单字母变量（循环计数器除外）

##### 格式规范
- 缩进：使用 4 个空格（不使用 Tab）
- 行宽：单行不超过 120 字符
- 函数定义：采用统一格式
  ```
  function_name() {
      local var_name="value"
      # 函数体
  }
  ```
- 条件判断：统一使用 `[[ ]]` 双括号
- 注释规范
  - 文件头部包含功能描述、使用方法
  - 每个函数前添加功能说明
  - 复杂逻辑添加行内注释

##### 结构规范
- 全局变量集中声明
- 函数按功能分类组织
- 工具函数置于前部
- 主函数置于文件末尾

#### 工具选择
- **shfmt**：Shell 脚本格式化工具
  - 自动格式化代码
  - 统一缩进风格
  - 检测格式不一致
- **自定义规则脚本**：针对项目特定规范的检查

#### 质量标准
- 所有代码符合格式规范
- 命名规范覆盖率 100%
- 关键函数注释完整率 100%

### 3.3 安全性检查

#### 检查目标
识别脚本中的安全漏洞和潜在风险，防止注入攻击和权限滥用。

#### 检查项

##### 命令注入防护
- 用户输入参数必须经过验证
- 禁止直接拼接外部输入到命令中
- 使用数组而非字符串传递参数

##### 路径遍历防护
- 文件路径参数必须验证合法性
- 禁止使用未清理的相对路径
- 限制文件访问范围

##### 权限控制
- 避免不必要的 root 权限操作
- 临时文件使用安全的权限设置
- 敏感信息不输出到日志

##### 资源安全
- 防止无限循环
- 防止递归调用栈溢出
- 临时文件清理机制

#### 工具选择
- **ShellCheck**：内置安全检查规则（SC2086, SC2155 等）
- **Bandit-baseline（Bash 适配版）**：安全扫描工具
- **自定义安全规则**：项目特定的安全检查

#### 质量标准
- 零高危安全问题
- 中危问题必须有风险评估记录
- 低危问题有改进计划

### 3.4 性能与效率检查

#### 检查目标
识别性能瓶颈和低效代码模式，优化脚本执行效率。

#### 检查项

##### 反模式检测
- 循环中的重复命令调用
- 不必要的子进程创建
- 管道滥用导致的性能损失
- 频繁的文件 I/O 操作

##### 资源使用
- 大文件处理方式
- 内存占用估算
- 并发执行可行性

##### 算法复杂度
- 避免 O(n²) 或更高复杂度的嵌套循环
- 优先使用内置命令而非外部工具

#### 工具选择
- **ShellCheck**：识别常见性能反模式（SC2013, SC2002 等）
- **自定义性能检查脚本**：检测项目特定的性能问题

#### 质量标准
- 无明显性能反模式
- 关键路径执行时间可控
- 资源占用在可接受范围

### 3.5 可移植性检查

#### 检查目标
确保脚本在不同 Linux 发行版和 Bash 版本中都能正常运行。

#### 检查项

##### Bash 版本兼容性
- 避免使用 Bash 4.0 以下版本不支持的特性
- 关联数组仅在确认支持的环境使用
- 检查脚本 Shebang 声明

##### 命令可用性
- 依赖的外部命令明确声明
- 启动时检查必需命令是否存在
- 提供命令缺失的友好提示

##### 路径与环境
- 避免硬编码绝对路径
- 环境变量使用的兼容性
- 临时目录创建的可移植性

#### 工具选择
- **ShellCheck**：POSIX 兼容性检查（使用 `-s bash` 参数）
- **checkbashisms**：检测非标准 Bash 特性
- **容器化测试**：在不同发行版镜像中验证

#### 质量标准
- 支持 Red Hat 7+、CentOS 7+、Ubuntu 18.04+ 等主流发行版
- 支持 Bash 4.0 及以上版本
- 所有依赖命令明确列出

### 3.6 测试覆盖率检查

#### 检查目标
确保关键功能有充分的测试覆盖，降低回归风险。

#### 检查项

##### 功能覆盖
- 每个异常检测函数有对应测试用例
- 边界条件测试
- 错误处理路径测试

##### 场景覆盖
- 正常诊断数据场景
- 不完整诊断数据场景
- 异常输入参数场景
- 空数据或无效数据场景

##### 输出验证
- 文本格式输出正确性
- JSON 格式输出合法性
- 退出码准确性

#### 工具选择
- **bats-core**：Bash 自动化测试框架
- **kcov**：代码覆盖率统计工具
- **自定义测试脚本**：集成测试框架

#### 质量标准
- 核心功能覆盖率 ≥ 80%
- 关键检测函数覆盖率 100%
- 所有公开函数有测试

### 3.7 文档与注释检查

#### 检查目标
确保代码有足够的文档支持，便于理解和维护。

#### 检查项

##### 文件级文档
- 脚本用途说明
- 使用方法示例
- 依赖项列表
- 版本信息

##### 函数级文档
- 函数功能描述
- 参数说明
- 返回值说明
- 使用示例（复杂函数）

##### 关键逻辑注释
- 复杂算法的实现思路
- 非显而易见的业务规则
- Workaround 的原因说明

#### 工具选择
- **自定义文档检查脚本**：验证注释完整性
- **markdownlint**：检查 Markdown 文档规范

#### 质量标准
- 所有公开函数有文档注释
- README 文档完整更新
- 复杂逻辑有清晰注释

## 4. 检查流程设计

### 4.1 开发阶段检查

#### 执行时机
开发人员在本地提交代码前

#### 检查方式
通过 Git Hook（pre-commit）自动触发

#### 检查流程
```
开发人员执行 git commit
         ↓
触发 pre-commit hook
         ↓
对变更的 .sh 文件执行快速检查
├─ ShellCheck 语法检查
├─ shfmt 格式检查
└─ 自定义规则检查
         ↓
检查结果判断
├─ 通过：允许提交
└─ 失败：阻止提交，显示错误信息
```

#### 检查内容
- 语法正确性（阻断）
- 格式规范（阻断）
- 基础安全检查（警告）

#### 执行时间要求
不超过 10 秒

### 4.2 持续集成检查

#### 执行时机
- 代码推送到远程仓库时
- Pull Request 创建或更新时
- 定期调度（每日）

#### 检查方式
通过 CI/CD 流水线（GitHub Actions / GitLab CI）自动执行

#### 检查流程
```
代码推送到远程仓库
         ↓
触发 CI 流水线
         ↓
环境准备
├─ 安装 ShellCheck
├─ 安装 shfmt
├─ 安装 bats-core
└─ 安装其他依赖工具
         ↓
并行执行检查任务
├─ 静态代码检查
│   ├─ ShellCheck 全面扫描
│   ├─ shfmt 格式验证
│   └─ 自定义规则检查
├─ 安全扫描
│   ├─ ShellCheck 安全规则
│   └─ 自定义安全检查
├─ 性能分析
│   └─ 性能反模式检测
└─ 自动化测试
    ├─ 单元测试（bats）
    ├─ 集成测试
    └─ 覆盖率统计
         ↓
生成质量报告
├─ 检查结果汇总
├─ 问题分类统计
├─ 趋势对比分析
└─ 上传制品
         ↓
质量门禁判断
├─ 通过：允许合并
└─ 失败：阻止合并，通知相关人员
```

#### 检查内容（全面）
- 所有维度的质量检查
- 完整的自动化测试
- 代码覆盖率统计
- 质量趋势分析

#### 执行时间要求
不超过 5 分钟

### 4.3 发布前检查

#### 执行时机
版本发布前的最终验证

#### 检查方式
手动触发或自动触发（版本标签创建时）

#### 检查流程
```
创建版本标签或手动触发
         ↓
执行完整质量检查
├─ CI 阶段所有检查
├─ 额外的兼容性测试
│   ├─ CentOS 7 环境测试
│   ├─ Ubuntu 18.04 环境测试
│   ├─ Aliyun Linux 环境测试
│   └─ Kylin 环境测试
└─ 文档完整性检查
         ↓
生成发布质量报告
         ↓
人工评审
         ↓
批准发布
```

#### 检查内容（最严格）
- CI 阶段所有检查
- 多发行版兼容性测试
- 文档同步性验证
- 版本号一致性检查

#### 执行时间要求
不超过 15 分钟

## 5. 质量门禁规则

### 5.1 阻断性问题（必须修复）

#### 语法错误
- ShellCheck Error 级别问题
- Bash 语法解析失败

#### 严重安全漏洞
- 命令注入风险
- 路径遍历漏洞
- 权限提升风险

#### 功能破坏性问题
- 核心功能测试失败
- 回归测试失败

### 5.2 警告性问题（建议修复）

#### 格式不规范
- 缩进不一致
- 命名不符合约定
- 代码行过长

#### 轻微安全隐患
- 变量未引用
- 临时文件权限过宽

#### 性能问题
- 明显的性能反模式
- 不必要的资源消耗

### 5.3 提示性问题（可选修复）

#### 代码改进建议
- 可读性优化建议
- 更简洁的实现方式

#### 文档完善
- 注释可以更详细
- 示例可以更丰富

## 6. 工具链配置

### 6.1 ShellCheck 配置

#### 配置文件：`.shellcheckrc`
```
# 启用所有检查
enable=all

# 排除特定规则（根据项目需要）
disable=SC2034  # 未使用的变量（某些全局配置变量）
disable=SC1090  # 无法跟踪动态source

# Shell 类型
shell=bash

# 严重级别阈值
severity=style
```

#### 集成方式
- 命令行工具：直接调用 `shellcheck kudig.sh`
- 编辑器集成：VSCode ShellCheck 插件
- CI 集成：在流水线中自动执行

### 6.2 shfmt 配置

#### 格式化规则
```
-i 4      # 缩进4个空格
-bn       # 二元操作符换行
-ci       # case 语句缩进
-sr       # 重定向符后加空格
-kp       # 保持列对齐
```

#### 集成方式
- 本地格式化：`shfmt -w -i 4 -bn -ci -sr kudig.sh`
- 格式检查：`shfmt -d -i 4 -bn -ci -sr kudig.sh`
- CI 集成：在流水线中验证格式

### 6.3 bats-core 测试框架

#### 测试文件组织
```
tests/
├── test_helper.bash           # 测试辅助函数
├── unit/                      # 单元测试
│   ├── test_system_check.bats
│   ├── test_network_check.bats
│   └── test_kernel_check.bats
├── integration/               # 集成测试
│   ├── test_full_analysis.bats
│   └── test_output_format.bats
└── fixtures/                  # 测试数据
    ├── normal_diagnose/
    └── abnormal_diagnose/
```

#### 测试执行方式
```bash
# 执行所有测试
bats tests/

# 执行单个测试文件
bats tests/unit/test_system_check.bats

# 生成 TAP 格式输出
bats --formatter tap tests/

# 生成覆盖率报告（结合 kcov）
kcov coverage bats tests/
```

### 6.4 自定义检查脚本

#### 脚本：`scripts/quality_check.sh`

##### 功能模块
- 命名规范检查
- 项目特定规则验证
- 文档完整性检查
- 退出码使用规范

##### 执行方式
```bash
# 检查单个文件
./scripts/quality_check.sh kudig.sh

# 检查所有 Shell 脚本
./scripts/quality_check.sh *.sh

# 输出详细报告
./scripts/quality_check.sh --verbose kudig.sh
```

## 7. CI/CD 集成方案

### 7.1 GitHub Actions 配置

#### 工作流文件：`.github/workflows/code-quality.yml`

##### 触发条件
- 推送到 main/develop 分支
- Pull Request 创建或更新
- 定时任务（每日）

##### 任务矩阵
```
Job 1: 静态代码检查
├─ 安装 ShellCheck
├─ 安装 shfmt
├─ 执行 ShellCheck 扫描
├─ 执行 shfmt 格式检查
└─ 上传检查报告

Job 2: 安全扫描
├─ ShellCheck 安全规则
├─ 自定义安全检查
└─ 生成安全报告

Job 3: 自动化测试
├─ 安装 bats-core
├─ 安装 kcov
├─ 执行单元测试
├─ 执行集成测试
├─ 生成覆盖率报告
└─ 上传覆盖率到 Codecov

Job 4: 多环境兼容性测试
├─ CentOS 7 容器测试
├─ Ubuntu 18.04 容器测试
└─ Aliyun Linux 容器测试

Job 5: 质量报告汇总
├─ 收集所有检查结果
├─ 生成质量趋势图
└─ 发布评论到 PR
```

##### 质量门禁规则
- Job 1-3 必须全部通过
- Job 4 至少通过 2 个环境
- 代码覆盖率不低于 80%

### 7.2 GitLab CI 配置（可选）

#### 配置文件：`.gitlab-ci.yml`

##### 阶段划分
- lint：静态检查
- security：安全扫描
- test：自动化测试
- report：报告生成

##### 并行执行
各阶段内的任务并行执行，加快反馈速度

## 8. 质量报告设计

### 8.1 报告内容

#### 综合评分
基于各维度检查结果计算综合质量分（0-100 分）

#### 问题分类统计
- 按严重级别分类：Error / Warning / Info
- 按类型分类：语法 / 格式 / 安全 / 性能 / 文档
- 按文件分布：各文件问题数量

#### 趋势分析
- 与上次检查对比
- 历史趋势图表
- 新增/修复问题统计

#### 详细问题清单
| 级别 | 类型 | 位置 | 问题描述 | 修复建议 |
|------|------|------|----------|----------|
| Error | 语法 | kudig.sh:123 | 变量未定义 | 添加变量声明 |

### 8.2 报告格式

#### HTML 格式（CI 生成）
- 可视化图表
- 可交互的问题列表
- 代码高亮显示

#### JSON 格式（API 集成）
```json
{
  "version": "1.0",
  "timestamp": "2024-01-07T13:00:00Z",
  "score": 95,
  "summary": {
    "total": 12,
    "error": 0,
    "warning": 8,
    "info": 4
  },
  "categories": {
    "syntax": { "error": 0, "warning": 2 },
    "format": { "error": 0, "warning": 5 },
    "security": { "error": 0, "warning": 1 }
  },
  "issues": [...]
}
```

#### Markdown 格式（PR 评论）
- 简洁的问题摘要
- 直接链接到代码位置
- 修复建议

### 8.3 报告分发

#### 自动通知
- CI 失败时通过邮件/IM 通知责任人
- PR 中自动评论质量报告
- 定期发送质量趋势报告

#### 报告存储
- CI 制品存储：保留最近 30 次报告
- 趋势数据库：长期存储质量指标
- 代码覆盖率平台：Codecov/Coveralls

## 9. 持续改进机制

### 9.1 规则演进

#### 规则评审
- 每季度评审一次检查规则
- 根据实际问题调整规则
- 引入新的最佳实践

#### 误报处理
- 收集误报反馈
- 调整规则配置减少误报
- 添加例外规则（需要审批）

### 9.2 度量与分析

#### 关键指标
- 代码质量分数趋势
- 问题修复周期
- 测试覆盖率变化
- 技术债务数量

#### 分析报告
- 月度质量分析报告
- 问题根因分析
- 改进建议

### 9.3 团队培训

#### 培训内容
- Bash 编码规范
- 常见问题与最佳实践
- 工具使用指南

#### 培训方式
- 文档知识库
- 内部分享会
- Code Review 反馈

## 10. 实施路线图

### 第一阶段：基础建设（1-2周）

#### 目标
建立基本的代码质量检查能力

#### 任务
- 配置 ShellCheck 和 shfmt
- 创建 pre-commit hook
- 编写自定义检查脚本
- 修复现有代码的明显问题

#### 交付物
- `.shellcheckrc` 配置文件
- `.git/hooks/pre-commit` 脚本
- `scripts/quality_check.sh` 检查脚本
- 初步的质量基线报告

### 第二阶段：CI 集成（2-3周）

#### 目标
在 CI/CD 流程中自动化质量检查

#### 任务
- 配置 GitHub Actions 工作流
- 集成 bats-core 测试框架
- 编写核心功能的单元测试
- 配置代码覆盖率上传

#### 交付物
- `.github/workflows/code-quality.yml` 工作流
- `tests/` 目录及测试用例
- Codecov 集成配置
- CI 质量门禁规则文档

### 第三阶段：完善与优化（3-4周）

#### 目标
完善测试覆盖，优化检查效率

#### 任务
- 补充集成测试用例
- 实现多环境兼容性测试
- 优化 CI 执行时间
- 完善质量报告

#### 交付物
- 完整的测试套件
- 多发行版测试镜像
- 优化后的 CI 流程
- HTML 格式质量报告

### 第四阶段：持续改进（长期）

#### 目标
建立质量持续改进机制

#### 任务
- 定期评审和更新规则
- 分析质量趋势
- 团队培训与分享
- 引入新的检查工具

#### 交付物
- 季度质量分析报告
- 规则更新日志
- 培训文档与视频
- 持续优化的检查体系

## 11. 风险与应对

### 11.1 工具兼容性风险

#### 风险描述
ShellCheck 等工具在不同环境中版本不一致，导致检查结果差异

#### 应对措施
- 在 CI 中使用固定版本的工具
- 提供 Docker 镜像统一检查环境
- 文档中明确工具版本要求

### 11.2 检查时间过长

#### 风险描述
完整的质量检查耗时过长，影响开发效率

#### 应对措施
- 本地 pre-commit 仅执行快速检查
- CI 中并行执行检查任务
- 增量检查：仅检查变更的代码
- 缓存机制：缓存依赖和测试数据

### 11.3 误报率过高

#### 风险描述
大量误报导致开发人员忽视检查结果

#### 应对措施
- 精细调整检查规则
- 允许在代码中添加例外标记（需要注释说明原因）
- 定期收集误报反馈并优化
- 区分阻断性问题和建议性问题

### 11.4 团队抵触情绪

#### 风险描述
严格的质量检查增加工作量，团队成员抵触

#### 应对措施
- 充分沟通质量检查的价值
- 提供详细的修复指南
- 逐步推进，不一次性全部启用
- 展示质量改进的成果

## 12. 成功标准

### 12.1 定量指标

- 代码质量综合评分 ≥ 90 分
- ShellCheck Error 级别问题数量 = 0
- 代码测试覆盖率 ≥ 80%
- CI 检查通过率 ≥ 95%
- 问题平均修复时间 ≤ 2 天

### 12.2 定性目标

- 团队成员认可质量检查的价值
- 代码 Review 效率显著提升
- 生产环境缺陷率明显下降
- 新人上手难度降低
- 代码库整体可维护性增强

## 13. 参考资源

### 13.1 工具文档

- ShellCheck Wiki：https://github.com/koalaman/shellcheck/wiki
- shfmt 文档：https://github.com/mvdan/sh
- bats-core 文档：https://bats-core.readthedocs.io/

### 13.2 编码规范

- Google Shell Style Guide：https://google.github.io/styleguide/shellguide.html
- Bash Best Practices：https://bertvv.github.io/cheat-sheets/Bash.html

### 13.3 社区资源

- ShellCheck 在线检查：https://www.shellcheck.net/
- Awesome Shell：https://github.com/alebcay/awesome-shell

---

## 附录 A：检查规则清单

### ShellCheck 关键规则

| 规则代码 | 类型 | 说明 | 级别 |
|---------|------|------|------|
| SC2086 | 安全 | 变量应该用引号包围 | Error |
| SC2046 | 安全 | 命令替换应该用引号包围 | Error |
| SC2155 | 安全 | 局部变量声明与赋值应分开 | Warning |
| SC2181 | 最佳实践 | 直接检查退出码而非 $? | Warning |
| SC2034 | 代码质量 | 未使用的变量 | Warning |
| SC2002 | 性能 | 不必要的 cat 使用 | Info |

### 自定义规则

| 规则名称 | 检查内容 | 级别 |
|---------|---------|------|
| NAMING_01 | 全局常量使用大写命名 | Error |
| NAMING_02 | 函数名使用小写加下划线 | Warning |
| COMMENT_01 | 公开函数必须有注释 | Warning |
| FORMAT_01 | 缩进必须为 4 个空格 | Error |
| FORMAT_02 | 单行长度不超过 120 字符 | Warning |

## 附录 B：质量检查工作流配置示例

### GitHub Actions 完整配置

```yaml
name: Code Quality Check

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]
  schedule:
    - cron: '0 2 * * *'  # 每日凌晨2点

jobs:
  shellcheck:
    name: ShellCheck静态分析
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 运行ShellCheck
        uses: ludeeus/action-shellcheck@master
        with:
          scandir: '.'
          severity: warning
          
      - name: 上传ShellCheck报告
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: shellcheck-report
          path: shellcheck.log

  format-check:
    name: 代码格式检查
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 安装shfmt
        run: |
          wget -O shfmt https://github.com/mvdan/sh/releases/download/v3.7.0/shfmt_v3.7.0_linux_amd64
          chmod +x shfmt
          sudo mv shfmt /usr/local/bin/
      
      - name: 检查代码格式
        run: |
          shfmt -d -i 4 -bn -ci -sr *.sh
  
  security-scan:
    name: 安全扫描
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: ShellCheck安全规则
        run: |
          shellcheck -S error *.sh
      
      - name: 自定义安全检查
        run: |
          ./scripts/security_check.sh
  
  unit-test:
    name: 单元测试
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 安装bats-core
        run: |
          sudo apt-get update
          sudo apt-get install -y bats
      
      - name: 安装kcov（覆盖率）
        run: |
          sudo apt-get install -y kcov
      
      - name: 运行单元测试
        run: |
          kcov coverage bats tests/unit/
      
      - name: 上传覆盖率报告
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage/cobertura.xml
          
  integration-test:
    name: 集成测试
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 准备测试数据
        run: |
          ./tests/fixtures/prepare_test_data.sh
      
      - name: 运行集成测试
        run: |
          bats tests/integration/
  
  compatibility-test:
    name: 兼容性测试 - ${{ matrix.os }}
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: 
          - centos:7
          - ubuntu:18.04
          - registry.cn-hangzhou.aliyuncs.com/alinux/aliyunlinux:3
    container:
      image: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      
      - name: 安装依赖
        run: |
          if [ -f /etc/redhat-release ]; then
            yum install -y bash grep gawk sed coreutils findutils
          else
            apt-get update && apt-get install -y bash grep gawk sed coreutils findutils
          fi
      
      - name: 运行基本测试
        run: |
          bash kudig.sh --help
          bash kudig.sh --version
      
      - name: 功能测试
        run: |
          ./tests/compatibility_test.sh
  
  quality-report:
    name: 生成质量报告
    needs: [shellcheck, format-check, security-scan, unit-test, integration-test]
    runs-on: ubuntu-latest
    if: always()
    steps:
      - uses: actions/checkout@v3
      
      - name: 下载所有报告
        uses: actions/download-artifact@v3
      
      - name: 生成综合报告
        run: |
          ./scripts/generate_quality_report.sh
      
      - name: 上传质量报告
        uses: actions/upload-artifact@v3
        with:
          name: quality-report
          path: quality-report.html
      
      - name: 评论PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const report = fs.readFileSync('quality-report.md', 'utf8');
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: report
            });
```

## 附录 C：质量检查脚本示例

### 自定义质量检查脚本（scripts/quality_check.sh）

```
（由于不能包含代码实现，此处用自然语言描述脚本功能）

该脚本主要包含以下检查功能模块：

1. 命名规范检查模块
   - 扫描脚本中的全局变量定义，验证是否使用大写命名
   - 提取所有函数名，验证是否使用小写加下划线
   - 检查是否使用了 Bash 保留字作为变量名

2. 注释完整性检查模块
   - 验证文件头部是否包含功能描述和使用说明
   - 检查每个函数定义前是否有注释说明
   - 统计注释覆盖率

3. 项目特定规则检查模块
   - 验证异常添加格式是否符合规范（严重级别|中文|英文|详情|位置）
   - 检查日志输出函数是否正确使用
   - 验证退出码使用是否符合约定（0/1/2）

4. 文档同步性检查模块
   - 检查 README 中的版本号是否与脚本一致
   - 验证使用示例是否与实际参数匹配
   - 检查异常检测规则表是否完整

5. 报告生成模块
   - 汇总所有检查结果
   - 按严重级别分类输出问题
   - 生成 Markdown 格式或 JSON 格式报告
```

## 附录 D：测试用例示例

### 单元测试用例（tests/unit/test_system_check.bats）

```
（用自然语言描述测试用例结构和验证逻辑）

该测试文件验证系统资源检查函数的正确性：

测试套件：check_system_resources 功能测试

测试用例 1：CPU负载检测 - 正常情况
- 准备：创建包含正常负载数据的模拟诊断目录
- 执行：调用 check_system_resources 函数
- 验证：确认未添加异常记录，日志输出显示"正常"

测试用例 2：CPU负载检测 - 偏高情况
- 准备：创建包含偏高负载数据（2-4倍CPU核心数）的诊断目录
- 执行：调用 check_system_resources 函数
- 验证：确认添加了警告级别的 ELEVATED_SYSTEM_LOAD 异常

测试用例 3：CPU负载检测 - 严重情况
- 准备：创建包含严重负载数据（>4倍CPU核心数）的诊断目录
- 执行：调用 check_system_resources 函数
- 验证：确认添加了严重级别的 HIGH_SYSTEM_LOAD 异常

测试用例 4：内存使用率检测 - 边界值测试
- 准备：创建内存使用率为 85%, 95%, 100% 的三组数据
- 执行：分别测试三种情况
- 验证：
  - 85% 触发 ELEVATED_MEMORY_USAGE 警告
  - 95% 触发 HIGH_MEMORY_USAGE 严重
  - 100% 触发 HIGH_MEMORY_USAGE 严重

测试用例 5：磁盘空间检测 - 多挂载点
- 准备：创建包含多个挂载点数据，部分超过阈值
- 执行：调用检测函数
- 验证：仅对超过阈值的挂载点添加异常记录

测试用例 6：缺失数据场景
- 准备：创建不包含 memory_info 文件的诊断目录
- 执行：调用检测函数
- 验证：函数正常退出，不崩溃，输出"跳过"状态
```
- ShellCheck Error 级别问题数量 = 0
- 代码测试覆盖率 ≥ 80%
- CI 检查通过率 ≥ 95%
- 问题平均修复时间 ≤ 2 天

### 12.2 定性目标

- 团队成员认可质量检查的价值
- 代码 Review 效率显著提升
- 生产环境缺陷率明显下降
- 新人上手难度降低
- 代码库整体可维护性增强

## 13. 参考资源

### 13.1 工具文档

- ShellCheck Wiki：https://github.com/koalaman/shellcheck/wiki
- shfmt 文档：https://github.com/mvdan/sh
- bats-core 文档：https://bats-core.readthedocs.io/

### 13.2 编码规范

- Google Shell Style Guide：https://google.github.io/styleguide/shellguide.html
- Bash Best Practices：https://bertvv.github.io/cheat-sheets/Bash.html

### 13.3 社区资源

- ShellCheck 在线检查：https://www.shellcheck.net/
- Awesome Shell：https://github.com/alebcay/awesome-shell

---

## 附录 A：检查规则清单

### ShellCheck 关键规则

| 规则代码 | 类型 | 说明 | 级别 |
|---------|------|------|------|
| SC2086 | 安全 | 变量应该用引号包围 | Error |
| SC2046 | 安全 | 命令替换应该用引号包围 | Error |
| SC2155 | 安全 | 局部变量声明与赋值应分开 | Warning |
| SC2181 | 最佳实践 | 直接检查退出码而非 $? | Warning |
| SC2034 | 代码质量 | 未使用的变量 | Warning |
| SC2002 | 性能 | 不必要的 cat 使用 | Info |

### 自定义规则

| 规则名称 | 检查内容 | 级别 |
|---------|---------|------|
| NAMING_01 | 全局常量使用大写命名 | Error |
| NAMING_02 | 函数名使用小写加下划线 | Warning |
| COMMENT_01 | 公开函数必须有注释 | Warning |
| FORMAT_01 | 缩进必须为 4 个空格 | Error |
| FORMAT_02 | 单行长度不超过 120 字符 | Warning |

## 附录 B：质量检查工作流配置示例

### GitHub Actions 完整配置

```yaml
name: Code Quality Check

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]
  schedule:
    - cron: '0 2 * * *'  # 每日凌晨2点

jobs:
  shellcheck:
    name: ShellCheck静态分析
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 运行ShellCheck
        uses: ludeeus/action-shellcheck@master
        with:
          scandir: '.'
          severity: warning
          
      - name: 上传ShellCheck报告
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: shellcheck-report
          path: shellcheck.log

  format-check:
    name: 代码格式检查
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 安装shfmt
        run: |
          wget -O shfmt https://github.com/mvdan/sh/releases/download/v3.7.0/shfmt_v3.7.0_linux_amd64
          chmod +x shfmt
          sudo mv shfmt /usr/local/bin/
      
      - name: 检查代码格式
        run: |
          shfmt -d -i 4 -bn -ci -sr *.sh
  
  security-scan:
    name: 安全扫描
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: ShellCheck安全规则
        run: |
          shellcheck -S error *.sh
      
      - name: 自定义安全检查
        run: |
          ./scripts/security_check.sh
  
  unit-test:
    name: 单元测试
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 安装bats-core
        run: |
          sudo apt-get update
          sudo apt-get install -y bats
      
      - name: 安装kcov（覆盖率）
        run: |
          sudo apt-get install -y kcov
      
      - name: 运行单元测试
        run: |
          kcov coverage bats tests/unit/
      
      - name: 上传覆盖率报告
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage/cobertura.xml
          
  integration-test:
    name: 集成测试
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 准备测试数据
        run: |
          ./tests/fixtures/prepare_test_data.sh
      
      - name: 运行集成测试
        run: |
          bats tests/integration/
  
  compatibility-test:
    name: 兼容性测试 - ${{ matrix.os }}
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: 
          - centos:7
          - ubuntu:18.04
          - registry.cn-hangzhou.aliyuncs.com/alinux/aliyunlinux:3
    container:
      image: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      
      - name: 安装依赖
        run: |
          if [ -f /etc/redhat-release ]; then
            yum install -y bash grep gawk sed coreutils findutils
          else
            apt-get update && apt-get install -y bash grep gawk sed coreutils findutils
          fi
      
      - name: 运行基本测试
        run: |
          bash kudig.sh --help
          bash kudig.sh --version
      
      - name: 功能测试
        run: |
          ./tests/compatibility_test.sh
  
  quality-report:
    name: 生成质量报告
    needs: [shellcheck, format-check, security-scan, unit-test, integration-test]
    runs-on: ubuntu-latest
    if: always()
    steps:
      - uses: actions/checkout@v3
      
      - name: 下载所有报告
        uses: actions/download-artifact@v3
      
      - name: 生成综合报告
        run: |
          ./scripts/generate_quality_report.sh
      
      - name: 上传质量报告
        uses: actions/upload-artifact@v3
        with:
          name: quality-report
          path: quality-report.html
      
      - name: 评论PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const report = fs.readFileSync('quality-report.md', 'utf8');
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: report
            });
```

## 附录 C：质量检查脚本示例

### 自定义质量检查脚本（scripts/quality_check.sh）

```
（由于不能包含代码实现，此处用自然语言描述脚本功能）

该脚本主要包含以下检查功能模块：

1. 命名规范检查模块
   - 扫描脚本中的全局变量定义，验证是否使用大写命名
   - 提取所有函数名，验证是否使用小写加下划线
   - 检查是否使用了 Bash 保留字作为变量名

2. 注释完整性检查模块
   - 验证文件头部是否包含功能描述和使用说明
   - 检查每个函数定义前是否有注释说明
   - 统计注释覆盖率

3. 项目特定规则检查模块
   - 验证异常添加格式是否符合规范（严重级别|中文|英文|详情|位置）
   - 检查日志输出函数是否正确使用
   - 验证退出码使用是否符合约定（0/1/2）

4. 文档同步性检查模块
   - 检查 README 中的版本号是否与脚本一致
   - 验证使用示例是否与实际参数匹配
   - 检查异常检测规则表是否完整

5. 报告生成模块
   - 汇总所有检查结果
   - 按严重级别分类输出问题
   - 生成 Markdown 格式或 JSON 格式报告
```

## 附录 D：测试用例示例

### 单元测试用例（tests/unit/test_system_check.bats）

```
（用自然语言描述测试用例结构和验证逻辑）

该测试文件验证系统资源检查函数的正确性：

测试套件：check_system_resources 功能测试

测试用例 1：CPU负载检测 - 正常情况
- 准备：创建包含正常负载数据的模拟诊断目录
- 执行：调用 check_system_resources 函数
- 验证：确认未添加异常记录，日志输出显示"正常"

测试用例 2：CPU负载检测 - 偏高情况
- 准备：创建包含偏高负载数据（2-4倍CPU核心数）的诊断目录
- 执行：调用 check_system_resources 函数
- 验证：确认添加了警告级别的 ELEVATED_SYSTEM_LOAD 异常

测试用例 3：CPU负载检测 - 严重情况
- 准备：创建包含严重负载数据（>4倍CPU核心数）的诊断目录
- 执行：调用 check_system_resources 函数
- 验证：确认添加了严重级别的 HIGH_SYSTEM_LOAD 异常

测试用例 4：内存使用率检测 - 边界值测试
- 准备：创建内存使用率为 85%, 95%, 100% 的三组数据
- 执行：分别测试三种情况
- 验证：
  - 85% 触发 ELEVATED_MEMORY_USAGE 警告
  - 95% 触发 HIGH_MEMORY_USAGE 严重
  - 100% 触发 HIGH_MEMORY_USAGE 严重

测试用例 5：磁盘空间检测 - 多挂载点
- 准备：创建包含多个挂载点数据，部分超过阈值
- 执行：调用检测函数
- 验证：仅对超过阈值的挂载点添加异常记录

测试用例 6：缺失数据场景
- 准备：创建不包含 memory_info 文件的诊断目录
- 执行：调用检测函数
- 验证：函数正常退出，不崩溃，输出"跳过"状态
```
- ShellCheck Error 级别问题数量 = 0
- 代码测试覆盖率 ≥ 80%
- CI 检查通过率 ≥ 95%
- 问题平均修复时间 ≤ 2 天

### 12.2 定性目标

- 团队成员认可质量检查的价值
- 代码 Review 效率显著提升
- 生产环境缺陷率明显下降
- 新人上手难度降低
- 代码库整体可维护性增强

## 13. 参考资源

### 13.1 工具文档

- ShellCheck Wiki：https://github.com/koalaman/shellcheck/wiki
- shfmt 文档：https://github.com/mvdan/sh
- bats-core 文档：https://bats-core.readthedocs.io/

### 13.2 编码规范

- Google Shell Style Guide：https://google.github.io/styleguide/shellguide.html
- Bash Best Practices：https://bertvv.github.io/cheat-sheets/Bash.html

### 13.3 社区资源

- ShellCheck 在线检查：https://www.shellcheck.net/
- Awesome Shell：https://github.com/alebcay/awesome-shell

---

## 附录 A：检查规则清单

### ShellCheck 关键规则

| 规则代码 | 类型 | 说明 | 级别 |
|---------|------|------|------|
| SC2086 | 安全 | 变量应该用引号包围 | Error |
| SC2046 | 安全 | 命令替换应该用引号包围 | Error |
| SC2155 | 安全 | 局部变量声明与赋值应分开 | Warning |
| SC2181 | 最佳实践 | 直接检查退出码而非 $? | Warning |
| SC2034 | 代码质量 | 未使用的变量 | Warning |
| SC2002 | 性能 | 不必要的 cat 使用 | Info |

### 自定义规则

| 规则名称 | 检查内容 | 级别 |
|---------|---------|------|
| NAMING_01 | 全局常量使用大写命名 | Error |
| NAMING_02 | 函数名使用小写加下划线 | Warning |
| COMMENT_01 | 公开函数必须有注释 | Warning |
| FORMAT_01 | 缩进必须为 4 个空格 | Error |
| FORMAT_02 | 单行长度不超过 120 字符 | Warning |

## 附录 B：质量检查工作流配置示例

### GitHub Actions 完整配置

```yaml
name: Code Quality Check

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]
  schedule:
    - cron: '0 2 * * *'  # 每日凌晨2点

jobs:
  shellcheck:
    name: ShellCheck静态分析
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 运行ShellCheck
        uses: ludeeus/action-shellcheck@master
        with:
          scandir: '.'
          severity: warning
          
      - name: 上传ShellCheck报告
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: shellcheck-report
          path: shellcheck.log

  format-check:
    name: 代码格式检查
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 安装shfmt
        run: |
          wget -O shfmt https://github.com/mvdan/sh/releases/download/v3.7.0/shfmt_v3.7.0_linux_amd64
          chmod +x shfmt
          sudo mv shfmt /usr/local/bin/
      
      - name: 检查代码格式
        run: |
          shfmt -d -i 4 -bn -ci -sr *.sh
  
  security-scan:
    name: 安全扫描
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: ShellCheck安全规则
        run: |
          shellcheck -S error *.sh
      
      - name: 自定义安全检查
        run: |
          ./scripts/security_check.sh
  
  unit-test:
    name: 单元测试
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 安装bats-core
        run: |
          sudo apt-get update
          sudo apt-get install -y bats
      
      - name: 安装kcov（覆盖率）
        run: |
          sudo apt-get install -y kcov
      
      - name: 运行单元测试
        run: |
          kcov coverage bats tests/unit/
      
      - name: 上传覆盖率报告
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage/cobertura.xml
          
  integration-test:
    name: 集成测试
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: 准备测试数据
        run: |
          ./tests/fixtures/prepare_test_data.sh
      
      - name: 运行集成测试
        run: |
          bats tests/integration/
  
  compatibility-test:
    name: 兼容性测试 - ${{ matrix.os }}
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: 
          - centos:7
          - ubuntu:18.04
          - registry.cn-hangzhou.aliyuncs.com/alinux/aliyunlinux:3
    container:
      image: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      
      - name: 安装依赖
        run: |
          if [ -f /etc/redhat-release ]; then
            yum install -y bash grep gawk sed coreutils findutils
          else
            apt-get update && apt-get install -y bash grep gawk sed coreutils findutils
          fi
      
      - name: 运行基本测试
        run: |
          bash kudig.sh --help
          bash kudig.sh --version
      
      - name: 功能测试
        run: |
          ./tests/compatibility_test.sh
  
  quality-report:
    name: 生成质量报告
    needs: [shellcheck, format-check, security-scan, unit-test, integration-test]
    runs-on: ubuntu-latest
    if: always()
    steps:
      - uses: actions/checkout@v3
      
      - name: 下载所有报告
        uses: actions/download-artifact@v3
      
      - name: 生成综合报告
        run: |
          ./scripts/generate_quality_report.sh
      
      - name: 上传质量报告
        uses: actions/upload-artifact@v3
        with:
          name: quality-report
          path: quality-report.html
      
      - name: 评论PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const report = fs.readFileSync('quality-report.md', 'utf8');
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: report
            });
```

## 附录 C：质量检查脚本示例

### 自定义质量检查脚本（scripts/quality_check.sh）

```
（由于不能包含代码实现，此处用自然语言描述脚本功能）

该脚本主要包含以下检查功能模块：

1. 命名规范检查模块
   - 扫描脚本中的全局变量定义，验证是否使用大写命名
   - 提取所有函数名，验证是否使用小写加下划线
   - 检查是否使用了 Bash 保留字作为变量名

2. 注释完整性检查模块
   - 验证文件头部是否包含功能描述和使用说明
   - 检查每个函数定义前是否有注释说明
   - 统计注释覆盖率

3. 项目特定规则检查模块
   - 验证异常添加格式是否符合规范（严重级别|中文|英文|详情|位置）
   - 检查日志输出函数是否正确使用
   - 验证退出码使用是否符合约定（0/1/2）

4. 文档同步性检查模块
   - 检查 README 中的版本号是否与脚本一致
   - 验证使用示例是否与实际参数匹配
   - 检查异常检测规则表是否完整

5. 报告生成模块
   - 汇总所有检查结果
   - 按严重级别分类输出问题
   - 生成 Markdown 格式或 JSON 格式报告
```

## 附录 D：测试用例示例

### 单元测试用例（tests/unit/test_system_check.bats）

```
（用自然语言描述测试用例结构和验证逻辑）

该测试文件验证系统资源检查函数的正确性：

测试套件：check_system_resources 功能测试

测试用例 1：CPU负载检测 - 正常情况
- 准备：创建包含正常负载数据的模拟诊断目录
- 执行：调用 check_system_resources 函数
- 验证：确认未添加异常记录，日志输出显示"正常"

测试用例 2：CPU负载检测 - 偏高情况
- 准备：创建包含偏高负载数据（2-4倍CPU核心数）的诊断目录
- 执行：调用 check_system_resources 函数
- 验证：确认添加了警告级别的 ELEVATED_SYSTEM_LOAD 异常

测试用例 3：CPU负载检测 - 严重情况
- 准备：创建包含严重负载数据（>4倍CPU核心数）的诊断目录
- 执行：调用 check_system_resources 函数
- 验证：确认添加了严重级别的 HIGH_SYSTEM_LOAD 异常

测试用例 4：内存使用率检测 - 边界值测试
- 准备：创建内存使用率为 85%, 95%, 100% 的三组数据
- 执行：分别测试三种情况
- 验证：
  - 85% 触发 ELEVATED_MEMORY_USAGE 警告
  - 95% 触发 HIGH_MEMORY_USAGE 严重
  - 100% 触发 HIGH_MEMORY_USAGE 严重

测试用例 5：磁盘空间检测 - 多挂载点
- 准备：创建包含多个挂载点数据，部分超过阈值
- 执行：调用检测函数
- 验证：仅对超过阈值的挂载点添加异常记录

测试用例 6：缺失数据场景
- 准备：创建不包含 memory_info 文件的诊断目录
- 执行：调用检测函数
- 验证：函数正常退出，不崩溃，输出"跳过"状态
```
