import { formatJobSpecType } from './formatJobSpecType'

describe('formatJobSpecType', () => {
  it("removes 'Spec' suffix as a default", () => {
    expect(formatJobSpecType('KeeperSpec')).toEqual('Keeper')
  })

  it('formats Direct Request as a special case', () => {
    expect(formatJobSpecType('DirectRequestSpec')).toEqual('Direct Request')
  })

  it('formats Flux Monitor as a special case', () => {
    expect(formatJobSpecType('FluxMonitorSpec')).toEqual('Flux Monitor')
  })
})
