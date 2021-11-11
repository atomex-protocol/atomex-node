# Atomex Protocol Node

[![Build Status](https://github.com/atomex-protocol/watch-tower/workflows/CI/badge.svg)](https://github.com/atomex-protocol/watch-tower/actions?query=branch%3Amaster+workflow%3A%22build%22)
[![made_with golang](https://img.shields.io/badge/made_with-golang-blue.svg)](https://golang.org/)

It's a project which contains Golang realization of modules of atomex protocol:

* watch tower
* market maker

## Build

You can build binaries from source

```bash
cd cmd/watch_tower && go run . -c ../../configs/production

cd cmd/market_maker && go run . -c ../../configs/production
```

or by docker

```bash
docker-compose up -d --build watch_tower market_maker

# or

make up
```

## Usage

### Watch tower

To run watch tower you should execute command

```bash
watch_tower -c configs
```

Watch tower has 1 argument:

* `c` - path to directory which contains configuration files. Default: `configs`.


### Market Maker

To run market maker you should execute command

```bash
market_maker -c configs
```

Market maker has 1 argument:

* `c` - path to directory which contains configuration files. Default: `configs`.

## Configuration

All configuration files are YAML files which located in configuration directory. Default configuration directory is `./configs`. It's located in application folder. For example, you can find configuration directories [here](https://github.com/atomex-protocol/watch-tower/configs).

Configuration directory contains following files:

* `assets.yml` - file which contains using currencies
* `atomex.yml` - atomex settings
* `binance.yml` - binance settings
* `chains.yml` - file which contains all supported chains settings
* `market_maker.yml` - market maker settings
* `symbols.yml` - file which contains instruments descriptions
* `watch_tower.yml` - watch tower settings.

You can find full description of configuration files below. **WARNING!** Don't use configs from examples. It may contains invalid contract addresses.

### Assets

In `assets.yml` you can add, remove or edit using currencies. File structure is:

```yaml
asset_id: # asset ID must be unique and it will be used in other configuration files
    name: <asset name>
    chain: <tezos or ethereum>
    contract: <contract address of asset if exists>
    atomex_contract: <atomex contract which interacted with asset>
    decimals: <decimals count of asset>

# =============================================================
# For example
# =============================================================

USDT:
  name: USDT
  chain: ethereum
  contract: 0xdac17f958d2ee523a2206206994597c13d831ec7
  atomex_contract: 0xAc5881D77Db9340c94F53300736a9cf9e61fA25E
  decimals: 18
XTZ:
  name: XTZ
  chain: tezos
  atomex_contract: KT1GyzWoSh9A2ACr1wQpBoGHiDrjT4oHJT2J
  decimals: 6
```

### Symbols

In `symbols.yml` you can add, remove or edit using traded instruments. It contains array of instruments. File structure is:

```yaml
- name: <symbol name which will be used in application and other configuration files>
  base: <base asset ID from assets config>
  quote: <quote asset ID from assets config>

# =============================================================
# For example
# =============================================================

- name: XTZ_USDT
  base: XTZ
  quote: USDT
```

### Chains

In `chains.yml` you can edit supported chains settings. Now supported only `tezos` and `ethereum`. File structure is:

```yaml
tezos:
  node: <URL to tezos node RPC>
  tzkt: <URL to TzKT API>
  ttl: <Time-to-live of sent operations in blocks>

ethereum:
  node:  <URL to ethereum node RPC>
  wss:  <URL to ethereum node websocket>

# =============================================================
# For example
# =============================================================

tezos:
  node: https://rpc.tzkt.io/mainnet
  tzkt: https://api.tzkt.io
  ttl: 2

ethereum:
  node: https://main-light.eth.linkpool.io/
  wss: wss://main-light.eth.linkpool.io/ws
```

### Atomex
In `atomex.yml` you can edit base atomex settings. File structure is:

```yaml
rest_api: <URL to Atomex API>
wss: <URL to Atomex Websocket>

settings:
  reward_for_redeem: <value of reward for redeem>
  lock_time: <value of lock time>

to_symbols:
  <symbol ID from symbols config>: <symbol name in atomex>

from_symbols:
  <symbol name in atomex>: <symbol ID from symbols config>

# =============================================================
# For example
# =============================================================

rest_api: https://api.atomex.me/v1
wss: wss://ws.api.atomex.me/ws

settings:
  reward_for_redeem: 1000
  lock_time: 3600

to_symbols:
  XTZ_USDT: XTZ/USDT

from_symbols:
  XTZ/USDT: XTZ_USDT
```

### Third-party quote providers

Market maker use third-party quote provider for sending limits to atomex. Now only Binance is supported. But list of supported providers can be extended. Each provider has own config in separate file. File is named like provider kind. For example: `binance.yml`.

Symbol can be absent at quote provider. That's why synthetics concept is reailized. There are two synthetic types: direct and divided. Direct is an one-to-one synthetics. Divided is a synthetics which use third currency for exchange. For example, pair `XTZ_ETH` is not exists at Binance. You can realized it via combination of pairs `XTZ_USDT` and `ETH_USDT`. So, if you sell `XTZ` on `XTZ_ETH` at Atomex, you have to sell `XTZ_USDT` at Binance and buy `ETH` on `ETH_USDT`.

#### Binance

In `binance.yml` you can edit binance settings. File structure is:

```yaml
to_symbols:
  <binance symbol name>: <symbol ID from symbols config>
from_symbols:
 <symbol ID from symbols config>:
    type: <synthetic type (direct | divided)>
    symbols: <array of binnace symbols>

# =============================================================
# For example
# =============================================================

to_symbols:
  XTZUSDT: XTZ_USDT
  ETHUSDT: ETH_USDT
from_symbols:
  XTZ_USDT:
    type: direct
    symbols:
      - XTZUSDT
  XTZ_ETH:
    type: divided
    symbols:
      - XTZUSDT
      - ETHUSDT
```

### Watch tower

In `watch_tower.yml` you can edit watch tower settings. File structure is:

```yaml
restore: <flag which is set for finding passed swaps without redeem or refund (*false* by default)>
types: <following for `redeem` or/and `refund` (both by default)>
retry_count_on_failed_tx: <retry count if transaction sending was failed (*0* by default)>

# =============================================================
# For example
# =============================================================

restore: true
types:
  - redeem
  - refund
retry_count_on_failed_tx: 2
```

### Market maker

In `market_maker.yml` you can edit market maker settings. File structure is:

```yaml
quote_provider:
  kind: <kind of quote provider. Only `binance`supported>

keys:
  file: <file which contains key data>
  kind: <kind of key file. Only `custom`supported>
  generate_if_not_exists: <flag which is set for generating keys if not exists (*false* by default)>

strategies:  # array of strategies
  - kind: <volatility | follow | one-by-one>
    symbol: <symbol ID from symbols config>
    spread:
      ask: <minimal half-spread for ask in percents>
      bid: <minimal half-spread for bid in percents>
    volume: <tradable volume>
    dist:
      min: <minimal offset for volatility strategy>
      max: <maximum offset for volatility strategy>
    width:  <count of standard deviation for volatility strategy>
    window: <rolling window for volatility strategy>

log_level: trace

# =============================================================
# For example
# =============================================================

quote_provider:
  kind: binance

keys:
  file: api_test.json
  kind: custom
  generate_if_not_exists: true

strategies:  
  - kind: volatility
    symbol: XTZ_ETH
    spread:
      ask: 0.03
      bid: 0.05
    volume: 0.01
    dist:
      min: 0.001
      max: 0.006
    width: 1.5
    window: 60

log_level: trace

```