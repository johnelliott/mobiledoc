package mobiledoc

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextRenderer(t *testing.T) {
	doc := Document{
		Version: Version,
		Markups: []Markup{
			{Tag: "b"},
			{Tag: "i"},
			{Tag: "a", Attributes: Map{"href": "http://example.com"}},
		},
		Atoms: []Atom{
			{Name: "atom1", Text: "__atom1_txt__", Payload: Map{"bar": 42}},
			{Name: "atom2", Text: "__atom2_txt__", Payload: Map{"bar": 24}},
		},
		Cards: []Card{
			{Name: "card1", Payload: Map{"foo": 42}},
			{Name: "card2", Payload: Map{"foo": 42}},
		},
	}
	doc.Sections = []Section{
		{Type: CardSection, Card: &doc.Cards[0]},
		{Type: MarkupSection, Tag: "p", Markers: []Marker{
			{Type: TextMarker, Text: " then text "},
			{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, ClosedMarkups: 1, Text: "then marked up text "},
			{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[1]}, Text: "then more "},
			{Type: TextMarker, ClosedMarkups: 1, Text: "then one close markup "},
			{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[1], &doc.Markups[2]}, ClosedMarkups: 1, Text: "then 2 open 1 closed markups "},
			{Type: TextMarker, ClosedMarkups: 1, Text: "then one close markup "},
		}},
		{Type: MarkupSection, Tag: "p", Markers: []Marker{
			{Type: AtomMarker, Atom: &doc.Atoms[0]},
			{Type: AtomMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, Atom: &doc.Atoms[1]},
			{Type: AtomMarker, ClosedMarkups: 1, Atom: &doc.Atoms[0]},
		}},
		{Type: ImageSection, Source: "http://example.com/foo.png"},
		{Type: ListSection, Tag: "ul", Items: [][]Marker{
			{
				{Type: TextMarker, ClosedMarkups: 0, Text: "then first ul item "},
				{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, ClosedMarkups: 1, Text: "then second ul item 1 open 1 close markup "},
			},
			{
				{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, Text: "then another open markup "},
				{Type: TextMarker, ClosedMarkups: 1, Text: "then escaped html stuff <foo> "},
			},
		}},
		{Type: 10, Card: &doc.Cards[1]},
	}

	r := NewTextRenderer()
	r.Atoms["atom1"] = func(w *bufio.Writer, text string, payload Map) error {
		_, err := w.WriteString(fmt.Sprintf("%s", text))
		return err
	}
	r.Atoms["atom2"] = func(w *bufio.Writer, text string, payload Map) error {
		_, err := w.WriteString(fmt.Sprintf("%s", text))
		return err
	}
	r.Cards["card1"] = func(w *bufio.Writer, payload Map) error {
		_, err := w.WriteString("**card1**")
		return err
	}
	r.Cards["card2"] = func(w *bufio.Writer, payload Map) error {
		_, err := w.WriteString("**card2**")
		return err
	}

	out := `**card1** then text then marked up text then more then one close markup then 2 open 1 closed markups then one close markup __atom1_txt____atom2_txt____atom1_txt__ http://example.com/foo.png then first ul item then second ul item 1 open 1 close markup then another open markup then escaped html stuff <foo> **card2**`

	buf := &bytes.Buffer{}
	err := r.Render(buf, doc)
	assert.NoError(t, err)
	assert.Equal(t, out, buf.String())
}
