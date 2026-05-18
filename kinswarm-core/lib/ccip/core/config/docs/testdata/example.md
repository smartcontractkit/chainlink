[//]: # (Generated - DO NOT EDIT.)

## Example

```toml
Bar = 7 # Required
```

## Global
```toml
FieldName = 'foo' # Default
```


This example demonstrates some of the features:

### FieldName
```toml
FieldName = 'foo' # Default
```
FieldName is a string with a default value. Every field **must** be documented with a comment that begins with the field name.

This is a loose comment.
Comments can span multiple lines.

## TableName
```toml
[TableName]
Bar = 10 # Example
TrickyField = true # Default
```
TableName holds settings that do something...
#### Details

We can include a long description here:
1. some
2. list
3. items

### Bar
```toml
Bar = 10 # Example
```
Bar doesn't have a default value, so an example **must** be included.

### TrickyField
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
TrickyField = true # Default
```
TrickyField should only be used by advanced users, so it includes the special line `**ADVANCED**` to include a common warning tag.

