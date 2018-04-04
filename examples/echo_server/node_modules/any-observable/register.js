'use strict';
module.exports = require('./loader')(global, loadImplementation);

function loadImplementation(implementation) {
	var impl;

	if (implementation === 'global.Observable') {
		// if no implementation or env specified use global.Observable
		impl = {
			Observable: global.Observable,
			implementation: 'global.Observable'
		};
	} else if (implementation) {
		// if implementation specified, require it
		var lib = require(implementation);

		impl = {
			Observable: lib.Observable || lib.default || lib,
			implementation: implementation
		};
	} else {
		// try to auto detect implementation. This is non-deterministic
		// and should prefer other branches, but this is our last chance
		// to load something without throwing error
		impl = tryAutoDetect();
	}

	if (!impl) {
		throw new Error('Cannot find any-observable implementation nor' +
			' global.Observable. You must install polyfill or call' +
			' require("any-observable/register") with your preferred' +
			' implementation, e.g. require("any-observable/register")(\'rxjs\')' +
			' on application load prior to any require("any-observable").');
	}

	return impl;
}

function tryAutoDetect() {
	var libs = [
		'rxjs/Observable',
		'zen-observable'
	];

	for (var i = 0; i < libs.length; i++) {
		try {
			return loadImplementation(libs[i]);
		} catch (err) {}
	}

	return null;
}
