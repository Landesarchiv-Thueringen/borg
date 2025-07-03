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
See [README.md](https://github.com/Landesarchiv-Thueringen/borg/blob/main/gui/src/app/features/file-analysis/README.md).

## Development

### Start a Frontend Development Server

In the directory `gui` run

```sh
pnpm install
pnpm start
```

### Go Workspaces

We use Go workspaces so any tooling (e.g., IDE integration) knows about directories that contain Go modules.

To add a new directory, run for example

```sh
go work use ./tools/new-tool
```

### Releasing a new version

- Choose a version tag based on semantic versioning. In most cases, this means incrementing the minor version when there are new features and otherwise, incrementing the patch version.
- Update `CHANGELOG.md` with the chosen version tag and any changes.
- Update the version env in `compose.yml`.
- Push any changes to `main`.
- Draft a new [release](https://github.com/Landesarchiv-Thueringen/borg/releases) on GitHub.
  - Include the release's section of `CHANGELOG.md` as description.

## Documentation

The documentation is generated with [MkDocs](https://www.mkdocs.org/). To change it, edit the markdown files under docs.

### Build

To build and server the documentation use

```sh
mkdocs serve
```

#### Prerequisites (Arch Linux)

- mkdocs
- python-pymdown-extensions
- mkdocs-glightbox (AUR)

### Deploy

To deploy the documentation to GitHub Pages use

```sh
mkdocs gh-deploy
```
