Go code generator for [Capnproto](https://capnproto.org) definitions.
It integrates MsgPack and Capn'proto code generators. Code generation options available as annotations. Also this tool adds support for possible data validation.

# HOW TO

1. Install `capnp` tool. [Instructions](https://capnproto.org/install.html)

2. Install `go-capnproto`
   
   ```sh
   go get github.com/tpukep/go-capnproto/...
   ```

3. Write Capn'proto schema

   For example `model.capnp`:
   ```capnp
   @0xad8d3cbc6db52a1d;
   
   using Go = import "/go.capnp";
   $Go.package("model");
   
   struct Book {
      title      @0:   Text;
      pageCount  @1:   Int32;
      authors    @2:   List(Text);
      content    @3:   Text;
   }
   ```

4. Generate Go plain code.

   ```sh
   go-capnproto -source model.capnp
   ```

# Annotations

## Codecs

   ```capnp
   using Codec = import "/codec.capnp";

   # By default only plain Go code will be generated
   # These annotations are enable extra features

   $Codec.msgp;  # Enables msgp code generation
   $Codec.json;  # Enables go json tags generation in plain Go code
   $Codec.capnp; # Enables Capn'proto code generation
  
   struct Person {
      name  @0 :Text;
      email @1 :Text;
      age   @2 :UInt8;
      phone @3 :Text;
   }
   ```

## Fields

   ```capnp
   using Field = import "/field.capnp";
   
   struct Person {
      firstName  @0 :Text $Field.required("name"); # Rename field to use more
                                                   # compact name in serialization code
      email @1 :Text;
      age   @2 :UInt8 $Field.optional("pc");      # Mark as optional
      phone @3 :Text $Field.ignored;              # This field will be ignored
   }
   ```

## Checks

   ```capnp
   using Check = import "/check.capnp";

   # Number checks
   $Check.multof(Num :UInt32); # Multiple of Num
   $Check.min(Num :UInt64);    # Great than Num
   $Check.max(Num :UInt64);    # Less than Num

   # Text check
   $Check.format(FORMAT_TYPE); # Predefined format
   $Check.pattern(REGEXP);     # Regular expression

   # List check
   $Check.unique;              # All List elements must be unique

   # Text & List check
   $Check.minlen(Len :UInt32); # Minimum length
   $Check.maxlen(Len :UInt32); # Maximum length

   struct Person {
      name      @0 :Text  $Check.maxlen(256) $Check.minlen(2);
      email     @1 :Text  $Check.format("email");
      age       @2 :UInt8 $Check.max(40);
      phone     @3 :Text  $Check.pattern("\\d+");
      addresses @4 :List(Text) $Check.unique;
   }
   ```

Examples of use located at `demo/annotations.capnp`.
