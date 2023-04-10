### 安装 fil-kms

1. Clone the repository:
    ```ssh
     git clone https://github.com/FilBoost/fil-kms.git
     cd fil-kms
    ```
2. Build fil-kms
   ```sh
    make clean
    make build
    ```
   
### 启动进程
1. 自定义aks，aks是自定义的一个密钥信息，用于客户远程访问时加密数据 
2. 启动进程
   ```sh
   fil-kms run --aks $AKS
   ```

### 导入私钥
   ```sh
  fil-kms import
   ```
### 

