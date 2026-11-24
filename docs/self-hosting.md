# Self-hosting

The service can run on a host or in a container. Persistent data is required for the compressed CVE dataset and its checksum. For public use, configure the zone as a delegated subdomain that you control; the default `cve.` zone is intended for networks with custom DNS routing.

The current implementation does not automatically fetch or refresh the dataset. Generate the data file separately and place it in the configured data directory. Keep DNS port 53 reachable over both UDP and TCP. The Compose configuration must publish both transports as well as the HTTP port.

Before exposing an instance to the internet, account for the missing DNSSEC signing, AXFR ACLs, response-rate limiting, HTTP rate limiting, admin authentication, and automated refresh features described in [Security](security.md).
