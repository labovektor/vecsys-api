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
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	iconPath := participant.Event.Icon
	if iconPath != "" {
		iconPath = strings.ReplaceAll(iconPath, "\\", "/")
		iconURL := fmt.Sprintf("http://127.0.0.1:8787/api/v1%s", iconPath)
		cleanURL, _, _ := strings.Cut(iconURL, "?")

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
	}

	pdf.Ln(50)
	pdf.SetFont("Helvetica", "B", 20)

	pdf.CellFormat(0, 10, tr("Kartu Peserta"), "", 1, "C", false, 0, "")

	pdf.CellFormat(0, 10, tr(participant.Event.Name), "", 1, "C", false, 0, "")

	pdf.Ln(10)

	pdf.SetFont("Helvetica", "", 14)
	pdf.CellFormat(50, 10, tr("Asal Sekolah"), "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, tr(": "+participant.Institution.Name), "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, tr("Jenjang/Kategori"), "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, tr(": "+participant.Category.Name), "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, tr("Region"), "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, tr(": "+participant.Region.Name), "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, tr("Nomor Peserta"), "", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(0, 10, tr(": "+participant.Id.String()), "", 1, "", false, 0, "")

	pdf.SetFont("Helvetica", "", 14)
	pdf.CellFormat(50, 10, tr("Nama"), "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, tr(": "+participant.Name), "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, tr("Email Peserta"), "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, tr(": "+participant.Email), "", 1, "", false, 0, "")

	pdf.Ln(5)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, 10, tr("Anggota:"))
	pdf.Ln(10)
	for i, m := range *participant.Biodata {
		pdf.SetFont("Helvetica", "", 14)
		pdf.CellFormat(0, 8, tr(fmt.Sprintf("%d. %s (%s)", i+1, m.Name, m.IdNumber)), "", 1, "", false, 0, "")
		pdf.Ln(5)
	}

	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(5)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, 10, tr("Contact Person:"))
	pdf.Ln(10)
	pdf.SetFont("Helvetica", "", 14)
	pdf.CellFormat(0, 8, tr(participant.Region.ContactNumber+fmt.Sprintf(" (%s)", participant.Region.ContactName)), "", 1, "", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
