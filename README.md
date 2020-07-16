# any-get

一个方便从服务器上下载文件的小工具

本工具监听一个http上传请求,远端机器需要传送文件的时使用curl post文件到该服务上.
该服务会将文件压缩成zip,并生成下载地址便于下载

## Usage

```shell
启动

./any-get -baseUrl=http://127.0.0.1:8080
```

```shell
单个文件

curl -F test=@$(pwd)/file  http://127.0.0.1:8080

多个文件

curl -F test=@$(pwd)/file1 -F test=@$(pwd)/file2   http://127.0.0.1:8080
 ```