pflow.eth - Web3 Metamodel Repository
=====================================

"We're gonna need some bigger Hyperstructures".  

Overview
--------

I want to build a Dapp that incentive-ized users to produce petri-net smart contracts
and registers them with a merkle tree.The merkle tree is used to gate access to airdrops and rewards to modelers.

The Dapp will be a web3 front end that allows users to upload their petri-net smart contracts and register them with a merkle tree.  
The merkle tree is used to gate access to registry and may be used in airdrops and rewards to modelers.

Model Authors
-------------

As a model author, I want to be able to:

* Access the DApp's front-end hosted on IPFS.
* Upload my JavaScript model code to the DApp.
* Receive notifications and updates regarding my uploaded models. (solidity event)

Model Users
-----------

As a model user, I want to be able to:

* Access the DApp's front-end hosted on IPFS.
* Browse and search for available JavaScript models.
* Use the model with a connect Web3 wallet button to make on-chain transactions.

Registry Admins
---------------
As a registry admin, I want to be able to:

* Access the DApp's front-end hosted on IPFS.
* Administer the model registry on-chain. (Using admin permission)
* Manage user access rights and permissions to the registry.
* Control proxy contract updates in the production environment.
* Receive notifications and updates on registry changes and proxy contract updates.

Front-End DApp
--------------

The front-end DApp hosted on IPFS should include the following:

* HTML templates for user interfaces.
* JavaScript code for interactivity and functionality.
* Web3.js library for interacting with the Ethereum blockchain.
* Event listeners to capture user actions and trigger API calls.
* Display notifications and updates to the respective users.
* Provide a user-friendly interface for model authors, model users, registry admins, and treasury signatories.
 
Backend Node
------------
* The backend node hosted by QuickNode should:

* Communicate with the front-end DApp and web server.
* Receive event updates from the Ethereum blockchain.
* Trigger API calls to the web server based on received events.
* Handle communication between the DApp and the web server.

Web Server (pflow.dev)
----------------------

The web server should:

* Listen for API requests from the backend node.
* Receive event-triggered API calls from the backend node.
* Perform necessary actions on-chain based on the received requests.
* Update the PostgreSQL database to store contract code and index registry data.
* Provide appropriate responses to the backend node for event-triggered API calls.

Indexer
-------

The Indexer database should:

* Be provisioned as a backend database for the web server.
* Store the contract code and related metadata.
* Index the on-chain registry data for efficient querying.
* Support necessary CRUD operations for managing contract code and registry data.
* Publish new build to IPFS