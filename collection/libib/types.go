package libib

type libibCSVEntries struct {
	ItemType        string `csv:"item_type"`
	Title           string `csv:"title"`
	Author          string `csv:"creators"`
	FirstName       string `csv:"first_name"`
	LastName        string `csv:"last_name"`
	ISBN            string `csv:"upc_isbn10"`
	ISBN13          string `csv:"ean_isbn13"`
	Description     string `csv:"description"`
	Publisher       string `csv:"publisher"`
	PublishDate     string `csv:"publish_date"`
	Group           string `csv:"group"`
	Tags            string `csv:"tags"`
	Notes           string `csv:"notes"`
	Length          int    `csv:"length"`
	Lexile          string `csv:"lexile"`
	CallNUmber      string `csv:"call_number"`
	DDC             string `csv:"ddc"`
	LCC             string `csv:"lcc"`
	LCCN            string `csv:"lccn"`
	OCLC            string `csv:"oclc"`
	NumberOfDisks   int    `csv:"number_of_disks"`
	NumberOfPlayers int    `csv:"number_of_players"`
	AgeGroup        string `csv:"age_group"`
	Ensemble        string `csv:"ensemble"`
	AspectRatio     string `csv:"aspect_ratio"`
	ESRB            string `csv:"esrb"`
	Rating          string `csv:"rating"`
	Review          string `csv:"review"`
	ReviewDate      string `csv:"review_date"`
	Status          string `csv:"status"`
	Begin           string `csv:"begin"`
	Completed       string `csv:"completed"`
	Added           string `csv:"added"`
	Copies          int    `csv:"copies"`
}
