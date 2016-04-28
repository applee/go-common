package card

import (
	"testing"
)

func Test_IdentitySingle(t *testing.T) {
	id := "44522119830409724x"
	i, err := NewIdentityCard(id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("年龄：%d", i.CalcAge())
	t.Logf("校验码：%+q", i.CalcChecksum())
	t.Logf("归属地：%s", i.GetAddress())
	ok, err := i.Validate()
	if !ok {
		t.Error(err)
	} else {
		t.Log("校验通过")
	}

}
