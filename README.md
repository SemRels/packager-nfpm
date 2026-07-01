# packager-nfpm

[![Latest Release](https://img.shields.io/github/v/release/SemRels/packager-nfpm?label=version\&color=blue)](https://github.com/SemRels/packager-nfpm/releases/latest)

nFPM packager plugin for semrel. It creates Linux packages (`deb`, `rpm`, `apk`) by invoking `nfpm package` during the release pipeline.

## Repository Layout

```text
cmd/plugin/              Plugin entry point
internal/plugin/         Business logic scaffold
internal/grpc/           gRPC transport scaffold
proto/v1                 Symlink to the SemRel protobuf contract
.github/workflows/       CI, release, and security automation
```

## Installation

Published binaries are distributed through releases and synchronized to `registry.semrel.io`.

## Development

```bash
go build ./cmd/plugin
go test ./...
```

## Configuration

Configure in `.semrel.yaml`:

```yaml
plugins:
	- uses: packager-nfpm
		args:
			config: packaging/nfpm.yaml
			target: dist/packages
			packagers: deb,rpm,apk
```

Runtime inputs:

- `SEMREL_VERSION` / `SEMREL_NEXT_VERSION` (required)
- `SEMREL_DRY_RUN` (`true` to print commands only)
- `SEMREL_PLUGIN_CONFIG` (default: `nfpm.yaml`)
- `SEMREL_PLUGIN_TARGET` (default: `dist`)
- `SEMREL_PLUGIN_PACKAGERS` (CSV, default: `deb,rpm,apk`)

Dependencies:

- `nfpm` must be installed and available on `PATH`.
