# jfsh

> **本仓库是 [hacel/jfsh](https://github.com/hacel/jfsh) 的复刻版本。**
> 相比上游的改动：
> - **手动跳过片段**：在 mpv 播放窗口中按 `Ctrl+s` 可跳过当前片段（片头、片尾等），不受 `skip_segments` 自动跳过配置的限制。

[English](README.md)

基于终端的 [Jellyfin](https://jellyfin.org) 客户端，通过 [mpv](https://mpv.io) 浏览媒体库并播放视频。
灵感来自 [jftui](https://github.com/Aanok/jftui)。

![演示](demo/demo.gif)

## 功能特性

- 使用**你自己的** mpv 配置
- 断点续播
- 追踪播放进度并同步到 Jellyfin
- 自动及手动跳过片段（片头、片尾等）
- 纯键盘操作，无需鼠标

## 安装

### 前置条件

- 运行中的 [Jellyfin](https://jellyfin.org) 服务器
- [mpv](https://mpv.io) 已安装且在 PATH 中
- [Go](https://go.dev) 1.23 或更高版本

```sh
git clone https://github.com/xiuusi/jfsh.git
cd jfsh
go build -o jfsh .
```

将编译好的二进制文件移动到 PATH 目录中，例如：

```sh
mv jfsh ~/.local/bin/
```

## 使用方法

1. **启动 jfsh**

   ```sh
   jfsh
   ```

2. **登录**

   首次启动时需要输入：

   - **Host**：Jellyfin 服务器地址，如 `http://localhost:8096`
   - **Username**：用户名
   - **Password**：密码

3. **播放媒体**

   - 选择一个条目，按 **Enter** 或 **Space** 播放。
   - mpv 将启动并开始串流播放。

4. **退出**

   - 按 **`q`** 退出 jfsh。

## 配置

配置文件默认存储在 `$XDG_CONFIG_HOME/jfsh/jfsh.yaml`。若 `$XDG_CONFIG_HOME` 未设置，则默认路径为：

- **Linux**：`~/.config/jfsh/jfsh.yaml`
- **macOS**：`~/Library/Application Support/jfsh/jfsh.yaml`
- **Windows**：`%APPDATA%/jfsh/jfsh.yaml`

```yaml
host: http://localhost:8096
username: me
password: hunter2
device: mycomputer # 上报给 Jellyfin 的设备名称（默认：主机名）
skip_segments: # 自动跳过的片段类型（默认：[]）
  - Recap
  - Preview
  - Intro
  - Outro
```

### 片段跳过

默认不自动跳过任何片段。要启用自动跳过，需在配置文件中添加 `skip_segments`。可选值为 Jellyfin 中的片段类型：`Unknown`、`Commercial`、`Preview`、`Recap`、`Outro` 和 `Intro`。

在 mpv 播放窗口中按 **`Ctrl+s`** 可手动跳过当前片段。手动跳过不受 `skip_segments` 配置限制，适用于所有片段类型——只要当前播放位置处于任意片段范围内，就会跳转到该片段的末尾。

## 计划

- 通过 TUI 进行配置
- 完整的媒体库浏览
- 排序功能
- 更好的搜索：按媒体类型、观看状态和元数据过滤
