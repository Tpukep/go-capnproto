@0xcbef7e7210d94a7e;

using Go = import "/go.capnp";
using Proto = import "vendor/gitlab.ostrovok.ru/hotcore/supplierd-proto/proto.capnp";

$Go.package("icepeak");

struct Tx {
	packages @0 :List(Proto.Package);
	supplier @1 :Text;
	feed     @2 :Text;
	created  @3 :UInt64;
	commited @4 :UInt64;
}
