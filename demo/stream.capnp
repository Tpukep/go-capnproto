@0xb3ccab48d0e1a1d7;

using Go = import "/go.capnp";
using Codec = import "../codec.capnp";

$Go.package("demo");

$Codec.capnp;

struct Stream {
    id      @0 :UInt64;
    seq     @1 :UInt64;
}

struct Session {
   syn     @0 :UInt64;
   ack     @1 :UInt64;
   sess    @2 :UInt64;

   streams @3 :List(Stream);
}
