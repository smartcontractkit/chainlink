const command = require('@oclif/command')
const fs = require('fs')
// Can't abort fetch on nodejs which results in a process that never exits
// on network request timeout. Use axios instead until support is added to fetch.
const axios = require('axios')

const CONTENT_TYPE_JSON = 'application/vnd.api+json'

class SaRequester extends command.Command {
  async run () {
    const { args, flags } = this.parse(SaRequester)
    const agreement = fs.readFileSync(flags.agreement, 'utf8')
    const oracleURLs = args.file.split(/\s+/)

    createServiceAgreements(agreement, oracleURLs)
      .then(signatures =>
        console.table(
          signatures.map(s => ({signature: s})),
          ['signature']
        )
      )
      .catch(e => console.log('Unable to create SA, got error:\n\n\t%s\n', e.message))
  }
}

async function createServiceAgreements (agreement, oracleURLs) {
  return Promise.all(
    oracleURLs.map(url => {
      return axios.post(url, agreement, { timeout: 5000 }).then(response => {
        if (response.status === 200) {
          const contentType = response.headers['content-type']
          if (contentType === CONTENT_TYPE_JSON) {
            return response.data.data
          } else {
            throw new Error(`Unexpected response content type: "${contentType}" expected: "${CONTENT_TYPE_JSON}"`)
          }
        }
        throw new Error(`Unexpected response: ${response.status} body: ${response.json()}`)
      }).then(data => data.attributes.signature)
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
