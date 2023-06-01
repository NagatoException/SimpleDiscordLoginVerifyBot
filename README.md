# Discord Login Verification Bot
[中文文档 - Chinese Document](https://github.com/NagatoException/SimpleDiscordLoginVerifyBot/blob/main/README.zh.md)
This project is a Discord login verification bot designed to handle user login requests and verification operations.



## Installation

There are two installation methods available:

#### 1. Building from source

```shell
git clone https://github.com/NagatoException/SimpleDiscordLoginVerifyBot.git
cd SimpleDiscordLoginVerifyBot
go build

./SimpleDiscordLoginVerifyBot
```

#### 2. Using pre-compiled binary

Download the latest pre-compiled binary from the Releases section.

## Dependencies

This project depends on the following Go packages:

- `github.com/bwmarrin/discordgo`: Discord API for Go.
- `github.com/google/uuid`: UUID generation library.
- `github.com/patrickmn/go-cache`: Cache library for login requests.

## Configuration

The project comes with a `config.yaml` file. Modify the file according to your needs:

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

## API Usage

#### 1. `/new` - Create a new login request

- Method: GET
- Parameters:
  - `username`: Username to be verified
  - `tag`: User's tag to be verified

Example request:

```SQL
GET /new?username=NagatoException&tag=6898
```

#### 2. `/verify` - Verify a previous login request

- Method: GET
- Parameters:
  - `uuid`: UUID of the previous login request

Example request:

```SQL
GET /verify?uuid=00000000-0000-0000-0000-000000000000
```

##### 3. `/login` - Process a login request

- Method: GET
- Parameters:
  - `uuid`: UUID of the login request

Example request:

```SQL
GET /login?uuid=00000000-0000-0000-0000-000000000000
```

## Contributing

If you have any ideas or improvements, feel free to open an issue or submit a pull request.

## License

This project is licensed under the GPL-3.0 License.