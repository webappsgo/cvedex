# Command-line client

The client binary is `cvedex-cli`. Shared flags are `--server`, `--token`, `--output`, `--debug`, `--color`, and `--lang`. Output is plain text by default; use `--output json` for machine-readable responses.

```sh
cvedex-cli get CVE-2024-3094
cvedex-cli get CVE-2024-3094 --output json
cvedex-cli search "xz utils" --server http://127.0.0.1:64580
cvedex-cli --shell bash
```

## Data generation

`data generate` builds the canonical compressed dataset from a local cvelistV5 checkout and the official KEV and EPSS feeds:

```sh
cvedex-cli data generate --source-dir ./cvelistV5 --output ./cvedex-data.json.gz
```

The command also writes the matching `.sha256` file. Feed URLs can be overridden with `--kev-url` and `--epss-url`.

`data sync` downloads and verifies an operator-specified HTTPS data URL and checksum URL before writing a local snapshot:

```sh
cvedex-cli data sync --data-url https://example.invalid/cvedex-data.json.gz \
  --checksum-url https://example.invalid/cvedex-data.json.gz.sha256 \
  --output ./cvedex-data.json.gz
```

## Zone download

`zone generate` downloads the BIND-format zone from the configured server:

```sh
cvedex-cli zone generate --output ./cvedex.zone
```

The downloaded zone file is signed when the server has DNSSEC enabled and a
verified dataset loaded. The CLI does not sign zones locally.
