# Hyperledger Fabric Insurance Application

This repository demonstrates the integration of a **Hyperledger Fabric** blockchain network with a REST API built using **Node.js** and **Express** for managing insurance policies. The application allows for registering insurance policies, querying policy details, and processing claims through smart contract interactions with the Hyperledger Fabric network.

## Overview

This application provides an API layer for interacting with an insurance management system built on **Hyperledger Fabric**. The system leverages Fabric's smart contracts (chaincode) to handle insurance policy registration, claim processing, and querying of policy information.

The API exposes three main functionalities:
1. **Registering a new insurance policy**.
2. **Querying an existing policy** by its ID.
3. **Processing a claim** against a policy.

## Technologies

- **Node.js**: JavaScript runtime used for backend development.
- **Express.js**: Web framework to handle HTTP requests and route them.
- **Hyperledger Fabric**: Permissioned blockchain framework for the smart contract and ledger.
- **Fabric Network SDK**: SDK to interact with Hyperledger Fabric blockchain from a Node.js application.
- **JavaScript (ES6)**: Scripting language used for backend logic.
- **Body-parser**: Middleware for parsing incoming request bodies in a middleware before your handlers.
- **Cors**: Middleware to enable cross-origin requests for the frontend.

## Setup and Installation

### Prerequisites

Before setting up the project, ensure the following are installed on your machine:

- **Node.js** (v14.x or above)
- **npm** (Node package manager)
- **Hyperledger Fabric** installed locally
- **Docker** (For peer nodes management)

