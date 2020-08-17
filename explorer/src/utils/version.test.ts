import { existsSync, unlinkSync } from 'fs'
import { Environment } from '../config'
import { getVersion, VERSION_FILE_NAME, writeVersion } from './version'

function removeVersionFile() {
  try {
    unlinkSync(VERSION_FILE_NAME)
    // eslint-disable-next-line no-empty
  } catch {}
}

beforeAll(removeVersionFile)

describe('version tests', () => {
  describe('in a production environment', () => {
    afterEach(removeVersionFile)
    it(`writes to a ${VERSION_FILE_NAME} file`, async () => {
      await writeVersion()
      const file = await getVersion(Environment.PROD)
      console.log(file)
      expect(file).toBeTruthy()
    })

    it(`fails to read a ${VERSION_FILE_NAME} file if it does not exist`, async () => {
      expect.assertions(1)
      try {
        await getVersion(Environment.PROD)
      } catch (e) {
        expect((e as Error).message).toMatchInlineSnapshot(
          `"Could not read VERSION.json: ENOENT: no such file or directory, open 'VERSION.json'"`,
        )
      }
    })
  })

  describe('in a development or test environment', () => {
    it('reads directly from the git repository and package.jsons', async () => {
      expect(existsSync(VERSION_FILE_NAME)).toBeFalsy()
      const file = await getVersion(Environment.DEV)
      console.log(file)
      expect(file).toBeTruthy()
    })
  })
})
