Plain Go code generator for [Capnproto](https://capnproto.org) definitions.

# HOW TO

1. Install

   - Install capnp tool
   
   - Install capnp plugins
   
   ```
   go get github.com/glycerine/bambam
   go get github.com/glycerine/go-capnproto
   go get github.com/tpukep/go-capnproto
   ```

   Copy go.capnp from glycerine/go-capnproto
   
   ```sh
   cp glycerine/go-capnproto/go.capnp /usr/local/include
   ```

2. Write Capn'proto schema

   `model.capnp`:
   ```
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

3. Generate Go plain code

   ```sh
   capnp compile -opgo model.capnp
   ```

4. Generate translation code

   ```sh
   bambam -o . -p model model.go
   ```

5. Generate Cap'n proto code

   Modify `schema.capnp`. Replace:
   ```
   using Go = import "go.capnp";
   ```
   
   With:
   ```
   using Go = import "/go.capnp";
   ```
   
   Generate
   
   ```sh
   capnp compile -ogo schema.capnp
   ```

