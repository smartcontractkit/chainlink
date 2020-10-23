context('End to end', function () {
  it('Archives a job', () => {
    cy.login()

    // Create Job
    cy.clickLink('New Job')
    cy.contains('h5', 'New Job').should('exist')
    cy.getJobJson().then((jobJson) => {
      cy.get('textarea[id=json]').paste(jobJson)
    })
    cy.clickButton('Create Job')
    cy.contains('p', 'Successfully created job')
      .children('a')
      .invoke('text')
      .as('jobId')

    // Archive Job
    cy.get('#created-job').click()
    cy.contains('h6', 'Job spec detail').should('exist')
    cy.clickButton('Archive')
    cy.contains('h5', 'Warning').should('exist')
    cy.get('@jobId').then((jobId) => {
      cy.contains('button', `Archive ${jobId}`).click()
    })
    cy.contains('p', 'Successfully archived job').should('exist')
  })
})
