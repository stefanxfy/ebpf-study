# hello-ebpf

一个简单的 eBPF 示例项目，演示如何使用 eBPF tracepoint 监控系统调用。

## 功能

本项目实现了一个基础的 eBPF 程序，用于监控 Linux 系统中的 `sys_enter_write` 系统调用。当任何进程执行 `write` 系统调用时，eBPF 程序会记录触发该调用的进程 PID，并通过内核日志输出。

**主要功能：**
- 监控 `sys_enter_write` tracepoint
- 记录触发系统调用的进程 PID
- 通过内核日志输出监控信息

## 代码架构

```
hello-ebpf/
├── cmd/
│   └── main.go          # Go 用户空间程序，负责加载和附加 eBPF 程序
├── ebpf/
│   ├── hello.ebpf.c     # eBPF 内核程序源码
│   ├── vmlinux.h        # Linux 内核类型定义
│   ├── vmlinux_flavors.h
│   ├── vmlinux_missing.h
│   └── CMakeLists.txt   # CMake 构建配置
├── target/              # 构建输出目录（自动创建）
│   ├── hello.o          # 编译后的 eBPF 对象文件
│   └── hello            # 编译后的 Go 可执行文件
├── go.mod               # Go 模块依赖配置
├── justfile             # Just 构建脚本
└── README.md            # 项目文档
```

### 核心组件说明

#### 1. eBPF 内核程序 (`ebpf/hello.ebpf.c`)

这是一个简单的 tracepoint 程序，附加到 `sys_enter_write` tracepoint：

```c
SEC("tp/syscalls/sys_enter_write")
int handle_tp(void *ctx)
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    bpf_printk("BPF triggered sys_enter_write from PID %d.\n", pid);
    return 0;
}
```

**功能：**
- 监听 `syscalls/sys_enter_write` tracepoint
- 获取当前进程的 PID
- 使用 `bpf_printk` 输出日志到内核缓冲区

#### 2. Go 用户空间程序 (`cmd/main.go`)

负责加载 eBPF 程序并附加到 tracepoint：

**主要流程：**
1. 移除内存锁定限制（`rlimit.RemoveMemlock()`）
2. 加载编译好的 eBPF 对象文件（`hello.o`）
3. 创建 eBPF collection
4. 获取 `handle_tp` 程序
5. 附加到 `sys_enter_write` tracepoint
6. 等待中断信号，优雅退出

## 构建流程

本项目使用 [Just](https://github.com/casey/just) 作为构建工具。

### 前置要求

- Linux 内核版本 >= 5.8（支持 eBPF tracepoint）
- Go 1.24 或更高版本
- Clang（用于编译 eBPF 程序）
- Just（构建工具）
- sudo 权限（运行程序需要）

### 构建步骤

#### 1. 编译 eBPF 程序

```bash
just ebpf
```

该命令会：
- 自动创建 `target` 目录（如果不存在）
- 根据系统架构（x86_64 或 aarch64）设置编译标志
- 使用 Clang 编译 `hello.ebpf.c` 生成 `target/hello.o`

**编译参数说明：**
- `-g -O2`: 启用调试信息并优化
- `-mcpu=v2`: 使用 eBPF v2 CPU 特性
- `-D__TARGET_ARCH_x86` 或 `-D__TARGET_ARCH_arm64`: 根据架构设置目标
- `-target bpf`: 指定目标为 BPF

#### 2. 构建 Go 程序

```bash
just build
```

该命令会：
- 自动创建 `target` 目录（如果不存在）
- 运行 `go mod tidy` 下载并整理依赖
- 编译 Go 程序生成 `target/hello` 可执行文件

#### 3. 一键构建并运行

```bash
just run
```

该命令会先执行 `build`，然后以 sudo 权限运行程序。

## 运行方式

### 基本运行

```bash
sudo ./target/hello
```

程序运行后会：
1. 加载 eBPF 程序到内核
2. 附加到 `sys_enter_write` tracepoint
3. 输出 "eBPF program attached successfully to sys_enter_write tracepoint"
4. 持续运行，监控系统调用
5. 按 `Ctrl+C` 优雅退出

### 使用 Just 运行

```bash
just run
```

这会自动构建并运行程序。

## 内核日志查看方式

eBPF 程序使用 `bpf_printk` 输出的日志会写入内核的 trace 缓冲区。可以通过以下方式查看：

### 方法 1: 使用 trace_pipe（实时查看）

```bash
sudo cat /sys/kernel/debug/tracing/trace_pipe
```

这是实时流式输出，会持续显示新的日志条目。

### 方法 2: 使用 trace 文件（查看快照）

```bash
sudo cat /sys/kernel/debug/tracing/trace
```

这会显示当前缓冲区中的所有日志。

### 日志输出示例

当程序运行时，如果某个进程执行了 `write` 系统调用，你会看到类似以下的日志：

```
hello-ebpf-1234  [001] d... 12345.678901: bpf_trace_printk: BPF triggered sys_enter_write from PID 1234.
```

### 清理日志

如果需要清理 trace 缓冲区：

```bash
sudo echo > /sys/kernel/debug/tracing/trace
```

## 依赖说明

### Go 依赖

- `github.com/cilium/ebpf`: eBPF 库，用于加载和管理 eBPF 程序
- `github.com/vishvananda/netlink`: 网络链接管理（当前未使用，但已包含在依赖中）

### 系统依赖

- Linux 内核头文件（通常位于 `/usr/include/linux`）
- eBPF 相关头文件（`bpf/bpf_helpers.h`, `bpf/bpf_tracing.h`）

## 参考资料

- [eBPF 官方文档](https://ebpf.io/)
- [Cilium eBPF 库文档](https://pkg.go.dev/github.com/cilium/ebpf)
- [Linux 内核 tracepoint 文档](https://www.kernel.org/doc/html/latest/trace/tracepoints.html)

