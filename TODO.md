# To Do List 

## GotBtc-002: Fix bug not allowing to add more chunks while the dispatcher has started
It seems that when the dispatcher has started, the queue is blocked and can not be filled anymore. 
Miners and filling the queue have to be working in parallele together of course.

## GoBtc-003: Re-evaluate the maximum number of Goroutines available for mining and update the stats of the project
With the time-out added, the mining stats have changed. Since we are using more goroutines per miner for the timeout, an update of the performances must be done. We have to try running more or less goroutines, see in which case we can find a gain of performances. Finally, update the main page of the project with the new average of operations per seconds.

## GoBtc-004: Chunk validation
When a successful chunk is found, we want to give it to a "validation entity" before share it on the network. The validation entity will perform basic testing to make sure the block found is correct and will then submit it on the Bitcoin network. In this ticket, we have to make sure that we can send back a chunk over a chan (which has not be done yet). 

## GotBtc-005: Websocket
We do not want to do ugly bashly curl requests anymore. This is highly inneficient. 
Implement a reliable and efficient websocket for the basic operations that this bitcoin miner needs.

## GotBtc-006: Build coinbasetxn and the Merkle Root
Use Bitcoin documentation to build coinbasetxn which will allow us to build the Merkle Root.

## GoBtc-007: Calculate the Target
Use the difficulty to build the target of the BlockHeader. 
