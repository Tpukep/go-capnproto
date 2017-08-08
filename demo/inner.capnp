# schema.capnp
@0xb42bbb69f69a82f8;
$import "/go.capnp".package("blockav");
$import "/caps/codec.capnp".json(void);
$import "/caps/codec.capnp".capnp(void);
struct BlockHotel @0x86476d39d95b39b0 {  # 0 bytes, 1 ptrs
  offers @0 :List(Offer) $import "/caps/field.capnp".required("offers");  # ptr[0]
  struct Offer @0xb38682c6fad804b2 {  # 0 bytes, 1 ptrs
    blockID @0 :Text $import "/caps/field.capnp".required("block_id");  # ptr[0]
  }
}
