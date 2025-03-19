## Tool for adding custom structure tags for proto message

https://morioh.com/a/65b8a7df8385/protoc-go-inject-tag-supercharge-golang-protobuf-with-custom-tags

```
go install github.com/favadi/protoc-go-inject-tag@latest
```

**Example**

1. Add comment line to your proto message field

```azure
// file: test.proto
syntax = "proto3";

package pb;
option go_package = "/pb";

message IP {
  // @gotags: valid:"ip"
  string Address = 1;

  // Or:
  string MAC = 2; // @gotags: validate:"omitempty"
}
```

2. Then regenerate pb.go file. Execute this

```azure
protoc-
go
-inject-tag -input="*.pb.go"
```