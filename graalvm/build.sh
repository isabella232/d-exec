#!/usr/bin/env bash

javac -d out smartcontract/SmartContractException.java smartcontract/SmartContract.java tcp/Main.java

# Run it with `java -classpath out Server`
