# kudig v1.2.0 - Bash 版本

> ✅ **生产可用** - 稳定的 Bash 脚本实现

## 简介

`kudig` v1.0 是一个高效的 Kubernetes 节点诊断日志分析工具，用于分析 `diagnose_k8s.sh` 收集的诊断数据。

## 核心特性

- **Bash 脚本实现**：轻量级、无依赖、开箱即用
- **离线分析模式**：分析 `diagnose_k8s.sh` 收集的诊断数据
- **120+项异常检测规则**：涵盖系统资源、进程服务、网络、内核、容器运行时、Kubernetes、存储、安全、性能、日志、容器Pod、网络性能、硬件健康、应用层、安全审计、容量规划、高可用性、备份恢复、合规性、监控告警等维度
- **生产环境就绪**：专为生产环境设计，包含硬件健康、安全审计、高可用性等生产级检测
- **中英文双语报告**：支持文本和 JSON 两种输出格式
- **智能排查建议**：针对每个异常提供详细的排查步骤和解决方案
- **逐项详细输出**：显示每个检查项的名称、状态和结果

## 快速开始

### 1. 安装

```bash
# 直接使用（无需安装）
chmod +x kudig
```

### 2. 收集诊断数据

```bash
# 在 Kubernetes 节点上收集诊断数据
sudo ./diagnose_k8s.sh
# 生成目录: /tmp/diagnose_1702468800
```

### 3. 分析诊断数据

```bash
# 基本使用
./kudig /tmp/diagnose_1702468800

# 详细模式（显示每一步检查项）
./kudig --verbose /tmp/diagnose_1702468800

# JSON 格式输出
./kudig --json /tmp/diagnose_1702468800

# 保存报告到文件
./kudig -o report.txt /tmp/diagnose_1702468800
```

## 命令行选项

| 选项 | 说明 |
|-----|------|
| `-h, --help` | 显示帮助信息 |
| `-v, --version` | 显示版本信息 |
| `--verbose` | 详细输出模式（显示每一步检查项） |
| `--json` | 输出 JSON 格式（默认文本格式） |
| `-o, --output <文件>` | 保存报告到指定文件 |
| `<诊断目录>` | diagnose_k8s.sh 生成的诊断目录路径 |

## 内置检测规则（120+ 项）

### 1. 系统资源检测（6 项）
- **CPU 负载检测**（2倍/4倍核心数阈值）
  - 排查命令：`uptime`、`top -bn1 | head -5`
  - 修复建议：检查高CPU进程 `ps aux --sort=-%cpu | head -10`
- **内存使用检测**（85%/95%阈值）
  - 排查命令：`free -h`、`cat /proc/meminfo`
  - 修复建议：检查内存占用 `ps aux --sort=-%mem | head -10`
- **磁盘空间检测**（90%/95%阈值）
  - 排查命令：`df -h`、`du -sh /* | sort -rh | head -10`
  - 修复建议：清理无用文件 `find /var/log -name "*.log" -delete`
- **文件句柄检测**（>50000警告）
  - 排查命令：`lsof | wc -l`、`cat /proc/sys/fs/file-nr`
  - 修复建议：增加文件句柄限制 `echo "* soft nofile 65536" >> /etc/security/limits.conf`
- **进程/线程数检测**（>5000/10000阈值）
  - 排查命令：`ps -elf | wc -l`、`ps -eo nlwp | sort -nr | head -10`
  - 修复建议：检查异常进程 `ps -eo pid,nlwp,cmd | sort -nrk2 | head -5`
- **磁盘使用检测**
  - 排查命令：`iostat -x 1 5`、`df -i`
  - 修复建议：检查 inode 使用 `df -i | sort -nrk5`

### 2. 进程与服务检测（6 项）
- **Kubelet 服务检测**
  - 排查命令：`systemctl status kubelet`、`journalctl -u kubelet -n 50`
  - 修复建议：重启服务 `systemctl restart kubelet`
- **容器运行时检测**（Docker/Containerd）
  - 排查命令：`systemctl status docker containerd`
  - 修复建议：重启运行时 `systemctl restart containerd`
