var dot = require('graphlib-dot')

// Works
var digraph = dot.read(
  `digraph { test [type=http resp="{\\"blah\\": \\"a\\"}"]; }`,
)
console.log(digraph)
// Fails
//var digraph = dot.read(`digraph { test [type=http resp=<{"blah": "a"}>]; }`);
let a = `digraph { test [type=http resp=<{"blah": "a"}>]; }`
//var digraph = dot.read(a.replace('/<[^>]*>/g', ""));

//var digraph = dot.read(a.replace('/<[^>]*>/g', ""));
console.log(a.replace(/\<[^\>]*\>/, '""'))
