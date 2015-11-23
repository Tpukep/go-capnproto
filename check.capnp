using Go = import "go.capnp";
@0xc93cdda73f575f51;

$Go.package("capn");

# struct Int64Range {
# 	start @0 :Int64;
# 	end @1 :Int64;
# }

# enum TextFormat {
#   datetime @0;
#   hostname @1;
#   email @2;
#   ipv4 @3;
#   ipv6 @4;
#   uri @5;
# }


# Number check
annotation multof(field) :UInt32;
annotation min(field) :Int64;
annotation max(field) :Int64;
# annotation range(field) :Int64Range;

# Text check
annotation format(field) :Text;
annotation pattern(field) :Text;

# List check
annotation unique(field) :Void;

# Text & List check
annotation minlen(field) :UInt32;
annotation maxlen(field) :UInt32;


