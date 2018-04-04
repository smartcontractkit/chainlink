'use strict';
var symbolObservable = require('symbol-observable');

module.exports = function (fn) {
	return Boolean(fn && fn[symbolObservable]);
};
