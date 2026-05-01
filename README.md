# code-code-contracts

Contract source and generated language bindings for the Code Code platform.

This repository owns:

- `packages/proto`: canonical protobuf contracts.
- `packages/go-contract`: generated Go contracts and small contract helpers.
- `packages/agent-contract`: generated TypeScript contracts.

Contract changes start here. Downstream platform, console, and deploy
repositories should consume released contract versions instead of editing
generated files directly.

Useful checks:

```bash
cd packages/proto && buf lint
cd packages/proto && buf generate
cd packages/go-contract && go test ./...
cd packages/agent-contract && pnpm install && pnpm build
```
