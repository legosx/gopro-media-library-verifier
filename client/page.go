package client

type Page struct {
	totalPages int
	medias     []Media
}

func NewPage(totalPages int, medias []Media) Page {
	return Page{
		totalPages: totalPages,
		medias:     medias,
	}
}

func (p Page) TotalPages() int {
	return p.totalPages
}

func (p Page) Medias() []Media {
	return p.medias
}
