const repl = require('repl');
const msg = 'message';

repl.start('> ').context.m = msg;
