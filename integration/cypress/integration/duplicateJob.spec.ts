context('End to end', function () {
  it('Duplicates a job', () => {
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
      .as('jobId1')
    cy.get('#created-job').click()
    cy.contains('h6', 'Job spec detail').should('exist')

    // Duplicate Job
    cy.clickLink('Duplicate')
    cy.clickButton('Create Job')
    cy.contains('p', 'Successfully created job')
      .children('a')
      .invoke('text')
      .as('jobId2')
    cy.get('#created-job').click()
    cy.contains('h6', 'Job spec detail').should('exist')

    // Ensure jobs IDs are different
    cy.get('@jobId1').then((jobId1) => {
      cy.get('@jobId2').then((jobId2) => {
        expect(jobId1).to.not.equal(jobId2)
      })
    })
  })
})
