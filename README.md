### 安装 fil-kms

1. Clone the repository:
   ```ssh
      https://github.com/ipfs-force-community/fil-kms.git
      cd fil-kms
   ```
2. Build fil-kms
   ```sh
      make clean
      make build
   ```
   
### 启动进程

> 自定义aks，aks是自定义的一个密钥信息，用于客户远程访问时加密数据 

```sh
   ./fil-kms --ip 0.0.0.0 run --aks xxxxx
```

启动后会在 `~/.filkms` 生成一个目录，里面包含一个配置文件 `config.toml`，需要手动调整。

配置文件格式如下：

```bash
WalletURL = "127.0.0.1:8000"
WalletToken = "token"

[[Filters]]
  Client = "t2i4llai5x72clnz643iydyplvjmni74x4vyme7ny"
  Miner = "t2i4llai5x72clnz643iydyplvjmni74x4vyme7ny"
  Limit = "10Tib"
  Used = 0
  Start = ""
  End = ""

[[Filters]]
  Client = "t2ydigmboo6yjmbyqxyy3eyxifnreamnhdjzybfuy"
  Miner = "t2ydigmboo6yjmbyqxyy3eyxifnreamnhdjzybfuy"
  Limit = "8Gib"
  Used = 0
  Start = ""
  End = ""
```
