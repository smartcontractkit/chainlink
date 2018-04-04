var merge = function merge(a) {
  return function (b) {
    var c = {};
    for (var key in a) {
      c[key] = a[key];
    }for (var _key in b) {
      c[_key] = b[_key];
    }return c;
  };
};

var remove = function remove(removeKey) {
  return function (a) {
    var b = {};
    for (var key in a) {
      if (key !== removeKey) b[key] = a[key];
    }return b;
  };
};

module.exports = {
  merge: merge,
  remove: remove
};