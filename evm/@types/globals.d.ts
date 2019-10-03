declare const web3: import('web3')
declare const assert: Chai.Assert
declare const contract: import('mocha').MochaGlobals['describe']
declare const artifacts: {
  require: (contract: string) => any
}
