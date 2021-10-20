interface ExecaMock extends jest.Mock {
  sync: typeof jest.fn
}

const execa = jest.genMockFromModule('execa') as ExecaMock
execa.sync = jest.fn().mockReturnValue({ stdout: '{}' })
export default execa
