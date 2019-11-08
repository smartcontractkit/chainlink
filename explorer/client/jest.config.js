module.exports = {
  preset: 'ts-jest/presets/js-with-ts',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
  transformIgnorePatterns: ['node_modules/(?!(@chainlink/json-api-client)/)'],
  testRegex: '(/__tests__/(?!support/*)|(\\.|/)(test|spec))\\.tsx?$',
}
