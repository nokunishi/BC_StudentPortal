package main

import (
	"encoding/json"
	"strings"
)

const GNAME_INDEX = 0
const FNAME_INDEX = 1
const AGE_INDEX = 2
const PHONE_INDEX = 3
const SCHOOL_INDEX = 4
const ADDR_INDEX = 5

type Student struct {
	Name    Name
	Age     json.Number
	Phone   string
	School  string
	Address Address
}

type Address struct {
	City     string
	Province string
	Zip      string
	Country  string
}

type Name struct {
	Given  string
	Family string
}

func NewModel(attr []string) Student {
	for i := range attr {
		attr[i] = string(strings.Trim(attr[i], " "))
	}

	addr := NewAddr(attr[ADDR_INDEX])
	name := Name{attr[GNAME_INDEX], attr[FNAME_INDEX]}

	return Student{
		name,
		json.Number(attr[AGE_INDEX]),
		attr[PHONE_INDEX],
		attr[SCHOOL_INDEX],
		addr,
	}
}

func NewAddr(addr string) Address {
	addr_ := strings.Split(addr, ".")

	if len(addr_) < 4 {
		addr_ = append(addr_, "Canada")
	}

	if !strings.Contains(addr_[2], "-") {
		tmp1 := addr_[2][:3]
		tmp2 := addr_[2][3:]
		addr_[2] = tmp1 + "-" + tmp2
	}

	return Address{
		addr_[0], addr_[1], addr_[2], addr_[3],
	}
}
