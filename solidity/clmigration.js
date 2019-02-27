// clmigration provides two key helpers for Chainlink development:
// 1. wraps migrations that to be skipped in the test environment, since we
// recreate every contract beforeEach test, and hit other APIs in our migration process.
// 2. Prepare plumbing for correct async/await behavior,
// in spite of https://github.com/trufflesuite/truffle/issues/501

module.exports = function(callback) {
  return function(deployer, network) {
    if (network == 'test') {
      console.log('===== SKIPPING MIGRATIONS IN TEST ENVIRONMENT =====')
    } else {
      deployer
        .then(async () => {
          return await callback(deployer, network)
        })
        .catch(console.log)
    }
  }
}
