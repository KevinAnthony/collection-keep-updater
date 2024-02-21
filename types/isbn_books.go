package types

type ISBNBooks []ISBNBook

func NewISBNBooks(length int) ISBNBooks {
	return make(ISBNBooks, 0, length)
}

func (b ISBNBooks) Diff(s ISBNBooks) (ISBNBooks, error) {
	diff := NewISBNBooks(0)

	for _, book := range s {
		if b.Contains(book) {
			continue
		}

		diff = append(diff, book)
	}

	return diff, nil
}

func (b ISBNBooks) Contains(book ISBNBook) bool {
	for _, l := range b {
		if l.Equals(book) {
			return true
		}
	}

	return false
}
