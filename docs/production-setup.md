# Permissionless node setup

This guide describes how to run a permissionless node that:

- Synchronizes the network
- Expose a JSON RPC interface, acting as an archive node

Note that sequencing and proving functionalities are not covered in this document.

## Requirements

- A machine to run the zkEVM permissionless node with the following requirements:
  - Hardware: 32GB RAM, 4 cores, 256GB Disk with high IOPS (as the network is super young the current disk requirements are quite low, but they will increase overtime. Also note that this requirement is true if the DBs run on the same machine, but it's recommended to run Postgres on dedicated infra). Currently ARM-based CPUs are not supported
  - Software: Ubuntu 22.04, Docker
- A L1 node: we recommend using geth, but what it's actually needed is access to a JSON RPC interface for the L1 network (Sepolia for Cardona zkEVM testnet, Goerli for zkEVM testnet, Ethereum mainnet for zkEVM mainnet)

## Setup

This is the most straightforward path to run a zkEVM node, and it's perfectly fine for most use cases, however if you are interested in providing service to many users it's recommended to do some tweaking over the default configuration. Furthermore, this is quite opinionated, feel free to run this software in a different way, for instance it's not needed to use Docker, you could use the Go and C++ binaries directly.

Follow next steps to deploy a permissionless node:

1. Define network `mainnet`, `testnet` or `cardona` (e.g `mainnet`):
```bash
ZKEVM_NET=mainnet
```
2. Create and define installation directory (e.g. `~/zkevm-node`):
```bash
mkdir ~/zkevm-node
ZKEVM_DIR=~/zkevm-node
```
3. Create and define config directory (e.g. `~/zkevm-node/config`):
```bash
mkdir ~/zkevm-node/config
ZKEVM_CONFIG_DIR=~/zkevm-node/config
```

4. In the following [link](https://github.com/0xPolygonHermez) check the node version for the network you are deploying and set the version in the environment variable (e.g. v0.5.0):

```bash
ZKEVM_VERSION=v0.5.0
```
> **NOTE:** It's recommended to source this environment variables in your `~/.bashrc`, `~/.zshrc` or whatever you're using

5. Download and extract the artifacts: 
```bash
curl -L https://github.com/0xPolygonHermez/zkevm-node/releases/download/$ZKEVM_VERSION/$ZKEVM_NET.zip > $ZKEVM_NET.zip 
unzip -o $ZKEVM_NET.zip -d $ZKEVM_DIR 
rm $ZKEVM_NET.zip 
```

6. Copy the file with the environment variables into config directory:
```bash
cp $ZKEVM_DIR/$ZKEVM_NET/example.env $ZKEVM_CONFIG_DIR/.env
```
7. Edit the `.env` file with your favourite editor (e.g nano) and set the variables:
```bash
nano $ZKEVM_CONFIG_DIR/.env
```
> **NOTE:** With the configuration done in step 7 it is enough to run the permissionless node with the default config parameters. If you want to customize some config parameters do the following steps:  
> 1. Copy the config files into the config directory: `cp $ZKEVM_DIR/$ZKEVM_NET/config/environments/$ZKEVM_NET/* $ZKEVM_CONFIG_DIR/`
> 2. Make sure the modify the `ZKEVM_ADVANCED_CONFIG_DIR` variable from `$ZKEVM_CONFIG_DIR/.env` with the correct path (same value as `$ZKEVM_CONFIG_DIR`)
> 3. Edit the different configuration files in the `$ZKEVM_CONFIG_DIR` directory and make the necessary changes
8. Run the node (you may need to run this command using `sudo` depending on your Docker setup): 
```bash
docker compose --env-file $ZKEVM_CONFIG_DIR/.env -f $ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml up -d
```
9. Make sure that all components are running:
```bash
docker compose --env-file $ZKEVM_CONFIG_DIR/.env -f $ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml ps
```
10. You should see the following containers:
   - zkevm-rpc
   - zkevm-sync
   - zkevm-state-db
   - zkevm-pool-db
   - zkevm-prover
11. If everything has gone as expected you should be able to run queries to the JSON RPC at `http://localhost:8545`. For instance you can run the following query that fetches the latest synchronized L2 block, if you call this every few seconds, you should see the number increasing:
```bash
curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' http://localhost:8545
```

## Troubleshooting

- It's possible that the machine you're using already uses some of the necessary ports. In this case you can change them directly at `$ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml`
- If one or more containers are crashing please check the logs using:
```bash
docker compose --env-file $ZKEVM_CONFIG_DIR/.env -f $ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml logs <cointainer_name>
```
## Stop
You can stop all the containers using:
```bash
docker compose --env-file $ZKEVM_CONFIG_DIR/.env -f $ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml down
```

## Updating

In order to update the software, you have to repeat the steps of the setup, but taking care of not overriding the config that you have modified. Basically, instead of running `cp $ZKEVM_DIR/$ZKEVM_NET/example.env $ZKEVM_CONFIG_DIR/.env`, check if the variables of `$ZKEVM_DIR/$ZKEVM_NET/example.env` have been renamed or there are new ones, and update `$ZKEVM_CONFIG_DIR/.env` accordingly.

## Advanced setup

> DISCLAIMER: right now this part of the documentation attempts to give ideas on how to improve the setup for better performance, but is far from being a detailed guide on how to achieve this. Please open issues requesting more details if you don't understand how to achieve something. We will keep improving this doc for sure!

There are some fundamental changes that can be done towards the basic setup, in order to get better performance and scale better.

### Database

In the basic setup, there are Postgres being instanciated as Docker containers. For better performance is recommended to:

- Run dedicated instances for Postgres. To achieve this you will need to:
  - Remove the Postgres services (`zkevm-pool-db` and `zkevm-state-db`) from the `docker-compose.yml`
  - Instantiate Postgres elsewhere (note that you will have to create credentials and run some queries to make this work, following the config files and docker-compose should give a clear idea of what to do)
  - Update the `node.config.toml` to use the correct URI for both DBs
  - Update `prover.config.json` to use the correct URI for the state DB
- Use a setup of Postgres that allows to have separated endpoints for read / write replicas

### JSON RPC

Unlike the synchronizer, that needs to have only one instance running (having more than one synchronizer running at the same time connected to the same DB can be fatal), the JSON RPC can scale horizontally.

There can be as many instances of it as needed, but in order to not introduce other bottlenecks, it's important to consider the following:

- Read replicas of the State DB should be used for the JSON RPCs instances
- Synchronizer should have an exclusive instance of `zkevm-prover`
- JSON RPCs should scale in correlation with instances of `zkevm-prover`. The most obvious way to do so is by having a dedicated `zkevm-prover` for each `zkevm-rpc` instance. But depending on the payload of your solution it could be worth to have `many zkevm-rpc : 1 zkevm-prover`
