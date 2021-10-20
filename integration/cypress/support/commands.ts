// Cypress doesn't always click buttons because it thinks they are hidden (when they are not)
Cypress.Commands.overwrite(
  'click',
  (originalFn: CypressChildFunction, subject: CypressChild, options: {}) => {
    const defaultOptions = {
      force: true,
    }
    options = Object.assign({}, defaultOptions, options)
    return originalFn(subject, options)
  },
)

Cypress.Commands.add('clickLink', (linkText) => {
  cy.contains('a', linkText).click()
})

Cypress.Commands.add('clickButton', (buttonText) => {
  cy.contains('button', buttonText).click()
})

// Used to fill a form by passing input fields as an argument in
// the format {name: value}
Cypress.Commands.add(
  'fill',
  { prevSubject: 'element' },
  (form: CypressChild, values: object) => {
    for (const field in values) {
      const value = values[field]
      cy.wrap(form).get(`input[name=${field}]`).clear().type(value)
    }
  },
)

Cypress.Commands.add('getJobJson', () => {
  cy.fixture('job').then((job) => {
    const host = Cypress.env('JOB_SERVER_HOST') || 'cypress-job-server'
    const port = Cypress.env('JOB_SERVER_PORT') || '6692'
    job.tasks[0].params.get = `http://${host}:${port}`
    cy.wrap(JSON.stringify(job, null, 4))
  })
})

// TODO - remove in future. Cypress potentially working on fix to 2 visit superdomain limit.
// or refactor ete tests to not share state b/t tests
// https://docs.cypress.io/guides/guides/web-security.html#One-Superdomain-per-Test
// https://github.com/cypress-io/cypress/issues/944
Cypress.Commands.add('forceVisit', (url) => {
  cy.get('body').then((body$) => {
    const appWindow = body$[0].ownerDocument?.defaultView
    const appIframe = appWindow?.parent.document.querySelector('iframe')
    return new Promise((resolve) => {
      if (appIframe) appIframe.onload = () => resolve()
      if (appWindow) appWindow.location = url
    })
  })
})

Cypress.Commands.add(
  'login',
  (email = 'notreal@fakeemail.ch', password = 'twochains') => {
    const url = Cypress.env('CHAINLINK_URL') || 'http://localhost:6688'
    cy.visit(url)
    cy.contains('Operator').should('exist')
    cy.get('form input[id=email]').type(email)
    cy.get('form input[id=password]').type(password)
    cy.get('form button').click()
    cy.contains('h5', 'Activity').should('exist')
  },
)

Cypress.Commands.add('paste', { prevSubject: true }, (subject, text) => {
  cy.wrap(subject).clear().invoke('val', text).type(' {backspace}')
})

// Command will continue to reload the page until a selector matches or
// the max # of reloads is reached
Cypress.Commands.add('reloadUntilFound', (selector, options = {}) => {
  const defaultOptions = {
    maxAttempts: 10,
    waitTime: Cypress.config('defaultCommandTimeout'),
  }
  options = Object.assign({}, defaultOptions, options)
  if (options.maxAttempts < 1) {
    throw `Unable to find ${selector} on page`
  }
  options.maxAttempts--
  const $element = Cypress.$(selector)
  if ($element.length > 0) {
    cy.wrap($element)
  } else {
    cy.reload()
    cy.wait(options.waitTime)
    cy.reloadUntilFound(selector, options)
  }
})

type CypressChild = Cypress.Chainable<JQuery<any>>
type CypressChildFunction = (subject: CypressChild, options: {}) => CypressChild
