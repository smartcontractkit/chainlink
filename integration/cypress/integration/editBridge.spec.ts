context('End to end', function() {
  it('Edits a bridge', () => {
    cy.login()

    // Navigate to New Bridge page
    cy.clickLink('Bridges')
    cy.contains('h4', 'Bridges').should('exist')
    cy.clickLink('New Bridge')
    cy.contains('h5', 'New Bridge').should('exist')
    // Create Bridge
    cy.get('form input[name=name]').type('create_test_bridge')
    cy.get('form input[name=url]').type('http://example_1.com')
    cy.get('form input[name=minimumContractPayment]').type('123')
    cy.get('form input[name=confirmations]').type('5')
    cy.clickButton('Create Bridge')
    // Check new bridge created successfuly
    cy.contains('p', 'Successfully created bridge')
      .children('a')
      .click()
    cy.contains('p', 'create_test_bridge').should('exist')
    cy.contains('p', 'http://example_1.com').should('exist')
    cy.contains('p', '123').should('exist')
    cy.contains('p', '5').should('exist')
  })
})
