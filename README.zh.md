# Discord 登录验证机器人

这个项目是适用于 Discord 的登录验证机器人，用于处理用户的登录请求和验证操作。



## 安装

我们共提供了两种安装方式

#### 1.编译安装

```shell
git clone https://github.com/NagatoException/SimpleDiscordLoginVerifyBot.git
cd SimpleDiscordLoginVerifyBot
go build

./SimpleDiscordLoginVerifyBot
```

#### 2.使用预先编译好的二进制文件

在Release中下载最新的二进制编译文件即可



## 依赖

该机项目依赖以下 Go 包：

- `github.com/bwmarrin/discordgo`：适用于Go的Discord API。
- `github.com/google/uuid`：用于生成 UUID 。
- `github.com/patrickmn/go-cache`：用于缓存登录请求。



## 配置

Release和源码都自带一个config.yaml，修改即可

```yaml
general:
  domain: "" # Your Domain
  server_id: "" # Your Discord ServerID
  ssl_enabled: false
  ssl:
    cert_file: "" # Your SSL Cert File Path
    key_file: "" # Your SSL Key File Path
bot:
  token: "" # Your Discord Bot Token
cache:
  expire: "" # Cache Expire Time (Seconds)
  clean_interval: "300" # Cache Clean Interval (Seconds)

```



## API 使用方法

#### 1. /new

- 方法：GET
- 参数：
  - `username`：要验证的用户名
  - `tag`：要验证的用户的代号

示例请求：

```SQL
/new?username=NagatoException&tag=6898
```

#### 2. `/verify` - 验证先前的登录请求

- 方法：GET
- 参数：
  - `uuid`：先前登录请求的 UUID

示例请求：

```SQL
GET /verify?uuid=00000000-0000-0000-0000-000000000000
```

##### 3. `/login` - 处理登录请求

- 方法：GET
- 参数：
  - `uuid`：登录请求的 UUID

示例请求：

```SQL
GET /login?uuid=00000000-0000-0000-0000-000000000000
```



## 贡献

如果你有好的点子和代码，欢迎发起Issue或pr



## 协议

该项目基于GPL-3.0协议开源