- **ps 命令挂起检测**
  - 排查命令：`timeout 5 ps -ef`、`ps aux | grep " D "`
  - 修复建议：检查 D 状态进程 `ps aux | grep " D "`
- **D 状态进程检测**
  - 排查命令：`ps aux | grep " D "`
  - 修复建议：检查 IO 设备状态 `iostat -x`
- **Runc 进程检测**
  - 排查命令：`ps aux | grep runc | wc -l`
  - 修复建议：检查容器状态 `crictl ps`
- **Firewalld 检测**
  - 排查命令：`systemctl status firewalld`、`iptables -L`
  - 修复建议：关闭 firewalld `systemctl stop firewalld && systemctl disable firewalld`

### 3. 网络检测（5 项）
- **连接追踪表检测**（80%/95%阈值）
  - 排查命令：`cat /proc/sys/net/netfilter/nf_conntrack_count`、`sysctl net.netfilter.nf_conntrack_max`
  - 修复建议：增加连接追踪表大小 `sysctl -w net.netfilter.nf_conntrack_max=262144`
- **网卡状态检测**
  - 排查命令：`ip link show`、`ifconfig -a`
  - 修复建议：启用网卡 `ip link set <interface> up`
- **默认路由检测**
  - 排查命令：`ip route show default`、`route -n`
  - 修复建议：添加默认路由 `ip route add default via <gateway> dev <interface>`
- **Kubelet 端口监听检测**
  - 排查命令：`netstat -tlnp | grep 10250`、`ss -tlnp | grep 10250`
  - 修复建议：检查 kubelet 配置 `cat /etc/kubernetes/kubelet.conf`
- **Iptables 规则检测**
  - 排查命令：`iptables -L -n | wc -l`、`iptables-save | wc -l`
  - 修复建议：清理无用规则 `iptables -F && iptables -X`

### 4. 内核检测（7 项）
- **内核 Panic 检测**
  - 排查命令：`dmesg | grep -i panic`、`cat /var/log/dmesg | grep -i panic`
  - 修复建议：检查内核日志 `journalctl -k -n 100`
- **OOM Killer 检测**
  - 排查命令：`dmesg | grep -i "out of memory"`、`grep -i oom /var/log/messages`
  - 修复建议：增加内存或调整 Pod 资源限制
- **Messages 日志 OOM 检测**
  - 排查命令：`grep -i oom /var/log/messages`、`grep -i "kill process" /var/log/messages`
  - 修复建议：检查被 kill 的进程 `grep -i killed /var/log/messages`
- **文件系统只读检测**
  - 排查命令：`mount | grep ro`、`touch /tmp/test && rm /tmp/test`
  - 修复建议：检查文件系统 `fsck -y /dev/<device>`
- **磁盘 IO 错误检测**
  - 排查命令：`dmesg | grep -i "io error"`、`smartctl -a /dev/<device>`
  - 修复建议：检查磁盘健康状态 `smartctl -H /dev/<device>`
- **内核模块加载检测**
  - 排查命令：`dmesg | grep -i "module.*failed"`、`lsmod`
  - 修复建议：检查模块依赖 `modinfo <module>`
- **NMI Watchdog 检测**
  - 排查命令：`dmesg | grep -i nmi`、`cat /proc/sys/kernel/nmi_watchdog`
  - 修复建议：调整 watchdog 设置 `sysctl -w kernel.nmi_watchdog=0`

### 5. 容器运行时检测（4 项）
- **Docker 启动检测**
  - 排查命令：`docker info`、`systemctl status docker`
  - 修复建议：重启 Docker `systemctl restart docker`
- **Docker 存储驱动检测**
  - 排查命令：`docker info | grep Storage`、`docker system df`
  - 修复建议：检查存储驱动配置 `cat /etc/docker/daemon.json`
- **Containerd 容器创建检测**
  - 排查命令：`crictl ps`、`crictl info`
  - 修复建议：重启 containerd `systemctl restart containerd`
- **镜像拉取检测**
  - 排查命令：`crictl pull <image>`、`docker pull <image>`
  - 修复建议：检查镜像仓库连接 `curl -v <registry-url>`

### 6. Kubernetes 组件检测（9 项）
- **PLEG 状态检测**
  - 排查命令：`grep -i "pleg" /var/log/kubelet.log`、`journalctl -u kubelet | grep -i pleg`
  - 修复建议：重启容器运行时 `systemctl restart containerd`
