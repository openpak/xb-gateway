# Xbox Gateway Server

OpenPak's Xbox Gateway Server.

| | |
| --- | --- |
| Network | Original Xbox |
| Role | Security gateway: terminates the native Xbox secure transport, validates xb-auth tickets, admits users through xb-account and dispatches service ids to feature backends |
| Depends on | [`xb-auth`](https://github.com/openpak/xb-auth) (ticket keyring), [`xb-account`](https://github.com/openpak/xb-account) (mTLS), feature backends such as `xb-social` and `xb-tou` (`xbfeature` gRPC stream), [`xb-common`](https://github.com/openpak/xb-common) (private Go module) |
| Image | `ghcr.io/openpak/xb-gateway` |
| Based on | [Xombie](https://github.com/XombieOnline/xombie) security gateway (AGPL-3.0) |
| Licence | AGPL-3.0-only |

## What it serves

| Endpoint | What |
| --- | --- |
| UDP 3074 | Native Xbox secure transport |
| `xbfeature` gRPC stream | Service-id dispatch to feature backends |

## Run

Runs as a container on the secondary host with host networking. It is configured from `deploy/xb-gateway.env.example`; the quadlet example is `deploy/openpak-xb-gateway.container`. Details are in [deploy.md](docs/deploy.md).

## Build

```sh
GOWORK=off go build ./... && go vet ./... && go test -race ./...
```

Needs `GOPRIVATE=github.com/openpak/*` and `gh auth setup-git`, because `xb-common` is a private module.

## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)
- [Notes](docs/notes.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
