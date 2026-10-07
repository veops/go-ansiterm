<h3 align="center">Simple VT-compatible Linux Terminal Emulator</h3>
<div align="center">

[![Licence](https://img.shields.io/badge/Licence-Apache-brightgreen)](https://github.com/veops/gin-api-template/blob/main/LICENSE)
[![Golang](https://img.shields.io/badge/go-1.18+-blue)](https://go.dev/dl/)
</div>

------------------------------

[中文](README_zh.md)


## Introduction
**`go-ansiterm`**  is a Linux terminal emulator similar to pyte, specifically designed for the Go language. While retaining the powerful features of pyte, it further adapts to the Go ecosystem.

## Use Cases
In scenarios involving jump servers and other situations that require strict management of user terminal command execution, it's necessary to filter commands before they are executed by the user. In the development stack based on Golang, we haven't found any related open-source libraries. Additionally, some open-source jump server projects do not support command extraction. Therefore, this tool is mainly developed based on practical needs, to solve the problem of extracting commands executed by terminal users.

## Core Features

- **Command Extraction**: The core functionality lies in extracting effective information from terminal outputs and echoes, mainly focusing on user command extraction. This feature is extremely useful in jump servers for the convenient extraction and effective filtering of executed commands.

- **Screen Simulation**: Includes a screen simulator capable of processing character streams on the screen, supporting operations like cursor movement and text scrolling.

- **ANSI Escape Sequences**: Handling various ANSI escape sequences for enhanced terminal interactions.

## How to Require
```shell
    go get github.com/veops/go-ansiterm
```

## Usage Guide

```shell
# First, create a screen object
screen := NewScreen(80, 24)

# Create a stream object
stream := InitByteStream(screen, false)

# Connect the stream object to the screen
stream.Attach(screen)
defer stream.Close()

# Feed input into the stream object
stream.Feed(input)

# output
output := screen.Display()

```

## Stream lifecycle and command checks

Call `Close()` when a session ends, including when it disconnects during an
escape sequence. `Close()` waits for the parser worker to exit and is idempotent.
`Feed`, `Attach`, `InitializeParser` and `Close` are serialized. Direct screen
access must be serialized by the caller with these operations.

The screen represents terminal output, including edits echoed by the remote
shell. This library does not execute commands or enforce authorization. A
bastion must check input before forwarding the submission to the target, use
the echoed screen to reconcile command edits, and maintain its own permissions,
command policy and audit trail. Screen output alone cannot safely authorize
multiline paste, shell scripts or commands that have already been forwarded.
