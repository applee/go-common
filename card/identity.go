package card

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// The up bound of person age.
const (
	MaxAge = 256
)

// Define errors
var (
	ErrInvalidFormat   = errors.New("invalid format of China identity card")
	ErrInvalidLength   = errors.New("invalid length of China identity card")
	ErrInvalidAge      = errors.New("invalid age(between 0-120)")
	ErrInvalidCheckSum = errors.New("invalid checksum")
)

// IdentityCard represents China identity card.
type IdentityCard struct {
	Original string
	Address  string
	Birthday string
	Order    string
	Checksum byte
}

// New initializes the IdentityCard with given string.
func New(s string) (*IdentityCard, error) {
	e := regexp.MustCompile("(^\\d{15}$)|(^\\d{17}(\\d|X|x)$)")
	if !e.Match([]byte(s)) {
		return nil, ErrInvalidFormat
	}
	i := &IdentityCard{
		Original: s,
		Address:  s[0:6],
	}
	if len(s) == 15 {
		i.Birthday = "19" + s[6:12]
		i.Order = s[12:]
	} else {
		i.Birthday = s[6:14]
		i.Order = s[14:17]
		i.Checksum = s[17]
		if i.Checksum == 'x' {
			i.Checksum = 'X'
		}
	}
	return i, nil
}

// CalcAge calculates the real age.
func (i *IdentityCard) CalcAge(t *time.Time) (age int) {
	birthday, err := time.Parse("20060102", i.Birthday)
	if err != nil {
		return -1
	}
	now := time.Now()
	if t != nil {
		now = *t
	}
	age = now.Year() - birthday.Year()
	if age > 1 && (birthday.Month() > now.Month() || (birthday.Month() == now.Month() && birthday.Day() > now.Day())) {
		age--
	}
	return
}

// CalcChecksum calculate the checksum number.
func (i *IdentityCard) CalcChecksum() byte {
	var sum int
	s := i.Original
	if len(s) == 15 {
		s = s[:6] + "19" + s[6:]
	}
	for index := 0; index < 17; index++ {
		sum += int(s[index]-'0') * IdentityWeightFactor[index]
	}
	return IdentifyChecksums[sum%11]
}

// Validate checks the card number.
func (i *IdentityCard) Validate() (ok bool, err error) {
	ok, err = i.ValidateProvince()
	if !ok {
		return
	}
	ok, err = i.ValidateAge()
	if !ok {
		return
	}
	ok, err = i.ValidateChecksum()
	return
}

// ValidateProvince checks the province.
func (i *IdentityCard) ValidateProvince() (ok bool, err error) {
	_, ok = IdentityCardProvince[i.Address[:2]]
	return
}

// ValidateAge checks the age.
func (i *IdentityCard) ValidateAge() (ok bool, err error) {
	age := i.CalcAge(nil)
	if age < 0 || age > MaxAge {
		return false, ErrInvalidAge
	}
	return true, nil
}

// ValidateChecksum checks the checksum number.
func (i *IdentityCard) ValidateChecksum() (bool, error) {
	if i.CalcChecksum() != i.Checksum {
		return false, ErrInvalidCheckSum
	}
	return true, nil
}

// GetAddress get address name.
func (i *IdentityCard) GetAddress() string {
	province := IdentityCardProvince[i.Address[:2]]
	cityCode, _ := strconv.Atoi(i.Address[:6])
	city := IdentityAddress[cityCode]
	if city == "" {
		city = "Unknown"
	}
	return province + "/" + city
}
