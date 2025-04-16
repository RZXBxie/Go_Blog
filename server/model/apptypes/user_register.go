package apptypes

import "encoding/json"

// Register 用户注册来源
type Register int

const (
	Email Register = iota
	QQ
)

func (r Register) MarshalJson() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Register) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*r = toRegister(str)
	return nil
}

func (r Register) String() string {
	var str string
	switch r {
	case Email:
		str = "邮箱"
	case QQ:
		str = "QQ"
	default:
		str = "未知注册方式"
	}
	return str
}

func toRegister(str string) Register {
	var r Register
	switch str {
	case "邮箱":
		r = Email
	case "QQ":
		r = QQ
	default:
		return -1
	}
	return r
}
