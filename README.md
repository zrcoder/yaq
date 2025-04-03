# yaq

A terminal game engine for learning programming, inspired by [spx](https://github.com/goplus/spx).

## Features

- Develop games with only config files
- Play games with users' inputted code

## Install

If go < 1.23

```sh
go install github.com/zrcoder/yaq/cmd/yaq@latest
```

Else
```sh
go install -ldflags="-checklinkname=0" github.com/zrcoder/yaq/cmd/yaq@latest
```


## Usage

The first thing is defining your game with yaml files.
 ```text
mygame-dir
 ├── Scene1
 │   └── index.yaml
 ├── Scene2
 │   ├── Level1.yaml
 │   ├── Level2.yaml
 │   └── index.yaml
 └── index.yaml
 ```

Run the game

```shell
cd mygame-dir
yaq
```

## Games powered by yaq

- [yaqs](https://github.com/zrcoder/yaqs)
