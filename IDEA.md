# cvedex Specification

## Project description

cvedex is an authoritative DNS server that exposes the entire CVE (Common Vulnerabilities and Exposures) corpus as a queryable zone under the custom `.cve` TLD. Every CVE becomes a hierarchy of DNS names, so any standard tool — `dig`, `host`, `drill`, `nslookup`, log enrichers, threat-intel pipelines — can look up CVE details with a plain DNS query. Type `2024-3094.cve` into a browser and land on the canonical CVE page. Run `dig TXT cvss.2024-3094.cve` and get the score. Slave the zone into your existing BIND/Knot/PowerDNS infrastructure via AXFR.

cvedex is built for security teams, threat-intel pipelines, log-enrichment workflows, homelabs, and any operator who controls their own DNS resolution and wants frictionless programmatic access to CVE data without hammering rate-limited APIs. The data is always fresh (rebuilt from the official CVE list on a daily schedule), enriched with the CISA KEV catalog and FIRST.org EPSS scores, signed with DNSSEC, and served from a single static binary with zero external runtime dependencies.

The zone is configurable. The default is `cve.` — a custom internal-use TLD for networks where the operator controls DNS resolution (split-horizon setups, internal resolvers, lab environments). It can also be set to a real subdomain the operator owns and delegates from a parent zone (`cve.mydomain.example`), in which case cvedex works as a publicly-resolvable service with standard DNSSEC delegation and public-CA HTTPS for the redirector.

## Project variables

```
project_name:  cvedex
project_org:   casapps
internal_name: cvedex
app_name:      cvedex
```

## Business logic

### Product scope & non-goals

**In scope:**

- Authoritative DNS for one configurable zone (default `cve.`), exposing the full CVE corpus (~350k records) as queryable names
- Per-CVE records covering description, CVSS, severity, vector, CWE, dates, state, vendors, products, references, KEV flag, EPSS score
- Zone-level index records for total counts, per-year counts, per-severity counts, KEV count, latest CVE, build metadata
- A URI record (RFC 7553) at each CVE name returning the canonical upstream URL
- HTTP redirector that lets `https://2024-3094.cve` route to the canonical CVE page on the upstream of the operator's choice
- JSON gateway exposing per-CVE detail, vendor reverse lookups, product reverse lookups, vendor/product enumeration, full-text search, KEV catalog dump, DNSSEC DS record, and zone-file download
- DNSSEC signing (offline by default, online optional)
- AXFR/IXFR/NOTIFY for slaving the zone into existing DNS infrastructure
- Three modes: long-running server, one-shot zone-file generator, one-shot sync-only refresh
- Daily refresh from the canonical sources with atomic zone swap and graceful degradation
- A landing page served at the zone apex (and at any unknown hostname under the zone) explaining what cvedex is and how to use it

**Non-goals:**

- Not a recursive resolver
- Not a general-purpose DNS server (it serves the CVE zone and nothing else; out-of-zone queries are REFUSED by default)
- No write API — the zone is read-only, rebuilt from canonical sources
- No operator-injected custom records — the zone is built exclusively from canonical sources, no manual additions, no per-deployment local extensions
- No web UI for browsing CVEs (the JSON gateway is for clients; humans use `dig`, `curl`, or downstream tools)
- No vendor-neutral reverse-CNAME-for-vendor scheme via DNS — vendor/product reverse lookups live in the JSON gateway because they don't fit a DNS-shaped query model
- No DNS-over-TLS or DNS-over-HTTPS — operators who want encrypted DNS transport front cvedex with a resolver that supports it

### Roles & permissions

| Role | Capabilities |
|------|--------------|
| **Anonymous DNS client** | Query any name in the zone (DNS read). Subject to RRL. |
| **Anonymous HTTP client** | Read any endpoint on the JSON gateway. Hit the redirector. Subject to standard HTTP rate limiting. |
| **AXFR secondary** | Pull the zone via AXFR/IXFR. ACL-gated; default localhost only. |
| **Server Admin** | Standard template-defined admin: manages refresh schedule, upstream redirect URL template, AXFR ACL, NOTIFY targets, DNSSEC key generation/rotation, enrichment toggles, RRL parameters, out-of-zone policy. |

There is no Regular User concept — the data is anonymous-public read-only. Multi-user features (PART 34) are not implemented.

### Data model & sensitivity

