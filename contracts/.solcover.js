module.exports = {
  mocha: {
    grep: "@skip-coverage", // Find everything with this tag
    invert: true               // Run the grep's inverse set.
  },
  skipFiles: [
    "src/v0.4",
    "src/v0.5"
  ]
}