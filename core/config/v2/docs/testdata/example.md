[//]: # (Generated - DO NOT EDIT.)

## Table of contents

- [Global](#Global)
- [TableName](#TableName)

## Global<a id='Global'></a>
```toml
FieldName = 'foo' # Default
```


This example demonstrates some of the features:

### FieldName<a id='FieldName'></a>
```toml
FieldName = 'foo' # Default
```
FieldName is a string with a default value. Every field **must** be documented with a comment that begins with the field name.

This is a loose comment.
Comments can span multiple lines.

## TableName<a id='TableName'></a>
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

### Bar<a id='TableName-Bar'></a>
```toml
Bar = 10 # Example
```
Bar doesn't have a default value, so an example **must** be included.

### TrickyField<a id='TableName-TrickyField'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
TrickyField = true # Default
```
TrickyField should only be used by advanced users, so it includes the special line `**ADVANCED**` to include a common warning tag.

