service "disk" "webservers" {
    apply_to = ["apache", "nginx"]
    check {
        name = "nagios.check_disk"
        args = [
            "-w", "10000",
            "-c", "200000",
            "-p", "/tmp"
        ]
    }
}

check "nagios" "check_disk" {
    runner = "nagios.debain"
}

// Runner directory for Debian/Ubuntu flavor of nagios
runner "nagios" "debian" {
    directory = "/usr/lib/nagios/plugins"
}

// Runner directory for RHEL/CentOS flavor of nagios
// Any checks using this runner will require the program
// to be passed in (ex: check_disk)
runner "nagios" "centos" {
    directory = "/usr/lib64/nagios/plugins"
}

// Runner directory and executable name for system-wide
// python3
runner "python3" "linux_system" {
    directory = "/usr/bin"
    executable = "python3"
}