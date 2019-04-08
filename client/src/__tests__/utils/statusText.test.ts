import statusText from '../../utils/statusText'

describe('utils/statusText', () => {
  it('is "Pending" for a status of "in_progress"', () => {
    expect(statusText('in_progress')).toEqual('Pending')
  })

  it('is "Failed" for a status of "error"', () => {
    expect(statusText('error')).toEqual('Failed')
  })

  it('is "Complete" for a status of "completed"', () => {
    expect(statusText('completed')).toEqual('Complete')
  })

  it('returns the status as titlecase by default', () => {
    expect(statusText('pending_confirmations')).toEqual('Pending Confirmations')
  })
})
