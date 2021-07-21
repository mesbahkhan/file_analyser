NAME:
   Go File Analysis Suite - gofiles [mode]

USAGE:
   file_analyser.exe [global options] command [command options] [arguments...]

VERSION:
   0.0.2

AUTHOR:
   Mesbah Khan <khanm@ontoledgy.io>

COMMANDS:
   hash, h    use it to create a hashtable for a directory. options: --hashAlgo {sha256, sha512, md5} sets hashing algorithms , --skipFiles {n} skips n files, --batchSize {n} writes n processed files.
   
   copy, c    Use it to copy files using a csv loader with source and destination paths
   
   unzip, u   use it to unzip all zips within a directory. Select recursivity using -recursive yes or no
   
   report, a  Use get an anlysis report on folder
   
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

COPYRIGHT:
   copyright 2020
