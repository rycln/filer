# Filer

Interactive REPL utility for quick file sorting through a "keep or delete" process.

## What is it?

`filer` helps you quickly sort through a pile of files. It shows you one file at a time, and you decide whether to **keep** it (move to a target directory) or **delete** it.

## Usage

```bash
filer [-s SOURCE_DIR] [-t TARGET_DIR] [-p REGEX_PATTERN]
```

## Arguments

- -s, --source SOURCE_DIR - Directory with files to sort (default: current directory)
- -t, --target TARGET_DIR - Directory where kept files will be moved (default: files remain in place)
- -p, --pattern REGEX_PATTERN - Regular expression to filter files (e.g., "\.jpg$", "^2024-", ".*\.(jpg|png)$")

## Controls

When running, you'll see:

```bash
[1/1000]

File: photo001.jpg

[K]eep, [D]elete, [S]kip, [Q]uit? 
```

- k - Keep the file (moves to target_dir if specified)
- d - Delete the file permanently
- s - Skip file 
- q - Exit the application

## Note

Deletion is permanent - use with caution!




