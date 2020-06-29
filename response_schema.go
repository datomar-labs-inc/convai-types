package ctypes

import "encoding/xml"

type Response struct {
	Messages []Message `json:"messages" mapstructure:"messages" msgpack:"messages"`
}

type ResponseConfig struct {
	Basic           MessageConfig            `json:"basic" mapstructure:"basic" msgpack:"basic"`
	SendNow         bool                     `json:"sendNow" mapstructure:"sendNow" msgpack:"sendNow"`
	CustomResponses map[string]MessageConfig `json:"customResponses" mapstructure:"customResponses" msgpack:"customResponses"`
}

type Message struct {
	ShouldBatch bool       `json:"shouldBatch" mapstructure:"shouldBatch" msgpack:"shouldBatch"`
	GraphID     *int64     `json:"graphId" mapstructure:"graphId" msgpack:"graphId"`
	NodeID      *int64     `json:"nodeId" mapstructure:"nodeId" msgpack:"nodeId"`
	Message     XMLMessage `json:"data" mapstructure:"data" msgpack:"data"`
	Seq         int        `json:"seq" mapstructure:"seq" msgpack:"seq"`
}

type MessageConfig struct {
	ResponseXML string `json:"responseXML" mapstructure:"responseXML" msgpack:"responseXML"`
}

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
	XMLName  xml.Name `xml:"qr" json:"-" msgpack:"-" mapstructure:"-"`
	Text     string   `xml:",innerxml" json:"text" msgpack:"text" mapstructure:"text"`
	Value    *string  `xml:"value,attr,omitempty" json:"value,omitempty" msgpack:"value,omitempty" mapstructure:"value,omitempty"`
	Phone    bool     `xml:"phone,attr,omitempty" json:"phone,omitempty" msgpack:"phone,omitempty" mapstructure:"phone,omitempty"`
	Email    bool     `xml:"email,attr,omitempty" json:"email,omitempty" msgpack:"email,omitempty" mapstructure:"email,omitempty"`
	Image    *string  `xml:"image,attr,omitempty" json:"image,omitempty" msgpack:"image,omitempty" mapstructure:"image,omitempty"`
	ImageURL *string  `xml:"imageURL,attr,omitempty" json:"imageUrl,omitempty" msgpack:"imageUrl,omitempty" mapstructure:"imageUrl,omitempty"`
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
	ID      string   `xml:"id,attr" json:"id" msgpack:"id"`
	Width   *uint64  `xml:"width,attr,omitempty" json:"width" msgpack:"width"`
	Height  *uint64  `xml:"height,attr,omitempty" json:"height" msgpack:"height"`
	X       *uint64  `xml:"x,attr,omitempty" json:"x" msgpack:"x"`
	Y       *uint64  `xml:"y,attr,omitempty" json:"y" msgpack:"y"`
	URL     string   `xml:"url,attr,omitempty" json:"url" msgpack:"url"`
}

type XMLButton struct {
	XMLName xml.Name `xml:"button" json:"-" msgpack:"-" mapstructure:"-"`
	Text    string   `xml:",innerxml" json:"text" msgpack:"text" mapstructure:"text"`
	Value   *string  `xml:"value,attr,omitempty" json:"value,omitempty" msgpack:"value,omitempty" mapstructure:"value,omitempty"`
	URL     *string  `xml:"url,attr,omitempty" json:"url,omitempty" msgpack:"url,omitempty" mapstructure:"url,omitempty"`
}

type XMLSender struct {
	XMLName xml.Name `xml:"sender" json:"-" msgpack:"-" mapstructure:"-"`
	Name    string   `xml:",innerxml" json:"name" msgpack:"name" mapstructure:"name"`
	Persona *string  `xml:"persona,attr,omitempty" json:"persona" msgpack:"persona"`
}

type XMLTextRandomizer struct {
	XMLName xml.Name `xml:"random" json:"-" msgpack:"-" mapstructure:"-"`
	Text    []string `xml:"text" json:"text" msgpack:"text"`
}
