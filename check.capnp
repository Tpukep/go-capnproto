@0xc00cd2272161d09c;

using Go = import "/go.capnp";

$Go.package("caps");

# Number check
annotation multof(field) :UInt32;
annotation min(field) :Int64;
annotation max(field) :Int64;

# Text check
annotation format(field) :Text;
annotation pattern(field) :Text;

# List check
annotation unique(field) :Void;

# Text & List check
annotation minlen(field) :UInt32;
annotation maxlen(field) :UInt32;
