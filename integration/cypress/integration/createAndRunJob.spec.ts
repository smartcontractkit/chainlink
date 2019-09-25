context('End to end', function() {
  it('Creates a job that runs', () => {
    cy.login()

    // Create Job
    cy.clickLink('New Job')
    cy.contains('h5', 'New Job').should('exist')
    cy.fixture('job').then(job => {
      const port = Cypress.env('JOB_SERVER_PORT')
      job.tasks[0].params.get = `http://localhost:${port}`
      cy.get('textarea[id=json]').paste(JSON.stringify(job, null, 4))
    })
    cy.clickButton('Create Job')
    cy.contains('p', 'Successfully created job').should('exist')

    // Run Job
    cy.get('#created-job').click()
    cy.contains('Job Spec Detail')
    cy.clickButton('Run')
    cy.contains('p', 'Successfully created job run')
      .children('a')
      .click()
      .invoke('text')
      .as('runId')
    cy.contains('a > p', 'JSON').click()

    // Wait for job run to complete
    cy.reloadUntilFound('h5:contains(Completed)', { waitTime: 500 })
    cy.contains('h5', 'Completed').should('exist')

    // Navigate to transactions page
    cy.contains('li > a', 'Transactions').click()
    cy.contains('h4', 'Transactions').should('exist')

    // Navigate to Explorer
    cy.forceVisit('http://localhost:8080')
    cy.get('@runId').then(runId => {
      cy.get('input[name=search]').type(runId)
    })
    cy.clickButton('Search')
    cy.get('@runId').then(runId => {
      cy.clickLink(runId)
    })
    cy.contains('h5', 'Complete').should('exist')
  })
})
