context('End to end', function() {
  it('creates a job that runs', () => {
    cy.visit('http://localhost:6688')
    cy.contains('Chainlink').should('exist')

    cy.login()

    // Create Job
    cy.contains('New Job').click({ force: true })
    cy.get('h5').should('contain', 'New Job')
    cy.fixture('job').then(job => {
      cy.get('textarea[id=json]').paste(JSON.stringify(job, null, 4))
    })
    cy.contains('Button', 'Create Job').click()
    cy.contains('p', 'Successfully created job').should('exist')

    // Run Job
    cy.get('#created-job').click({ force: true })
    cy.contains('Job Spec Detail')
    cy.contains('Button', 'Run').click({ force: true })
    cy.contains('p', 'Successfully created job run')
      .children('a')
      .click({ force: true })
      .invoke('text')
      .as('runId')
    cy.contains('a > p', 'JSON').click({ force: true })

    // Wait for job run to complete
    cy.refreshUntilFound('h5:contains(Completed)', { waitTime: 500 })
    cy.contains('h5', 'Completed').should('exist')

    // Navigate to transactions page
    cy.contains('li > a', 'Transactions').click({ force: true })
    cy.contains('h4', 'Transactions').should('exist')

    // Navigate to Explorer
    cy.forceVisit('http://localhost:8080')
    cy.get('@runId').then(runId => {
      cy.get('input[name=search]').type(runId)
    })
    cy.contains('Button', 'Search').click({ force: true })
    cy.get('@runId').then(runId => {
      cy.contains('a', runId).click({ force: true })
    })
    cy.contains('h5', 'Complete').should('exist')
  })
})
