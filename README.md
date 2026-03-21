# jfsh

> **This is a fork of [hacel/jfsh](https://github.com/hacel/jfsh).**
> Changes in this fork:
> - **Media library browsing**: Added a new "Libraries" tab to browse your Jellyfin media libraries directly. (Implemented by minimax-m2.7)
> - **Manual segment skipping**: Press `Ctrl+s` in the mpv window to skip the current segment (intro, outro, etc.), independent of the `skip_segments` auto-skip configuration. (Implemented by opus 4.6)
> - **Media info display**: Shows resolution, codec, HDR info for movies and episodes in the description line. (Implemented by minimax-m2.7)

[中文说明](README_zh.md)

A terminal-based client for [Jellyfin](https://jellyfin.org) that lets you browse your media library and play videos via [mpv](https://mpv.io).
Inspired by [jftui](https://github.com/Aanok/jftui).

![Demo](demo/demo.gif)

## Features

- Uses _your_ mpv config!
- Resumes playback!
- Tracks playback progress and updates jellyfin!
- Automatic and manual segment (intro, etc.) skipping!
- **Library browsing**: Browse your media libraries via the Libraries tab!
- **Media info display**: Shows resolution, codec, HDR info for movies and episodes!
- No mouse required!

## Installation

### Prerequisites

- A running [Jellyfin](https://jellyfin.org) instance.
- [mpv](https://mpv.io) available in PATH.
- [Go](https://go.dev) 1.23 or later.

```sh
git clone https://github.com/xiuusi/jfsh.git
cd jfsh
go build -o jfsh .
```

Move the built binary to a directory in your PATH, e.g.:

```sh
mv jfsh ~/.local/bin/
```

## Usage

1. **Start jfsh**

   ```sh
   jfsh
   ```

2. **Login**

   On first launch, you'll be prompted to enter:

   - **Host**: e.g., `http://localhost:8096`
   - **Username**
   - **Password**

3. **Play Media**

   - Select an item and press **Enter** or **Space** to play it.
   - `mpv` will launch and begin streaming.

4. **Navigate tabs**

   - **Resume**: Continue watching (incomplete media)
   - **Next Up**: Next episode of TV series
   - **Recently Added**: Recently added media
   - **Libraries**: Media library list
   - **Search**: Search

5. **Browse libraries**

   - Switch to the **Libraries** tab to view all media libraries
   - Select a library and press **Enter** to browse its contents
   - Select a specific item (movie, series) and press **Enter** to play
   - Press **Esc** or **Backspace** to go back

6. **Quit**

   - Press **`q`** at any time to exit jfsh.

## Configuration

By default, the configuration file is stored in `$XDG_CONFIG_HOME/jfsh/jfsh.yaml`. If `$XDG_CONFIG_HOME` is not set it defaults to:

- **Linux**: `~/.config/jfsh/jfsh.yaml`
- **macOS**: `~/Library/Application Support/jfsh/jfsh.yaml`
- **Windows**: `%APPDATA%/jfsh/jfsh.yaml`

```yaml
host: http://localhost:8096
username: me
password: hunter2
device: mycomputer # Device name to report to jellyfin (default: hostname)
skip_segments: # Segments to automatically skip (default: [])
  - Recap
  - Preview
  - Intro
  - Outro
```

### Segment skipping

By default, no segments are automatically skipped. To enable skipping segments you must add `skip_segments` to the configuration file. Possible values for `skip_segments` are the segment types in Jellyfin which are: `Unknown`, `Commercial`, `Preview`, `Recap`, `Outro` and `Intro`.

You can also manually skip the current segment by pressing **`Ctrl+s`** in the mpv window. This works for all segment types regardless of `skip_segments` configuration — if the current playback position is inside any segment, it will seek to the end of that segment.

## Plans

- Configuration through TUI
- Sorting
- Better search: Filter by media type, watched status, and metadata
