# CPU Bitcoin Miner in Go


### Disclaimer: this is only a prototype of a bitcoin miner. Do not use it for official mining.

## Synopsis

Implementation of a Bitcoin mining client in Go to take advantage of the built-in concurrency primitives of the language. 
This implementation integrates _getblocktemplate_ used for pool mining instead of _getwork_ as it is getting the standard with BIP 0023. 

## Motivation

This project is for **educational** purpose at start for a deep and better understanding of Bitcoin _under the hood_ as well as the discovery of Golang.

## Install & Run
After cloning this repo, get the information that you need from your BitCoin Core Client and set your own configuration in config.json:
```
{
    "user" : "user_name",
    "password" : "Pa$$w0rd",
    "host" : "127.0.0.1:8332",
    "account" : "MyBTCAccount",
    "log": {
        "Activated": true,
        "Level": "info",
        "File": ""
    }
}

```
Just start the program to begin minning:
```
go run src/main.go
```

## Implementation

The interactions with the Bitcoin Core client are currently done with curl calls, that have to be re-written with net/http or with a websocket implementation to maintain long-help connections. 

The mining side is implemented with a Job/Worker/Dispatcher pattern (see reference).

A first step is to send the raw data from bitcoin client to the dispatcher which is formating and spliting it into 'Chunks'. Those Chunks are sent to the miners. When a successful block is found, it is sent back to the Dispatcher for verification and back to the WebSocket for submitting to the network (last part to be done). 

The overall Bitcoin implementation, but especially the Block Header, tries to stay as close as possible to the reference.

## Performance

The current version is able to perform around 350.000 operations (double sha256 hashs) per second on a 2.4 GHz Intel Core i5.

## Dependencies

None, beside your Bitcoin client. 
The btcsuite/btcrpcclient library was used for the JSON-RPC call at the beginning but was recently removed due to its deprecated mining implementation (BIP 0023 again).

## Contributors

- Mathieu Mailhos (@mathieumailhos)

## References
- [Start with Bitcoin and download your client](https://bitcoin.org/en/)
- [BIP 0023](https://en.bitcoin.it/wiki/BIP_0023)
- [Block Hashing Algorithm](https://en.bitcoin.it/wiki/Block_hashing_algorithm)
- [Thread Pool in Golang](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)

