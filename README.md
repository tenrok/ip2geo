# ip2geo

A tool with HTTP server to get the geodata by IP.

```
usage: ip2geo [<flags>] <command> [<args> ...]

A tool with HTTP server to get the geodata by IP.

Flags:
  -h, --help     Show context-sensitive help (also try --help-long and --help-man).
  -d, --dbpath="./SxGeoCity.dat"  
                 Path to SxGeoCity.dat file.
  -v, --version  Show application version.

Commands:
  help [<command>...]
    Show help.


  find <ip>
    Find IP address in database.


  serve [<flags>]
    Start HTTP server, listen and serve.

    -a, --addr=":8080"  Address to listen for HTTP requests on.
```
