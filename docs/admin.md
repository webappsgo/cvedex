# Administration

The server provides first-run administrator setup at `/server/administration/setup`, sign-in at `/server/auth/login`, a protected dashboard at `/server/administration`, and POST sign-out. The default dashboard path can be changed with `ADMIN_PATH`. Runtime configuration remains environment and YAML based as described in [Configuration](configuration.md).

On first start, the server prints a one-time setup token and a one-time operator token to standard output. The setup token is stored as a hash in the private config directory. The operator token is stored in `secrets.yml` with mode `0600`; only its SHA-256 digest is held in memory for request authentication. Create the first administrator before exposing the panel. Passwords use Argon2id; browser sessions use HttpOnly, SameSite=Strict cookies with absolute and idle expiry; login failures are rate-limited. Setup, login, and sign-out forms use CSRF tokens.

The dashboard links to `/server/{admin_path}/settings`, where an administrator
can change the refresh interval, upstream redirect template, AXFR/IXFR client
ACL, NOTIFY targets, DNSSEC signing, and DNS response rate limiting. Settings
are validated and persisted in `server.yml`; the running service applies
supported changes immediately. The page also provides DNSSEC key rotation.

The matching JSON API is `/api/v1/server/{admin_path}/settings` (GET and PATCH),
with POST `/api/v1/server/{admin_path}/security/dnssec/rotate` for key rotation.
All admin API requests require `Authorization: Bearer <operator-token>`; browser
session cookies are not API credentials. Bearer-authenticated API requests do
not use CSRF tokens.

The admin security page at `/server/{admin_path}/security` shows the operator
token masked. Regeneration requires the current administrator password, writes
the replacement to `secrets.yml`, invalidates the previous token immediately,
and records `security.operator_token_regenerated` in `LOG_DIR/audit.log`. The
replacement is displayed once and is not returned by an API.
