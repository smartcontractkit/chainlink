const njsp = require("./../dist/index.js");

// Callback is called when there is a complete JSON
const parse = njsp((json) => console.log(json));

parse('[1,2,3,4]');

parse('[1,2');
parse(',3,4]');

parse('[::invalid_json_is_ignored::]');

parse('{"pos": {"x":');
parse('1.70, "y": 2.');
parse('49, "z": 2e3}}');

parse('[ "aaaa\\"abcd\\u0123\\\\aa\\/aa" ]')

parse('[true,0,false,0')
parse(',false,0,true]');
