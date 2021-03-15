const bridgeProperties = {
  name: 'edit_test_bridge',
  url: 'http://example_1.com',
  minimumContractPayment: '123',
  confirmations: '5',
}

const newBridgeProperties = {
  url: 'http://example_2.com',
  minimumContractPayment: '321',
  confirmations: '10',
}

context('End to end', function () {
  it('Edits a bridge', () => {
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

    // Edit Bridge
    cy.clickLink('Edit')
    cy.get('form').fill(newBridgeProperties)
    cy.clickButton('Save Bridge')
    cy.contains('p', 'Successfully updated')
      .should('exist')
      .children('a')
      .click()
    cy.contains('td', bridgeProperties.name).should('exist')
    cy.contains('td', newBridgeProperties.url).should('exist')
    cy.contains('td', newBridgeProperties.minimumContractPayment).should(
      'exist',
    )
    cy.contains('td', newBridgeProperties.confirmations).should('exist')
  })
})
