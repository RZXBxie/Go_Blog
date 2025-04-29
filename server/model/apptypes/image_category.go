package apptypes

import "encoding/json"

type Category int

const (
	Null         Category = iota
	System                // 系统
	Carousel              // 背景
	Cover                 // 封面
	Illustration          //插图
	AdImage               // 广告
	Logo                  // 友链
)

func (c Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Category) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*c = ToCategory(s)
	return nil
}

func (c Category) String() string {
	var str string
	switch c {
	case Null:
		str = "未使用"
	case System:
		str = "系统"
	case Carousel:
		str = "背景"
	case Cover:
		str = "封面"
	case Illustration:
		str = "插图"
	case AdImage:
		str = "广告"
	case Logo:
		str = "友链"
	default:
		str = "未知类别"
	}
	return str
}

func ToCategory(s string) Category {
	var c Category
	switch s {
	case "未使用":
		c = Null
	case "系统":
		c = System
	case "背景":
		c = Carousel
	case "封面":
		c = Cover
	case "插图":
		c = Illustration
	case "广告":
		c = AdImage
	case "友链":
		c = Logo
	default:
		return -1
	}
	return c
}
