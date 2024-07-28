# pflow-eth WebUI

## Introduction

This web service is designed to facilitate interactions with Ethereum blockchain,
focusing on transaction management and synchronization with a local hardhat test network.

It provides a backend for the pflow-eth developer UI, enabling efficient testing and development workflows.

## Features

- **Transaction Management**: Supports querying, inserting, and managing blockchain transactions.
- **Synchronization**: Offers synchronization capabilities with a local hardhat test network for development and testing
  purposes.
- **Developer UI**: Includes a minimal API (without authentication) to interact with the contract, complemented by a
  developer UI for enhanced testing.
- **New Relic Integration**: Utilizes New Relic for application performance monitoring, logging, and custom event
  tracking.
- **Dynamic Configuration**: Features dynamic configuration through environment variables, allowing for flexible
  deployment and testing scenarios.
- **Materialized Views**: Leverages PostgreSQL materialized views for efficient data management and query optimization.
- **Contract Event Monitoring**: Provides infrastructure to experiment with multiple streams of contract events and
  their storage.

## Technology Stack

- **Backend**: Go (Gorilla Mux for routing, go-ethereum for Ethereum client, go.rice for embedding assets)
- **Database**: PostgreSQL
- **Frontend**: React (for the developer UI)
- **Monitoring**: New Relic
- **Other Tools**: npm for JavaScript package management

## Getting Started

To get started with the pflow-eth web service, clone the repository and ensure you have Go, npm, and PostgreSQL
installed. Follow the setup instructions in the `webui/README.md` for the frontend and refer to the `go.mod` file for
backend dependencies.

1. Clone the repository: `git clone <repository-url>`
2. Install backend dependencies: `go mod tidy`
3. Set up the database: `psql -f setup.sql`
4. Start the backend service: `go run main.go`
5. Navigate to the `webui` directory and install frontend dependencies: `npm install`
6. Start the developer UI: `npm start`

For detailed configuration options and environment variables, refer to the `config/vars.go` file.