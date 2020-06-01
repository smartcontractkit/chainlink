# Chainlink Heartbeat ETH/USD network graph

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

## Available Scripts

In the project directory, you can run:

### `npm start:dev`

Runs the app in the development mode.<br>
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.<br>
You will also see any lint errors in the console.

### `npm test`

Launches the test runner in the interactive watch mode.<br>
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `npm run build`

Builds the app for production to the `build` folder.<br>
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.<br>
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

### `npm run start`

Launches the expressjs server that serve the `/build` folder

## Staging Deploy

This repo is configured to automatically create a draft staging environment
on Netlify when a PR is opened. You can access the URL by viewing the PR checks

![](./docs/pr-checks-deploy-url.png)

## Production Deploy

This repo is configured to automatically deploy to production when there is a
new push to the `master-feeds` branch.

https://www.netlify.com/blog/2015/10/01/a-step-by-step-guide-middleman-on-netlify/#step-2-link-to-your-github

## Hidden Features

### Display Offchain Comparison Links, Contract Health

On the landing page you can enable a link to an offchain price comparison resource and health checks with the `details` query parameter.

```
https://feeds.chain.link?details=true
```

![display-offchain-comparison](./docs/feed-landing-compare-offchain.png)

#### Checks

- Within price threshold
- Answer is 0

Color Codes

- Red: A check has failed (hover for tooltip that includes failed checks)
- Yellow: Unknown status (when the price health check is not configured)
- Green: Ok

### Override feeds & nodes JSON endpoints

Override the urls with a URI encoded query param

```
https://feeds.chain.link?feeds-json=https%3A%2F%2Fweiwatchers.com%2Ffeeds.json
https://feeds.chain.link?nodes-json=https%3A%2F%2Fweiwatchers.com%2Fnodes.json
```

You can use the browser console to encode your URI:

```javascript
> encodeURIComponent('https://weiwatchers.com/feeds.json')
"https%3A%2F%2Fweiwatchers.com%2Ffeeds.json"
```

![reference-contract-health](./docs/reference-contract-health.png)

## Available env vars

```
REACT_APP_INFURA_KEY              - infura key
REACT_APP_GA_ID                   - google analytics key
REACT_APP_FEEDS_JSON              - URL to load reference contract feeds
REACT_APP_NODES_JSON              - URL to load oracle nodes
REACT_APP_HOSTNAME_WHITELIST      - Comma separated hosts that can override the feed & node source via query parameters ([localhost,weiwatchers.com], etc...)
REACT_APP_DEV_HOSTNAME_WHITELIST  - Comma separated hosts that can use dev features like /custom ([localhost,feeds.web.dev], etc...)
REACT_APP_DEV_PROVIDER            - development web3 provider
```
