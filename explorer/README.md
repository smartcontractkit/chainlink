# Explorer

## Environment variables

See [./src/config.ts](./src/config.ts) for the available list of environment variables.

## Deployment

### Build Docker image

The Explorer application is part of our yarn workspace monorepo. To ensure that `yarn` can
resolve dependencies correctly you must build the Docker image from the root of the monorepo.

```bash
docker build . -f explorer/Dockerfile -t smartcontract/explorer
```

## Setup

### Install packages

```
yarn install
cd client && yarn install && cd -
```

### Database Configuration

##### Creation

[TypeORM](https://typeorm.io/#/) requires you to create the db manually,
and in our case that will involve leverage postgresql's `createdb`:

```
createdb explorer_dev
yarn migration:run
```

##### Deletion

```
dropdb explorer_dev
```

##### Connection

[TypeORM](https://typeorm.io/#/migrations) has been configured to load
`ormconfig/<env>.json`. Therefore, if in development, it loads `ormconfig/development.json`,
if production, `ormconfig/production.json`.

##### Running alongside Chainlink Node (dev)

```
$ EXPLORER_URL=ws://localhost:8080 cldev node
$ yarn run dev # in another terminal
```

##### Migrations

Please see [TypeORM's migration guide](https://typeorm.io/#/migrations).

##### Local admin

Run `yarn seed:admin <username> <password>` to set up local admin credentials during development.

## Running on seperate origins locally

The client is able to run on a different origin than the server. The steps below outline a
quick way of testing this locally via [ngrok.](https://ngrok.com/)

### Configure ngrok

In a terminal pane:

```sh
# Setup ngrok to proxy the default server settings.
ngrok http 8080
```

In a separate terminal pane:

```sh
# Setup ngrok to proxy the default client settings.
ngrok http 3001
```

### Configuring the server

```sh
# replace http://1b623c12.ngrok.io with the forwarded url that the previous step gave you for
# forwarding the client via ngrok
EXPLORER_CLIENT_ORIGIN=http://1b623c12.ngrok.io yarn dev:server
```

### Configuring the client

```sh
# replace http://03045a9a.ngrok.io with the forwarded url that the previous step gave you for
# forwarding the server via ngrok
DANGEROUSLY_DISABLE_HOST_CHECK=true  REACT_APP_EXPLORER_BASEURL=http://03045a9a.ngrok.io yarn start
```

Note the usage of `DANGEROUSLY_DISABLE_HOST_CHECK`, it is described here: https://create-react-app.dev/docs/proxying-api-requests-in-development/#invalid-host-header-errors-after-configuring-proxy
Using the safe `HOST` variable does not work with ngrok, so unforunately this is the only way of using ngrok with
create react app. Consider running the client dev server in a VM or remote machine that is sandboxed.

Another way of testing a separate domain is to not use ngrok to forward the client, and to just use it locally via
`localhost:3001`. Make sure to set `EXPLORER_CLIENT_ORIGIN` to `http://localhost:3001` if so.

You should now be able to visit the client via browser by using the forwarded ngrok url, or localhost.
Observe network requests using the api having a different origin than the client, and successfully returning data.

### Configuring the client environment variables

Set Google Analytics tracking ID:

```
GA_ID="UA-128878871-10"
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
