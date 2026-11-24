# Self-hosting

The service can run on a host or in a container. Persistent storage holds the compressed CVE dataset, its checksum, and SQLite state. For public use, configure the zone as a delegated subdomain that you control; the default `cve.` zone is intended for networks with custom DNS routing.

The service downloads the latest checksummed data asset on first start and refreshes daily. Operators may mirror the asset and set `DATA_URL` and `CHECKSUM_URL`. Keep DNS port 53 reachable over both UDP and TCP. The Compose configurations publish both transports as well as the HTTP port.

Before exposing an instance to the internet, account for the missing DNSSEC signing, AXFR ACLs, response-rate limiting, HTTP rate limiting, admin authentication, and automated refresh features described in [Security](security.md).
