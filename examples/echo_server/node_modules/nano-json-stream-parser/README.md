# nano-json-stream-parser

A complete, pure JavaScript, streamed JSON parser in about `750 bytes` (gzipped). It is similar to [Oboe.js](https://github.com/jimhigson/oboe.js/), a streaming JSON micro-library with a size of `4.8kb` (gzipped). While that alone isn't much, sizes add up quickly when you stack many libs. This lib achieves a 85% size reduction, while still offering the same main functionality. Uses ES6 arrows.

## Install

    npm i nano-json-stream-parser

## Usage

Usage is self explanatory:

```javascript
const njsp = require("nano-json-stream-parser");

// Callback is called when there is a complete JSON
const parse = njsp((json) => console.log(json));

parse('[1,2,3,4]');

parse('[1,2');
parse(',3,4]');

parse("[::invalid_json_is_ignored::]");

parse('{"pos": {"x":');
parse('1.70, "y": 2.');
parse('49, "z": 2e3}}');

parse('[ "aaaa\\"abcd\\u0123\\\\aa\\/aa" ]')
```

Output:

```
[ 1, 2, 3, 4 ]
[ 1, 2, 3, 4 ]
{ pos: { x: 1.7, y: 2.49, z: 2000 } }
[ 'aaaa"abcdģ\\aa/aa' ]
```

## Disclaimer

This library has no tests yet and could contain buggy edge-cases.
