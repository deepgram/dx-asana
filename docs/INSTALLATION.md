# Installation

## Prerequisites

- Go 1.22 or newer
- An Asana personal access token from [app.asana.com/0/my-apps](https://app.asana.com/0/my-apps)

## From source

```bash
git clone https://github.com/deepgram/dx-asana.git
cd dx-asana
make install
```

`make install` runs `go install` with the version, commit, and build-date ldflags baked in. The binary lands in `$GOBIN` (or `$GOPATH/bin` if `GOBIN` is unset). Make sure that directory is on your `PATH`.

To build into the working directory instead:

```bash
make build      # produces ./dx-asana
./dx-asana --version
```

## Authentication

```bash
export ASANA_TOKEN=your-token-here
# or persist it to ~/.dx-asana/config.json
dx-asana config set --token your-token-here
```

## Homebrew

Paused for now. A `deepgram/homebrew-tap` formula will be wired back in once the binary name and release cadence settle on this fork.

## Migrating from the upstream `asana-cli`

If you previously used the upstream binary, your config lives at `~/.asana-cli/config.json`. `dx-asana` reads from `~/.dx-asana/config.json` instead. Either copy your token over:

```bash
mkdir -p ~/.dx-asana && cp ~/.asana-cli/config.json ~/.dx-asana/config.json
```

…or just `dx-asana config set --token …` from scratch.
