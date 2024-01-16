<h3 align="center">Simple VT-compatible Linux Terminal Emulator</h3>
<div align="center">

[![Licence](https://img.shields.io/badge/Licence-Apache-brightgreen)](https://github.com/veops/gin-api-template/blob/main/LICENSE)
[![Golang](https://img.shields.io/badge/go-1.18+-blue)](https://go.dev/dl/)
</div>

------------------------------
[English](README.md)

## 介绍

**`go-ansiterm`** 是一款基于类似pyte的linux终端仿真器，专为Go语言设计。它在保留pyte强大功能的基础上，更进一步地适应了Go生态系统


## 适用场景

在跳板机和其他需要严格管理用户终端执行命令的场景中，需要在用户输入执行进行执行之前对命令进行过滤，基于golang的开发技术栈没有找到相关的开源库，在某些开源的跳板机中也并不支持命令的提取，因此，本工具主要是从实际需求出发，解决终端用户执行命令的提取。


## 核心功能

- **`命令提取`** 本项目的核心功能在于通过开发通过用户终端的输出和回显等信息提取有效信息，主要为用户的命令的提取，该功能在跳板机中将非常有用，可以非常便捷的提取出用户执行的命令，从而对命令进行有效筛选和过滤。
- **`屏幕模拟`**  包含一个屏幕模拟器，可以处理屏幕上的字符流，支持诸如光标移动、文本滚动等操作。
- **`ANSI转义序列`** 处理各种ANSI转义序列以增强终端交互



## 下载
```shell
    go get github.com/veops/go-ansiterm
```

### 使用
```shell
# 首先创建一个屏幕对象
screen := NewScreen(80, 24)

# 创建一个流对象
stream := InitByteStream(screen, false)

# 将流对象连接到屏幕
stream.Attach(screen)

# 将流对象连接到屏幕
stream.Feed(input)

# 获取屏幕输出
output := screen.Display()

# 重置屏幕
screen.Reset()
```