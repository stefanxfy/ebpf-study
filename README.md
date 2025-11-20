# eBPF Study

这是一个 eBPF（Extended Berkeley Packet Filter）学习项目集合，包含多个 eBPF 示例和实验项目。

## 项目简介

eBPF 是 Linux 内核中的一个强大功能，允许在内核空间中运行沙箱程序，而无需修改内核源代码或加载内核模块。本项目旨在通过实践学习 eBPF 技术，包含从基础到进阶的各种示例。

## 项目结构

```
ebpf-study/
├── hello-ebpf/          # eBPF 入门示例：监控系统调用
│   ├── cmd/             # Go 用户空间程序
│   ├── ebpf/            # eBPF 内核程序源码
│   ├── go.mod           # Go 模块配置
│   ├── justfile         # 构建脚本
│   └── README.md        # 项目详细文档
└── LICENSE              # Apache 2.0 许可证
```

## 子项目

### hello-ebpf

一个简单的 eBPF 入门示例，演示如何使用 tracepoint 监控 Linux 系统调用。

**主要特性：**
- 监控 `sys_enter_write` 系统调用
- 记录触发系统调用的进程 PID
- 通过内核日志输出监控信息

**快速开始：**
```bash
cd hello-ebpf
just run
```

详细文档请参考 [hello-ebpf/README.md](./hello-ebpf/README.md)

## 前置要求

### 系统要求

- Linux 内核版本 >= 5.8（推荐 >= 5.15）
- 支持 eBPF 的 Linux 发行版（Ubuntu 20.04+, Debian 11+, Fedora 33+, 等）

### 开发工具

- **Go**: 1.24 或更高版本
- **Clang**: 用于编译 eBPF 程序
- **Just**: 构建工具（可选，但推荐）
- **bpftool**: eBPF 工具集（可选，用于调试）

### 安装依赖

#### Ubuntu/Debian

```bash
sudo apt-get update
sudo apt-get install -y \
    clang \
    llvm \
    libbpf-dev \
    golang-go \
    just
```

#### Fedora/RHEL

```bash
sudo dnf install -y \
    clang \
    llvm \
    libbpf-devel \
    golang \
    just
```

#### macOS (使用 Homebrew)

```bash
brew install clang llvm go just
```

**注意**: eBPF 程序只能在 Linux 系统上运行，macOS 只能进行编译，无法运行。

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd ebpf-study
```

### 2. 运行示例

```bash
cd hello-ebpf
just run
```

### 3. 查看日志

在另一个终端中查看内核日志：

```bash
sudo cat /sys/kernel/debug/tracing/trace_pipe
```

## 学习路径

1. **hello-ebpf**: 从最简单的 tracepoint 示例开始，了解 eBPF 的基本概念和工作流程
2. **更多示例**: （待添加）逐步学习更复杂的 eBPF 应用场景

## 常见问题

### Q: 为什么需要 root 权限？

A: eBPF 程序需要加载到内核空间，这需要 root 权限。这是 Linux 内核的安全机制。

### Q: 可以在 macOS 上运行吗？

A: 不可以。eBPF 是 Linux 内核的特性，只能在 Linux 系统上运行。macOS 上只能编译代码，无法实际运行 eBPF 程序。

### Q: 如何调试 eBPF 程序？

A: 
- 使用 `bpf_printk` 输出调试信息到内核日志
- 使用 `bpftool` 查看加载的程序和 maps
- 使用 `perf` 工具进行性能分析

### Q: 编译错误：找不到 vmlinux.h

A: 可以使用 `bpftool` 生成：
```bash
bpftool btf dump file /sys/kernel/btf/vmlinux format c > vmlinux.h
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

本项目采用 [Apache 2.0 许可证](./LICENSE)。

## 参考资料

- [eBPF 官方文档](https://ebpf.io/)
- [Cilium eBPF 库文档](https://pkg.go.dev/github.com/cilium/ebpf)
- [Linux 内核 eBPF 文档](https://www.kernel.org/doc/html/latest/bpf/)
- [eBPF 和 XDP 参考指南](https://github.com/xdp-project/xdp-tutorial)

## 相关资源

- [bcc 工具集](https://github.com/iovisor/bcc)
- [bpftrace](https://github.com/iovisor/bpftrace)
- [eBPF Summit](https://ebpf.io/summit/)

