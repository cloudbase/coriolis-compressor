# Coriolis compressor

This is a simple web app that accepts a POST of binary data, and echoes back a compressed version of that body. The compression is done using a [parallel implementation of gzip/zlib](https://github.com/klauspost/pgzip), so for large chunks, there is an added benefit of using all cores available to the system, to do the actual compression.

The purpose of this is to be used in Coriolis to speed up replica transfers. The built in ```zlib``` module is single threaded, and takes longer to compress bigger chunks of data.

This component will probably be replaced with a native python module that leverages 7zip or some other implementation of ```zlib``` and ```gzip```.

## Build it

```bash
go install -ldflags="-s -w" -mod vendor ./...
```

## API

```bash
POST /
```

Headers:

| Name                 | Type   | Optional | Description                                            |
| -------------------- | ------ | -------- | ------------------------------------------------------ |
| X-Compression-Format | string |   true   | Compression format to use. Possible values: gzip, zlib |

## Example usage

```bash
$ echo 'Hi there!' > /tmp/hello-world.txt
$ curl -s -X POST -H 'X-Compression-Format: zlib' \
    --data-binary @/tmp/hello-world.txt \
    http://127.0.0.1:7766/ | zlib-flate -uncompress
Hi there!
```