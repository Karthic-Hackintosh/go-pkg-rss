package feeder

import (
	"crypto/md5"
	"encoding/json"
	"io"
	"time"
)

type Item struct {
	// RSS and Shared fields
	Provider string
	FeedUrl  string //provider and feedurl are the meta inforamtion added before storing to kafka

	Title       string
	Links       []*Link
	Description string
	Author      Author
	Categories  []*Category
	Comments    string
	Enclosures  []*Enclosure
	Guid        *string
	PubDate     string
	Source      *Source

	// Atom specific fields
	Id           string
	Generator    *Generator
	Contributors []string
	Content      *Content
	Updated      string

	Extensions map[string]map[string][]Extension

	//for satisfying sarama encoder
	encoded []byte
	err     error
}

func (ale *Item) ensureEncoded() {
	if ale.encoded == nil && ale.err == nil {
		ale.encoded, ale.err = json.Marshal(ale)
	}
}

func (ale *Item) Length() int {
	ale.ensureEncoded()
	return len(ale.encoded)
}

func (ale *Item) Encode() ([]byte, error) {
	ale.ensureEncoded()
	return ale.encoded, ale.err
}

//the data being sent to kafka insertion shoudl satisfy the sarama.Encoder from https://github.com/Shopify/sarama/blob/master/utils.go,
func (i *Item) ParsedPubDate() (time.Time, error) {
	return parseTime(i.PubDate)
}

func (i *Item) Key() string {
	switch {
	case i.Guid != nil && len(*i.Guid) != 0:
		return *i.Guid
	case len(i.Id) != 0:
		return i.Id
	case len(i.Title) > 0 && len(i.PubDate) > 0:
		return i.Title + i.PubDate
	default:
		h := md5.New()
		io.WriteString(h, i.Description)
		return string(h.Sum(nil))
	}
}
