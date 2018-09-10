# protoc-gen-gotag (PGGT)

PGGT is a protoc plugin used to add/replace struct tags on generated protobuf messages.
Get  it using ```go get github.com/srikrsna/protoc-gen-gotag ```It supports the following features,

## Add/Replace Tags

New tags like xml, sql, bson etc... can be added to struct messages of protobuf. Example
```proto
syntax = "proto3";

package example;

import "tagger/tagger.proto";

message Example {
    string with_new_tags = 1 [(tagger.tags) = "graphql:\"withNewTags,optional\"" ];
    string with_new_multiple = 2 [(tagger.tags) = "graphql:\"withNewTags,optional\" xml:\"multi,omitempty\"" ];
    
    string replace_default = 3 [(tagger.tags) = "json:\"replacePrevious\""] ; 

    oneof one_of {
        option (tagger.oneof_tags) = "graphql:\"withNewTags,optional\"";
        string a = 5 [(tagger.tags) = "json:\"A\""];
        int32 b_jk = 6 [(tagger.tags) = "json:\"b_Jk\""];
    }
}

message SecondMessage {
    string with_new_tags = 1 [(tagger.tags) = "graphql:\"withNewTags,optional\"" ];
    string with_new_multiple = 2 [(tagger.tags) = "graphql:\"withNewTags,optional\" xml:\"multi,omitempty\"" ];
    
    string replace_default = 3 [(tagger.tags) = "json:\"replacePrevious\""] ; 
}
``` 

Then struct tags can be added by running this command **after** the regular protobuf generation command.
```bash
    protoc -I /usr/local/include \
    	-I . \
    	--gotag_out=:. example/example.proto
```

In the above example tags like graphql and xml will be added whereas existing tags such as json are replaced with the supplied values. 

## Add tags to XXX* fields

It is very useful to ignore XXX* fields in protobuf generated messages. The go protocol buffer compiler adds ```json:"-"``` tag to all XXX* fields. Additional tags can be added to these fields using the 'xxx' option of PGGT. It can be done like this. All '+' characters will be replaced with ':'.

```bash
    protoc -I /usr/local/include \
    	-I . \
    	--gotag_out=xxx="graphql+\"-\" bson+\"-\"":. example/example.proto
```

### Note
 
 This should always run after protocol buffer compiler has run. The command such as the one below will fail/produce unexpected results.
 ```bash
    protoc -I /usr/local/include \
        	-I . \
        	--go_out=:. \
        	--gotag_out=xxx="graphql+\"-\" bson+\"-\"":. example/example.proto
``` 
