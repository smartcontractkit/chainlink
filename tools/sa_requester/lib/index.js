const command = require('@oclif/command')
const fs = require('fs')
// Can't abort fetch on nodejs which results in a process that never exits
// on network request timeout. Use axios instead until support is added to fetch.
const axios = require('axios')

const CONTENT_TYPE_JSON = 'application/vnd.api+json'
const FETCH_TIMEOUT = 5000
const SERVICE_AGREEMENTS_PATH = '/v2/service_agreements'
const ACCOUNT_BALANCE_PATH = '/v2/account_balance'

function urlWithPath (t, path) {
  const u = new URL(t)
  u.pathname = path
  return u.toString()
}

class SaRequester extends command.Command {
  async run () {
    const { args, flags } = this.parse(SaRequester)
    const agreement = fs.readFileSync(flags.agreement, 'utf8')
    const oracleURLs = args.file.split(/\s+/)
    const addresses = await getOracleAddresses(oracleURLs)

    createServiceAgreements(agreement, addresses, oracleURLs)
      .then(signatures =>
        console.table(
          signatures,
          ['address', 'signature']
        )
      )
      .catch(e => console.log('Unable to create SA, got error:\n\n\t%s\n', e.message))
  }
}

const parseResponse = response => {
  if (response.status === 200) {
    const contentType = response.headers['content-type']
    if (contentType === CONTENT_TYPE_JSON) {
      return response.data.data
    } else {
      throw new Error(`Unexpected response content type: "${contentType}" expected: "${CONTENT_TYPE_JSON}"`)
    }
  }
  throw new Error(`Unexpected response: ${response.status} body: ${response.json()}`)
}

async function getOracleAddresses (oracleURLs) {
  return Promise.all(
    oracleURLs.map(baseURL => {
      const url = urlWithPath(baseURL, ACCOUNT_BALANCE_PATH)
      return axios.get(url, { timeout: FETCH_TIMEOUT })
        .then(parseResponse)
        .then(data => data.id)
    })
  )
}

async function createServiceAgreements (agreement, addresses, oracleURLs) {
  return Promise.all(
    oracleURLs.map((u, i) => {
      const url = urlWithPath(u, SERVICE_AGREEMENTS_PATH)
      return axios.post(url, agreement, { timeout: FETCH_TIMEOUT })
        .then(parseResponse)
        .then(data => ({signature: data.attributes.signature, address: addresses[i]}))
    })
  )
}

SaRequester.description = 'Collect the signatures for a service agreement from multiple chainlink nodes'
SaRequester.flags = {
  version: command.flags.version({ char: 'v' }),
  help: command.flags.help({ char: 'h' }),
  agreement: command.flags.string({ char: 'a', description: 'Location of agreement' })
}
SaRequester.args = [{ name: 'file' }]

module.exports = SaRequester
