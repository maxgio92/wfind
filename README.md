[![Latest release](https://img.shields.io/github/v/release/maxgio92/wfind?style=for-the-badge)](https://github.com/maxgio92/wfind/releases/latest)
[![License](https://img.shields.io/github/license/maxgio92/wfind?style=for-the-badge)](COPYING)
![Go version](https://img.shields.io/github/go-mod/go-version/maxgio92/wfind?style=for-the-badge)

# ![](./logo.svg) like find but for web sites

`wfind` (world wide web find) search for files in a web site directory hierarchy over HTTP, through HTML references.

The tool is inspired by GNU `find(1)` and `wget(1)`.

### Usage

```
wfind URL [flags]
```

For details please read the CLI [documentation](./docs/wfind.md).

### In action

```shell
$ wfind https://mirrors.edge.kernel.org/debian/dists/ -t f -n Release
https://mirrors.edge.kernel.org/debian/dists/bullseye/Release
https://mirrors.edge.kernel.org/debian/dists/buster/Release
https://mirrors.edge.kernel.org/debian/dists/sid/Release
https://mirrors.edge.kernel.org/debian/dists/stretch/Release
...
```
