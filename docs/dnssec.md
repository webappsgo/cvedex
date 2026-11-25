# DNSSEC

When DNSSEC is enabled (the default), the server creates and persists ECDSA
P-256 KSK and ZSK key pairs under `DATA_DIR/security/dnssec/`, signs each
verified zone snapshot offline, and serves DNSKEY, RRSIG, NSEC3, and denial
proof records. The signed BIND-format snapshot is available from
`/api/v1/zone`; `/api/v1/ds` returns the DS record for the active KSK. Both
endpoints remain unavailable until a verified dataset has been loaded.

Publish the returned DS at the parent zone before enabling validating
delegation. Use the administration settings page to rotate keys; rotation
publishes the new signed zone before replacing the active keys. Keep the data
directory persistent and protected because it contains the private keys.
