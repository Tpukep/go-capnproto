package demo

import (
  capn "github.com/glycerine/go-capnproto"
  "io"
)




func (s *DirectoryEntry) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	DirectoryEntryGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *DirectoryEntry) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootDirectoryEntryCapn(capMsg)
      DirectoryEntryCapnToGo(z, s)
   return nil
}



func DirectoryEntryCapnToGo(src DirectoryEntryCapn, dest *DirectoryEntry) *DirectoryEntry {
  if dest == nil {
    dest = &DirectoryEntry{}
  }
  dest.Name = src.Name()

  return dest
}



func DirectoryEntryGoToCapn(seg *capn.Segment, src *DirectoryEntry) DirectoryEntryCapn {
  dest := AutoNewDirectoryEntryCapn(seg)
  dest.SetName(src.Name)

  return dest
}
