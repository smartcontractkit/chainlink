module.exports = {
  skipFiles: [
    'v0.7/',
    'v0.8/mocks',
    'v0.8/interfaces',
    'v0.8/vendor',
    'v0.8/dev/interfaces',
    'v0.8/dev/vendor',
    'v0.8/dev/Keeper2_0/interfaces',
    'v0.8/dev/transmission',
    'v0.8/tests',
  ],
  istanbulReporter: ['text', 'text-summary', 'json'],
  mocha: {
    grep: '@skip-coverage', // Find everything with this tag
    invert: true, // Run the grep's inverse set.
  },
  configureYulOptimizer: true,
  solcOptimizerDetails: {
    peephole: false,
    jumpdestRemover: false,
    orderLiterals: true,
    deduplicate: false,
    cse: false,
    constantOptimizer: false,
    yul: true,
  },
}
