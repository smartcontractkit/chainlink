module.exports = function (fork) {
  fork.use(require("./babel-core"));
  fork.use(require("./flow"));
};