- **CNI 插件检测**
  - 排查命令：`ls /etc/cni/net.d/`、`grep -i cni /var/log/kubelet.log`
  - 修复建议：检查 CNI 配置 `cat /etc/cni/net.d/*`
- **Kubelet 证书检测**
  - 排查命令：`openssl x509 -in /var/lib/kubelet/pki/kubelet.crt -text -noout`
  - 修复建议：更新证书 `kubeadm certs renew all`
- **API Server 连接检测**
  - 排查命令：`curl -k https://<apiserver>:6443/healthz`、`kubectl cluster-info`
  - 修复建议：检查 kubeconfig `cat ~/.kube/config`
- **Kubelet 认证检测**
  - 排查命令：`grep -i "unauthorized" /var/log/kubelet.log`
  - 修复建议：检查认证配置 `cat /etc/kubernetes/kubelet.conf`
- **Pod 驱逐检测**
  - 排查命令：`grep -i evict /var/log/kubelet.log`、`kubectl get events | grep Evict`
  - 修复建议：增加节点资源或调整驱逐阈值
- **节点状态检测**
  - 排查命令：`kubectl get node <node> -o yaml | grep -A 10 status:`
  - 修复建议：检查节点条件 `kubectl describe node <node>`
- **磁盘压力检测**
  - 排查命令：`kubectl describe node <node> | grep -A 5 "DiskPressure"`
  - 修复建议：清理磁盘空间 `df -h`
- **内存压力检测**
  - 排查命令：`kubectl describe node <node> | grep -A 5 "MemoryPressure"`
  - 修复建议：检查内存使用 `free -h`

### 7. 时间同步检测（1 项）
- **NTP/Chrony 状态检测**
  - 排查命令：`timedatectl status`、`chronyc sources`
  - 修复建议：启用时间同步 `systemctl start chronyd && systemctl enable chronyd`

### 8. 配置检测（5 项）
- **Swap 配置检测**
  - 排查命令：`swapon -s`、`cat /proc/swaps`
  - 修复建议：禁用 swap `swapoff -a && sed -i '/swap/d' /etc/fstab`
- **IP 转发检测**
  - 排查命令：`sysctl net.ipv4.ip_forward`、`cat /proc/sys/net/ipv4/ip_forward`
  - 修复建议：启用 IP 转发 `sysctl -w net.ipv4.ip_forward=1`
- **bridge-nf-call-iptables 检测**
  - 排查命令：`sysctl net.bridge.bridge-nf-call-iptables`
  - 修复建议：启用桥接转发 `sysctl -w net.bridge.bridge-nf-call-iptables=1`
- **ulimit open files 检测**
  - 排查命令：`ulimit -n`、`cat /etc/security/limits.conf | grep nofile`
  - 修复建议：增加限制 `echo "* soft nofile 65536" >> /etc/security/limits.conf`
- **SELinux 检测**
  - 排查命令：`getenforce`、`sestatus`
  - 修复建议：设置为 Permissive `setenforce 0`

### 9. 存储和文件系统检测（4 项）
- **磁盘 IO 等待检测**（>20%/50%阈值）
  - 排查命令：`iostat -x 1 5`、`sar -d 1 5`
  - 修复建议：检查 IO 密集进程 `iotop -o`
- **文件系统挂载选项检测**（noatime）
  - 排查命令：`mount | grep noatime`、`cat /etc/fstab`
  - 修复建议：添加 noatime 选项到 /etc/fstab
- **EXT4 文件系统错误检测**
  - 排查命令：`dmesg | grep -i "ext4.*error"`、`fsck -n /dev/<device>`
  - 修复建议：修复文件系统 `fsck -y /dev/<device>`
- **NFS 挂载状态检测**
  - 排查命令：`mount | grep nfs`、`showmount -e <nfs-server>`
  - 修复建议：检查 NFS 服务 `systemctl status nfs-server`

### 10. 安全配置检测（4 项）
- **密码策略检测**
  - 排查命令：`cat /etc/login.defs | grep PASS`、`cat /etc/pam.d/system-auth | grep pam_cracklib`
  - 修复建议：调整密码策略 `vim /etc/login.defs`
