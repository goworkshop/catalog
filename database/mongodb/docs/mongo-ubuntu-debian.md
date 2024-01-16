# MongoDB on Ubuntu and Debian

Follow these steps to install MongoDB on Ubuntu/Debian.

## Prerequisites

- Ubuntu or Debian OS

## Installation Steps

1. Update the package list:

    ```bash
    sudo apt update
    ```

2. Install MongoDB:

    ```bash
    sudo apt install -y mongodb
    ```

3. Start the MongoDB service:

    ```bash
    sudo systemctl start mongod
    ```

4. Verify the MongoDB service status:

    ```bash
    sudo systemctl status mongod
    ```

## MongoDB Container

The MongoDB container is running at `localhost:27017` with the following credentials:

- Username: root
- Password: example
