# BorgFormat

Identifies and validates file formats.

## Run

Build or rebuild images as required and start backend and frontend services:

```sh
docker compose up --build
```

## Development

### Frontend

#### Start a Development Server

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
