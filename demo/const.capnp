@0xd40674aff4d0ce00;

using Go = import "/go.capnp";

$Go.package("demo");

struct Person {
	name @0 :Text;
	email @1 :Text;
}

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

const en :Rfc3092Variable = foo;
const in8 :Int8 = 3;
const uin8 :UInt8 = 49;
const in16 :Int16 = 159;
const uin16 :UInt16 = 59;
const in32 :Int32 = 59;
const uin32 :UInt32 = 159;
const in64 :Int64 = 159;
const uin64 :UInt64 = 4159;
const f32 :Float32 = 3.14159;
const f64 :Float64 = 3.14159;

const em :Text = "lisa@example.com";
const nm :Text = "Lisa";

const secret :Data = 0x"9f98739c2b53835e 6720a00907abd42f";

const bob :Person = (email = "bob@example.com", name = "Bob");
const liz :Person = (email = .em, name = .nm);

const seqp :List(Person) = [.bob, .liz];

enum Rfc3092Variable {
  foo @0;
  bar @1;
  baz @2;
  qux @3;
}
