const bridgeProperties = {
  name: 'create_test_bridge',
  url: 'http://example.com',
  minimumContractPayment: '123',
  confirmations: '5',
}

context('End to end', function () {
  it('Creates a bridge', () => {
    cy.login()

    // Navigate to New Bridge page
    cy.clickLink('Bridges')
    cy.contains('h4', 'Bridges').should('exist')
    cy.clickLink('New Bridge')
    cy.contains('h5', 'New Bridge').should('exist')

    // Create Bridge
    cy.get('form').fill(bridgeProperties)
    cy.clickButton('Create Bridge')

    // Check new bridge created successfully
    cy.contains('p', 'Successfully created bridge')
      .should('exist')
      .children('a')
      .click()
    cy.contains('td', bridgeProperties.name).should('exist')
    cy.contains('td', bridgeProperties.url).should('exist')
    cy.contains('td', bridgeProperties.minimumContractPayment).should('exist')
    cy.contains('td', bridgeProperties.confirmations).should('exist')
  })
})
