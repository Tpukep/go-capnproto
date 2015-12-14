@0xcdc2bb7bf9b89132;
using Go = import "/go.capnp";
using Codec = import "/caps/codec.capnp";
$Go.package("protocol");
$Codec.capnp;

enum Status {
  syn @0;
  ack @1;
  psh @2;
}

struct Stream {
    id      @0 :UInt64;
    seq     @1 :UInt64;
}

struct Frame {
    session @0 :UInt64;
    status  @1 :Status;
	stream  @2 :Stream;
    seq     @3 :UInt64;
    payload @4 :Data;
}
