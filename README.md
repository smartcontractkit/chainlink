# LINK Stats

## Deployment

TODO...

## Setup

### Creating the database

There is an odd behavior of `db-migrate` when creating the database using the 
default `yarn db-migrate db:create mydatabase_name` command. Use the following workaround:

```bash
yarn db-migrate db:create --config database_create.json linkstats_dev
```

More information can be found [here](https://github.com/db-migrate/node-db-migrate/issues/393)

### Migrations

https://db-migrate.readthedocs.io/en/latest/Getting%20Started/commands/#commands

```
# migrate create
yarn db-migrate create createMyTable

# migrate up
yarn db-migrate up

# migrate down
yarn db-migrate down
```
