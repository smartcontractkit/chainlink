'use strict';

var expectedImpl
if(process.env.ANY_PROMISE){
  // should load registration regardless
  expectedImpl = process.env.ANY_PROMISE
  require('any-promise/register')(expectedImpl)
} else {
  expectedImpl = 'global.Promise'
}

var impl = require('any-promise/implementation')

if(impl !== expectedImpl){
  throw new Error('Expecting '+expectedImpl+' got '+impl)
}
