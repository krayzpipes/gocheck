
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
    executable = "/usr/lib/nagios/plugins/check_disk"
}
