Go code generator for [Capnproto](https://capnproto.org) definitions.
It integrates MsgPack and Capn'proto code generators. Code generation options available as annotations. Also this tool adds support for possible data validation.

# HOW TO

You need Go 1.5 version. Set environment variable: `GO15VENDOREXPERIMENT=1`.

1. Install `capnp` tool. See [Instructions](https://capnproto.org/install.html)

2. Install `caps`
   
   ```sh
   go get github.com/tpukep/caps/...
   ```

3. Write Capn'proto schema `model.capnp`:
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
   caps -source model.capnp
   ```

# Annotations

## Codecs

   ```capnp
   using Codec = import "/caps/codec.capnp";

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
   using Field = import "/caps/field.capnp";
   
   struct Person {
      firstName  @0 :Text $Field.required("name"); # Rename field to use more
                                                   # compact name in serialization code
      email @1 :Text;
      age   @2 :UInt8 $Field.optional("pc");      # Mark as optional
      phone @3 :Text $Field.ignored;              # This field will be ignored
   }
   ```

## Checks

You can use [Go-playground Validator](https://github.com/go-playground/validator) expressions to generate tags in plain Go code. Note that `$Field` tags also generate corresponding `validate` tags.

### Example of various checks
  
   ```capnp
   using Check = import "/caps/check.capnp";

   struct Person $Go.doc("Some Person") {
      name  @0 :Text  $Check.value("max=256,min=2");
      email @1 :Text  $Check.value("email");
      age   @2 :UInt8 $Check.value("max=40");
      phone @3 :Text;
   }
   ```

Examples of use located at `demo/annotations.capnp`.
