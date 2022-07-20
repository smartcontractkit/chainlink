import titleize from 'utils/titleize'

describe('Titleizes strings', () => {
  it('Capitalizes first words of the sentence', () => {
    const brokenCase = 'cApiTAlS hErE, LoweRcaseS There'
    const brokenCaseTitleized = 'Capitals Here, Lowercases There'
    expect(titleize(brokenCase)).toEqual(brokenCaseTitleized)
  })
  it('Converts underscores into spaces', () => {
    const underscoreString = 'Pending_Run_Success'
    const underscoreStringTitleized = 'Pending Run Success'
    expect(titleize(underscoreString)).toEqual(underscoreStringTitleized)
  })
  it('Capitalizes first words and converts underscores into spaces', () => {
    const brokenCaseWithUnderscores = 'job_error_now'
    const brokenCaseWithUnderscoresCorrect = 'Job Error Now'
    expect(titleize(brokenCaseWithUnderscores)).toEqual(
      brokenCaseWithUnderscoresCorrect,
    )
  })
  it('Does not converts non string values ', () => {
    const date = new Date()
    expect(titleize(1)).toEqual(1)
    expect(titleize(undefined)).toEqual(undefined)
    expect(titleize(date)).toEqual(date)
  })
})
