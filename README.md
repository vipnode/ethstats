# ethstats

Go implementation of an ethstats collection server.


## Challenges with the netstats/ethstats protocol

- Undocumented. Must be reverse engineered from existing implementations (only 2 at present, linked in appendix).
- No reusable libraries/schemas/SDKs.
- Inconsistent implementations.
  (Example: Node version sends `clientTime` timestamps in a different format than the Geth built-in implementation.)
- Unused components of protocol.
  (Example: Node and Geth implementations ignore the `clientTime` field from the response, which is probably for the best since they're incompatible.)
- Lacking node authentication (only has authorization). Would be nice if the auth handshake included a signed message from the enodeID.
- Non-standard framework-specific websocket payload format (`{"emit": ["topic", {payload}}`). Non-homogeneous array types are unnecessarily frustrating to work with.
- No ability for the server to throttle the rate from the clients.
- Lacking metrics about peers (only peer count).
- Lacking metrics about the server/runtime.


## Ethstats v2 Protocol Wishlist

- Clearly defined request/response schemas which are easy to use across programming languages. Perhaps GRPC or plain Protobufs? The goal should be to get native support from every major Ethereum client.
- Connections identified and authenticated by EnodeID (challenge signed by Enode private key?).
- Support for sharing the full peer list. This is useful for validating bi-direcitonal peer serving claims (such as for vipnode).
- Support for sharing node runtime metrics. This is useful for maintainers of large fleets of nodes, and debugging platform-specific performance quirks.
- Designed to be easily integrated with mainstream timeseries tooling like Prometheus/InfluxDB.


## Appendix

### References

- https://github.com/cubedro/eth-net-intelligence-api/blob/bdc192ebb76fc9964ef0da83ee88bc86ba69c052/lib/node.js
- https://github.com/ethereum/go-ethereum/blob/6286c255f16a914b39ffd3389cba154a53e66a13/ethstats/ethstats.go

## License

MIT.
