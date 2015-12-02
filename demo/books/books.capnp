using Go = import "/go.capnp";
using Check = import "/check.capnp";

@0x85d3acc39d94e0f8;

$Go.package("books");

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
}
