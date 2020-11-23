# gocheck
Distributed monitoring platform written in Go.

## Initial

This is the agent that runs on the remote or local hosts.

To test:
1. `docker-compose build`
2. `docker-compose run --rm app`
3. Once in the shell, you can run:
    ```bash
    ./gocheck run --config_file /app/etc/testFile.hcl
    ```
    and the output should look something like:
    ```bash
    root@1fb51aa30b85:/app# ./gocheck run --config_file /app/etc/testFile.hcl
    2020/11/23 18:23:41 disk.webservers_1 - HEALTHY: "DISK OK - free space: / 157784 MB (36% inode=95%);| /=278652MB;449865;359865;0;459865\n"
    2020/11/23 18:23:41 disk.webservers_2 - WARNING: "DISK WARNING - free space: / 157784 MB (36% inode=95%);| /=278652MB;259865;359865;0;459865\n"
    2020/11/23 18:23:41 disk.webservers_3 - CRITICAL: "DISK CRITICAL - free space: / 157784 MB (36% inode=95%);| /=278652MB;359865;259865;0;459865\n"
    2020/11/23 18:23:41 disk.webservers_4 - No valid check found.
    ```


The included config looks like:

```hcl
// Test healthy
service "disk" "webservers_1" {
    apply_to = ["apache", "nginx"]
    check {
        name = "nagios.check_disk"
        args = [
            "-w", "10000",
            "-c", "100000",
            "-p", "/tmp"
        ]
    }
}

// Test warning
service "disk" "webservers_2" {
    apply_to = ["apache", "nginx"]
    check {
        name = "nagios.check_disk"
        args = [
            "-w", "200000",
            "-c", "100000",
            "-p", "/tmp"
        ]
    }
}
// Test critical

service "disk" "webservers_3" {
    apply_to = ["apache", "nginx"]
    check {
        name = "nagios.check_disk"
        args = [
            "-w", "100000",
            "-c", "200000",
            "-p", "/tmp"
        ]
    }
}

// Test invalid check name
service "disk" "webservers_4" {
    apply_to = ["apache", "nginx"]
    check {
        name = "nagios.nope"
        args = [
            "-w", "100000",
            "-c", "200000",
            "-p", "/tmp"
        ]
    }
}

check "nagios" "check_disk" {
    executable = "/usr/lib/nagios/plugins/check_disk"
}

```