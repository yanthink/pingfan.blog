## 运行环境

- Go 1.20+
- MySQL 8.0+
- Redis 5.0+
- Zincsearch 0.4.9+

## 开发环境部署/安装

#### 1. 安装模块依赖

```shell
go mod tidy
```

#### 2. 配置

根据情况修改 [config.toml](config.toml) 配置文件。
可通过 go run key_generate.go 生成 key。

#### 3. 数据库迁移

```shell
go run db_migrate.go
```

#### 4. 安装 zincsearch

https://zincsearch-docs.zinc.dev/installation/

运行 zincsearch 服务并修改 config.toml 的 zincsearch 配置

初始化 zincsearch
```shell
go run zincsearch_init.go
```

> 如果提示 fatal error: stdatomic.h: No such file or directory
> #include <stdatomic.h> 错误，检查一下 gcc 版本过低，尝试更新下 gcc 版本再重试。

**更新 gcc**
```shell
yum install centos-release-scl
yum install devtoolset-8
```

#### 5. 运行服务

```shell
go run main.go
```

## FAQ

### fatal error: stdatomic.h: No such file or directory

检查一下 gcc 版本是否过低，尝试更新下 gcc 版本再重试。

centos7 更新 gcc 版本

#### 1. 安装 Developer Toolset：

```shell
sudo yum install centos-release-scl
sudo yum install devtoolset-9
```

#### 2. 启用 Developer Toolset：

```shell
scl enable devtoolset-9 bash
```