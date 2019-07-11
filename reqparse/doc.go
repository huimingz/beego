// Copyright 2019 huimingz
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
http请求时的参数解析器，需要搭配beego框架使用

Useage
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
	    parser := reqparse.RequestParser{}
	    stu := student{}
	    err := parser.ParseArgs(&t.Controller, &stu)
	    if err != nil {
	        t.Ctx.WriteString(err.Error() + "\n")
	        return
	    }
	    t.Ctx.WriteString(fmt.Sprintf("%+v\n", stu))
	}

more docs: https://github.com/huimingz/beego/blob/master/reqparse/README.md
*/
package reqparse
