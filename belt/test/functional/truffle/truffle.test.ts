import { join } from 'path'
import { ls, rm } from 'shelljs'
import compile from '../../../src/commands/compile'

const TEST_PATH = 'test/functional/truffle/'
const TEST_FS_PATH = join(TEST_PATH, 'testfs')
const FIXTURES_PATH = join(TEST_PATH, 'fixtures')

describe('compileAll', () => {
  beforeEach(() => {
    rm('-r', TEST_FS_PATH)
  })

  it('should produce a truffle contract abstraction via a json artifact produced by sol-compiler', async () => {
    await compile.run([
      `--config=${join(FIXTURES_PATH, 'test.config.json')}`,
      'truffle',
    ])

    expect([...ls('-R', TEST_FS_PATH)]).toMatchSnapshot()
  })
})
