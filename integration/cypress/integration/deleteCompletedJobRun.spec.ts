context('End to end', function() {
  it('Deletes a completed job', () => {
    cy.login()

    // Create Job
    cy.clickLink('Configuration')
    cy.contains('h5', 'Configuration').should('exist')
    cy.clickButton('Delete Completed Jobs')
    cy.contains('span', 'Confirm delete all completed job runs').should('exist')
  })
})
