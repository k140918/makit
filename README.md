# 🗂️ makit

A CLI tool that allows you to create files and directories with a single command  
![MIT License](https://img.shields.io/badge/license-MIT-blue "MIT License")
[![Go Report Card](https://goreportcard.com/badge/github.com/k140918/makit)](https://goreportcard.com/report/github.com/k140918/makit)
[![Coverage Status](https://coveralls.io/repos/github/k140918/makit/badge.svg?branch=main)](https://coveralls.io/github/k140918/makit?branch=main)

## 👀 Overview

This is a simple yet powerful CLI tool that lets you **create files and directories simultaneously** with a single command.  
No more typing `mkdir` followed by `touch`—this utility streamlines your workflow and boosts productivity, especially for developers and script writers.

## 🥞 Usage

`makit` is a powerful CLI tool that simplifies file and directory creation. Once installed, you can use it from any directory in your terminal.

```
makit [OPTION] <FILES|DIRS...>

OPTION
    -p, --parents         Create parent directories as needed (default; --parents=false to disable)
    -m, --mode <mode>     Set file/directory permissions (e.g., -m 755)
    -d, --date <stamp>    Use specific timestamp (e.g., 202504181200)
    -c, --no-create       Do not create file if it doesn't exist
    -v, --verbose         Enable verbose output
        --version         Print version
    -h, --help            Print this message
```

A path that ends with `/` or has no extension is created as a directory; otherwise it is created as an empty file.
New files are created with mode `644` and new directories with `755` unless `-m` is given.

**Examples:**

```bash
makit my_new_file.txt                      # create a file
makit my_new_directory/                    # create a directory
makit project/src/main.go                  # create a file with its parent directories
makit client/index.html server/app.js      # create multiple paths at once
makit -m 700 private_data/                 # create a directory with permissions 700
makit -d 202506201800 meeting_notes.txt    # create a file with a specific timestamp
makit -c existing_file.txt                 # only update existing paths, never create
```

## 🍈 Installation

To use `makit` from any directory in your terminal, place the `makit` executable in a directory that is included in your system's `PATH` environment variable.

1.  **Build the executable** (in the project root, where `go.mod` is located):

    ```bash
    go build -o makit ./cmd/main          # macOS / Linux
    go build -o makit.exe ./cmd/main      # Windows
    ```

    If you have [just](https://github.com/casey/just), `just build` runs the tests and builds the binary.

2.  **Move the executable to a `PATH` directory:**

    ```bash
    sudo mv makit /usr/local/bin/         # macOS / Linux
    ```

    On Windows, move `makit.exe` to any folder listed in the `Path` environment variable.

3.  **Verify the installation:** open a new terminal and run

    ```bash
    makit --version
    ```

Prebuilt binaries for Linux / macOS / Windows (amd64, arm64) are attached to each [release](https://github.com/k140918/makit/releases).

## 🐼 About

### Developer Name

    K140918

### License

[MIT License](LICENSE)

### Icon

<img src="file_dir.png" alt="ファイルとディレクトリが半分のアイコン" width="200">

### Origin of Name

    The name makit is a blend of:

      make – representing the action of creating files and directories, inspired by traditional UNIX tools like make, mkdir, and touch.

      it – referring to “it”, the thing you want to create. Whether it’s a file, a folder, or an entire structure—just makit.

### Version

0.1.3
