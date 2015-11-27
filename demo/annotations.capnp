@0xeac197c12d74cbbb;

using Go = import "/go.capnp";
using Check = import "../check/check.capnp";
using Json = import "../jsontag/json.capnp";
using Msgp = import "../msgptag/msgp.capnp";

$Go.package("demo");

struct Person $Go.doc("Some Person") {
	name @0 :Text $Check.maxlen(256) $Check.minlen(2);
	email @1 :Text $Check.format("email");
	age @2 :UInt8 $Check.max(40);
	phone @3 :Text $Check.pattern("\\d+");
}

struct Book {
	title @0 :Text $Json.required("title"); # Title of the book.
	pageCount @1 :Int32 $Json.required("page_count") $Msgp.field("-"); # Number of pages in the book.
	authors @2 :List(Person) $Json.optional("authors"); # Authors of the book	
	content @3 :Text $Check.maxlen(256) $Json.required("content"); # Book content
}
