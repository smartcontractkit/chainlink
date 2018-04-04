// The RLP format
// Serialization and deserialization for the BytesTree type, under the following grammar:
// | First byte | Meaning                                                                    |
// | ---------- | -------------------------------------------------------------------------- |
// | 0   to 127 | HEX(leaf)                                                                  |
// | 128 to 183 | HEX(length_of_leaf + 128) + HEX(leaf)                                      |
// | 184 to 191 | HEX(length_of_length_of_leaf + 128 + 55) + HEX(length_of_leaf) + HEX(leaf) |
// | 192 to 247 | HEX(length_of_node + 192) + HEX(node)                                      |
// | 248 to 255 | HEX(length_of_length_of_node + 128 + 55) + HEX(length_of_node) + HEX(node) |

var encode = function encode(tree) {
  var padEven = function padEven(str) {
    return str.length % 2 === 0 ? str : "0" + str;
  };

  var uint = function uint(num) {
    return padEven(num.toString(16));
  };

  var length = function length(len, add) {
    return len < 56 ? uint(add + len) : uint(add + uint(len).length / 2 + 55) + uint(len);
  };

  var dataTree = function dataTree(tree) {
    if (typeof tree === "string") {
      var hex = tree.slice(2);
      var pre = hex.length != 2 || hex >= "80" ? length(hex.length / 2, 128) : "";
      return pre + hex;
    } else {
      var _hex = tree.map(dataTree).join("");
      var _pre = length(_hex.length / 2, 192);
      return _pre + _hex;
    }
  };

  return "0x" + dataTree(tree);
};

var decode = function decode(hex) {
  var i = 2;

  var parseTree = function parseTree() {
    if (i >= hex.length) throw "";
    var head = hex.slice(i, i + 2);
    return head < "80" ? (i += 2, "0x" + head) : head < "c0" ? parseHex() : parseList();
  };

  var parseLength = function parseLength() {
    var len = parseInt(hex.slice(i, i += 2), 16) % 64;
    return len < 56 ? len : parseInt(hex.slice(i, i += (len - 55) * 2), 16);
  };

  var parseHex = function parseHex() {
    var len = parseLength();
    return "0x" + hex.slice(i, i += len * 2);
  };

  var parseList = function parseList() {
    var lim = parseLength() * 2 + i;
    var list = [];
    while (i < lim) {
      list.push(parseTree());
    }return list;
  };

  try {
    return parseTree();
  } catch (e) {
    return [];
  }
};

module.exports = { encode: encode, decode: decode };