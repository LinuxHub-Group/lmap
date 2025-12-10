# lmap

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

<h3 align="center">LinuxHub's Nmap - 下一代网络扫描工具</h3>

<p align="center">
  比 nmap 更强大、更灵活的网络探测与扫描能力
  <br/>
  <a href="#使用方法"><strong>探索 lmap 的功能 »</strong></a>
  <br/>
  <br/>
  <a href="https://github.com/LinuxHub-Group/lmap/issues">报告 Bug</a>
  ·
  <a href="https://github.com/LinuxHub-Group/lmap/issues">请求新功能</a>
  ·
  <a href="#贡献">贡献</a>
</p>

## 目录

- [关于项目](#关于项目)
- [功能特性](#功能特性)
- [快速开始](#快速开始)
  - [安装](#安装)
  - [使用方法](#使用方法)
- [构建选项](#构建选项)
- [贡献](#贡献)
- [许可证](#许可证)

## 关于项目

lmap (LinuxHub's Nmap) 是 LinuxHub 团队开发的下一代网络扫描工具，被称为 nmap 的 Pro Plus Max 版本。它是一个现代化的网络扫描工具，旨在提供比传统 nmap 更强大、更灵活的网络探测与扫描能力。

该工具专为网络安全工程师、系统管理员和渗透测试人员设计，帮助他们快速识别网络中的活动主机、开放端口和服务。

### 为什么选择 lmap？

- **快速扫描** - 利用 Go 语言的并发特性，实现高速网络扫描
- **跨平台支持** - 支持 Windows、Linux、macOS 等多种操作系统
- **易于使用** - 简洁的命令行界面，直观的参数设置
- **高度可定制** - 支持灵活的排除规则和扫描配置
- **开源免费** - 基于 GPL-3.0 许可证，完全开源

## 功能特性

- [x] 🚀 **网络扫描** - 快速扫描指定网络段中的活动主机
- [x] ✅ **IP检查** - 验证IP地址有效性并检查主机是否在线
- [x] 🔍 **子网解析** - 解析CIDR格式的子网并生成IP地址列表
- [x] 📡 **ICMP监听** - 监听网络中的ICMP数据包
- [x] 🏓 **Ping探测** - 使用ICMP协议探测主机是否在线
- [x] ⚡ **并发处理** - 支持高并发扫描，提高扫描效率
- [x] 🚫 **排除规则** - 支持排除特定IP或子网，避免扫描不必要目标
- [x] 📋 **详细输出** - 提供详细的扫描过程信息

## 快速开始

### 安装

#### 预编译二进制文件

从 [Releases](https://github.com/LinuxHub-Group/lmap/releases) 页面下载适用于您系统的预编译二进制文件。

#### 从源码构建

**要求**:
- Go 1.16 或更高版本

```bash
# 克隆项目
git clone https://github.com/LinuxHub-Group/lmap.git
cd lmap

# 构建
make

# 或者直接使用Go构建
go build ./cmd/lmap
```

构建完成后，您将在项目根目录下获得 [lmap](file:///D:/works/lmap/lmap/lmap.exe) 可执行文件。

### 使用方法

#### 基本扫描

```bash
# 扫描单个子网
./lmap -subnet 192.168.1.0/24

# 扫描多个子网
./lmap -subnet 192.168.1.0/24 -subnet 10.0.0.0/16

# 详细输出模式
./lmap -subnet 192.168.1.0/24 -v
```

#### 排除特定IP或子网

```bash
# 排除单个IP
./lmap -subnet 192.168.1.0/24 -exclude 192.168.1.1

# 排除多个IP或子网
./lmap -subnet 192.168.1.0/24 -exclude 192.168.1.1 -exclude 192.168.1.10/32
```

#### 命令行选项

| 选项 | 描述 | 示例 |
|------|------|------|
| `-subnet` | 要扫描的网络段，CIDR格式 (可多次指定) | `-subnet 192.168.1.0/24` |
| `-exclude` | 要排除的IP或子网 (可多次指定) | `-exclude 192.168.1.1` |
| `-v` | 详细输出模式 | `-v` |

## 构建选项

```bash
# 默认构建当前平台版本
make

# 构建所有平台版本
make all

# 格式化代码
make fmt

# 运行测试
make test

# 清理构建产物
make clean
```

## 贡献

欢迎任何形式的贡献！如果您想为 lmap 做出贡献，请遵循以下步骤：

1. Fork 项目
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

### 开发环境设置

```bash
# 1. 克隆项目
git clone https://github.com/LinuxHub-Group/lmap.git

# 2. 进入项目目录
cd lmap

# 3. 安装依赖
make install

# 4. 运行测试
make test
```

## 社区和支持

- [报告 Bug](https://github.com/LinuxHub-Group/lmap/issues)
- [请求新功能](https://github.com/LinuxHub-Group/lmap/issues)
- [查看已知问题](https://github.com/LinuxHub-Group/lmap/issues)

## 许可证

```
Copyright (C) <2021>  <LinuxHub-Group>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```