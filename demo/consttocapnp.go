package demo

import (
  capn "github.com/glycerine/go-capnproto"
  "io"
)




func (s *Person) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	PersonGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *Person) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootPersonCapn(capMsg)
      PersonCapnToGo(z, s)
   return nil
}



func PersonCapnToGo(src PersonCapn, dest *Person) *Person {
  if dest == nil {
    dest = &Person{}
  }
  dest.Name = src.Name()
  dest.Email = src.Email()

  return dest
}



func PersonGoToCapn(seg *capn.Segment, src *Person) PersonCapn {
  dest := AutoNewPersonCapn(seg)
  dest.SetName(src.Name)
  dest.SetEmail(src.Email)

  return dest
}