- **SSH 配置检测**（root登录）
  - 排查命令：`grep PermitRootLogin /etc/ssh/sshd_config`
  - 修复建议：禁用 root 登录 `sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config && systemctl restart sshd`
- **防火墙状态检测**
  - 排查命令：`systemctl status firewalld iptables`、`iptables -L`
  - 修复建议：配置防火墙规则 `iptables -A INPUT -p tcp --dport 6443 -j ACCEPT`
- **AppArmor/SELinux 状态检测**
  - 排查命令：`apparmor_status`、`getenforce`
  - 修复建议：根据需要调整安全模块配置

### 11. 性能指标检测（4 项）
- **上下文切换率检测**（>100000/s）
  - 排查命令：`vmstat 1 5 | awk '{print $12}'`、`sar -w 1 5`
  - 修复建议：检查高上下文切换进程 `pidstat -w 1`
- **系统调用率检测**（>500000/s）
  - 排查命令：`strace -c -p <pid>`、`sar -c 1 5`
  - 修复建议：分析系统调用 `strace -p <pid>`
- **缓存使用检测**（>80%）
  - 排查命令：`free -h`、`cat /proc/meminfo | grep -E 'Cached|Buffers'`
  - 修复建议：清理缓存 `sync && echo 3 > /proc/sys/vm/drop_caches`
- **中断率检测**（>50000/s）
  - 排查命令：`vmstat 1 5 | awk '{print $11}'`、`sar -I SUM 1 5`
  - 修复建议：检查中断源 `cat /proc/interrupts | sort -nrk2 | head -10`

### 12. 日志深度分析检测（5 项）
- **Kubelet 错误统计**（>100）
  - 排查命令：`grep -i error /var/log/kubelet.log | wc -l`、`journalctl -u kubelet | grep -i error | wc -l`
  - 修复建议：查看详细错误 `grep -i error /var/log/kubelet.log | tail -20`
- **容器重启检测**（>10次）
  - 排查命令：`grep -i "container.*restart" /var/log/kubelet.log | wc -l`
  - 修复建议：检查容器日志 `crictl logs <container-id>`
- **镜像拉取超时检测**（>5次）
  - 排查命令：`grep -i "timeout" /var/log/kubelet.log | grep -i pull | wc -l`
  - 修复建议：检查网络连接 `ping <registry>`
- **存储空间不足检测**
  - 排查命令：`grep -i "no space left" /var/log/kubelet.log`、`df -h`
  - 修复建议：清理磁盘空间 `du -sh /* | sort -rh | head -10`
- **内核错误统计**（>50）
  - 排查命令：`grep -i "kernel:" /var/log/messages | wc -l`、`dmesg | grep -i error | wc -l`
  - 修复建议：查看内核错误详情 `dmesg | grep -i error | tail -20`

### 13. 容器和 Pod 状态检测（5 项）
- **容器创建失败检测**（>10次）
  - 排查命令：`grep -i "failed to create container" /var/log/kubelet.log | wc -l`
  - 修复建议：检查容器运行时 `crictl info`
- **容器启动超时检测**（>5次）
  - 排查命令：`grep -i "startcontainer.*timed out" /var/log/kubelet.log | wc -l`
  - 修复建议：检查容器启动时间 `crictl inspect <container-id> | grep StartedAt`
- **Pod 挂载失败检测**（>5次）
  - 排查命令：`kubectl get events | grep FailedMount | wc -l`
  - 修复建议：检查 PV/PVC 状态 `kubectl get pv,pvc`
- **Pod 沙箱创建失败检测**（>5次）
  - 排查命令：`grep -i "failed to create pod sandbox" /var/log/kubelet.log | wc -l`
  - 修复建议：检查 CNI 配置 `cat /etc/cni/net.d/*`
- **容器运行时健康检测**
  - 排查命令：`crictl info`、`docker info`
  - 修复建议：重启容器运行时 `systemctl restart containerd`

### 14. 网络性能检测（5 项）
- **网络错误检测**（>100）
  - 排查命令：`ifconfig <interface> | grep -E 'errors|dropped'`、`netstat -i`
  - 修复建议：检查物理连接和网卡驱动
