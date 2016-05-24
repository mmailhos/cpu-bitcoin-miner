# CPU Bitcoin Miner in Go

## Synopsis

Implementation of a Bitcoin mining client in Go to take advantage of the built-in concurrency primitives of the language. 
This implementation integrates _getblocktemplate_ used for pool mining instead of _getwork_ as it is getting the standard with BIP 0023. 

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

## Motivation

This project is for educational purpose at start for a deep and better understanding of Bitcoin _under the hood_ as well as the discovery of Go.

## Implementation

The interactions with the Bitcoin Core client are currently done with __ugly__ curl calls but we are soon going to move to our own Web Socket (coming days/weeks). This will allow us to avoid executing external Unix commands, to maintain long-held connections with the client and to reduce the overhead related to HTTP.

The mining side is implemented with a Job/Worker/Dispatcher pattern.

A WebSocket interacts with the Bitcoin Core Client. It sends the data to the Dispatcher which is formating it and spliting into 'Chunks'. Those Chunks are sent to the miners. When a successful block is found, it is sent back to the Dispatcher for verification and back to the WebSocket for submitting. 

The overall Bitcoin implementation, but especially the Block Header, stay as close as possible to the Reference. 

## Performance

The current version is able to perform around __350.000__ operations (double sha256 hashs) per second on a 2.4 GHz Intel Core i5.

## Dependencies

None, beside your Bitcoin client. 
The btcsuite/btcrpcclient library was used for the JSON-RPC call at the beginning but was recently removed due to its deprecated mining implementation (BIP 0023 again).

## Roadmap

Checkout the TODO list to see the incoming features and what is needed to meet the Bitcoin standard. 

## Suggestions - Contact

Any suggestions, comments, questions, ideas of features or advices for improvements are great welcome! 
Please send me a message on Twitter @mathieumailhos.

## Contributors

- Mathieu Mailhos

## References
- [Start with Bitcoin and download your client](https://bitcoin.org/en/)
- [BIP 0023](https://en.bitcoin.it/wiki/BIP_0023)
- [Block Hashing Algorithm](https://en.bitcoin.it/wiki/Block_hashing_algorithm)
- [Thread Pool in Golang](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)
