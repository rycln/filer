# Filer

Interactive REPL utility for quick file sorting through a "keep or delete" process.

## What is it?

`filer` helps you quickly sort through a pile of files. It shows you one file at a time, and you decide whether to **keep** it (move to a target directory) or **delete** it.

## Usage

```bash
filer [source_dir] [target_dir]
```

## Arguments

- source_dir - Directory with files to sort (default: current directory)
- target_dir - Directory where kept files will be moved (default: files remain in place)

## Controls

When running, you'll see:

```bash
File: photo001.jpg
[K]eep, [D]elete, [Q]uit? _
```

- K or k - Keep the file (moves to target_dir if specified)
- D or d - Delete the file permanently
- Q or q - Exit the application

## Note

Files are processed in alphabetical order. Deletion is permanent - use with caution!




