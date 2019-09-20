context('state sharing test', function() {
  before(() => {})
  it('sets some state', () => {
    cy.wrap(1).as('foo')
  })

  it('retrieves state', () => {
    cy.get('@foo').should('equal', 1)
  })
})
