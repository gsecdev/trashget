# trashget

a tool which serves virtual big files on a http server for intrusion detection purposes.

'''
Usage:
  trashget [OPTIONS]

Application Options:
  -p, --port=     port to listen at
  -i, --ip=       IP to listen at
  -f, --filename= filename to serve (default: full_backup.zip)
  -s, --size=     virtual size to server (in MB)
  -u, --uri=      URI to serve at (default: /)

Help Options:
  -h, --help      Show this help message
'''

## ToDo
- implement throttling
- implement option to abort the download after a specified transmission length
- implement serving of pseudo random data
- implement the simulation of valid file structures and headers, e.g. zip files