- **网络丢包检测**（>1000）
  - 排查命令：`ping -c 100 <ip> | grep packet`、`mtr <ip>`
  - 修复建议：检查网络设备和线缆
- **TCP 连接状态检测**（TIME_WAIT>10000）
  - 排查命令：`netstat -ant | grep TIME_WAIT | wc -l`、`ss -s`
  - 修复建议：调整 TCP 参数 `sysctl -w net.ipv4.tcp_fin_timeout=30`
- **网络带宽使用检测**
  - 排查命令：`sar -n DEV 1 5`、`iftop -i <interface>`
  - 修复建议：检查带宽占用 `ps aux | grep -E '(wget|curl|scp)'`
- **DNS 配置检测**
  - 排查命令：`cat /etc/resolv.conf`、`nslookup kubernetes.default.svc.cluster.local`
  - 修复建议：配置正确的 DNS 服务器 `echo "nameserver 8.8.8.8" > /etc/resolv.conf`

### 15. 硬件健康检测（5 项）
- **CPU温度检测**（>70°C/85°C阈值）
  - 排查命令：`sensors`、`cat /sys/class/thermal/thermal_zone*/temp`
  - 修复建议：检查散热系统，清理灰尘
- **内存ECC错误检测**
  - 排查命令：`dmesg | grep -i ecc`、`journalctl -k | grep -i mce`
  - 修复建议：更换故障内存条
- **磁盘SMART健康检测**
  - 排查命令：`smartctl -H /dev/<device>`、`smartctl -a /dev/<device>`
  - 修复建议：备份并更换故障磁盘
- **硬件错误日志检测**
  - 排查命令：`dmesg | grep -i "hardware error"`、`journalctl -k -p err`
  - 修复建议：检查硬件设备状态
- **CPU频率降频检测**
  - 排查命令：`cat /proc/cpuinfo | grep MHz`、`cpupower frequency-info`
  - 修复建议：检查电源管理和散热

### 16. 应用层检测（5 项）
- **Pod健康检查失败检测**（>10次）
  - 排查命令：`kubectl get pods`、`kubectl describe pod <pod-name>`
  - 修复建议：检查应用健康检查配置
- **应用启动超时检测**（>5次）
  - 排查命令：`kubectl logs <pod-name>`、`kubectl describe pod <pod-name>`
  - 修复建议：优化应用启动时间或增加超时时间
- **容器资源限制检测**
  - 排查命令：`kubectl describe pod <pod-name> | grep -A 10 Limits`
  - 修复建议：为容器设置合理的资源限制
- **应用崩溃检测**（>10次）
  - 排查命令：`kubectl logs <pod-name> --previous`、`kubectl describe pod <pod-name>`
  - 修复建议：检查应用日志，修复崩溃原因
- **应用日志错误检测**（>100条）
  - 排查命令：`kubectl logs <pod-name>`、`journalctl -u kubelet`
  - 修复建议：分析应用错误日志，修复问题

### 17. 安全审计检测（5 项）
- **失败登录尝试检测**（>10次/50次阈值）
  - 排查命令：`lastb`、`grep "Failed password" /var/log/auth.log`
  - 修复建议：配置fail2ban或限制SSH访问
- **特权容器检测**
  - 排查命令：`kubectl get pods --all-namespaces -o jsonpath='{.items[*].spec.containers[*].securityContext.privileged}'`
  - 修复建议：移除特权模式，使用最小权限原则
- **root用户容器检测**（>5个）
  - 排查命令：`kubectl get pods --all-namespaces -o jsonpath='{.items[*].spec.containers[*].securityContext.runAsUser}'`
  - 修复建议：使用非root用户运行容器
- **敏感挂载检测**
  - 排查命令：`kubectl describe pod <pod-name> | grep -A 5 Mounts`
  - 修复建议：避免挂载敏感路径
- **权限提升操作检测**（>20次）
  - 排查命令：`grep sudo /var/log/auth.log | wc -l`、`last -n 20`
  - 修复建议：审计权限提升操作，配置sudo日志

