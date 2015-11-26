package demo

// AUTO GENERATED - DO NOT EDIT

type Book struct {
	Title       string
	PageCount   int32
	Authors     []Person
	Content     []byte
	Description struct {
		Genre    uint32
		Review   string
		Glossary string
	}
}

type Person struct {
	Name       string
	Email      string
	Age        uint8
	Phone      string
	Address    PersonAddress
	Employment struct {
		/* Not implemented */ Employer string
		School                         string
		/* Not implemented */
	}
}

type PersonEmployment_Which uint16

const (
	PERSONEMPLOYMENT_UNEMPLOYED   PersonEmployment_Which = 0
	PERSONEMPLOYMENT_EMPLOYER     PersonEmployment_Which = 1
	PERSONEMPLOYMENT_SCHOOL       PersonEmployment_Which = 2
	PERSONEMPLOYMENT_SELFEMPLOYED PersonEmployment_Which = 3
)

type PersonAddress struct {
	HouseNumber uint32
	Street      string
	City        string
	Country     string
}

type PhoneNumber struct {
	Number string
	Type   PhoneNumberType
}

type PhoneNumberType uint16

const (
	PHONENUMBERTYPE_MOBILE PhoneNumberType = 0
	PHONENUMBERTYPE_HOME   PhoneNumberType = 1
	PHONENUMBERTYPE_WORK   PhoneNumberType = 2
)

func (c PhoneNumberType) String() string {
	switch c {
	case PHONENUMBERTYPE_MOBILE:
		return "mobile"
	case PHONENUMBERTYPE_HOME:
		return "home"
	case PHONENUMBERTYPE_WORK:
		return "work"
	default:
		return ""
	}
}

func PhoneNumberTypeFromString(c string) PhoneNumberType {
	switch c {
	case "mobile":
		return PHONENUMBERTYPE_MOBILE
	case "home":
		return PHONENUMBERTYPE_HOME
	case "work":
		return PHONENUMBERTYPE_WORK
	default:
		return 0
	}
}
