#!/usr/bin/env node

var echoes = 0;
var express = require('express');
var body_parser = require('body-parser');
var app = express();
var PORT = process.argv[2];
app.use(body_parser.json())

app.get('/count', function (req, res) {
  res.json(echoes);
});

app.all('*', function (req, res) {
  echoes += 1;
  console.log({headers: req.headers, body: req.body});
  res.json({headers: req.headers, body: req.body});
});

if(isNaN(parseFloat(PORT))){
  console.log("defaulting to 6690");
  PORT = 6690;
}

app.listen(PORT, function () {
  console.log("listening on port " + PORT);
});
