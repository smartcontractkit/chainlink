var Benchmark = require('benchmark')
var SHA3 = require('sha3');
var jsSha3 = require('js-sha3')
var cryptojsSha3 = require('crypto-js/sha3');
var encHex = require("crypto-js/enc-hex");
var crypto = require('crypto') 

var suite = new Benchmark.Suite;

// add tests 
suite.add('node sha3', function() {
  var d = new SHA3.SHA3Hash()
  d.update(crypto.randomBytes(64))
  d.digest('hex')

})
.add('cryptojs', function() {
  var data = encHex.parse(crypto.randomBytes(64).toString('hex'))
  cryptojsSha3(data, {
    outputLength: 512
  }).toString();
})
.add('js-sha3', function(){
  jsSha3.sha3_512(crypto.randomBytes(64));
})

// add listeners 
.on('cycle', function(event) {
  console.log(String(event.target));
})
.on('complete', function() {
  console.log('Fastest is ' + this.filter('fastest').pluck('name'));
})
// run async 
.run({ 'async': true });
