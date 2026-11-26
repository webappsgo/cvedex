# Getting started

1. Install cvedex using the [installation guide](installation.md).
2. Set `DATA_URL` and `CHECKSUM_URL` if you use a data mirror. The defaults point to the project's latest published data asset.
3. Start the service with DNS port 53 available, or set `DNS_PORT` to a development port. The host HTTP port is selected once and persisted; set `CVEDEX_PORT` if you want a fixed port.
4. Query the configured zone with `dig` or use the JSON API.

```sh
dig TXT severity.2024-3094.cve @127.0.0.1
curl -q -LSsf -H 'Accept: application/json' http://127.0.0.1:64580/api/v1/cves/CVE-2024-3094
```

The sample HTTP command uses port 64580; set `CVEDEX_PORT=64580` when starting cvedex to use it. On first start cvedex downloads and verifies the current dataset. Until a valid snapshot is available, DNS data queries return `SERVFAIL` and readiness reports unhealthy. Later refreshes run daily by default; a failed refresh leaves the previous verified dataset serving.
