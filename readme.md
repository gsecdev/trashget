## trashget

a simple tool which serves virtual big files on a http server for intrusion detection purposes.

### Installation
install go on your operating system. eg. on debian `sudo apt install golang` or follow instructions on `https://go.dev/doc/install`.

then just run `go install github.com/gsecdev/trashget@latest`

you may proxy the server behind you webserver. eg proxying behind an nginx server by using a similar configuration:

```
       location /trash/ful_backup.zip {
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_buffering off;
                proxy_cache off;
                proxy_pass http://localhost:8000;
        }
```

### Usage
```
trashget [OPTIONS]

Application Options:
  -p, --port=       port to listen at (default: 8000)
  -i, --ip=         IP to listen at (defaults to all IPs)
  -f, --filename=   filename to serve (default: full_backup.zip)
  -s, --size=       virtual size of file (in MB) (default: 1000)
  -u, --uri=        URI to serve at (default: /)
  -t, --throttle=   throttle bandwith (in Mbit/s) (default: -1)
  -a, --abortAfter= abort transmission after given % (default: 100)

Help Options:
  -h, --help        Show this help message
```

### ToDo
- [x] throttling
- [x] option to abort the download after a specified transmission length
- [x] serving of pseudo random data
- [ ] the simulation of valid file structures and headers, e.g. zip files
