## wfind

Find folders and files in web sites using HTTP or HTTPS

```
wfind URL [flags]
```

### Options

```
      --async                               Whether to scrape with asynchronous jobs. (default true)
      --connection-pool-size int            The maximum number of idle connections across all hosts. (default 1000)
      --connection-pool-size-per-host int   The maximum number of idle connections across for each host. (default 1000)
      --connection-timeout int              The maximum amount of time in milliseconds a dial will wait for a connect to complete. (default 180000)
  -h, --help                                help for wfind
      --idle-connection-timeout int         The maximum amount of time in milliseconds a connection will remain idle before closing itself. (default 120000)
      --keep-alive-interval int             The interval between keep-alive probes for an active network connection. (default 30000)
      --max-body-size int                   The maximum size in bytes a response body is read for each request. (default 524288)
  -n, --name string                         Base of file name (the path with the leading directories removed) exact pattern. (default ".+")
  -r, --recursive                           Whether to examine entries recursing into directories. Disable to behave like GNU find -maxdepth=0 option. (default true)
      --tls-handshake-timeout int           The maximum amount of time in milliseconds a connection will wait for a TLS handshake. (default 30000)
  -t, --type string                         The file type
  -v, --verbose                             Enable verbosity to log all visited HTTP(s) files
```

