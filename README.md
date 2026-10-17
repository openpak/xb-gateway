# xb-gateway

Original Xbox security gateway (SG): terminates the native Xbox secure transport on UDP 3074, validates xb-auth tickets (shared keyring), admits users through xb-account (mTLS) and dispatches service ids to feature backends (`xb-social`, `xb-tou`, ...) over the `xbfeature` gRPC stream. Runs on the secondary host (host network). `XBGW_BIND` (local interface) and `XBGW_ADVERTISE_IP` (public IP baked into tickets) are separate on purpose. Health: `XBGW_HEALTH_ADDR` (loopback) serves `/healthz`, `/readyz`, `/metricsz`; `xb-gateway healthcheck` is the listener-liveness probe. Configuration: `deploy/xb-gateway.env.example`; unit example `deploy/openpak-xb-gateway.container`. The `sgwire` package (framing/crypto of the SG protocol) is public so xb-e2e can drive a console-side client.

## Repository

- Module `github.com/openpak/xb-gateway`, image `ghcr.io/openpak/xb-gateway` (static, scratch, uid 10001).
- Shared packages: `github.com/openpak/xb-common` (private; see its README for access). Service map: workspace docs `xbox/native-service-architecture.md`.
- Build and test: `GOWORK=off go build ./... && go vet ./... && go test -race ./...` (with `GOPRIVATE=github.com/openpak/*` and `gh auth setup-git`).
- Image: `podman build --secret id=netrc,src=$HOME/.netrc -t ghcr.io/openpak/xb-gateway:dev .`
- Release: push a `vX.Y.Z` tag; `.github/workflows/release.yml` (tags and manual dispatch only) tests and publishes the multi-arch image. CI reads xb-common with the job's `GITHUB_TOKEN` through a netrc build secret; nothing secret is in git or in the image.
- `deploy/` holds the env example and a Podman quadlet example (digests in it are examples; production pins its own).
- Provenance and licence: AGPL-3.0-only, see `LICENSE` and `NOTICE`. Derived from the Xombie project in part; the legacy Rust workspace remains in https://github.com/openpak/xbox-live.
