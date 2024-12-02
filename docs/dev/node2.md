# Node 2 Orientation

## Repo Layout

```text
kwil-db
├── app                     shared components of CLI apps
│   ├── custom
│   ├── key
│   ├── node
│   │   └── conf
│   ├── setup
│   └── shared
│       ├── bind
│       ├── display
│       ├── generate
│       └── version
├── cmd                     CLI apps, composed from app pkgs
│   ├── kwil-cli
│   │   ├── client
│   │   ├── cmds
│   │   ├── config
│   │   ├── csv
│   │   ├── generate
│   │   └── helpers
│   └── kwild
│       ├── generate
│       └── internal
├── common                  context and data structs for extensions
├── config                  the main Config struct with toml tags
├── contrib                 helper scripts, docker files, non-code, etc.
│   ├── docker
│   │   └── compose
│   ├── scripts
│   │   ├── build
│   │   ├── kuneiform
│   │   ├── mods
│   │   └── publish
│   └── systemd
├── core                    top-level utils, structs, logger, and client types
│   ├── client
│   │   └── types
│   ├── crypto
│   │   └── auth
│   ├── gatewayclient
│   ├── log
│   ├── rpc
│   │   ├── client
│   │   ├── json
│   │   └── transport
│   ├── types
│   │   ├── admin
│   │   ├── decimal
│   │   ├── serialize
│   │   └── validation
│   └── utils
│       ├── json
│       ├── order
│       └── random
├── docs
│   ├── dev
│   ├── release-notes
│   └── sql
├── extensions              extension packages for binary customization
│   ├── auth
│   ├── consensus
│   ├── listeners
│   ├── precompiles
│   └── resolutions
├── node                    the node (p2p and all the dependencies inc. consensus)
│   ├── accounts
│   ├── admin
│   ├── consensus           the consensus engine (CE) where decisions are made
│   ├── engine              interpreter of kuneiform and user dataset SQL
│   │   ├── execution
│   │   ├── generate
│   │   ├── integration
│   │   └── testdata
│   ├── ident               end-user identity based core/crypto{,auth}
│   ├── mempool
│   ├── meta                chain metadata (current state)
│   ├── peers               node peer manager (PEX etc.)
│   ├── pg
│   ├── services
│   │   └── jsonrpc
│   ├── store               block store / index / tx index
│   ├── txapp
│   ├── types
│   │   └── sql
│   ├── utils
│   ├── versioning
│   └── voting              voting store a.k.a. validator+event store
├── parse
│   ├── gen
│   ├── grammar
│   ├── planner
│   ├── postgres
│   └── wasm
├── test
│   ├── acceptance
│   ├── integration
│   └── stress
├── testing                 kuneiform testing framework
└── version                 global (sem)versioning
```
