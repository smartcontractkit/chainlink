'use strict';
var REGISTRATION_KEY = '@@any-observable/REGISTRATION';
var registered = null;

module.exports = function (global, loadImplementation) {
	return function register(implementation, opts) {
		opts = opts || {};

		// global registration unless explicitly  {global: false} in options (default true)
		var registerGlobal = opts.global !== false;

		// load any previous global registration
		if (registerGlobal && !registered) {
			registered = global[REGISTRATION_KEY];
		}

		if (registered && implementation && registered.implementation !== implementation) {
			throw new Error('any-observable already defined as "' + registered.implementation +
				'".  You can only register an implementation before the first ' +
				' call to require(\'any-observable\') and an implementation cannot be changed');
		}

		if (!registered) {
			// use provided implementation
			if (implementation && opts.Observable) {
				registered = {
					Observable: opts.Observable,
					implementation: implementation
				};
			} else {
				// require implementation if implementation is specified but not provided
				registered = loadImplementation(implementation || null);
			}

			if (registerGlobal) {
				// register preference globally in case multiple installations
				global[REGISTRATION_KEY] = registered;
			}
		}

		return registered;
	};
};
