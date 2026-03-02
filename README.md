# POC: Binary Dispatcher

## Inspiration

- git (git-shell) -> `git shell`
- docker (docker-compose) -> `docker compose`
- app (app-sub, app-add) -> `app sub` and `app add`

## Step

1. `task build` to build modules and core binary
2. `PATH=$PWD:$PATH ./app add 1 1` to calculate additional
3. `PATH=$PWD:$PATH ./app sub 1 1` to calculate subtraction
