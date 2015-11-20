using Go = import "../../go.capnp";
using Check = import "../../check.capnp";
# using Json = import "../../json.capnp";

@0x85d3acc39d94e0f8;

$Go.package("books");
$Go.import("fmt");


const eseq :List(Rfc3092Variable) = [foo, bar, baz];

const seq8 :List(Int8) = [3, 5, 7, 9];
const sequ8 :List(UInt8) = [3, 5, 7, 9];
const seq16 :List(Int16) = [3, 5, 7, 9];
const sequ16 :List(UInt16) = [3, 5, 7, 9];
const seq32 :List(Int32) = [3, 5, 7, 9];
const sequ32 :List(UInt32) = [3, 5, 7, 9];
const seq64 :List(Int64) = [3, 5, 7, 9];
const sequ64 :List(UInt64) = [3, 5, 7, 9];
const seqf32 :List(Float32) = [3.50, 5.34, 7.0, 9.0];
const seqf64 :List(Float64) = [3.50, 5.34, 7.0, 9.0];
const ans :List(Bool) = [true, false, true];
const let :List(Text) = ["a", "b", "c"];
const letd :List(Data) = [0x"9f98739c2b53835e 6720a00907abd42f", 0x"9f98739c2b53835e 6720a00907abd42f", 0x"9f98739c2b53835e 6720a00907abd42f"];

struct Book {
	title @0 :Text; # Title of the book.
	pageCount @1 :Int32; # Number of pages in the book.
	authors @2 :List(Person); # Authors of the book	
	content @3 :Data; # Book content
}

struct Person $Go.doc("Some Person") {
	name @0 :Text $Check.maxlen(256) $Check.minlen(2);
	email @1 :Text $Check.format("email");
	age @2 :UInt8 $Check.max(40);
	phone @3 :Text $Check.pattern("\\d+");

	# age @2 :UInt8 $Json.field("email") $Json.omitempty;

	# address @2 :Address;

	# struct Address {
	#     houseNumber @0 :UInt32;
	#     street @1 :Text;
	#     city @2 :Text;
	#     country @3 :Text;
	# }

	# address :group {
	#     houseNumber @2 :UInt32;
	#     street @3 :Text;
	#     city @4 :Text;
	#     country @5 :Text;
	# }
}

# const pi :Float32 = 3.14159;
# const secret :Data = 0x"9f98739c2b53835e 6720a00907abd42f";
# const bob :Person = (email = "bob@example.com", name = "Bob", age = 40);

# const foo :Text = "rss";
# const bar :Text = "Hello";
# const baz :Person = (name = .foo, email = .bar);

# struct Contact {
# 	address @0 :Person.Address;
# }

enum Rfc3092Variable {
  foo @0;
  bar @1;
  baz @2;
  qux @3;
}
