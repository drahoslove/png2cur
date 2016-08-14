# png2cur
Convert PNG to CUR

Written in [golang](https://golang.org/).

## Build
```
go get -u github.com/Drahoslav7/png2cur
```

## Examples
simple:
```
png2cur arrow.png
```
custom output name and shifted hotspot:
```
png2cur arrow.png -x 2 -y 2 -o cursor.cur
```

### Notes
Only support images with RGBA color model and with maximal dimensions of 256×256 px.
