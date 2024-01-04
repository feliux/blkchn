# blkchn

Blockchain project

**Features**

- Proof of Work: reward a miner if the new block hash starts by a MINING_DIFFICULTY times 0, such as '000ewfbweg5w1g61we1'.
- Consensus algorithm: the longest chain rule.
- Addresses: version 1 Bitcoin addresses.

**TODO**

Firsts goes first...

1. ~~First it must works~~
    - ~~Blocks, transactions, hashing method, nodes, etc~~
    - ~~Mining logic (PoW)~~
    - ~~Wallets, blockchain addresses, transaction verification~~
    - ~~APIs~~
    - ~~Structure blockchain network: nodes, wallets and consensus algo~~
2. Then make it better
    - Build a CLI
    - Make structs private
    - Checks "to review" comments
    - Docker, k8s, etc
    - CI/CD
    - Implement ICO for wallet creation
    - Implement fees.
    - Improve logging and error handling
        - Not decrease wallet amount when a transaction fails
    - Improve APIs with handlers and routers
    - Save blockchain state (despite cache)
3. Then make it faster
    - Synchronizing transactions with P2P

## Usage

### Serve Blockchain

```bash
make server-build
make wallet-build

# 3 nodes
make server-run ARGS="-port=5000"
make server-run ARGS="-port=5001"
make server-run ARGS="-port=5002"

# 2 wallets
make wallet-run ARGS="-port 8080 -gateway=http://127.0.0.1:5001"
make wallet-run ARGS="-port 8081 -gateway=http://127.0.0.1:5002"
```

## References

[Technical background of version 1 Bitcoin addresses](https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses)

1. Creating ECDSA private key (32 bytes) and public key (64 bytes).
2. Perform sha256 hashing on the public key -> 32 bytes
3. Perform RIPEMD-160 hashing on the step2 result -> 20 bytes
4. add version byte in front of RIPEMD-160 hash -> 0x00 for Main Network
5. Perform sha256 hash on the extended RIPEMD-160 result
6. Perform sha256 hash on the result of the previous step5 sha256 hash
7. Take the first 4 bytes of the second step6 sha256 hash for checksum
8. Add the 4 checksum bytes from step7 at the end of extended RIPEMD-160 hash from step4 -> 25 bytes
9. Convert the result from a byte string into base58
