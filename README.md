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
  Client = "f3tginto55sulyzzjoy3byof7u7t5vicwa372z4hn4zdhjkdwdkhalsn6mytk53paoc63uf575rfdvqetyvpka"
  Miner = "t01000"
  Limit = "10Tib"
  Used = 0

[[Filters]]
  Client = "f3tginto55sulyzzjoy3byof7u7t5vicwa372z4hn4zdhjkdwdkhalsn6mytk53paoc63uf575rfdvqetyvpka"
  Miner = "t01002"
  Limit = "8Gib"
  Used = 0
```

### curl

1. 列出钱包地址

```bash
curl 192.168.200.18:10025/walletList
```

### 签名

```bash
curl 192.168.200.18:10025/sign -X POST -H "Content-Type: application/json" -H "Authorization: xxxxx:615d6db49f07e58dc4bab0358bdd3c6cbeebee62" -d '{"addr":"f3tginto55sulyzzjoy3byof7u7t5vicwa372z4hn4zdhjkdwdkhalsn6mytk53paoc63uf575rfdvqetyvpka","msg":"i9gqWCgAAYHiA5IgIIPglc5o7ABVX2jLTVKHZmhsuiR4zioG4NoyVwxO6joVGgABhqD0WDEDmZDZu72VF4zlLsbDhxf0/PtUCsDf9Z4dvMjOlQ7DUcC5N8zE1d28Dhe3Qvf9iUdYQwDoB2R0ZXN0AABAQEA="}'
```
