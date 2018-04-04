'use strict';
const isStream = require('is-stream');
const isObservable = require('is-observable');

exports.isListr = obj => Boolean(obj && obj.setRenderer && obj.add && obj.run);
exports.isObservable = obj => isObservable(obj);
exports.isStream = obj => isStream(obj) && !isObservable(obj);
