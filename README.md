# go-auth-otp-service

## Minimum Required Infrastructure (LOM)

- 16 GB Memory
- 8 Core Cpu
- 64 GB Storage
- At least one IPv4 pointing to a domain

## Deployment Instructions

- Clone the project from CVS:
    ```shell
    git clone <repo> && cd <app_dir>
    ```
- Copy <code>.env.example</code> to <code>.env</code> And modify it as you wish.
    ```shell
    cp .env.example .env
    ```
- Download and install dependencies
    ```shell
    go mod download
    ```
    ```
- Handle the database migrations via:
    ```shell
    ./ database migrate up
    ```
- Handle the database seeders via:
    ```shell
    ./ database seed run
    ```
- Bootstrap the application via:
    ```shell
    ./  app bootstrap
    ```

### Via Docker:
Docker and docker compose are available as well:
```shell
docker compose -f docker-compose.yml up -d
```
### Important:
specifying docker-compose.yml is necessary due to port exposure management.
<hr>

