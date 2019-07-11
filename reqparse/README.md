# ReqParser

Go语言Beego框架参数解析器。

## Usage

```go
type student struct {
    Name   string `parser:"username; Required"`
    Grade  string `parser:"grade; Required; Choices(A, B, C)"`
    Number int    `parser:";Required"`
}

func (u student) ValidateNumber(number int) error {
    if number < 0 {
        return errors.New("value must be greater than 0")
    }
    return nil
}

type TestConroller struct {
    beego.Controller
}

func (t *TestConroller) Post() {
    parser := reqparse.RequerstParser{}
    stu := student{}
    err := parser.ParseArgs(&t.Controller, &stu)
    if err != nil {
        t.Ctx.WriteString(err.Error() + "\n")
        return
    }
    t.Ctx.WriteString(fmt.Sprintf("%+v\n", stu))
}
```



```
$ curl -X POST -d "username=dd&grade=A" "http://localhost:8080/v1/test"
number: missing required parameter 'number'

$ curl -X POST -d "username=dd&grade=1" "http://localhost:8080/v1/test"
'1' is not a valid choice

$ curl -X POST -d "username=dd&grade=A&number=-2" "http://localhost:8080/v1/test"
number: value must be greater than 0

$ curl -X POST -d "username=dd&grade=A&number=2" "http://localhost:8080/v1/test"
{Name:dd Grade:A Number:2}
```

## Feature

- 可自定义参数key值
- 提供基本标签验证
- 支持自定义参数验证

## Development Plan

- 多字段错误信息返回
- 更多的内置验证
- 自定义数据来源


## Supported Tag

- Required
- Range(min, max)
- Choices(1,2,3...)

## Custom Validator

如何添加：给需要验证的结构体添加方法，方法名为`Validate` + 字段名，接收的参数是一个当前字段类型的值，
返回类型是error。

自定义验证将会在标签验证之后执行，也就是作为参数验证的最后一环。

example:
```go
type student struct {
    Age int
}

func (s student) ValidateAge(age int) error {}
```
	

字段验证顺序为从左到右，会返回第一个错误。第一个字段为字段名key名，用户获取值，
留空时默认使用该结构体字段名作为key名，默认全部小写。

## License

```
Copyright 2019 huimingz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
