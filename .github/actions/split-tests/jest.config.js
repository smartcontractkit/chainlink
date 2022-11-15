/** @type {import('ts-jest').JestConfigWithTsJest} */
const jestConfig = {
  preset: "ts-jest/presets/default-esm",
  resolver: "<rootDir>/mjs-resolver.ts",
  transform: {
    "^.+\\.mts?$": [
      "ts-jest",
      {
        useESM: true,
      },
    ],
  },
  testEnvironment: "node",
};
export default jestConfig;
