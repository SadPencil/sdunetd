# sdunetd
Embedded SRUN3000 Portal Client for SDU-Qingdao
_________
Ver 1.1 is suitable for Shandong University, Qingdao Campus, up to 2018.

Ver 2.0-RC is probably suitable for Shandong University, Qingdao Campus, since March 2019.

## Copyright
Copyright © 2018-2019 Sad Pencil &lt;me@pencil.live&gt;

MIT License

## Get the executable

Static builds are available [here](https://github.com/SadPencil/sdunetd/releases).

You can also compile by your self.

```bash
go get -v github.com/SadPencil/sdunetd
go build github.com/SadPencil/sdunetd
```

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
ExecStart=/usr/local/bin/sdunetd -config /etc/sdunetd/config.json
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

  - Note: It might take a few minutes to copy a large file if `squashfs` is enabled on the router.


2. `chmod +x /usr/local/bin/sdunetd`
3. Create a configuration file at `/etc/sdunetd/config.json`
4. `touch /etc/init.d/sdunetd`
5. `chmod +x /etc/init.d/sdunetd`
6. `vi /etc/init.d/sdunetd`

```shell
#!/bin/sh /etc/rc.common

START=60
 
start() { 
/usr/local/bin/sdunetd -config /etc/sdunetd/config.json >/dev/null 2>&1 &
}

stop() { 
killall sdunetd
}
```

7. `/etc/init.d/sdunetd enable`
8. `/etc/init.d/sdunetd start`

## Installation on Windows
Although it is okay to create a shortcut at `startup` folder, it is better to create a service. `srvany` is a 32-bit program provided by Microsoft to let any program become a service, and you can get a 64-bit implementation at repo [birkett/srvany-ng](https://github.com/birkett/srvany-ng.git).
