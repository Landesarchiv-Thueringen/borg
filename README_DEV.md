# BorgFormat

Identifies and validates file formats.

## Configuration

Copy `.env.example` to `.env` and make changes as required:

```sh
cp .env.example .env
$EDITOR .env
```

## Run

Build or rebuild images as required and start backend and frontend services:

```sh
docker compose up --build
```

## Embedding

Borg can be embedded in other Angular applications.
See [./gui/src/app/features/file-analysis/README.md](./gui/src/app/features/file-analysis/README.md).

## Development

### Start a Frontend Development Server

In the directory `gui` run

```sh
npm install
npm start
```

### Go Workspaces

We use Go workspaces so any tooling (e.g., IDE integration) knows about directories that contain Go modules.

To add a new directory, run for example

```sh
go work use ./tools/new-tool
```
