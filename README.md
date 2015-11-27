Plain Go code generator for [Capnproto](https://capnproto.org) definitions.
Depends on [go-capnproto v1](github.com/glycerine/go-capnproto).

# Annotations

There are available annotations:
   - For generating json tags. See `jsontag/json.capnp` file.
   - For generating tags for possible validation. See `check.capnp` file.
   - For generating msgpack tags. See `msgptag/msgp.capnp` file.

Examples of use located at `demo/annotations.capnp`.

# HOW TO

1. Install

   - Install capnp tool. [Instructions](https://capnproto.org/install.html)
   - Install capnp plugins
   
   ```sh
   go get github.com/glycerine/go-capnproto
   go get github.com/glycerine/bambam
   go get github.com/tpukep/go-capnproto
   ```

   Copy `go.capnp` from glycerine/go-capnproto to `/usr/local/include`
   
   ```sh
   cp glycerine/go-capnproto/go.capnp /usr/local/include
   ```

2. Write Capn'proto schema

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

3. Generate Go plain code.

   ```sh
   capnp compile -opgo model.capnp
   ```

4. Generate translation code.

   ```sh
   bambam -o . -p model model.go
   ```

5. Generate Cap'n proto code.

   Modify `schema.capnp`. Replace:
   ```capnp
   using Go = import "go.capnp";
   ```
   
   With:
   ```capnp
   using Go = import "/go.capnp";
   ```
   
   Run capnp tool.
   
   ```sh
   capnp compile -ogo schema.capnp
   ```
