@0xbdc890c7fe41570c;

using Go = import "/go.capnp";
using Codec = import "/caps/codec.capnp";

$Go.package("demo");

$Codec.capnp;

struct Message {
	union {
		void               @0 :Void;
	    revokedPackages    @1 :List(UInt32);
	    userSourceChanged  @2 :Text;
	    errorsCount        @3 :UInt64;
	    # badPackage         @4 :BadPackage;
	    # badPackages        @4 :List(BadPackage);
	    # endpointClosed     @4 :Endpoint;
	    endpointsClosed     @4 :List(Endpoint);
	}
}

struct BadPackage {
	id    @0 :UInt32;
	error @1 :Text;
}

enum Endpoint {
	supplier @0;
	stats @1;
}



struct Static {
	matching @0 :List(Text);
}

struct Instance {
	static             @0 :Static;
    revokedPackages    @1 :List(UInt32);
    userSourceChanged  @2 :Text;
}




# struct Message {
# 	type :union {
# 		void               @0 :Void;
# 	    revokedPackages    @1 :List(UInt32);
# 	    userSourceChanged  @2 :Text;
# 	}
# }

# struct Message {
#   id @0 :UInt64;

#   union {
#     revokedPackages :group {
#       ids     @1 :UInt64;
#     }
#     userSourceChanged :group {
#       source  @2 :Text;
#     }
#   }
# }