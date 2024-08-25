# Lightblocks Project

This project is a simple command processing application that uses RabbitMQ for communication between clients and the server. The client sends commands to a RabbitMQ queue, and the server processes these commands by interacting with a concurrent ordered map data structure. The processed results are written to a file.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Building the Project](#building-the-project)
- [Assumptions](#assumptions)
- 

## Prerequisites

Ensure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/) (for containerization)
- [Go](https://golang.org/dl/) (version 1.16 or later)
- [RabbitMQ](https://www.rabbitmq.com/download.html) (Ensure the server is running)
- [Git](https://git-scm.com/) (for version control)

## Building the Project

1. Clone the repository:

   ```bash
   git clone https://github.com/kanuku/lbs.git
   cd lbs
   ```

2. Build the project:

   ```bash
   go build -o server ./app.go
   go build -o client ./client/client.go
   ```

 3. Start application using Docker:

    ```bash
    docker-compose up
    ```

4. Send commands using filemode:

    ```bash
    ./client commands.txt
    ```

5. Send commands using interactive mode:

    ```bash
    ./client
    ```

    * Available Commands
        * Add Item: `addItem <key> <value>`
        * Delete Item: `deleteItem <key>`
        * Get Item: `getItem <key>`
        * Get All Items: `getAllItems`
        * Exit Client: `exit`


### Assumptions
 * The RabbitMQ server is running on `localhost` with the default port 5672.
 * The RabbitMQ credentials are guest/guest.
 * The server and client will communicate using the `lightblocks` queue.
 * Commands processed by the server are written to `server-output.txt`.
 * The client will handle simple command types (`addItem`, `deleteItem`, `getItem`, `getAllItems`).
 * The server gracefully handles shutdown signals like `SIGINT`, `SIGKILL`, or `SIGTERM`.

### Troubleshooting
 * **RabbitMQ Connection Issues**: If the client or server cannot connect to RabbitMQ, ensure RabbitMQ is running and accessible. Check the credentials and hostname.
 * **File Access Errors**: Ensure the server has write permissions for server-output.txt.
 * **Command Errors**: If the client reports invalid commands, check that commands follow the correct format.
