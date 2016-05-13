@0xeb24122846d052e1;

using Go = import "/go.capnp";
using Codec = import "../codec.capnp";

$Go.package("demo");

# $Codec.capnp;

struct Map(Key, Value) {
  entries @0 :List(Entry);
  struct Entry {
    key @0 :Key;
    value @1 :Value;
  }
}

struct People {
  byName @0 :Map(Text, Person);
}

struct Person {
  name      @0 :Text;
  birthdate @1 :Int64;
}
