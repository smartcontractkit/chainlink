var times = function times(n, f, x) {
  for (var i = 0; i < n; ++i) {
    x = f(x);
  }
  return x;
};

module.exports = {
  times: times
};