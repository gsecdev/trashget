## trashget

a simple tool which serves virtual big files on a http server for intrusion detection purposes.

### Installation
install go on your operating system. eg. on debian `sudo apt install golang`.

then just run `go install github.com/gsecdev/trashget`

### Usage
```
trashget [OPTIONS]

Application Options:
  -p, --port=     port to listen at
  -i, --ip=       IP to listen at (defaults to all IPs)
  -f, --filename= filename to serve (default: full_backup.zip)
  -s, --size=     virtual size to server (in MB)
  -u, --uri=      URI to serve at (default: /)

Help Options:
  -h, --help      Show this help message
```

### ToDo
- [ ] throttling
- [ ] option to abort the download after a specified transmission length
- [ ] serving of pseudo random data
- [ ] the simulation of valid file structures and headers, e.g. zip files