// sample of keys to test for
const CONFIG_KEYS = [
  'ACCOUNT_ADDRESS',
  'CHAINLINK_TLS_REDIRECT',
  'CHAINLINK_TLS_PORT',
  'ETH_CHAIN_ID',
  'ETH_GAS_PRICE_DEFAULT',
  'LOG_SQL_STATEMENTS',
  'MINIMUM_CONTRACT_PAYMENT',
  'REAPER_EXPIRATION',
  'TX_ATTEMPT_LIMIT',
]

context('End to end', function () {
  it('Deletes a completed job', () => {
    cy.login()

    // Create Job
    cy.clickLink('Configuration')
    cy.contains('h5', 'Configuration').should('exist')
    CONFIG_KEYS.forEach((key) => {
      cy.contains(key).should('exist')
    })
  })
})
