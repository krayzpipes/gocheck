// Test healthy
check "disk" "webservers_1" {
    apply_to = ["apache", "nginx"]
    exec {
        path = "/usr/lib/nagios/plugins/check_disk"
        args = [
            "-w", "10000",
            "-c", "100000",
            "-p", "/tmp"
        ]
    }
    cron = "TZ=America/New_York 0/5 * * * ? *"
}

// Test warning
check "disk" "webservers_2" {
    apply_to = ["apache", "nginx"]
    exec {
        path = "/usr/lib/nagios/plugins/check_disk"
        args = [
            "-w", "200000",
            "-c", "100000",
            "-p", "/tmp"
        ]
    }
}

// Test critical
check "disk" "webservers_3" {
    apply_to = ["apache", "nginx"]
    exec {
        path = "/usr/lib/nagios/plugins/check_disk"
        args = [
            "-w", "100000",
            "-c", "200000",
            "-p", "/tmp"
        ]
    }
}

// Test invalid check name
check "disk" "webservers_4" {
    apply_to = ["apache", "nginx"]
    exec {
        path = "/usr/lib/nagios/plugins/check_disk"
        args = [
            "-w", "100000",
            "-c", "200000",
            "-p", "/tmp"
        ]
    }
}

