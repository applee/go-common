package card

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	AGE_MAX = 120

	ERR_IDENTITY_FORMAT_INVALID   = "身份证号格式不正确"
	ERR_IDENTITY_LENGTH_INVALID   = "身份证长度错误"
	ERR_IDENTITY_AGE_INVALID      = "身份证年龄不合法(0-%d岁之间)"
	ERR_IDENTITY_CHECKSUM_INVALID = "身份证校验位不正确"
)

type IdentityCard struct {
	Original string //原始
	Address  string //地址
	Birthday string //生日
	Order    string //顺序码
	Checksum byte   //校验码
}

func NewIdentityCard(s string) (*IdentityCard, error) {
	e := regexp.MustCompile("(^\\d{15}$)|(^\\d{17}(\\d|X|x)$)")
	if !e.Match([]byte(s)) {
		return nil, errors.New(ERR_IDENTITY_FORMAT_INVALID)
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

//计算年龄
func (i *IdentityCard) CalcAge() (age int) {
	birthday, err := time.Parse("20060102", i.Birthday)
	if err != nil {
		return -1
	}
	now := time.Now()
	age = now.Year() - birthday.Year()
	if age > 1 && (birthday.Month() > now.Month() || (birthday.Month() == now.Month() && birthday.Day() < now.Day())) {
		age--
	}
	return
}

//计算验证码
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

//验证
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

//验证省份
func (i *IdentityCard) ValidateProvince() (ok bool, err error) {
	_, ok = IdentityCardProvince[i.Address[:2]]
	return
}

//验证年龄
func (i *IdentityCard) ValidateAge() (ok bool, err error) {
	age := i.CalcAge()
	if age < 0 || age > AGE_MAX {
		return false, fmt.Errorf(ERR_IDENTITY_AGE_INVALID, AGE_MAX)
	}
	return true, nil
}

//验证校验码
func (i *IdentityCard) ValidateChecksum() (bool, error) {
	if i.CalcChecksum() != i.Checksum {
		return false, errors.New(ERR_IDENTITY_CHECKSUM_INVALID)
	}
	return true, nil
}

//归属地
func (i *IdentityCard) GetAddress() string {
	province := IdentityCardProvince[i.Address[:2]]
	cityCode, _ := strconv.Atoi(i.Address[:6])
	city := IdentityAddress[cityCode]
	if city == "" {
		city = "未知"
	}
	return province + "/" + city
}
