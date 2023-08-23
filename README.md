# wfind: like find but for web sites

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
