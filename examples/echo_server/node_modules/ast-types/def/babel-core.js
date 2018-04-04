module.exports = function (fork) {
  fork.use(require("./es7"));

  var types = fork.use(require("../lib/types"));
  var defaults = fork.use(require("../lib/shared")).defaults;
  var def = types.Type.def;
  var or = types.Type.or;

  def("Noop")
    .bases("Statement")
    .build();

  def("DoExpression")
    .bases("Expression")
    .build("body")
    .field("body", [def("Statement")]);

  def("Super")
    .bases("Expression")
    .build();

  def("BindExpression")
    .bases("Expression")
    .build("object", "callee")
    .field("object", or(def("Expression"), null))
    .field("callee", def("Expression"));

  def("Decorator")
    .bases("Node")
    .build("expression")
    .field("expression", def("Expression"));

  def("Property")
    .field("decorators",
           or([def("Decorator")], null),
           defaults["null"]);

  def("MethodDefinition")
    .field("decorators",
           or([def("Decorator")], null),
           defaults["null"]);

  def("MetaProperty")
    .bases("Expression")
    .build("meta", "property")
    .field("meta", def("Identifier"))
    .field("property", def("Identifier"));

  def("ParenthesizedExpression")
    .bases("Expression")
    .build("expression")
    .field("expression", def("Expression"));

  def("ImportSpecifier")
    .bases("ModuleSpecifier")
    .build("imported", "local")
    .field("imported", def("Identifier"));

  def("ImportDefaultSpecifier")
    .bases("ModuleSpecifier")
    .build("local");

  def("ImportNamespaceSpecifier")
    .bases("ModuleSpecifier")
    .build("local");

  def("ExportDefaultDeclaration")
    .bases("Declaration")
    .build("declaration")
    .field("declaration", or(def("Declaration"), def("Expression")));

  def("ExportNamedDeclaration")
    .bases("Declaration")
    .build("declaration", "specifiers", "source")
    .field("declaration", or(def("Declaration"), null))
    .field("specifiers", [def("ExportSpecifier")], defaults.emptyArray)
    .field("source", or(def("Literal"), null), defaults["null"]);

  def("ExportSpecifier")
    .bases("ModuleSpecifier")
    .build("local", "exported")
    .field("exported", def("Identifier"));

  def("ExportNamespaceSpecifier")
    .bases("Specifier")
    .build("exported")
    .field("exported", def("Identifier"));

  def("ExportDefaultSpecifier")
    .bases("Specifier")
    .build("exported")
    .field("exported", def("Identifier"));

  def("ExportAllDeclaration")
    .bases("Declaration")
    .build("exported", "source")
    .field("exported", or(def("Identifier"), null))
    .field("source", def("Literal"));

  def("CommentBlock")
    .bases("Comment")
    .build("value", /*optional:*/ "leading", "trailing");

  def("CommentLine")
    .bases("Comment")
    .build("value", /*optional:*/ "leading", "trailing");

  def("Directive")
    .bases("Node")
    .build("value")
    .field("value", def("DirectiveLiteral"));

  def("DirectiveLiteral")
    .bases("Node", "Expression")
    .build("value")
    .field("value", String, defaults["use strict"]);

  def("BlockStatement")
    .bases("Statement")
    .build("body")
    .field("body", [def("Statement")])
    .field("directives", [def("Directive")], defaults.emptyArray);

  def("Program")
    .bases("Node")
    .build("body")
    .field("body", [def("Statement")])
    .field("directives", [def("Directive")], defaults.emptyArray);

  // Split Literal
  def("StringLiteral")
    .bases("Literal")
    .build("value")
    .field("value", String);

  def("NumericLiteral")
    .bases("Literal")
    .build("value")
    .field("value", Number)
    .field("raw", or(String, null), defaults["null"])
    .field("extra", {
      rawValue: Number,
      raw: String
    }, function getDefault() {
      return {
        rawValue: this.value,
        raw: this.value + ""
      }
    });

  def("BigIntLiteral")
    .bases("Literal")
    .build("value")
    // Only String really seems appropriate here, since BigInt values
    // often exceed the limits of JS numbers.
    .field("value", or(String, Number))
    .field("extra", {
      rawValue: String,
      raw: String
    }, function getDefault() {
      return {
        rawValue: String(this.value),
        raw: this.value + "n"
      };
    });

  def("NullLiteral")
    .bases("Literal")
    .build()
    .field("value", null, defaults["null"]);

  def("BooleanLiteral")
    .bases("Literal")
    .build("value")
    .field("value", Boolean);

  def("RegExpLiteral")
    .bases("Literal")
    .build("pattern", "flags")
    .field("pattern", String)
    .field("flags", String)
    .field("value", RegExp, function () {
      return new RegExp(this.pattern, this.flags);
    });

  var ObjectExpressionProperty = or(
    def("Property"),
    def("ObjectMethod"),
    def("ObjectProperty"),
    def("SpreadProperty"),
    def("SpreadElement")
  );

  // Split Property -> ObjectProperty and ObjectMethod
  def("ObjectExpression")
    .bases("Expression")
    .build("properties")
    .field("properties", [ObjectExpressionProperty]);

  // ObjectMethod hoist .value properties to own properties
  def("ObjectMethod")
    .bases("Node", "Function")
    .build("kind", "key", "params", "body", "computed")
    .field("kind", or("method", "get", "set"))
    .field("key", or(def("Literal"), def("Identifier"), def("Expression")))
    .field("params", [def("Pattern")])
    .field("body", def("BlockStatement"))
    .field("computed", Boolean, defaults["false"])
    .field("generator", Boolean, defaults["false"])
    .field("async", Boolean, defaults["false"])
    .field("accessibility", // TypeScript
           or(def("Literal"), null),
           defaults["null"])
    .field("decorators",
           or([def("Decorator")], null),
           defaults["null"]);

  def("ObjectProperty")
    .bases("Node")
    .build("key", "value")
    .field("key", or(def("Literal"), def("Identifier"), def("Expression")))
    .field("value", or(def("Expression"), def("Pattern")))
    .field("accessibility", // TypeScript
           or(def("Literal"), null),
           defaults["null"])
    .field("computed", Boolean, defaults["false"]);

  var ClassBodyElement = or(
    def("MethodDefinition"),
    def("VariableDeclarator"),
    def("ClassPropertyDefinition"),
    def("ClassProperty"),
    def("ClassMethod")
  );

  // MethodDefinition -> ClassMethod
  def("ClassBody")
    .bases("Declaration")
    .build("body")
    .field("body", [ClassBodyElement]);

  def("ClassMethod")
    .bases("Declaration", "Function")
    .build("kind", "key", "params", "body", "computed", "static")
    .field("kind", or("get", "set", "method", "constructor"))
    .field("key", or(def("Literal"), def("Identifier"), def("Expression")))
    .field("params", [def("Pattern")])
    .field("body", def("BlockStatement"))
    .field("computed", Boolean, defaults["false"])
    .field("static", Boolean, defaults["false"])
    .field("generator", Boolean, defaults["false"])
    .field("async", Boolean, defaults["false"])
    .field("decorators",
           or([def("Decorator")], null),
           defaults["null"]);

  var ObjectPatternProperty = or(
    def("Property"),
    def("PropertyPattern"),
    def("SpreadPropertyPattern"),
    def("SpreadProperty"), // Used by Esprima
    def("ObjectProperty"), // Babel 6
    def("RestProperty") // Babel 6
  );

  // Split into RestProperty and SpreadProperty
  def("ObjectPattern")
    .bases("Pattern")
    .build("properties")
    .field("properties", [ObjectPatternProperty])
    .field("decorators",
           or([def("Decorator")], null),
           defaults["null"]);

  def("SpreadProperty")
    .bases("Node")
    .build("argument")
    .field("argument", def("Expression"));

  def("RestProperty")
    .bases("Node")
    .build("argument")
    .field("argument", def("Expression"));

  def("ForAwaitStatement")
    .bases("Statement")
    .build("left", "right", "body")
    .field("left", or(
      def("VariableDeclaration"),
      def("Expression")))
    .field("right", def("Expression"))
    .field("body", def("Statement"));

  // The callee node of a dynamic import(...) expression.
  def("Import")
    .bases("Expression")
    .build();
};
