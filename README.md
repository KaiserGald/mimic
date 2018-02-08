# Mimic

##### current version: 0.2.3

Mimic is a simple file cloner. Any changes made to the source directory will be reflected in the destination directory. It has
handling for all system events and will make changes accordingly.

## Getting Started

Using mimic is simple:

Install it using the makefile.
```bash
make all
```

Basic usage:
```bash
mimic -w "sourcedir:destinationdir"
```
This will immediately clone all files and directories in 'sourcedir' into 'destinationdir'. It will then watch 'sourcedir' for any
changes.

## Configuration

#### Color

Mimic supports colored output. Simply start it with the color flag.
```bash
mimic -c -w "sourcedir:destinationdir"
```
*(only tested in GNOME terminal, but should work in most bash)

#### Output

By default mimic only outputs error and file creation events. For more information try the ```-v``` flag.
```bash
mimic -v -w "sourcedir:destinationdir"
```

To have it only output errors, use the ```-q``` flag.
```bash
mimic -q -w "sourcedir:destinationdir"
````
