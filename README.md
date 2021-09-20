# Watch Tower

Watch tower is atomex module which sends client's redeem  or refund transactions.

## Usage

To run watch tower you should execute command

```bash
watch_tower -c config.yml
```

Watch tower has 1 argument:

* `c` - config file name.


## Private keys

To set private keys you should pass environment variables to docker container or set it on OS. `ETHEREUM_PRIVATE` - for ethereum wallet and `TEZOS_PRIVATE` - for tezos wallet.
By default `.env` file is expected in docker-compose.


### Build

You can build binary from source

```bash
cd cmd/watch_tower && go build .
```

or by docker

```bash
docker-compose up -d --build watch_tower

# or

make up
```

### Config

Example:

```yaml
tezos:
  node: https://rpc.tzkt.io/mainnet
  tzkt: https://api.tzkt.io
  contract: KT1VG2WtYdSWz5E7chTeAdDPZNy2MpP8pTfL
  tokens:
    - KT1EpQVwqLGSH7vMCWKJnq6Uxi851sEDbhWL
    - KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr
ethereum:
  node: https://main-light.eth.linkpool.io/
  wss: wss://main-light.eth.linkpool.io/ws
  eth_address: 0xE9C251cbB4881f9e056e40135E7d3EA9A7d037df
  erc20_address: 0xAc5881D77Db9340c94F53300736a9cf9e61fA25E
  user_address: YOUR_WALLET_ADDRESS_HERE
restore: true
types:
  - redeem
  - refund
```

* `types` - following for `redeem` or/and `refund` (both by default)
* `restore` - flag which is set for finding passed swaps without redeem or refund (*false* by default)
* `retry_count_on_failed_tx` - retry count if transaction sending was failed (*0* by default)
* `tezos` - tezos settings. Consists of:

    * `node` - URL to tezos node RPC (**required**)
    * `tzkt` - URL to TzKT API (**required**)
    * `contract` - address of main Atomex contract in Tezos (**required**)
    * `tokens` - list of Atomex addresses of token contracts (**required**)
    * `min_payoff` - minimal pay off for processing in microtez (*0* by default.  Type: **string**)

* `ethereum` - ethereum settings. Consists of:

    * `node` - URL to ethereum node RPC (**required**)
    * `tzkt` - URL to ethereum websocket (**required**)
    * `eth_address` - address of ETH Atomex contract in ethereum (**required**)
    * `erc20_address` - address of ERC20 Atomex contract in ethereum (**required**)
    * `user_address` - user wallet address (**required**)
    * `min_payoff` - minimal pay off for processing in wei (*0* by default. Type: **string**)


### Architecture

  Logically program consists of 2 parts: blockchain listeners and watch tower.
  Blockchain listeners get new transactions, prepare it and send events by channels.
  Watch tower subscribes on blockchain listeners channels and processes its events.
  After processing watch tower send to listener Redeem or Refund command.