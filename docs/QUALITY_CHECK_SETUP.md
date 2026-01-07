# 代码质量检查工具安装指南

本文档说明如何安装和配置 kudig.sh 项目的代码质量检查工具链。

## 目录

- [工具清单](#工具清单)
- [Linux 安装](#linux-安装)
- [macOS 安装](#macos-安装)
- [Windows WSL 安装](#windows-wsl-安装)
- [编辑器集成](#编辑器集成)
- [Git Hook 配置](#git-hook-配置)
- [验证安装](#验证安装)

## 工具清单

| 工具 | 用途 | 必需 |
|------|------|------|
| ShellCheck | Shell 脚本静态分析 | ✓ |
| shfmt | Shell 脚本格式化 | ✓ |
| bats-core | Bash 自动化测试框架 | 推荐 |
| kcov | 代码覆盖率统计 | 可选 |

## Linux 安装

### CentOS / RHEL / AliyunLinux

```bash
# ShellCheck
sudo yum install -y epel-release
sudo yum install -y ShellCheck

# shfmt
wget -O /tmp/shfmt https://github.com/mvdan/sh/releases/download/v3.7.0/shfmt_v3.7.0_linux_amd64
chmod +x /tmp/shfmt
sudo mv /tmp/shfmt /usr/local/bin/shfmt

# bats-core
git clone https://github.com/bats-core/bats-core.git /tmp/bats-core
cd /tmp/bats-core
sudo ./install.sh /usr/local

# kcov (可选)
sudo yum install -y cmake gcc-c++ binutils-devel elfutils-libelf-devel
git clone https://github.com/SimonKagstrom/kcov.git /tmp/kcov
cd /tmp/kcov
mkdir build && cd build
cmake ..
make
sudo make install
```

### Ubuntu / Debian

```bash
# ShellCheck
sudo apt-get update
sudo apt-get install -y shellcheck

# shfmt
sudo snap install shfmt
# 或者手动安装
wget -O /tmp/shfmt https://github.com/mvdan/sh/releases/download/v3.7.0/shfmt_v3.7.0_linux_amd64
chmod +x /tmp/shfmt
sudo mv /tmp/shfmt /usr/local/bin/shfmt

# bats-core
sudo apt-get install -y bats

# kcov (可选)
sudo apt-get install -y kcov
```

## macOS 安装

使用 Homebrew 安装所有工具：

```bash
# 安装 Homebrew (如果尚未安装)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装工具
brew install shellcheck
brew install shfmt
brew install bats-core
brew install kcov  # 可选
```

## Windows WSL 安装

首先安装 WSL（Windows Subsystem for Linux）：

```powershell
# 在 PowerShell 管理员模式下运行
wsl --install
```

然后在 WSL 中按照 Linux 安装步骤进行。

## 编辑器集成

### Visual Studio Code

安装以下扩展：

1. **ShellCheck** (timonwong.shellcheck)
   ```
   ext install timonwong.shellcheck
   ```

2. **shell-format** (foxundermoon.shell-format)
   ```
   ext install foxundermoon.shell-format
   ```

配置 VSCode 设置（`.vscode/settings.json`）：

```json
{
    "shellcheck.enable": true,
    "shellcheck.run": "onSave",
    "shellcheck.exclude": ["SC2034", "SC1090", "SC1091"],
    "shellformat.flag": "-i 4 -bn -ci -sr",
    "files.associations": {
        "*.sh": "shellscript"
    },
    "editor.formatOnSave": true,
    "[shellscript]": {
        "editor.defaultFormatter": "foxundermoon.shell-format"
    }
}
```

### Vim / Neovim

使用 ALE (Asynchronous Lint Engine)：

```vim
" 安装 vim-plug
call plug#begin()
Plug 'dense-analysis/ale'
call plug#end()

" 配置 ALE
let g:ale_linters = {'sh': ['shellcheck']}
let g:ale_fixers = {'sh': ['shfmt']}
let g:ale_sh_shfmt_options = '-i 4 -bn -ci -sr'
let g:ale_fix_on_save = 1
```

### IntelliJ IDEA / WebStorm

1. 安装 **Shell Script** 插件
2. 在设置中配置 ShellCheck 路径：
   - Settings → Tools → Shell Script → ShellCheck
   - 指定 ShellCheck 可执行文件路径

## Git Hook 配置

### 自动安装 Pre-commit Hook

在项目根目录运行：

```bash
# 创建 hooks 目录（如果不存在）
mkdir -p .git/hooks

# 复制 pre-commit 脚本
cp scripts/pre-commit .git/hooks/pre-commit

# 添加执行权限
chmod +x .git/hooks/pre-commit
```

### 手动配置

创建 `.git/hooks/pre-commit` 文件并添加以下内容：

```bash
#!/usr/bin/env bash
exec bash scripts/pre-commit
```

然后添加执行权限：

```bash
chmod +x .git/hooks/pre-commit
```

### 禁用 Hook（临时）

如果需要临时跳过检查：

```bash
git commit --no-verify -m "Your commit message"
```

## 验证安装

### 检查工具版本

```bash
# ShellCheck
shellcheck --version

# shfmt
shfmt --version

# bats
bats --version

# kcov (可选)
kcov --version
```

### 运行质量检查

```bash
# 检查 kudig.sh
shellcheck kudig.sh
shfmt -d -i 4 -bn -ci -sr kudig.sh
bash scripts/quality_check.sh kudig.sh

# 运行所有检查
bash scripts/quality_check.sh --verbose kudig.sh
```

### 测试 Pre-commit Hook

```bash
# 修改一个文件并尝试提交
echo "# test" >> kudig.sh
git add kudig.sh
git commit -m "test commit"

# 应该看到质量检查运行
```

## 持续集成检查

CI 环境会自动安装工具并运行检查，无需额外配置。

参见 `.github/workflows/code-quality.yml` 了解 CI 配置详情。

## 常见问题

### Q: ShellCheck 报告 SC2034 "变量未使用"

**A**: 某些全局变量在声明时可能未立即使用，这是正常的。项目配置已排除此规则。

### Q: shfmt 格式化改变了我的代码风格

**A**: shfmt 按照项目统一的格式规范进行格式化。建议接受自动格式化以保持一致性。

### Q: Pre-commit Hook 检查失败

**A**: 
1. 检查工具是否正确安装
2. 查看具体错误信息并修复
3. 如果必要，可以使用 `git commit --no-verify` 跳过（不推荐）

### Q: 在 Windows 上如何使用？

**A**: 推荐使用 WSL (Windows Subsystem for Linux) 环境，完全支持所有工具。

## 更多帮助

- [ShellCheck Wiki](https://github.com/koalaman/shellcheck/wiki)
- [shfmt 文档](https://github.com/mvdan/sh/blob/master/cmd/shfmt/shfmt.1.scd)
- [bats-core 文档](https://bats-core.readthedocs.io/)
- [项目 README](../README.md)

---

如有问题，请提交 Issue 或联系项目维护者。
