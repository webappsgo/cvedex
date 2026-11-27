# Self-hosting

The service can run on a server machine or in a container. Persistent storage holds the compressed CVE dataset, its checksum, and SQLite state. For public use, configure the zone as a delegated subdomain that you control; the default `cve.` zone is intended for networks with custom DNS routing.

The service downloads the latest checksummed data asset on first start. The
`cve_update` scheduler task refreshes it on the configured schedule, which
defaults to 05:00 daily in the scheduler timezone; configured retries use
exponential delays. Send `SIGHUP` for an immediate refresh. Operators may
mirror the asset and set `DATA_URL` and `CHECKSUM_URL`. Keep DNS port 53
reachable over both UDP and TCP. The Compose configurations publish both
transports as well as the HTTP port.

Before exposing an instance to the internet, complete first-run administrator
setup, review transfer ACLs and NOTIFY targets, and publish the active DNSSEC
DS record at the parent zone. DNS UDP and HTTP response rate limiting are
enabled by default; tune them for the expected query volume.

## HTTPS and overlay networks

HTTPS is available with a manually supplied certificate and key using the
`server.ssl` YAML settings. The listener enforces TLS 1.2 or newer. Certificate
issuance, ACME challenges, and automated renewal are not yet implemented.

The service starts an overlay manager with dedicated loopback HTTP backends.
When a Tor binary is available, the manager negotiates control authentication,
subscribes to bootstrap and hidden-service descriptor events before enabling
network access, and waits up to 60 seconds for descriptor upload. Control
events and 10-minute health checks trigger process restart after failure. I2P
starts when enabled and uses `i2pd` or a reachable SAM bridge to serve the
eepsite through its backend. Overlay enabled, running, status, provider, and
hostname fields are included in JSON and plain-text health output. I2P
provider/session recovery and retry backoff remain incomplete. See
[Integrations](integrations.md#tor-and-i2p-overlays).

## DNS secondary and DNSSEC

The server signs zones offline when DNSSEC is enabled, serves AXFR/IXFR to
clients allowed by `DNS_TRANSFER_ACL`, and sends best-effort NOTIFY after
successful refreshes. Persist `DATA_DIR` across restarts to retain DNSSEC
keys. Obtain the active DS from `/api/v1/ds` and publish it at the parent
zone. See [DNSSEC](dnssec.md) and [Integrations](integrations.md).
