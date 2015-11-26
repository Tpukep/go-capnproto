package demo

import (
  capn "github.com/glycerine/go-capnproto"
  "io"
)




func (s *Book) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	BookGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *Book) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootBookCapn(capMsg)
      BookCapnToGo(z, s)
   return nil
}



func BookCapnToGo(src BookCapn, dest *Book) *Book {
  if dest == nil {
    dest = &Book{}
  }
  dest.Title = src.Title()
  dest.PageCount = src.PageCount()
  dest.Authors = *[]PersonCapnToGo(src.Authors(), nil)
  dest.Content = src.Content()

  return dest
}



func BookGoToCapn(seg *capn.Segment, src *Book) BookCapn {
  dest := AutoNewBookCapn(seg)
  dest.SetTitle(src.Title)
  dest.SetPageCount(src.PageCount)
  dest.SetAuthors([]PersonGoToCapn(seg, &src.Authors))
  dest.SetContent(src.Content)

  return dest
}



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
  dest.Age = src.Age()
  dest.Phone = src.Phone()

  return dest
}



func PersonGoToCapn(seg *capn.Segment, src *Person) PersonCapn {
  dest := AutoNewPersonCapn(seg)
  dest.SetName(src.Name)
  dest.SetEmail(src.Email)
  dest.SetAge(src.Age)
  dest.SetPhone(src.Phone)

  return dest
}
