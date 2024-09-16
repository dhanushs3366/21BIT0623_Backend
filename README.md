# 21BIT0623_Backend
### Trademarkia Assignment

## Overview

This is a backend API built with Go, designed to allow users to upload files to AWS S3 and PostgreSQL with file-sharing capabilities. The project uses the Echo framework for the REST API, and is structured into three key modules, each responsible for specific tasks:

* **Handlers**: Manage incoming user requests and route them to the appropriate services.
* **Models**: Define the schema for data entities such as `User`, `File`, `Metadata`, etc.
* **Services**: The core logic that interacts with the database for CRUD operations, S3 for file storage, and Redis for caching purposes.

## Getting Started

1. **Install Golang**: Download and install Go from the official [Go installation guide](https://go.dev/doc/install).
   
2. **Install Docker**: Ensure Docker is installed by following the instructions [here](https://docs.docker.com/engine/install/).

3. **Clone the Repository**: Get the project code by running:
   ```bash
   git clone https://github.com/dhanushs3366/21BIT0623_Backend.git

4) **Navigate into the project**: 
    ```bash 
    cd 21BIT0623_Backend

5) **Configure your docker-compose file**: Follow [docker-compose.yaml](./docker-compose.yaml) set your preferred credentials for this project

6) **Configure your `.env` file**: Make a `.env` file by following [.sample.env](./.sample.env) and use your credentials for AWS,Redis and postgres DB

7) **Build and Run it**:Build the docker container using 
    ```bash
    docker compose up --build

