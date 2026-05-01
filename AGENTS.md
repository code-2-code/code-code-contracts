# Agent Rules

- This repository owns public protobuf contracts and generated Go/TypeScript contract bindings.
- Change protobuf source in `packages/proto` before generated output.
- Do not hand-edit generated files under `packages/go-contract` or `packages/agent-contract/src/gen`.
- Do not add runtime, UI, Helm, deployment, or service implementation code here.
- When a contract changes, regenerate all current outputs and verify the generated packages.
