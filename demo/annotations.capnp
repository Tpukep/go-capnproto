@0xeac197c12d74cbbb;

using Go = import "/go.capnp";
using Codec = import "../annotations/codec/codec.capnp";
using Check = import "../annotations/check/check.capnp";
using Field = import "../annotations/field/field.capnp";

$Go.package("demo");

$Codec.msgp;
$Codec.json;
$Codec.capnp;

struct Person $Go.doc("Some Person") {
	name  @0 :Text  $Check.value("max=256,min=2");
	email @1 :Text  $Check.value("email");
	age   @2 :UInt8 $Check.value("max=40");
	phone @3 :Text;
}

struct Book {
	title     @0 :Text;                                          # Title of the book.
	pageCount @1 :Int32        $Field.required("pc");            # Number of pages in the book.
	authors   @2 :List(Person) $Field.optional("authors");       # Authors of the book	
	content   @3 :Text         $Field.ignored;                   # Book content
}
