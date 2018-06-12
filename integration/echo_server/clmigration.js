// clmigration prepares plumbing for correct async/await behavior, in spite of
// https://github.com/trufflesuite/truffle/issues/501

module.exports = function (callback) {
  return function (deployer, network) {
    deployer.then(async () => {
      return await callback(deployer, network)
    }).catch(console.log)
  }
}
