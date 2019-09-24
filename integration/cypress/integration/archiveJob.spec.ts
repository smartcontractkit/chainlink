context('End to end', function() {
  it('Archives a job', () => {
    cy.login()

    // Create Job
    cy.contains('New Job').click()
    cy.get('h5').should('contain', 'New Job')
    cy.fixture('job').then(job => {
      cy.get('textarea[id=json]').paste(JSON.stringify(job, null, 4))
    })
    cy.contains('Button', 'Create Job').click()
    cy.contains('p', 'Successfully created job')
      .children('a')
      .invoke('text')
      .as('jobId')

    // Archive Job
    cy.get('#created-job').click()
    cy.contains('h6', 'Job Spec Detail').should('exist')
    cy.contains('Button', 'Archive').click()
    cy.contains('h5', 'Warning').should('exist')
    cy.get('@jobId').then(jobId => {
      cy.contains('button', `Archive ${jobId}`).click()
    })
    cy.contains('p', 'Successfully archived job').should('exist')
  })
})