### 18. 容量规划检测（5 项）
- **CPU使用率趋势检测**（>80%）
  - 排查命令：`kubectl top nodes`、`kubectl describe node <node-name>`
  - 修复建议：增加节点或优化应用资源使用
- **内存使用率趋势检测**（>80%）
  - 排查命令：`kubectl top nodes`、`free -h`
  - 修复建议：增加节点内存或优化应用内存使用
- **磁盘使用率趋势检测**（>80%）
  - 排查命令：`df -h`、`du -sh /* | sort -rh | head -10`
  - 修复建议：清理磁盘或增加存储容量
- **网络带宽趋势检测**（RX>80GB/s）
  - 排查命令：`sar -n DEV 1 5`、`iftop -i <interface>`
  - 修复建议：增加网络带宽或优化网络使用
- **inode使用率趋势检测**（>80%）
  - 排查命令：`df -i`、`find / -xdev -printf '%h\n' | sort | uniq -c | sort -k 1 -n | tail -10`
  - 修复建议：清理小文件或增加inode数量

### 19. 高可用性检查（5 项）
- **节点NotReady状态检测**
  - 排查命令：`kubectl get nodes`、`kubectl describe node <node-name>`
  - 修复建议：检查节点状态，修复NotReady原因
- **API Server连接检测**（>5次错误）
  - 排查命令：`kubectl cluster-info`、`curl -k https://<apiserver>:6443/healthz`
  - 修复建议：检查API Server状态和网络连接
- **etcd健康检测**
  - 排查命令：`kubectl get cs`、`etcdctl endpoint health`
  - 修复建议：检查etcd集群状态
- **控制器管理器检测**
  - 排查命令：`kubectl get pods -n kube-system | grep controller-manager`
  - 修复建议：重启控制器管理器
- **调度器检测**
  - 排查命令：`kubectl get pods -n kube-system | grep scheduler`
  - 修复建议：重启调度器

### 20. 备份和恢复检测（5 项）
- **etcd备份检测**
  - 排查命令：`ls -lh /etc/kubernetes/pki/etcd/`、`crontab -l | grep etcd`
  - 修复建议：配置etcd定期备份
- **PV备份检测**
  - 排查命令：`kubectl get pv`、`kubectl get storageclass`
  - 修复建议：配置PV快照或备份策略
- **备份失败检测**
  - 排查命令：`journalctl -u backup-service`、`grep backup /var/log/syslog`
  - 修复建议：检查备份配置和存储空间
- **恢复测试检测**
  - 排查命令：`ls -lh /backup/`、`cat /var/log/restore-test.log`
  - 修复建议：定期执行恢复测试
- **备份存储检测**
  - 排查命令：`df -h /backup`、`mount | grep backup`
  - 修复建议：配置外部备份存储（S3、NFS等）

### 21. 合规性检查（5 项）
- **RBAC配置检测**
  - 排查命令：`kubectl auth can-i --list`、`kubectl get clusterrole`
  - 修复建议：启用并配置RBAC
- **网络策略检测**
  - 排查命令：`kubectl get networkpolicy --all-namespaces`
  - 修复建议：配置网络策略限制Pod间通信
- **Pod安全策略检测**
  - 排查命令：`kubectl get psp`、`kubectl get podsecuritypolicy`
  - 修复建议：配置Pod安全策略
- **镜像扫描检测**
  - 排查命令：`trivy image <image-name>`、`kubectl get image`
  - 修复建议：配置镜像漏洞扫描
- **审计日志检测**
  - 排查命令：`kubectl get pods -n kube-system | grep audit`、`grep audit /var/log/kube-audit/`
  - 修复建议：启用审计日志

### 22. 监控和告警检测（5 项）
- **Prometheus检测**
  - 排查命令：`kubectl get pods -n monitoring | grep prometheus`、`kubectl get svc -n monitoring`
  - 修复建议：部署Prometheus监控系统
- **Grafana检测**
  - 排查命令：`kubectl get pods -n monitoring | grep grafana`、`kubectl get svc -n monitoring`
  - 修复建议：部署Grafana可视化面板
- **告警规则检测**
  - 排查命令：`kubectl get prometheusrule --all-namespaces`、`kubectl get alertmanager`
  - 修复建议：配置告警规则
