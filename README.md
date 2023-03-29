# wfind: like find but for web sites

`wfind` (world wide web find) search for files in a web site directory hierarchy over HTTP and HTTPS, through HTML references.

The tool is inspired by GNU `find(1)` and `wget(1)`.

### Usage

```
wfind URL [flags]
```

#### Options

```
  -h, --help          help for wfind
  -n, --name string   Base of file name (the path with the leading directories removed) pattern.
  -r, --recursive     Whether to examine entries recursing into directories. Disable to behave like GNU find -maxdepth=0 option. (default true)
  -t, --type string   The file type
  -v, --verbose       Enable verbosity to log all visited HTTP(s) files
```

### In action

```shell
$ wfind https://mirrors.edge.kernel.org/debian/dists/ -t f -n Release
https://mirrors.edge.kernel.org/debian/dists/bullseye/Release
https://mirrors.edge.kernel.org/debian/dists/buster/Release
https://mirrors.edge.kernel.org/debian/dists/sid/Release
https://mirrors.edge.kernel.org/debian/dists/stretch/Release
...
```
