# addlicense

The program ensures source code files have copyright license headers
by scanning directory patterns recursively.

It modifies all source files in place and avoids adding a license header
to any file that already has one.

## install

    go get -u github.com/google/addlicense

## usage

    addlicense [flags] pattern [pattern ...]
    
    -c copyright holder (default "Google Inc.")
    -l license type: apache, mit (default "apache")
    -y year (default 2016)

The pattern argument can be provided multiple times, and may also refer
to single files.

## custom license

modify my.tmpl as your license template.

	addlicense -l my partern...

note: tmpl file name(witch out .tmpl extention) always same as template name


## license

Apache 2.0

This is not an official Google product.
