consul {
    address = "127.0.0.1:8500"
    retry {
        enabled = true
        attempts = 4
        backoff = "200ms"
        max_backoff = "1m"
    }
}

 wait {
    min = "3s"
    max = "10s"
}

# exec command to manage lifecycle of application operator
exec {
    command = "./appoperator"
    splay = "5s"
    reload_signal = "SIGUSR2"
    kill_signal = "SIGINT"
}

template {
    source = "./template/file/template1.ctmpl"
    destination = "./template1.yaml"
    create_dest_dirs = true
    error_on_missing_key = true
    perms = 0600
    # prevent file to be rendered very fast
    # at the very least, cannot update template less than 5s
    # PLEASE DO NOT REMOVE
    wait {
        min = "5s"
        max = "20s"
    }
}