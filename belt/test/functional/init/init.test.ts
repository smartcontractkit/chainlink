import { join } from 'path'
import { cat, rm, mkdir } from 'shelljs'
import Init from '../../../src/commands/init'

const TEST_PATH = 'test/functional/init/'
const TEST_FS_PATH = join(TEST_PATH, 'testfs')
const RUNTIME_CONFIG_PATH = `${TEST_FS_PATH}/.beltrc`

describe('init tests', () => {
  beforeAll(() => {
    mkdir(TEST_FS_PATH)
  })

  afterAll(() => {
    rm('-r', TEST_FS_PATH)
  })

  it('should write to runtime config', async () => {
    await Init.run([
      "-m 'raise clutch area heavy horn course filter farm deny solid finger sudden' -c 4 -p fdf38d85d15e434e9b2ca152b7b1bc6f",
      TEST_FS_PATH,
    ])

    expect(cat(RUNTIME_CONFIG_PATH).stdout).toMatchSnapshot()
  })

  it('should support partial updates to runtime config', async () => {
    await Init.run(['-p gdf38d85d15e434e9b2ca152b7b1bc6f', TEST_FS_PATH])

    expect(cat(RUNTIME_CONFIG_PATH).stdout).toMatchSnapshot()
  })
})