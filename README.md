# yaq

A terminal game engine for learning programming, inspired by [spx](https://github.com/goplus/spx).

## Features

- Develop games with only config files
- Play games with users' inputted code

## Install

```shell
go install github.com/zrcoder/yaq/cmd/yaq@latest
```

## Usage

The first thing is defining your game with toml files.
 ```text
mygame-dir
 ├── Scene1
 │   └── index.toml
 ├── Scene2
 │   ├── Level1.toml
 │   ├── Level2.toml
 │   └── index.toml
 └── index.toml
 ```

Run the game

```shell
cd mygame-dir
yaq
```

## Games powered by yaq

- [yaqs](https://github.com/zrcoder/yaqs)
