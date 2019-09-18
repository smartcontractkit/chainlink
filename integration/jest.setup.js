/* eslint-disable @typescript-eslint/no-var-requires */

const pupExpect = require('expect-puppeteer')

jest.setTimeout(30000)
pupExpect.setDefaultOptions({ timeout: 3000 })
