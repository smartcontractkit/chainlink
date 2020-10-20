import { join } from 'path'
import { cat, cp, rm } from 'shelljs'
import Box from '../../../src/commands/box'
import {
  getJavascriptFiles,
  getPackageJson,
  getSolidityFiles,
  getTruffleConfig,
} from '../../../src/services/truffle-box'

const TEST_PATH = 'test/functional/box/'
const TEST_FS_PATH = join(TEST_PATH, 'testfs')
const FIXTURES_PATH = join(TEST_PATH, 'fixtures')

describe('truffle box tests', () => {
  function assertSnapshotOutput() {
    const packageJson = getPackageJson(TEST_FS_PATH)
    const truffleConfig = getTruffleConfig(TEST_FS_PATH)
    const solFiles = getSolidityFiles(TEST_FS_PATH)
    const jsFiles = getJavascriptFiles(TEST_FS_PATH)
    const allFiles = solFiles
      .concat(jsFiles, [truffleConfig, packageJson])
      .sort()

    expect(cat(allFiles).stdout).toMatchSnapshot()
  }

  beforeEach(() => {
    rm('-r', TEST_FS_PATH)
    cp('-r', FIXTURES_PATH, TEST_FS_PATH)
  })

  it('should properly convert to v0.4', async () => {
    await Box.run(['-s=0.4', TEST_FS_PATH])
    assertSnapshotOutput()
  })

  it('should properly convert to v0.5', async () => {
    await Box.run(['-s=0.5', TEST_FS_PATH])
    assertSnapshotOutput()
  })

  it('should throw when trying to select v0.6 (non-public)', async () => {
    expect(Box.run(['-s=0.6'])).rejects.toBeInstanceOf(Error)
  })

  it('should throw when passed invalid version', async () => {
    expect(Box.run(['-s=0.99'])).rejects.toBeInstanceOf(Error)
  })

  it('should output nothing on a dry run', async () => {
    const mockFn: any = () => {}

    const mockLog = jest.spyOn(console, 'log').mockImplementation(mockFn)
    mockLog.mockName('mockLog')

    const mockErr = jest.spyOn(console, 'error').mockImplementation(mockFn)
    mockErr.mockName('mockErr')

    await Box.run(['-d', '-s=0.5', TEST_FS_PATH])
    expect(mockErr).not.toHaveBeenCalled()
    expect(mockLog).toHaveBeenCalled()
    mockLog.mockRestore()
    mockErr.mockRestore()
    assertSnapshotOutput()
  })

  it.each([['0.5'], ['v0.5'], ['0.5.0']])(
    'should handle different version type inputs like %s',
    async (testVersion) => {
      expect(
        async () => await Box.run([`-s=${testVersion}`, TEST_FS_PATH]),
      ).not.toThrow()
    },
  )
})
