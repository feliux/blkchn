# blkchn

## Usage

### Init

```bash
go mod edit -replace github.com/feliux/block=./block
go mod edit -replace github.com/feliux/blockchain=./blockchain
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
