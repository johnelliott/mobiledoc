package mobiledoc

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
)

// TextRenderer implements a basic text renderer.
type TextRenderer struct {
	Atoms map[string]func(*bufio.Writer, string, Map) error
	Cards map[string]func(*bufio.Writer, Map) error
}

// NewTextRenderer creates a new TextRenderer.
func NewTextRenderer() *TextRenderer {
	return &TextRenderer{
		Atoms: make(map[string]func(*bufio.Writer, string, Map) error),
		Cards: make(map[string]func(*bufio.Writer, Map) error),
	}
}

// Render will render the document to the provided writer.
func (r *TextRenderer) Render(w io.Writer, doc Document) error {
	// wrap writer
	bw := bufio.NewWriter(w)

	// render sections
	for _, section := range doc.Sections {
		err := r.renderSection(bw, section)
		if err != nil {
			return err
		}
	}

	// flush buffer
	err := bw.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (r *TextRenderer) renderSection(w *bufio.Writer, section Section) error {
	// select sub renderer based on type
	switch section.Type {
	case MarkupSection:
		return r.renderMarkupSection(w, section)
	case ImageSection:
		return r.renderImageSection(w, section)
	case ListSection:
		return r.renderListSection(w, section)
	case CardSection:
		return r.renderCardSection(w, section)
	}

	return nil
}

func (r *TextRenderer) renderMarkupSection(w *bufio.Writer, section Section) error {
	// render markers
	err := r.renderMarkers(w, section.Markers)
	if err != nil {
		return err
	}

	return nil
}

func (r *TextRenderer) renderImageSection(w *bufio.Writer, section Section) error {
	// parse url
	src, err := url.Parse(section.Source)
	if err != nil {
		return err
	}

	// write tag
	_, err = w.WriteString(fmt.Sprintf(" %s ", src.String()))
	if err != nil {
		return err
	}

	return nil
}

func (r *TextRenderer) renderListSection(w *bufio.Writer, section Section) error {
	// write all items
	for _, item := range section.Items {
		// render markers
		err := r.renderMarkers(w, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TextRenderer) renderCardSection(w *bufio.Writer, section Section) error {
	// get card renderer
	renderer, ok := r.Cards[section.Card.Name]
	if !ok {
		return fmt.Errorf("missing card renderer")
	}

	// call renderer
	err := renderer(w, section.Card.Payload)
	if err != nil {
		return err
	}

	return nil
}

func (r *TextRenderer) renderMarkers(w *bufio.Writer, markers []Marker) error {
	// prepare stack
	stack := markupStack{}

	// write all markers
	for _, marker := range markers {
		// write opening markups
		for _, markup := range marker.OpenMarkups {
			// push markup
			stack.push(markup)
		}

		// write marker
		switch marker.Type {
		case TextMarker:
			// write text
			_, err := w.WriteString(marker.Text)
			if err != nil {
				return err
			}
		case AtomMarker:
			// just print atom text
			_, err := w.WriteString(marker.Atom.Text)
			if err != nil {
				return err
			}
		}

		// close markups
		// nothing to do because it's all string writing and no state management
	}

	return nil
}
