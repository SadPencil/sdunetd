# sdunetd

A very simple cross-platform srun3000 portal client.

Ver 1.1 is suitable for Shandong University, Qingdao Campus, up to 2018.

Ver 2.0-RC is probably suitable for Shandong University, Qingdao Campus, since March 2019.

Static builds are available [here](https://github.com/SadPencil/sdunetd/releases).

## Copyright
Copyright © 2018-2019 Sad Pencil &lt;me@pencil.live&gt;

MIT License

## Installation on OpenWRT

1. Copy the executable to `/usr/bin`, and rename it to `sdunetd`

  - Note: You MUST choose proper builds according to `/proc/cpuinfo`.

  - Note: It might take a few minutes to copy a large file.

  - Note: If the static builds are too large for a tiny ROM, you can either use `upx` to compress the executeable, or compiles by yourself with environment variable `CGO_ENABLED=1`, which relys on `libc`.

2. `chmod +x /usr/bin/sdunetd`

3. `touch /etc/init.d/sdunetd`

4. `chmod +x /etc/init.d/sdunetd`

5. `vi /etc/init.d/sdunetd`

```shell
#!/bin/sh /etc/rc.common 

START=60
 
start() { 
/usr/bin/sdunetd -s SERVER_IP -u USERNAME -p PASSWORD -interface INTERFACE -log /tmp/sdunetd.log &
} 

stop() { 
killall sdunetd
}
```

Replace `SERVER_IP` `USERNAME` `PASSWORD` `INTERFACE` with real ones.

`INTERFACE` is the name of your network interface, e.g, `eth0` `eth0.1` `eno1` `Local Area Connection`.

6. `/etc/init.d/sdunetd enable`
7. `/etc/init.d/sdunetd start`

## Installation on Linux (based on systemd)
- Step 1 to 4 is the same with last section.

5. `vi /etc/systemd/system/sdunetd.service`

```ini
[Unit]
Description=sdunetd
After=network.target
Wants=network.target

[Service]
Type=simple
PrivateTmp=true
ExecStart=/usr/bin/sdunetd -s SERVER_IP -u USERNAME -p PASSWORD -interface INTERFACE -log /tmp/sdunetd.log
Restart=always

[Install]
WantedBy=multi-user.target
```
Replace `SERVER_IP` `USERNAME` `PASSWORD` `INTERFACE` with real ones.

`INTERFACE` is the name of your network interface, e.g, `eth0` `eth0.1` `eno1` `Local Area Connection`.

6. `systemctl daemon-reload`
7. `systemctl enable sdunetd`
8. `systemctl start sdunetd`

## Installation on Windows
Although it is okay to create a shortcut at `startup` folder, it is better to create a service. `srvany` is a 32-bit program provided by Microsoft to let any program become a service, and you can get a 64-bit implementation at repo [birkett/srvany-ng](https://github.com/birkett/srvany-ng.git).
