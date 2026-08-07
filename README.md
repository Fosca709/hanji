# Hanji

Tiny desktop memo app.

## Setup

Install [uv](https://docs.astral.sh/uv/), then install the dependencies:

```bash
uv sync
```

Run the app:

```bash
uv run python main.py
```

## Shortcut

- `Ctrl + Q`: toggle `Always on Top`

## Build

Build the app binary:

```bash
bash scripts/build.sh
```

Make a Debian package:

```bash
bash scripts/package-deb.sh
```
