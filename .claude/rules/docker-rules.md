# Docker Rules (PART 27)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Put Dockerfile at the project root — it belongs in docker/Dockerfile
- ❌ Put docker-compose.yml at the project root — it belongs in docker/
- ❌ Use .env files for container config — hardcode sane defaults in docker-compose
- ❌ Use EXPOSE 8080 or any non-80 port as default inside container
- ❌ Omit tini — it is required as PID 1
- ❌ Omit tor package — required for optional Tor hidden service
- ❌ Use CMD instead of ENTRYPOINT
- ❌ Skip STOPSIGNAL directive
- ❌ Put runtime data in the image — use volumes for {data_dir} and {config_dir}
- ❌ Map to 0.0.0.0 in production — use 172.17.0.1:{randomport}:80

## CRITICAL - ALWAYS DO
- ✅ Dockerfile at docker/Dockerfile (multi-stage)
- ✅ docker-compose.yml in docker/
- ✅ Builder stage: casjaysdev/go:latest
- ✅ Runtime stage: alpine:latest
- ✅ ENTRYPOINT: ["tini", "-p", "SIGTERM", "--", "/usr/local/bin/entrypoint.sh"]
- ✅ STOPSIGNAL SIGRTMIN+3
- ✅ Default port 80 inside container
- ✅ Default TZ: America/New_York
- ✅ docker/rootfs/ overlay copied into image at build time
- ✅ Production port mapping: -p 172.17.0.1:{randomport}:80

## DOCKERFILE STRUCTURE
```dockerfile
# syntax=docker/dockerfile:1
FROM casjaysdev/go:latest AS builder
# ... build steps ...

FROM alpine:latest
RUN apk add --no-cache git curl bash tini tor
COPY --from=builder /app/cvedex /usr/local/bin/cvedex
COPY docker/rootfs/ /
STOPSIGNAL SIGRTMIN+3
ENTRYPOINT ["tini", "-p", "SIGTERM", "--", "/usr/local/bin/entrypoint.sh"]
```

## REQUIRED PACKAGES (runtime image)
| Package | Reason |
|---------|--------|
| git | Version metadata, update checks |
| curl | Health checks, update downloads |
| bash | entrypoint.sh shell |
| tini | PID 1 signal handling |
| tor | Optional Tor hidden service |

## DOCKER ROOTFS OVERLAY
- Location: docker/rootfs/
- Copied into image at build time (COPY docker/rootfs/ /)
- Contains: entrypoint.sh, any runtime config defaults

## DOCKER COMPOSE (docker/docker-compose.yml)
```yaml
services:
  cvedex:
    image: ghcr.io/casapps/cvedex:latest
    ports:
      - "172.17.0.1:8080:80"
    environment:
      TZ: America/New_York
    volumes:
      - cvedex-data:/var/lib/casapps/cvedex
      - cvedex-config:/etc/casapps/cvedex
volumes:
  cvedex-data:
  cvedex-config:
```

No .env files. All configuration via environment variables with sane defaults baked in.

---
For complete details, see AI.md PART 27
