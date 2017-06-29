package sm

import (
	unipdf "github.com/jiusanzhou/unidoc/pdf"
	"os"
)

func readPdf(f string) (pages []*unipdf.PdfPage, err error) {
	fio, err := os.Open(f)
	defer fio.Close()
	reader, err := unipdf.NewPdfReader(fio)
	pages = reader.PageList
	return
}

func writePdf(f string, pages []*unipdf.PdfPage) error {

	writer := unipdf.NewPdfWriter()

	for _, p := range pages {
		writer.AddPage(p.GetPageAsIndirectObject())
	}

	fio, err := os.Create(f)
	if err != nil {
		return err
	}
	return writer.Write(fio)
}

func getSize(page *unipdf.PdfPage) (float64, float64, int64, error) {
	box, err := page.GetMediaBox()
	if err != nil {
		return 0, 0, 0, err
	}
	var rotate int64
	if page.Rotate != nil {
		rotate = *page.Rotate
	} else {
		rotate = 0
	}
	return box.Urx - box.Llx, box.Ury - box.Lly, rotate, nil
}
