#身份证校验工具

## Usage

````
i, err := NewIdentityCard("44522119830409724x")
//获取年龄
i.CalcAge()
//获取归属地
i.GetAddress()
//校验合法性
i.Validate()

````