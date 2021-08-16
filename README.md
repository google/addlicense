# addlicense

The program ensures source code files have copyright license headers
by scanning directory patterns recursively.

It modifies all source files in place and avoids adding a license header
to any file that already has one.

addlicense requires go 1.16 or later.

## install

    go get -u github.com/google/addlicense

## usage

    addlicense [flags] pattern [pattern ...]

    -c copyright holder (defaults to "Google LLC")
    -f custom license file (no default)
    -l license type: apache, bsd, mit, mpl (defaults to "apache")
    -y year (defaults to current year)
    -check check only mode: verify presence of license headers and exit with non-zero code if missing
    -ignore file patterns to ignore, for example: -ignore **/*.go -ignore vendor/**

The pattern argument can be provided multiple times, and may also refer
to single files.

The `-ignore` flag can use any pattern [supported by
doublestar](https://github.com/bmatcuk/doublestar#patterns).

## Running in a Docker Container

The simplest way to get the addlicense docker image is to pull from GitHub
Container Registry:

```bash
docker pull ghcr.io/google/addlicense:latest
```

Alternately, you can build it from source yourself:

```bash
docker build -t ghcr.io/google/addlicense .
```

Once you have the image, you can test that it works by running:

```bash
docker run -it ghcr.io/google/addlicense -h
```

Finally, to run it, mount the directory you want to scan to `/src` and pass the
appropriate addlicense flags:

```bash
docker run -it ghcr.io/google/addlicense -v ${PWD}:/src -c "Google LLC" *.go
```

## license

Apache 2.0

This is not an official Google product.
