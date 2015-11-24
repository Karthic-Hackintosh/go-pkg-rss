package feeder

import (
	"encoding/json"
)

type Channel struct {
	Provider       string
	FeedUrl        string //provider and feedurl are the meta inforamtion added before storing to kafka
	Title          string
	Links          []Link
	Description    string
	Language       string
	Copyright      string
	ManagingEditor string
	WebMaster      string
	PubDate        string
	LastBuildDate  string
	Docs           string
	Categories     []*Category
	Generator      Generator
	TTL            int
	Rating         string
	SkipHours      []int
	SkipDays       []int
	Image          Image
	Items          []*Item
	Cloud          Cloud
	TextInput      Input
	Extensions     map[string]map[string][]Extension

	// Atom fields
	Id       string
	Rights   string
	Author   Author
	SubTitle SubTitle

	//for satisyfing sarama encoder
	encoded []byte
	err     error
}

func (ale *Channel) ensureEncoded() {
	if ale.encoded == nil && ale.err == nil {
		ale.encoded, ale.err = json.Marshal(ale)
	}
}

func (ale *Channel) Length() int {
	ale.ensureEncoded()
	return len(ale.encoded)
}

func (ale *Channel) Encode() ([]byte, error) {
	ale.ensureEncoded()
	return ale.encoded, ale.err
}

//the data being sent to kafka insertion shoudl satisfy the sarama.Encoder from https://github.com/Shopify/sarama/blob/master/utils.go,
func (c *Channel) Key() string {
	switch {
	case len(c.Id) != 0:
		return c.Id
	default:
		return c.Title
	}
}
