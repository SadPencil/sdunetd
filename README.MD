# sdunetd
Embedded SRUN3000 Portal Client for SDU-Qingdao
_________

# 适用院校
如有未列出的适用院校、科研单位或其他单位，请通过 issue 补充
- 山东大学（济南无线网、青岛有线网、青岛无线网）
- 中科院计算所（有线网、无线网）

Ver 1.1 is once suitable for Shandong University, Qingdao Campus, up to 2018. No longer supported.

Ver 2.0+ is suitable for Shandong University, Qingdao Campus, since March 2019.

## Copyright
Copyright © 2018-2019 Sad Pencil &lt;me@pencil.live&gt;

MIT License

## Get the executable

Static builds are available [here](https://github.com/SadPencil/sdunetd/releases).

You can also compile by your self.

Rename the executable to `sdunetd`.

## Generate a configuration file

Run the program without any parameters and it will guide you to create a configuration file.

```bash
./sdunetd
```

## Installation on Linux (based on systemd)

1. Copy the executable to `/usr/local/bin`, and rename it to `sdunetd`
2. `chmod +x /usr/local/bin/sdunetd`
3. Create a configuration file at `/etc/sdunetd/config.json`
4. `vi /etc/systemd/system/sdunetd.service`

```ini
[Unit]
Description=sdunetd
After=network.target
Wants=network.target

[Service]
Type=simple
PrivateTmp=true
ExecStart=/usr/local/bin/sdunetd -c /etc/sdunetd/config.json -m
Restart=always

[Install]
WantedBy=multi-user.target
```
6. `systemctl daemon-reload`
7. `systemctl enable sdunetd`
8. `systemctl start sdunetd`

## Installation on OpenWRT

1. Copy the executable to `/usr/local/bin`, and rename it to `sdunetd`

  - Note: You MUST choose proper builds according to `/proc/cpuinfo`.

  - Note: It might take a few minutes to copy a large file to the router.


2. `chmod +x /usr/local/bin/sdunetd`
3. Create a configuration file at `/etc/sdunetd/config.json`
4. `touch /etc/init.d/sdunetd`
5. `chmod +x /etc/init.d/sdunetd`
6. `vi /etc/init.d/sdunetd`

```shell
#!/bin/sh /etc/rc.common

START=60
 
start() { 
/usr/local/bin/sdunetd -c /etc/sdunetd/config.json >/dev/null 2>&1 &
}

stop() { 
killall sdunetd
}
```

7. `/etc/init.d/sdunetd enable`
8. `/etc/init.d/sdunetd start`

## Installation on Windows
Although it is okay to create a shortcut at `startup` folder, it is better to create a service. `srvany` is a 32-bit program provided by Microsoft to let any program become a service, and you can get a 64-bit implementation at repo [birkett/srvany-ng](https://github.com/birkett/srvany-ng.git).

Example:

0. Suppose `sdunetd.exe` `config.json` and `srvany.exe` are all placed at `C:\Program Files\sdunetd`

1. Create a service named `sdunetd`

```winbatch
sc create "sdunetd" start= auto binPath= "C:\Program Files\sdunetd\srvany.exe"
```

2. Import the following to the registry

```ini
Windows Registry Editor Version 5.00
[HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\sdunetd\Parameters]
"Application"="C:\\Program Files\\sdunetd\\srvany.exe"
"AppDirectory"="C:\\Program Files\\sdunetd"
"AppParameters"="-c \"C:\\Program Files\\sdunetd\\config.json\""
```


## Dynamic DNS
We recommend [TimothyYe/GoDNS](https://github.com/TimothyYe/godns). In the configuration file, set `ip_interface` to your network interface to help GoDNS get the real IPv4 address. Click [here](https://github.com/TimothyYe/godns#get-an-ip-address-from-the-interface) to get detailed help.

However, you can't use GoDNS behind a NAT router because the Internet traffic at SDU-Qingdao is being masqueraded, so that online services can't determine your real IPv4 address.

`sdunetd` is able to detect your real IPv4 address at SDU-Qingdao no matter you are under a router or not. So, if you do need this feature, open an issue at [sdunetd](https://github.com/SadPencil/sdunetd/issues) so that we can fork a special version of GoDNS which is suitable for SDU-Qingdao.

## How to compile sdunetd

Go 1.13 or higher version is **required**.

This project uses **go module**. If you live in Mainland China, you might need to configure a proxy like [goproxy.cn](https://github.com/goproxy/goproxy.cn) to execute the following code.

```bash
git clone https://github.com/SadPencil/sdunetd
cd sdunetd
# go get -v # the dependencies will be automatically downloaded with `go build` as go module is enabled
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags=all="-s -w" -o build/sdunetd-linux-arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags=all="-s -w" -o build/sdunetd-linux-arm64
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags=all="-s -w" -o build/sdunetd-linux-386
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=all="-s -w" -o build/sdunetd-linux-amd64
CGO_ENABLED=0 GOMIPS=softfloat GOOS=linux GOARCH=mips go build -ldflags=all="-s -w" -o build/sdunetd-linux-mips-softfloat
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags=all="-s -w" -o build/sdunetd-linux-mips64
CGO_ENABLED=0 GOMIPS=softfloat GOOS=linux GOARCH=mipsle go build -ldflags=all="-s -w" -o build/sdunetd-linux-mipsle-softfloat
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags=all="-s -w" -o build/sdunetd-linux-mips64le
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags=all="-s -w" -o build/sdunetd-linux-arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags=all="-s -w" -o build/sdunetd-linux-arm64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags=all="-s -w" -o build/sdunetd-windows-amd64.exe
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags=all="-s -w" -o build/sdunetd-windows-386.exe
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags=all="-s -w" -o build/sdunetd-darwin-amd64
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags=all="-s -w" -o build/sdunetd-freebsd-amd64
upx --best --ultra-brute build/*
```

## the cURL version of sdunetd
For unknown reasons, the Linux version's “strict mode" might not work properly when multiple default route through different network interfaces are set, while curl works well in this situation. Thus, we provided an additional cURL version.

Note, this version is not build for Windows.
### How to compile the cURL version
Add `-tags="curl"` to the end of `go build`. 

Example: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=all="-s -w" -o build/sdunetd-curl-linux-amd64 -tags="curl"`
