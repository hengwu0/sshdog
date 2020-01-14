# SSHDog

SSHDog is your go-anywhere lightweight SSH server.  Written in Go, it aims
to be a portable SSH server that you can drop on a system and use for remote
access without any additional configuration.

Useful for:

* Tech support
* Backup SSHD
* Authenticated remote bind shells

Supported features:

* Windows & Linux
* Configure port, host key, authorized keys
* Pubkey, passwords authentication
* Port forwarding
* SCP (but no SFTP support)

如果希望在单板环境运行，最好在go目录执行：
find -name "*.go" |xargs sed -i 's|/dev/random|/dev/urandom|'
rm -rf ./pkg/linux_amd64 ./pkg/linux_amd64_race
以免长时间不响应。

Example usage:

```
% go build
% ssh-keygen -t rsa -b 2048 -N '' -f config/ssh_host_rsa_key
% echo 2222 > config/port
% cp ~/.ssh/id_rsa.pub config/authorized_keys
% ./sshd
[DEBUG] Adding hostkey file: ssh_host_rsa_key
[DEBUG] Adding authorized_keys.
[DEBUG] Listening on :2222
[DEBUG] Waiting for shutdown.
[DEBUG] select...
```


Test ok for WinSCP v5.9.6 (buildin 7601)

Author: hengwu0 <wu.heng@zte.com.cn>
Author: David Tomaschik <dwt@google.com>

*This is not a Google product, merely code that happens to be owned by Google.*



