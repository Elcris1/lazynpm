# Lazynpm

Lazynpm is a tui interface for managing packages and scripts inside a node managed environment.
It allows users to check which scripts exists, what do they relate to, install/uninstall dependencies, look for actual dependencies.., pending updates and executing scripts.

## Dev useful commands

Module init

```bash
go mod init github.com/Elcris1/lazynpm
```

Adding dependencies

```bash
go get -u github/username/pacakge@version
```

Synchronizyng dependencies

```bash
# sync and clean module cache
go mod tidy

# Download explicit packages
go mod download
```

Building the packages

```bash
go buid -o bin/lazynpm .cmd/lazynpm
```

Running the built in project

```bash
./bin/lazynpm 
```

System wide installation of package

```bash
go install ./cmd/lazynpm
```

## Dev Notes

This will be written as we go, learning things

# How to use lazynpm

This wil be written whenever there is something to explain
