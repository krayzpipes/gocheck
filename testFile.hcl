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
