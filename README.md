# LINK Stats [![CircleCI](https://circleci.com/gh/smartcontractkit/linkstats.svg?style=shield)](https://circleci.com/gh/smartcontractkit/linkstats)

## Deployment

TODO...

## Setup

### Install packages

```
yarn install
```

### Database Configuration

##### Creation

[TypeORM](https://typeorm.io/#/) requires you to create the db manually,
and in our case that will involve leverage postgresql's `createdb`:

```
createdb linkstats_dev
```

##### Connection

[TypeORM](https://typeorm.io/#/migrations) has been configured to load
`ormconfig.<env>.json`. Therefore, if in development, it loads `ormconfig.development.json`,
if production, `ormconfig.production.json`.

##### Migrations

Please see [TypeORM's migration guide](https://typeorm.io/#/migrations).

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
