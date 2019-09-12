#!/usr/bin/env node

if (process.argv.length <= 3) {
  console.log(
    `Usage: ${__filename} <twilio_account_sid> <twilio_auth_token> <from_number> <to_number>`,
  )
  process.exit(-1)
}

const express = require('express')
const bodyParser = require('body-parser')
const app = express()
const PORT = 6691

const [accountSid, authToken, from, to] = process.argv.slice(2)
const client = require('twilio')(accountSid, authToken)

app.use(bodyParser.json())
app.all('*', function(req, res) {
  const log = req.body
  const amount = parseInt(log.topics[1], 16)
  const message =
    'Hello Chainlink! You just sent received ' +
    amount +
    ' wei at ' +
    log.address
  console.log(
    'Sending text from ' + from + ' to ' + to + ' with message: ' + message,
  )
  client.messages
    .create({
      from: from,
      to: to,
      body: message,
    })
    .then(null, console.log)

  res.json({ body: { status: 'ok' } })
})

app.listen(PORT, function() {
  console.log('listening on port ' + PORT + ' for incoming ether')
})
