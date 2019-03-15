# LINK Stats [![CircleCI](https://circleci.com/gh/smartcontractkit/linkstats.svg?style=shield)](https://circleci.com/gh/smartcontractkit/linkstats)

## Deployment

TODO...

## Setup

### Install packages

```
yarn install
```

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

## Typescript

If you would like to add an npm package that doesn't have Typescript support you will need
to add the type definition yourself and make the Typescript transpiler aware of it. Custom
types are stored in the `types` directory.

```
echo "declare module 'my-package'" > types/my-package.d.ts
```

To make the the transpiler aware of the new type definition you will also need to add it to
the `"files": ...` section of `tsconfig.json`.

```
{
  // ...

  "files": [
    // ...
    "types/my-package.d.ts"
  ]
}
```
