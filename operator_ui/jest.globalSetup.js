// The normal setupFiles did not work since they run too late (jest: ^23.5.0). So it is mandatory to use the globalSetup file.
// https://stackoverflow.com/questions/56261381/how-do-i-set-a-timezone-in-my-jest-config
module.exports = async () => {
  process.env.TZ = 'UTC'
}
