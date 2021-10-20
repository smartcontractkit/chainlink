context('End to end', function () {
  it('Deletes a completed job', () => {
    cy.login()

    // Can only delete job runs that are > 1 week old...
    // This test only ensures that the UI for deleting a job run is available

    cy.clickLink('Configuration')
    cy.contains('h5', 'Configuration').should('exist')
    cy.get('[data-cy=delete-completed-job-runs]').click()
    cy.contains('h5', 'Delete completed jobs runs').should('exist')
    cy.contains('span', 'Confirm').should('exist')
  })
})
