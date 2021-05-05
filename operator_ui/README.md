# Operator UI

This package is responsible for rendering the UI of the chainlink node, which allows interactions wtih node jobs, jobs runs, configuration and any other related tasks.

## Development

Assuming you already have a local chainlink node listening on port 6688, run:

```
CHAINLINK_BASEURL=http://localhost:6688 CHAINLINK_VERSION='1@1' NODE_ENV=development yarn start
```

Now navigate to http://localhost:3000.

If sign-in doesn't work, check your network console, it's probably a CORS issue. You may need to run your chainlink node with `ALLOW_ORIGINS=http://localhost:3000` set.
