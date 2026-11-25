# Administration

The server provides first-run administrator setup at `/server/administration/setup`, sign-in at `/server/auth/login`, a protected dashboard at `/server/administration`, and POST sign-out. The default dashboard path can be changed with `ADMIN_PATH`. Runtime configuration remains environment and YAML based as described in [Configuration](configuration.md).

On first start, the server prints a one-time setup token to standard output and stores only its hash in the private config directory. Create the first administrator before exposing the panel. Passwords use Argon2id; sessions use HttpOnly, SameSite=Strict cookies with absolute and idle expiry; login failures are rate-limited. Setup, login, and sign-out forms use CSRF tokens.

The dashboard links to `/server/{admin_path}/settings`, where an administrator
can change the refresh interval, upstream redirect template, AXFR/IXFR client
ACL, NOTIFY targets, DNSSEC signing, and DNS response rate limiting. Settings
are validated and persisted in `server.yml`; the running service applies
supported changes immediately. The page also provides DNSSEC key rotation.

The matching JSON API is `/api/v1/server/{admin_path}/settings` (GET and PATCH),
with POST `/api/v1/server/{admin_path}/security/dnssec/rotate` for key rotation.
These endpoints currently require an authenticated administrator session.
