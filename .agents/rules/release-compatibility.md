# Release Compatibility

Applies to any change visible to a plugin or the modality host.

## Handshake

- `ProtocolVersion` is `1` in the current contract.
- `MagicCookieKey` is `GOPN_PLUGIN`; `MagicCookieValue` is
  `gopn-plugin-handshake`, matching modality `library/gopn`.
- A protocol bump requires a matching change and negotiation path in modality
  `library/gopn`, plus a dated note.

## Service Registry

- Business services register through `Service.Register`, a
  `func(grpc.ServiceRegistrar)`, and must remain compatible with
  `protoc-gen-go-grpc` output.
- A breaking RPC signature requires a new service name or protocol version.

## Module

- Do not change the module path without recording a migration decision and
  coordinating the host documentation.
- Upgrade `hashicorp/go-plugin` and `google.golang.org/grpc` together with a
  handshake and broker test.
