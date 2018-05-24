# ethstats

Go implementation of an ethstats collection server.

## Endpoints

- `/api` is a WebSocket endpoint for collecting ethstats.
- `/` returns a JSON response with the current connected nodes. Details can be
    fetched by specifying the node ID with the `node` param.

## Quickstart

```
$ make run
2018/05/24 16:34:18 listening on :8080
```

```
$ geth --ethstats "somenodeid:somesecret@localhost:8080"
```

```
$ curl "http://localhost:8080/"
{"nodes":["somenodeid"]}
$ curl "http://localhost:8080/?node=somenodeid"
{"id":"somenodeid","info":{"name":"somenodeid","node":"Geth/testmooch/v1.8.3-unstable/linux-amd64/go1.10","port":30303,"net":"1","protocol":"les/2","api":"No","os":"linux","os_v":"amd64","client":"0.1.1","canUpdateHistory":true},"latency":"","block":{"number":null,"hash":"0x0000000000000000000000000000000000000000000000000000000000000000","parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000","timestamp":null,"miner":"0x0000000000000000000000000000000000000000","gasUsed":0,"gasLimit":0,"difficulty":"","totalDifficulty":"","transactions":null,"transactionsRoot":"0x0000000000000000000000000000000000000000000000000000000000000000","stateRoot":"0x0000000000000000000000000000000000000000000000000000000000000000","uncles":[]},"pending":{"pending":0},"status":{"active":false,"syncing":false,"mining":false,"hashrate":0,"peers":0,"gasPrice":0,"uptime":0},"last_seen":"2018-05-24T16:35:35.116948798-04:00"}
```

## Appendix

### Challenges with the netstats/ethstats protocol

- Undocumented. Must be reverse engineered from existing implementations (only 2 at present, linked in appendix).
- No reusable libraries/schemas/SDKs.
- Inconsistent implementations.
  (Example: Node version sends `clientTime` timestamps in a different format than the Geth built-in implementation.)
- Unused components of protocol.
  (Example: Node and Geth implementations ignore the `clientTime` field from the response, which is probably for the best since they're incompatible.)
- Inconsistent containers. Some responses have an extra redundant object container, others do not.
  (Example: `hello`, `node-ping`, `latency` are contained immediately in the payload, while `block` is contained under another `{"block": {payload}}` layer, and both `pending` and `stats` are contained under `{"stats": {payload}}`)
- Lacking node authentication (only has authorization). Would be nice if the auth handshake included a signed message from the enodeID.
- Non-standard framework-specific websocket payload format (`{"emit": ["topic", {payload}}`). Non-homogeneous array types are unnecessarily frustrating to work with.
- No ability for the server to throttle the rate from the clients.
- Lacking metrics about peers (only peer count).
- Lacking metrics about the server/runtime.


### Ethstats v2 Protocol Wishlist

- Clearly defined request/response schemas which are easy to use across programming languages. Perhaps GRPC or plain Protobufs? The goal should be to get native support from every major Ethereum client.
- Connections identified and authenticated by EnodeID (challenge signed by Enode private key?).
- Support for sharing the full peer list. This is useful for validating bi-direcitonal peer serving claims (such as for vipnode).
- Support for sharing node runtime metrics. This is useful for maintainers of large fleets of nodes, and debugging platform-specific performance quirks.
- Designed to be easily integrated with mainstream timeseries tooling like Prometheus/InfluxDB.
- Support for relaying signed ethstats reports.


### References

- https://github.com/cubedro/eth-net-intelligence-api/blob/bdc192ebb76fc9964ef0da83ee88bc86ba69c052/lib/node.js
- https://github.com/ethereum/go-ethereum/blob/6286c255f16a914b39ffd3389cba154a53e66a13/ethstats/ethstats.go

## License

MIT.
