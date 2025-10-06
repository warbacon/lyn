# Lyn ☁️

> [!IMPORTANT]
> Work in progress!

Lyn is a minimal, fast, and cross-platform web server for serving static files
and browsing directories, inspired by Python's `http.server`.

## 🚀 Usage

```bash
lyn
```

This will start a server in the current directory.

You can also customize its behavior using additional options — run:

```bash
lyn -h
```

for a full list of available flags.

## ⚙️ Installation

Currently, manual compilation is required, but it's very straightforward with
go. Simply run:

```bash
go build .
```

inside this repository.

This will generate a binary named `lyn` (or `lyn.exe` on Windows) that you can
run directly.

## ⚠️ Disclaimer

**Lyn is not intended for production use.** Use it for development, testing, or
learning purposes only.
