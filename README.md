# securerm-go
一个简单的通过三次覆盖删除文件的工具，Go 实现

# 三次覆盖
第一次使用随机数覆盖，随机数来源于操作系统，如 Linux、FreeBSD 使用 `getrandom(2)`，Windows 使用 ProcessPrng API 等

第二次使用全1覆盖（每个字节都是 0xFF）

第三次使用全0覆盖（每个字节都是 0x00）

# 注意
本软件不是 srm（secure-delete）的替代品

# 用法
```shell
srmgo [flags] [file] [file2] [file3]...
```

## flags
* -h, --help: print help menu
* -r, --recursion: recursion remove
* -v, --verbose: print verbose