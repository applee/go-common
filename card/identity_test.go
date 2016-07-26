package card

import (
	"testing"
	"time"
)

func Test_IdentitySingle(t *testing.T) {
	id := "44522119830409724x"
	i, err := NewIdentityCard(id)
	if err != nil {
		t.Fatal(err)
	}
	start, _ := time.Parse("2006-01-02", "2015-04-08")
	t.Logf("年龄：%d", i.CalcAge(&start))
	t.Logf("校验码：%+q", i.CalcChecksum())
	t.Logf("归属地：%s", i.GetAddress())
	ok, err := i.Validate()
	if !ok {
		t.Error(err)
	} else {
		t.Log("校验通过")
	}

}
