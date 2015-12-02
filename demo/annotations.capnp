@0xeac197c12d74cbbb;

using Go = import "/go.capnp";
using Check = import "../check.capnp";
using Field = import "../field.capnp";
using Codec = import "../codec.capnp";

$Go.package("demo");

$Codec.msgp;
$Codec.json;
$Codec.capnp;

struct Person $Go.doc("Some Person") {
	name  @0 :Text  $Check.maxlen(256) $Check.minlen(2);
	email @1 :Text  $Check.format("email");
	age   @2 :UInt8 $Check.max(40);
	phone @3 :Text  $Check.pattern("\\d+");
}

struct Book {
	title     @0 :Text;                                           # Title of the book.
	pageCount @1 :Int32        $Field.required("pc");             # Number of pages in the book.
	authors   @2 :List(Person) $Field.optional("authors");        # Authors of the book	
	content   @3 :Text         $Check.maxlen(256) $Field.ignored; # Book content
}
