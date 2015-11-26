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
  dest.Content = *[]byteCapnToGo(src.Content(), nil)
  dest.Genre = src.Genre()
  dest.Review = src.Review()
  dest.Glossary = src.Glossary()

  return dest
}



func BookGoToCapn(seg *capn.Segment, src *Book) BookCapn {
  dest := AutoNewBookCapn(seg)
  dest.SetTitle(src.Title)
  dest.SetPageCount(src.PageCount)
  dest.SetAuthors([]PersonGoToCapn(seg, &src.Authors))
  dest.SetContent([]byteGoToCapn(seg, &src.Content))
  dest.SetGenre(src.Genre)
  dest.SetReview(src.Review)
  dest.SetGlossary(src.Glossary)

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
  dest.Address = *PersonAddressCapnToGo(src.Address(), nil)
  dest.Employer = src.Employer()
  dest.School = src.School()

  return dest
}



func PersonGoToCapn(seg *capn.Segment, src *Person) PersonCapn {
  dest := AutoNewPersonCapn(seg)
  dest.SetName(src.Name)
  dest.SetEmail(src.Email)
  dest.SetAge(src.Age)
  dest.SetPhone(src.Phone)
  dest.SetAddress(PersonAddressGoToCapn(seg, &src.Address))
  dest.SetEmployer(src.Employer)
  dest.SetSchool(src.School)

  return dest
}



func (s *PersonAddress) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	PersonAddressGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *PersonAddress) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootPersonAddressCapn(capMsg)
      PersonAddressCapnToGo(z, s)
   return nil
}



func PersonAddressCapnToGo(src PersonAddressCapn, dest *PersonAddress) *PersonAddress {
  if dest == nil {
    dest = &PersonAddress{}
  }
  dest.HouseNumber = src.HouseNumber()
  dest.Street = src.Street()
  dest.City = src.City()
  dest.Country = src.Country()

  return dest
}



func PersonAddressGoToCapn(seg *capn.Segment, src *PersonAddress) PersonAddressCapn {
  dest := AutoNewPersonAddressCapn(seg)
  dest.SetHouseNumber(src.HouseNumber)
  dest.SetStreet(src.Street)
  dest.SetCity(src.City)
  dest.SetCountry(src.Country)

  return dest
}



func (s *PhoneNumber) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	PhoneNumberGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *PhoneNumber) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootPhoneNumberCapn(capMsg)
      PhoneNumberCapnToGo(z, s)
   return nil
}



func PhoneNumberCapnToGo(src PhoneNumberCapn, dest *PhoneNumber) *PhoneNumber {
  if dest == nil {
    dest = &PhoneNumber{}
  }
  dest.Number = src.Number()
  dest.Type = *PhoneNumberTypeCapnToGo(src.Type(), nil)

  return dest
}



func PhoneNumberGoToCapn(seg *capn.Segment, src *PhoneNumber) PhoneNumberCapn {
  dest := AutoNewPhoneNumberCapn(seg)
  dest.SetNumber(src.Number)
  dest.SetType(PhoneNumberTypeGoToCapn(seg, &src.Type))

  return dest
}