| Data | Source | Classification | Notes |
|------|--------|----------------|-------|
| CVE records | cvelistV5 (CVE Program) | **Public** | Canonical vulnerability data |
| KEV flags & metadata | CISA KEV catalog | **Public** | Authoritative exploitation status |
| EPSS scores | FIRST.org | **Public** | Advisory exploit-prediction scores |
| Build metadata (timestamp, source SHA) | Generated locally | **Public** | Surfaced via index records |
| DNSSEC private keys (KSK, ZSK) | Generated locally | **Sensitive** | Stored in `{data_dir}`, readable only by service account |
| Server config | Operator-provided | **Sensitive** | Standard template handling |
| DNS query logs (when enabled) | Generated locally | **Sensitive** | Reveals which CVEs a network is researching — itself threat-intel-grade information. Default off. |
| HTTP access logs | Generated locally | **Standard** | Template-handled retention |
| Refresh / build logs | Generated locally | **Standard** | Operational only |

### Trust boundaries & external services

| Service | Trust | Failure mode |
|---------|-------|--------------|
| Data-file distribution endpoint (default: project's GitHub Releases) | **Trusted source of CVE data for the running app.** TLS authenticates the host; SHA-256 checksum verifies content integrity. | If unreachable: keep serving previous build; surface staleness via index records, metrics, and logs. Never serve no zone. |
| Configured upstream redirect target (default `cve.org`) | **Untrusted at runtime** — only used as a redirect destination, never fetched or parsed. Operator can reconfigure. | Not validated; if dead, redirector still issues 302 (browser surfaces the error). |
| AXFR secondaries | Trusted by ACL. Operator-controlled list. | If a secondary is unreachable, NOTIFY is best-effort; the secondary catches up on next refresh poll. |

The data-generation workflow has its own upstream dependencies (cvelistV5 git repository, CISA KEV JSON, FIRST.org EPSS CSV, the GitHub Releases publishing endpoint) — all over HTTPS, all canonical authoritative sources. Those are not direct dependencies of running cvedex instances; running instances only depend on the data-file distribution endpoint above.

cvedex never accepts inbound writes, never executes upstream-supplied content, and never proxies upstream queries.

### Threat model & abuse cases

**Primary assets being protected:**

1. Zone integrity (clients get authentic, unaltered CVE data)
2. Zone availability (queries continue serving even when sources hiccup)
3. Data freshness (the zone reflects the canonical CVE corpus within the refresh window)
4. Operator privacy (DNS queries reveal what an org is researching)

**Trusted vs untrusted inputs:**

| Input | Trust |
|-------|-------|
| cvelistV5 git repo content | Trusted (HTTPS to github.com authenticates) |
| CISA KEV JSON | Trusted (HTTPS) |
| FIRST.org EPSS CSV | Trusted (HTTPS) |
| Operator config | Trusted |
| Inbound DNS queries | Untrusted — validate, rate-limit, never parse for code paths |
| Inbound HTTP queries | Untrusted — validate Host header, validate query params, standard template handling |
| Inbound AXFR requests | Untrusted by default — ACL-gated |

**Attacker / abuser goals and defenses:**

- **DNS amplification abuse** — defended by Response Rate Limiting (RRL) on by default, REFUSED for out-of-zone queries, RFC 8482 minimal response for ANY queries, EDNS(0) buffer size capped.
- **Zone tampering / spoofing** — defended by DNSSEC signing of every successful build (NSEC3 no-opt-out, ECDSA P-256). Operators delegating from a parent zone publish the DS record exposed by cvedex.
- **Denial of service via malformed queries** — defended by strict RFC compliance (NXDOMAIN/NOERROR per spec), short-circuit on out-of-zone, RRL, standard template-handled HTTP rate limits.
- **Resource exhaustion via build failure** — defended by atomic staging-then-swap rebuild (failed builds never affect the live zone), graceful degradation when enrichment sources fail, "previous good zone keeps serving" as the default failure behavior.
- **Source poisoning (compromised upstream)** — partially defended by HTTPS to authenticated sources at the workflow tier. The CVE Program does not cryptographically sign individual CVE records. If an upstream source is compromised, the workflow ingests the compromised data and the resulting data file propagates it — this is a documented limitation, not a defect.
- **Compromised data-file distribution** — defended by TLS to the configured distribution endpoint and SHA-256 checksum verification on every fetch. If the publishing infrastructure itself is compromised (the attacker can recompute valid checksums), operators with high-trust requirements can self-host the workflow and serve from their own endpoint. Cryptographic signing of the data file (Sigstore, GPG) is a future enhancement; current trust model relies on TLS + checksum + operator's choice of distribution origin.
- **Query-log harvesting / privacy leak** — defended by DNS query logging being **off by default**. Operators who opt in must understand they're collecting threat-intel-grade data. DNS-over-TLS / DNS-over-HTTPS is not a goal for cvedex itself; if the operator wants encrypted transport they front cvedex with a resolver that supports it.
- **Slow-loris / connection exhaustion** — defended by template-handled HTTP timeouts and connection limits.
- **AXFR data scraping** — defended by ACL (default localhost only). Operators who expose AXFR publicly do so deliberately.

**Project-specific abuse cases:**

- **Spam / scraping the JSON gateway** — defended by template-standard rate limiting. The data is public, so scraping isn't an attack so much as a load concern.
- **Malicious upload** — N/A; no upload surface.
- **SSRF** — N/A; cvedex never fetches operator-supplied URLs.
- **Credential stuffing** — N/A; no end-user authentication. Server Admin auth is template-handled.

### Security decisions & exceptions

- **Zone is configurable; deployment-mode caveats follow.** See the "Deployment modes" subsection. Custom-TLD mode requires operator-controlled DNS resolution and limits HTTPS/DNSSEC posture; delegated-subdomain mode unlocks public-CA HTTPS and standard DNSSEC delegation. Choice is the operator's; defaults are conservative for the custom-TLD case.
- **Redirector HTTPS is mode-dependent.** In custom-TLD mode, HTTPS requires an internal CA (operator-supplied). In delegated-subdomain mode, standard public-CA HTTPS is supported (template-handled ACME). HTTP-only is the safe default in both modes; the redirect destination is always HTTPS regardless, so users land on a real cert.
- **DNS query logging off by default.** Privacy decision: per-CVE query logs are threat-intel-grade. Opt-in only. Documented.
- **AXFR localhost-only by default.** Operators who want remote secondaries explicitly extend the ACL. Documented in the operator guide.
- **Out-of-zone REFUSED by default.** Authoritative-only posture. Operators who want forward-to-upstream behavior opt in. Documented.
- **No upstream signature verification.** cvelistV5 isn't cryptographically signed at the record level by the CVE Program. cvedex trusts HTTPS to GitHub. Documented as a limitation.
- **Online DNSSEC signing is opt-in.** Default is offline signing at build time (signatures fresh on every refresh). Online signing burns CPU per query and is only useful for operators who need sub-refresh-cycle signature freshness.
- **DNSSEC signature validity windows.** RRSIGs are issued on each rebuild with a validity window comfortably longer than the refresh interval (default 30 days, with rebuilds expected daily). If rebuilds fail repeatedly past the validity window, the zone goes BOGUS and DNSSEC-validating resolvers will refuse cached answers. Build failures are surfaced via metrics so operators are alerted long before signatures expire.

### Deployment modes

cvedex supports two deployment modes, distinguished by the zone name:

- **Custom-TLD mode (default).** Zone is `cve.` — a custom, ICANN-undelegated TLD. cvedex serves the zone authoritatively, but the records are only resolvable on networks that route `.cve` queries to cvedex (split-horizon DNS, internal resolvers, lab environments, dedicated security-tooling networks). Implications:
  - Public-CA TLS certificates for `*.cve` are not obtainable. The redirector runs HTTP-only by default; operators who want HTTPS bring their own internal CA.
  - DNSSEC has no parent zone to publish the DS to. Validating resolvers must install cvedex's KSK as a trust anchor (or run with DNSSEC validation off for the `.cve` tree). The DS record is exposed at `ds.<zone>` and via the JSON gateway for operator convenience.
  - cvedex can be deployed anywhere on the operator's internal network without coordinating with any external party.

- **Delegated-subdomain mode.** Zone is a real subdomain the operator owns and delegates from a parent zone they control — e.g., `cve.mydomain.example`. cvedex is the authoritative server for that delegated zone. Implications:
  - The redirector can run HTTPS with standard public-CA certificates for `*.cve.mydomain.example` (typically via ACME with a DNS-01 challenge — handled by the casapps Go server template's ACME integration).
  - DNSSEC works in the standard delegated way: the parent zone publishes the DS record cvedex exposes at `ds.<zone>`. Validating resolvers across the public internet validate signatures normally, no trust-anchor wrangling required.
  - cvedex can be publicly resolvable on the internet — anyone can `dig 2024-3094.cve.mydomain.example` from anywhere with no special configuration.

Both modes use exactly the same record schema, refresh behavior, JSON gateway, and HTTP redirector logic. The only differences are the implications above. All examples in this document use the custom-TLD form (`2024-3094.cve`) for brevity; substitute the operator's chosen zone for delegated-subdomain deployments.

### Data distribution model

cvedex uses a two-tier data distribution architecture rather than each running instance pulling raw upstream data:

**Tier 1 — Data-generation workflow.** A scheduled CI workflow (running in the cvedex project's repository, or self-hosted by operators who prefer it) processes the canonical upstream sources (cvelistV5 git repository, CISA KEV catalog, FIRST.org EPSS scores) into a single compressed, checksummed, schema-versioned data file. The workflow runs daily and publishes the resulting data file as a release asset on GitHub Releases. Data releases are marked as pre-release so they do not appear as the `latest` release — `latest` is reserved for binary releases. End users do not interact with the workflow directly; it is project-maintenance infrastructure.

**Tier 2 — Running cvedex instances.** Each running instance fetches the latest published data file on its configured refresh cadence, verifies the checksum, parses it, and atomically swaps the new zone into place. End-user instances do not clone the cvelistV5 git repository, do not fetch KEV or EPSS directly, and do not parse 350,000 individual JSON records on every refresh — that work is centralized in the workflow and reused by every consumer.

Why this architecture:

- End-user instances refresh in seconds rather than tens of minutes, with a single HTTP fetch instead of a multi-GB git pull plus per-record parsing
- End-user instances have no git tooling requirement and minimal disk needs
- Data-quality fixes (parser bugs, schema decisions, edge-case handling) ship once via the workflow rather than per-deployment
- Operators who need independence from the project's release pipeline can run their own copy of the workflow and configure cvedex to fetch from their own URL — same data file format, different distribution origin

The data-file URL is operator-configuraribution endpoint, verify SHA-256 checksum, parse, build a fresh in-memory zone in a staging structure, sign with DNSSEC, atomically swap into place.
- On successful swap: SOA serial advances, NOTIFY is sent to configured secondaries.
- On any failure (fetch, checksum mismatch, schema-version mismatch, parse, sign, disk): the previous good zone keeps serving, the failure is logged and surfaced via metrics, the SOA serial does NOT advance.
- If the distribution endpoint is unreachable: the running zone keeps serving previous data; staleness is exposed via the `built.cve` and `source.cve` index records and via metrics.
- A fetched data file that fails checksum verification is rejected and discarded; the previous zone keeps serving and an error is logged.
- KEV and EPSS data are baked into the data file by the workflow. Running instances do not fetch or refresh those upstream sources separately. Maximum-staleness handling for KEV and EPSS is the workflow's responsibility (default 7-day threshold; past that, the workflow drops the field rather than emitting long-stale values into the data file).
- The data file's freshness (its `built` timestamp and source SHA) is exposed via the zone's `built.cve` and `source.cve` records. Operators alert on stale data via their existing monitoring infrastructure.
- First run on an empty data dir: fetch the latest data file from the distribution endpoint. Until the first successful fetch and build, readiness is false.

### HTTP redirector

cvedex's HTTP layer (handled by the casapps Go server template) inspects the `Host` header on every request:

- **Host matches `YYYY-NNNNN.<zone>` and the CVE exists**: 302 redirect to the configured upstream URL template. Default template targets `cve.org`. Configurable; available substitution variables are the year, the sequence number, and the full CVE ID. 302 chosen by default so changing the upstream template doesn't burn long-lived caches; 301 is selectable for operators who want it.
- **Host is the bare zone apex** (e.g., `https://cve`): serve the landing page (200).
- **Host matches `YYYY-NNNNN.<zone>` but the CVE doesn't exist, or the hostname is malformed, or the hostname is otherwise out-of-zone**: 404 status with a human-readable HTML body explaining what cvedex is, the `YYYY-NNNNN.cve` naming convention, a working example link, instructions for using the DNS interface (`dig TXT 2024-3094.cve`), a pointer to the JSON gateway, and an operator-customizable footer.

The URI DNS record and the HTTP redirector resolve to the same upstream URL — they're driven by the same template, so changing the template updates both on the next build.

### JSON gateway data surface

The HTTP gateway (routing/middleware/auth/CORS/rate-limiting/health/readiness/metrics all template-handled) exposes:

- **Per-CVE detail** — every field surfaced via DNS plus richer detail that doesn't fit DNS cleanly (full reference URL list, full affected-products list, all CVSS metric variants, full KEV metadata, full EPSS data)
- **Vendor reverse lookups** — given a vendor, list of CVE IDs with severity counts and most-recent timestamp; paginated
- **Product reverse lookups** — given a vendor and product, list of CVE IDs; paginated
- **Vendor and product enumeration** — paginated lists of all known vendors (with counts) and products for a given vendor
- **Search** — keyword search across descriptions, filterable by severity, year, KEV flag; cursor-paginated
- **Latest CVEs** — most recently published, configurable count
- **Full KEV catalog** — current CISA KEV list as JSON
- **DNSSEC DS record** — for parent-zone delegation
- **Zone-file download** — current zone in BIND master-file format (same as Generate mode output)

### Modes of operation

- **Server (default)**: long-running daemon. Listens on UDP+TCP for DNS and on the template-defined HTTP listener. Fetches the data file on its refresh schedule. Atomic zone swap on rebuild — no query downtime.
- **Generate-zone**: one-shot. Fetches the latest data file, writes a signed BIND-format zone file to disk, exits. For operators wble. Default points at the project's GitHub Releases. Operators can mirror, self-host, or air-gap as needed. The data file format is schema-versioned so that newer cvedex binaries can refuse to load incompatible older data files (and vice versa).

### Zone naming and per-CVE records

The zone apex is configurable. The default is `cve.` (custom-TLD mode); operators delegating from a parent zone use their chosen subdomain, e.g., `cve.mydomain.example`. Every CVE ID maps to a label of the form `YYYY-NNNNN` under the zone — the `CVE-` prefix is dropped because the zone name already conveys the context. So `CVE-2024-3094` becomes `2024-3094.cve` in custom-TLD mode, or `2024-3094.cve.mydomain.example` in delegated-subdomain mode. All record-name examples below use the custom-TLD form for brevity.

For each CVE, the following names are queryable:

| Name | Type | Content |
|------|------|---------|
| `2024-3094.cve` | A/AAAA | IP of the cvedex HTTP redirector (so browsers reach the landing/redirect) |
| `2024-3094.cve` | TXT | One-line summary: severity, CVSS score, truncated description |
| `2024-3094.cve` | URI | RFC 7553 record returning the canonical upstream URL |
| `desc.2024-3094.cve` | TXT | Full English description (chunked per RFC 1035 if >255 bytes) |
| `cvss.2024-3094.cve` | TXT | Numeric CVSS base score (highest version available — v4 > v3.1 > v3.0 > v2) |
| `severity.2024-3094.cve` | TXT | One of `none / low / medium / high / critical` |
| `vector.2024-3094.cve` | TXT | Full CVSS vector string |
| `cwe.2024-3094.cve` | TXT | Comma-separated CWE IDs |
| `published.2024-3094.cve` | TXT | ISO-8601 publish date |
| `modified.2024-3094.cve` | TXT | ISO-8601 last-modified date |
| `state.2024-3094.cve` | TXT | `PUBLISHED` / `REJECTED` / `RESERVED` |
| `vendors.2024-3094.cve` | TXT | Comma-separated vendor list |
| `products.2024-3094.cve` | TXT | Comma-separated product list |
| `refs.2024-3094.cve` | TXT | Reference URLs (one per RR or chunked per spec) |
| `kev.2024-3094.cve` | TXT | `yes` / `no`, plus KEV metadata when `yes` |
| `epss.2024-3094.cve` | TXT | EPSS score and percentile |

Index records at the zone apex:

| Name | Content |
|------|---------|
| `count.cve` | Total CVE count |
| `count.year-2024.cve` | Per-year count |
| `count.severity-critical.cve` | Per-severity count |
| `count.kev.cve` | KEV total |
| `latest.cve` | Most recent CVE ID |
| `built.cve` | Build timestamp |
| `source.cve` | Source git commit SHA |
| `version.cve` | cvedex version (also via `version.bind` CHAOS/TXT) |
| `ds.cve` | DNSSEC DS record for parent-zone delegation |

Names that don't match the `YYYY-NNNNN` pattern and aren't in the reserved label set above return NXDOMAIN. Existing names with no record of the queried type return NOERROR with empty answer plus SOA.

### Edge cases for CVE record handling

- **REJECTED / RESERVED CVEs** are still served. `state` reflects the actual state. `desc` may be a rejection reason or empty.
- **Multiple CVSS versions present** — the canonical `cvss.*` / `severity.*` / `vector.*` records use the highest version available. All variants are exposed in the JSON gateway.
- **Multiple CVSS scores from multiple CNAs (assigner + ADP)** — pick the CNA-of-record's score for the canonical records; expose the rest in the JSON gateway.
- **Empty fields** return TXT with empty string, not NXDOMAIN — the name is structurally valid.
- **Multilingual descriptions** — prefer English (`lang=en`); fall back to the first available.
- **Malformed individual records** — skip that record, increment a "malformed records" metric, continue the build.
- **TXT content sanitization** — free-text CVE fields are sanitized before insertion into TXT records: non-printable bytes stripped, control characters removed, values chunked per RFC 1035 when over 255 bytes per string. The full unsanitized text remains available via the JSON gateway.

### Refresh behavior

- Refresh runs on a configurable schedule (default daily, with jitter) and on `SIGHUP`.
- Each refresh: fetch the latest published data file from the configured distho want to plug the data into their own DNS server.
- **Generate-data**: one-shot. Processes upstream sources directly (cvelistV5 git repo, CISA KEV, FIRST.org EPSS), writes the canonical data file to disk, exits. Used by the data-generation workflow and by operators who self-host the pipeline. This is the only mode that touches upstream sources directly; it requires git tooling and network access to GitHub, CISA, and FIRST.org.
- **Sync-only**: one-shot. Fetches the latest data file, builds the zone, exits without serving. Useful for cron-driven workflows that control timing externally.

### DNS-server-specific behavior

cvedex behaves as a standard RFC-1035 authoritative server: UDP+TCP on port 53, EDNS(0) for larger UDP responses, proper SOA at the zone apex (serial = Unix timestamp of last successful build), configurable NS records, TC bit + TCP retry per spec, NXDOMAIN with SOA in authority for nonexistent names, NOERROR with empty answer + SOA for existing names with no record of the queried type, RFC 8482 minimal response for ANY queries.

DNSSEC is enabled by default. Default mode is offline signing at build time (every successful rebuild produces a freshly-signed zone). NSEC3 is used for authenticated denial with opt-out disabled. Default algorithm is ECDSA P-256. KSK and ZSK live in `{data_dir}` and are readable only by the service account. Operators can enable online signing if they need fresher signatures.

AXFR/IXFR is supported and ACL-gated (default localhost only). NOTIFY is sent to configured secondaries on every successful rebuild.

Response Rate Limiting (RRL) is on by default with standard token-bucket parameters.

Out-of-zone queries are REFUSED by default. Optional forward-to-upstream is available for operators who want a single binary acting as both authoritative-for-CVE and recursive-for-everything-else.

### High availability and clustering

cvedex is a clusterable DNS server (per the template's app-type matrix: cluster nodes ✓, managed nodes ✗, HA ✓). Each instance is independent and stateless beyond `{data_dir}` plus DNSSEC keys. Two deployment patterns are supported:

- **Independent instances**: each instance refreshes from cvelistV5 directly on its own schedule. Simplest topology. Slight zone-content drift between instances is acceptable for this use case (the data is the same, the SOA serials may differ by minutes).
- **Primary / secondary**: one instance acts as primary, refreshes from cvelistV5, sends NOTIFY to secondaries; secondaries pull via AXFR/IXFR. Zone-consistent across instances. DNSSEC keys must be shared across primaries if running active-active primaries (operator concern).

### Documentation that ships with the project

- README and operator guide (template-handled location and structure)
- A `dig` cookbook with copy-paste queries for every record type
- A `curl` cookbook for the JSON gateway
- An operator guide covering both deployment modes, DNSSEC key generation and rotation, parent-zone delegation for delegated-subdomain mode, trust-anchor distribution for custom-TLD mode, and AXFR secondary setup
- Operator notes on configuring the data-file URL, mirroring or self-hosting the data file, and self-hosting the data-generation workflow
- Sample configuration snippets for slaving the cvedex zone into BIND, Knot, and PowerDNS as secondaries

