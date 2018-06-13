#!/usr/bin/env node

const express = require('express')
const bodyParser = require('body-parser')
const app = express()
app.use(bodyParser.json())

let echoes = 0
app.get('/count', function (req, res) {
  res.json(echoes)
})

app.all('*', function (req, res) {
  echoes += 1
  console.log({headers: req.headers, body: req.body})
  res.json({headers: req.headers, body: req.body})
})

let port = process.argv[2]
if (isNaN(parseFloat(port))) {
  port = 6690
  console.log(`defaulting to ${port}`)
}

app.listen(port, function () {
  console.log(`echo_server listening on port ${port}`)
})
