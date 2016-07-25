# Fetch
Fetch is what I use to retreive files from the internet (I find other tools like wget and curl too complicated; besides, there's a coolness factor in using your own tools to do stuff).

I started out with [The Go Programming Language's](//gopl.io) 'fetch' program and branched out from there, adding shortcuts that I found useful (including the '-x' flag) as well as other common features.

## Usage
```
usage: fetch [flags] [urls]
flags:
    -d  download contents to a file instead of printing to stdout
    -e  ignore errors and fetch page on status code != 200
    -h  show usage
    -r  don't follow redirects
    -x  if -d is set, mark the file as executable
```

## Known issues
* **Crappy filename handling** downloaded files have weird filenames when when filename in Content-Disposition is unspecified; I have no idea how to fix this, if you know something I don't, send a pull request
