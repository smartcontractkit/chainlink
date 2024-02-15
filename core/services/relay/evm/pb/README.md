# Generating protobufs

Generate protobufs with the following command:

```bash 
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/$PROTO_FILE_NAME
```

For example, to generate the automation protobuf, run:

```bash 
protoc -I=automation/ --go_out=automation/ automation/offchainconfig.proto
```