# 安全政策

## 支持的版本

以下版本正在接收安全更新：

| 版本 | 支持状态 |
|------|----------|
| v2.0.x | ✅ 支持 |
| v1.2.x | ✅ 支持（仅关键安全修复）|
| v1.1.x 及更早 | ❌ 不再支持 |

## 报告漏洞

**请不要在公开的 GitHub issue 中报告安全漏洞。**

如果您发现了安全漏洞，请通过以下方式私下报告：

- 邮箱: security@kudig.io
- GitHub Security Advisories: [创建安全公告](https://github.com/kudig/kudig/security/advisories/new)

请在报告中包含以下信息：

- **漏洞描述**: 清晰描述漏洞
- **影响版本**: 受影响的版本
- **复现步骤**: 详细的复现步骤
- **影响范围**: 漏洞可能造成的危害
- **建议修复**: 如果有，提供修复建议
- **您的联系信息**: 以便我们联系您获取更多信息

## 响应时间线

我们承诺在收到安全报告后：

- **24 小时内**: 确认收到报告
- **72 小时内**: 完成初步评估
- **7 天内**: 提供修复计划或请求更多信息
- **30 天内**: 发布修复（取决于漏洞严重程度）

## 披露政策

我们遵循负责任的披露原则：

1. 报告者私下向维护者报告漏洞
2. 维护者确认并修复漏洞
3. 维护者发布修复版本
4. 维护者公开披露漏洞（通常在修复后 30 天）

## 安全最佳实践

### 对于用户

1. **保持更新**: 始终使用最新版本
2. **权限最小化**: 以最小必要权限运行 kudig
3. **审查配置**: 仔细检查配置文件和规则
4. **监控日志**: 定期检查诊断日志

### 对于开发者

1. **依赖管理**: 定期更新依赖项
   ```bash
   cd v2-go
   go mod tidy
   go list -u -m all
   ```

2. **安全扫描**: 使用安全扫描工具
   ```bash
   # 检查依赖漏洞
   govulncheck ./...
   
   # 静态安全分析
   gosec ./...
   ```

3. **代码审查**: 所有代码更改必须经过审查
4. **秘密管理**: 不要将凭据提交到代码仓库
5. **输入验证**: 验证所有用户输入

## 已知漏洞

| 漏洞 ID | 描述 | 影响版本 | 修复版本 | 披露日期 |
|---------|------|----------|----------|----------|
| 暂无 | - | - | - | - |

## 安全相关配置

### 容器安全

使用 Docker 时：

```dockerfile
# 使用非 root 用户
USER 1000:1000

# 只读根文件系统
readOnlyRootFilesystem: true

# 禁止特权模式
privileged: false
```

### Kubernetes 安全

部署 kudig 到 Kubernetes 时：

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
  capabilities:
    drop:
      - ALL
```

## 安全更新通知

订阅安全更新：

- 关注 [GitHub Security Advisories](https://github.com/kudig/kudig/security/advisories)
- 关注 Releases 页面的安全更新

## 致谢

感谢以下安全研究人员对 kudig 安全性的贡献：

- 待添加

## 许可证

本安全政策遵循项目 [Apache 2.0](LICENSE) 许可证。
