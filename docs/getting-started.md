# Getting started

1. Install cvedex using the [installation guide](installation.md).
2. Place a valid compressed dataset and matching SHA-256 file in the configured data directory.
3. Start the service with DNS port 53 available, or set `DNS_PORT` to a development port.
4. Query the configured zone with `dig` or use the JSON API.

```sh
dig TXT severity.2024-3094.cve @127.0.0.1
curl -H 'Accept: application/json' http://127.0.0.1:64580/api/v1/cves/CVE-2024-3094
```

The initial snapshot must be generated separately. Automatic first-run download and daily refresh are not implemented yet.
