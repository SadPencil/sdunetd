# Update log of sdunetd

## [v2.4.0](https://github.com/SadPencil/sdunetd/releases/tag/v2.4.0)
- The network section is re-added in the configuration file.
The strict mode is now re-supported, behaves like curl, but it is now a Linux-specific feature. See [this page](https://stackoverflow.com/a/73295452/7774607) for technical details.
- Breaking change: the configuration file is now updated. Please re-generate the config file.
- Breaking change: the version parameter is changed to `-V` instead of `-v`.
- Breaking change: `-v` now stands for verbose output.
- Breaking change: the log is now by default written to stderr, instead of stdout.
- Retries are now supported, specified in the configuration file. 

## [v2.3.1](https://github.com/SadPencil/sdunetd/releases/tag/v2.3.1)
- The network section is removed in the configuration file
- Support detecting the Internet via either the auth server or the online service. In the configuration file, set `online_detection_method` to `auth` for the auth server, or `ms` for the detection url by Microsoft

## [v2.2.2](https://github.com/SadPencil/sdunetd/releases/tag/v2.2.2)

The 2.2.2 version. Suitable for Shandong University, Qingdao Campus, since March 2019.

Update log:

- add `-l` flag to logout
- add some tips on the guide

## [v2.2](https://github.com/SadPencil/sdunetd/releases/tag/v2.2)

The 2.2 version. Suitable for Shandong University, Qingdao Campus, since March 2019.

Update log:

- add jsonp callback for disgusting

## [v2.1](https://github.com/SadPencil/sdunetd/releases/tag/v2.1)

The 2.1 version. Suitable for Shandong University, Qingdao Campus, since March 2019.

Update log:

- `-a` flag will only output the IP address
- add `-t` flag
- will check whether the authenticate server reject the login, and output the given error message if so.
- the return value when `-a`, `-t`, `-f` is enabled is meaningful. Returning a zero value means there is no error occurs, and vise versa.

## [v2.0](https://github.com/SadPencil/sdunetd/releases/tag/v2.0)

The 2.0 version. Suitable for Shandong University, Qingdao Campus, since March 2019.

## [v1.1](https://github.com/SadPencil/sdunetd/releases/tag/v1.1)

used to be suitable for SDU-Qingdao, until March 2019
