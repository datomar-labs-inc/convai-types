package ctypes

import "encoding/xml"

type XMLResponse struct {
	XMLName  xml.Name     `xml:"response" json:"-" msgpack:"-" mapstructure:"-"`
	Messages []XMLMessage `xml:"message" json:"messages" msgpack:"messages"`
}

type XMLMessage struct {
	XMLName    xml.Name `xml:"message" json:"-" msgpack:"-" mapstructure:"-"`
	TypingTime *float64 `xml:"typing,attr,omitempty" json:"typing,omitempty" msgpack:"typing,omitempty"`
	Text       *string  `xml:"text,omitempty" json:"text,omitempty" msgpack:"text,omitempty"`

	Sender         *XMLSender         `xml:"sender,omitempty" json:"sender,omitempty" msgpack:"sender,omitempty"`
	QuickReplies   []XMLQR            `xml:"qr" json:"quickReplies,omitempty" msgpack:"quickReplies,omitempty" mapstructure:"quickReplies,omitempty"`
	CardCollection *XMLCardCollection `xml:"cards,omitempty" json:"cardCollection,omitempty" msgpack:"cardCollection,omitempty"`
	Image          *XMLImage          `xml:"image,omitempty" json:"image,omitempty" msgpack:"image,omitempty"`
}

type XMLQR struct {
	XMLName xml.Name `xml:"qr" json:"-" msgpack:"-" mapstructure:"-"`
	Text    string   `xml:",innerxml" json:"text" msgpack:"text" mapstructure:"text"`
	Value   *string  `xml:"value,attr,omitempty" json:"value,omitempty" msgpack:"value,omitempty" mapstructure:"value,omitempty"`
}

type XMLPhone struct {
	XMLName xml.Name `xml:"phone" json:"-" msgpack:"-" mapstructure:"-"`
	Number  string   `xml:",innerxml" json:"number" msgpack:"number" mapstructure:"number"`
	Display *string  `xml:"display,omitempty" json:"display,omitempty" msgpack:"display,omitempty"`
}

type XMLCardCollection struct {
	XMLName xml.Name  `xml:"cards" json:"-" msgpack:"-" mapstructure:"-"`
	Cards   []XMLCard `xml:"card" json:"cards" msgpack:"cards"`
}

type XMLCard struct {
	XMLName  xml.Name    `xml:"card" json:"-" msgpack:"-" mapstructure:"-"`
	Title    string      `xml:"title" json:"title" msgpack:"title" mapstructure:"title"`
	Subtitle *string     `xml:"subtitle,omitempty" json:"subtitle,omitempty" msgpack:"subtitle,omitempty" mapstructure:"subtitle,omitempty"`
	Image    *XMLImage   `xml:"image,omitempty" json:"image,omitempty" msgpack:"image,omitempty"`
	Buttons  []XMLButton `xml:"button" json:"buttons" msgpack:"button"`
}

type XMLImage struct {
	XMLName xml.Name `xml:"image" json:"-" msgpack:"-" mapstructure:"-"`
	URL     string   `xml:"url,attr,omitempty" json:"url" msgpack:"url"`
}

type XMLButton struct {
	XMLName xml.Name `xml:"button" json:"-" msgpack:"-" mapstructure:"-"`
	Text    string   `xml:",innerxml" json:"text" msgpack:"text" mapstructure:"text"`
	Value   *string  `xml:"value,attr,omitempty" json:"value,omitempty" msgpack:"value,omitempty" mapstructure:"value,omitempty"`
	URL     *string  `xml:"url,attr,omitempty" json:"url,omitempty" msgpack:"url,omitempty" mapstructure:"url,omitempty"`
}

type XMLSender struct {
	XMLName  xml.Name `xml:"sender" json:"-" msgpack:"-" mapstructure:"-"`
	Name     string   `xml:",innerxml" json:"name" msgpack:"name" mapstructure:"name"`
	ImageURL *string  `xml:"image-url,attr,omitempty" json:"image_url" msgpack:"name"`
}
