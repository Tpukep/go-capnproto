@0x9adec29d8e2e1d77;

using Go = import "/go.capnp";

$Go.package("demo");

struct Book {
	title @0 :Text; # Title of the book.
	pageCount @1 :Int32; # Number of pages in the book.
	authors @2 :List(Person); # Authors of the book	
	content @3 :Data; # Book content

	description :group {
	    genre @4 :UInt32;
	    review @5 :Text;
	    glossary @6 :Text;
	}
}

struct Person {
	name @0 :Text;
	email @1 :Text;
	age @2 :UInt8;
	phone @3 :Text;

	address @4 :Address;

	struct Address {
	    houseNumber @0 :UInt32;
	    street @1 :Text;
	    city @2 :Text;
	    country @3 :Text;
	}

	employment :union {
	    unemployed @5 :Void;
	    employer @6 :Text;
	    school @7 :Text;
	    selfEmployed @8 :Void;
	    # We assume that a person is only one of these.
	}
}

struct PhoneNumber {
    number @0 :Text;
    type @1 :Type;

    enum Type {
      mobile @0;
      home @1;
      work @2;
    }
}