- **日志收集检测**
  - 排查命令：`kubectl get pods -n logging | grep -E 'fluentd|filebeat|logstash'`
  - 修复建议：部署日志收集系统
- **告警通知检测**
  - 排查命令：`kubectl get configmap -n monitoring | grep alertmanager`、`kubectl get secret -n monitoring | grep slack`
  - 修复建议：配置告警通知渠道

## 输出示例

### 文本格式（默认）

```
================================================================
  kudig v1.2.0 - Kubernetes节点诊断分析工具
================================================================

诊断目录: /tmp/diagnose_1702468800
分析时间: 2026-01-16 10:50:01

开始诊断检查...

========== 系统资源检查 ==========
  [✓] CPU负载: 正常 (15min负载: 0.34, CPU核心: 8)
  [✓] 内存使用: 正常 (使用率: 32%)
  [✓] 磁盘空间: 正常 (所有挂载点使用率<90%)
  [✓] 文件句柄: 正常 (最大: 273)
  [✓] 进程/线程数: 正常 (最大线程数: 33)
  [✓] 磁盘使用: 正常 (所有挂载点<90%)

========== 进程与服务检查 ==========
  [✓] Kubelet服务: running
  [✓] 容器运行时: docker=unknown, containerd=running
  ...

========== 网络状态检查 ==========
  [✓] 连接跟踪表: 正常 (517/262144, 0%)
  [!] 网卡状态: 部分down (kube-ipvs0,nodelocaldns)
  ...

================================================================
  诊断结果汇总
================================================================
=== Kubernetes节点诊断异常报告 ===

-------------------------------------------
【警告级别】异常项
-------------------------------------------
[警告] 网卡接口down | NETWORK_INTERFACE_DOWN
  详情: 以下网卡处于down状态: kube-ipvs0,nodelocaldns
  位置: network_info

-------------------------------------------
异常统计
-------------------------------------------
严重: 0 项
警告: 1 项
提示: 0 项
总计: 1 项
```

## 测试

参见 [TESTING.md](TESTING.md) 了解详细测试说明。

## 文件结构

```
v1-bash/
├── kudig                 # 主脚本
├── README.md             # 本文档
├── TESTING.md            # 测试文档
└── reference/            # 参考诊断数据
    └── diagnose_k8s/
        └── diagnose_*/   # 示例诊断目录
```

## 版本历史

- **v1.2.0** (2026-01-27)
  - 新增 40 项异常检测规则，总计 120+ 项
  - 新增硬件健康诊断（5项）：CPU温度、内存ECC、磁盘SMART、硬件错误、CPU降频
  - 新增应用层诊断（5项）：Pod健康、启动超时、资源限制、应用崩溃、日志错误
  - 新增安全审计诊断（5项）：失败登录、特权容器、root容器、敏感挂载、权限提升
  - 新增容量规划诊断（5项）：CPU/内存/磁盘/网络/inode趋势
  - 新增高可用性检查（5项）：节点状态、API Server、etcd、控制器、调度器
  - 新增备份和恢复检测（5项）：etcd备份、PV备份、备份失败、恢复测试、备份存储
  - 新增合规性检查（5项）：RBAC、网络策略、Pod安全策略、镜像扫描、审计日志
  - 新增监控和告警检测（5项）：Prometheus、Grafana、告警规则、日志收集、告警通知
  - 生产环境就绪，适用于生产环境诊断

- **v1.1.0** (2026-01-27)
  - 新增 37 项异常检测规则，总计 80+ 项
  - 新增存储和文件系统诊断（4项）
  - 新增安全配置诊断（4项）
  - 新增性能指标诊断（4项）
  - 新增日志深度分析（5项）
  - 新增容器和Pod状态诊断（5项）
  - 新增网络性能诊断（5项）
  - 增强诊断覆盖范围，符合 kusheet 诊断标准

- **v1.0.0** (2024-12)
  - 初始 Bash 脚本版本
  - 43 项异常检测规则
  - 支持文本和 JSON 输出格式
  - 逐项详细输出每个检查项状态
  - 为每个异常提供排查建议

## 许可证

Apache License 2.0

## 相关链接

- [v2.0 Go 版本](../v2-go/) - 🚧 开发中
- [项目主页](../)
