# EasyConnect for Cisco AnyConnect VPN

## Install EasyConnect

```bash
go install github.com/fadinflame/easyconnect@latest
```

## Usage

```bash
easyconnect
```

## Configuration

At the first run, EasyConnect will prompt you to create a configuration file.
Or you can create one manually.
To configure EasyConnect, create a configuration file named `config.json` in the same directory as the executable
The configuration file is located at `~/.easyconnect/config.json` on Linux and `%AppData%/easyconnect/config.json` on Windows.

```json
{
    "server": "Your VPN Server",
    "group": "Your VPN Group",
    "username": "Your VPN Username",
    "password": "Your VPN Password",
    "cisco_logs": false
}
```

## Cisco AnyConnect VPN Logs

You can enable Cisco AnyConnect VPN logs by setting the `cisco_logs` field to `true` in the configuration file.

```json
{
    "server": "Your VPN Server",
    "group": "Your VPN Group",
    "username": "Your VPN Username",
    "password": "Your VPN Password",
    "cisco_logs": true
}
```

## License

This project is released under the [MIT License](https://opensource.org/licenses/MIT).

