package util

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/httpimg"
	"github.com/labovector/vecsys-api/entity"
)

func GenerateCard(participant *entity.Participant) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	iconPath := participant.Event.Icon
	iconPath = strings.ReplaceAll(iconPath, "\\", "/")
	iconURL := fmt.Sprintf("http://127.0.0.1:8787/api/v1%s", iconPath)
	cleanURL := strings.Split(iconURL, "?")[0]

	httpimg.Register(pdf, cleanURL, "")
	imgFormat := ""
	if strings.HasSuffix(strings.ToLower(cleanURL), ".jpg") {
		imgFormat = "JPG"
	} else if strings.HasSuffix(strings.ToLower(cleanURL), ".jpg") || strings.HasSuffix(strings.ToLower(cleanURL), ".jpeg") {
		imgFormat = "JPG"
	} else if strings.HasSuffix(strings.ToLower(cleanURL), ".gif") {
		imgFormat = "GIF"
	} else {
		imgFormat = "PNG"
	}

	pdf.Image(cleanURL, 90, 12, 30, 30, false, imgFormat, 0, "")
	pdf.Ln(50)
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 10, "Kartu Peserta", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 10, participant.Event.Name, "", 1, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 14)
	pdf.CellFormat(50, 10, "Asal Sekolah", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Institution.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Jenjang/Kategori", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Category.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Region", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Region.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Nomor Peserta", "", 0, "", false, 0, "")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, ": "+participant.Id.String(), "", 1, "", false, 0, "")

	pdf.SetFont("Arial", "", 14)
	pdf.CellFormat(50, 10, "Nama", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Email Peserta", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Email, "", 1, "", false, 0, "")

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Anggota:")
	pdf.Ln(10)
	for i, m := range *participant.Biodata {
		pdf.SetFont("Arial", "", 14)
		pdf.CellFormat(0, 8, fmt.Sprintf("%d. %s (%s)", i+1, m.Name, m.IdNumber), "", 1, "", false, 0, "")
		pdf.Ln(5)
	}

	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Contact Person:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 14)
	pdf.CellFormat(0, 8, participant.Region.ContactNumber+fmt.Sprintf(" (%s)", participant.Region.ContactName), "", 1, "", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	pdfBytes := buf.Bytes()

	return pdfBytes, nil
}
