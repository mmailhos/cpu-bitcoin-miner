# To Do List 
Ordered by Priority

## GotBtc-006: Websocket
We do not want to do ugly bashly curl requests anymore. This is highly inneficient. 
Implement a reliable and efficient websocket for the basic operations that this bitcoin miner needs. Progrma must exists if the Bitcoin Core Client can not be reached.

## GoBtc-004: Re-evaluate the maximum number of Goroutines available for mining and update the stats of the project
With the time-out added, the mining stats have changed. Since we are using more goroutines per miner for the timeout, an update of the performances must be done. We have to try running more or less goroutines, see in which case we can find a gain of performances. Finally, update the main page of the project with the new average of operations per seconds.
Must be done after GoBtc-005 as it mights slightly modify the performances.

## GotBtc-007: Build coinbasetxn and the Merkle Root
Use Bitcoin documentation to build coinbasetxn which will allow us to build the Merkle Root.

## GoBtc-008: Calculate the Target
Use the difficulty to build the target of the BlockHeader. 

## GoBtc-009: PDF Documentation
Good looking documentation on a PDF presenting the goal of the project, the language, Bitcoin, the architecture of the projects, the crucial points, the roadmap... Make the project 'sellable'/'marketeable'.
