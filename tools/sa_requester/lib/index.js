const command = require('@oclif/command')
const fs = require('fs')
const URL = require('url').URL
// Can't abort fetch on nodejs which results in a process that never exits
// on network request timeout. Use axios instead until support is added to fetch.
const axios = require('axios')
// This overrides console.table in nodejs >= 9.x
require('console.table')

const CONTENT_TYPE_JSON = 'application/vnd.api+json'
const FETCH_TIMEOUT = 5000
const SERVICE_AGREEMENTS_PATH = '/v2/service_agreements'
const ACCOUNT_BALANCE_PATH = '/v2/user/balances'

function urlWithPath(t, path) {
  const u = new URL(t)
  u.pathname = path
  return u.toString()
}

class SaRequester extends command.Command {
  async run() {
    const { args, flags } = this.parse(SaRequester)
    const agreement = JSON.parse(fs.readFileSync(flags.agreement, 'utf8'))
    const oracleURLs = args.file.split(/\s+/)
    const addresses = await getOracleAddresses(oracleURLs)

    createServiceAgreements(agreement, addresses, oracleURLs)
      .then(signatures => console.table(['address', 'signature'], signatures))
      .catch(e =>
        console.log('Unable to create SA, got error:\n\n\t%s\n', e.message),
      )
  }
}

const parseError = ({ response }) => {
  if (response.status === 422) {
    throw new Error(response.data.errors)
  }

  throw new Error(`Unexpected response: ${response.status}`)
}

const parseResponse = response => {
  if (response.status === 200) {
    const contentType = response.headers['content-type']
    if (contentType === CONTENT_TYPE_JSON) {
      return response.data.data
    } else {
      throw new Error(
        `Unexpected response content type: "${contentType}" expected: "${CONTENT_TYPE_JSON}"`,
      )
    }
  }

  throw new Error(`Unexpected response: ${response.status}`)
}

async function getOracleAddresses(oracleURLs) {
  return Promise.all(
    oracleURLs.map(baseURL => {
      const url = urlWithPath(baseURL, ACCOUNT_BALANCE_PATH)
      return axios
        .get(url, { timeout: FETCH_TIMEOUT })
        .then(parseResponse)
        .then(data => data.id)
        .catch(parseError)
    }),
  )
}

async function createServiceAgreements(baseAgreement, addresses, oracleURLs) {
  return Promise.all(
    oracleURLs.map((u, i) => {
      const url = urlWithPath(u, SERVICE_AGREEMENTS_PATH)
      const serviceAgreementRequest = Object.assign({}, baseAgreement, {
        oracles: addresses,
      })

      return axios
        .post(url, serviceAgreementRequest, { timeout: FETCH_TIMEOUT })
        .then(parseResponse)
        .then(data => [addresses[i], data.attributes.signature])
        .catch(parseError)
    }),
  )
}

SaRequester.description =
  'Collect the signatures for a service agreement from multiple chainlink nodes'
SaRequester.flags = {
  version: command.flags.version({ char: 'v' }),
  help: command.flags.help({ char: 'h' }),
  agreement: command.flags.string({
    char: 'a',
    description: 'Location of agreement',
  }),
}
SaRequester.args = [{ name: 'file' }]

module.exports = SaRequester
