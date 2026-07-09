package util

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/httpimg"
	"github.com/labovector/vecsys-api/entity"
)

type PdfDiag struct {
	Ok        bool
	Err       string
	Bytes     int
	FontCount int
	StartHex  string
}

func checkPdf(pdf *gofpdf.Fpdf, label string) {
	if !pdf.Ok() {
		log.Printf("PDF DIAG [%s]: NOT OK - err=%v", label, pdf.Err())
	} else {
		log.Printf("PDF DIAG [%s]: OK", label)
	}
}

func GenerateCard(participant *entity.Participant) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	checkPdf(pdf, "after New")

	pdf.AddPage()
	checkPdf(pdf, "after AddPage")

	iconPath := participant.Event.Icon
	if iconPath != "" {
		log.Printf("PDF DIAG: iconPath = %q", iconPath)
		iconPath = strings.ReplaceAll(iconPath, "\\", "/")
		iconURL := fmt.Sprintf("http://127.0.0.1:8787/api/v1%s", iconPath)
		cleanURL, _, _ := strings.Cut(iconURL, "?")

		httpimg.Register(pdf, cleanURL, "")
		checkPdf(pdf, "after httpimg.Register")
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
		checkPdf(pdf, "after pdf.Image")
	} else {
		log.Printf("PDF DIAG: iconPath is empty, skipping image block")
	}

	pdf.Ln(50)
	pdf.SetFont("Helvetica", "B", 20)
	checkPdf(pdf, "after SetFont Helvetica B 20")

	pdf.CellFormat(0, 10, "Kartu Peserta", "", 1, "C", false, 0, "")
	checkPdf(pdf, "after CellFormat Kartu Peserta")

	pdf.CellFormat(0, 10, participant.Event.Name, "", 1, "C", false, 0, "")
	checkPdf(pdf, "after CellFormat Event.Name")

	pdf.Ln(10)

	pdf.SetFont("Helvetica", "", 14)
	pdf.CellFormat(50, 10, "Asal Sekolah", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Institution.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Jenjang/Kategori", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Category.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Region", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Region.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Nomor Peserta", "", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(0, 10, ": "+participant.Id.String(), "", 1, "", false, 0, "")

	pdf.SetFont("Helvetica", "", 14)
	pdf.CellFormat(50, 10, "Nama", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Email Peserta", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Email, "", 1, "", false, 0, "")

	pdf.Ln(5)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, 10, "Anggota:")
	pdf.Ln(10)
	for i, m := range *participant.Biodata {
		pdf.SetFont("Helvetica", "", 14)
		pdf.CellFormat(0, 8, fmt.Sprintf("%d. %s (%s)", i+1, m.Name, m.IdNumber), "", 1, "", false, 0, "")
		pdf.Ln(5)
	}

	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(5)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, 10, "Contact Person:")
	pdf.Ln(10)
	pdf.SetFont("Helvetica", "", 14)
	pdf.CellFormat(0, 8, participant.Region.ContactNumber+fmt.Sprintf(" (%s)", participant.Region.ContactName), "", 1, "", false, 0, "")

	checkPdf(pdf, "before Output")

	// Write debug file to /tmp for container-level inspection
	f, fErr := os.Create("/tmp/vecsys-debug.pdf")
	if fErr == nil {
		var tmpBuf bytes.Buffer
		tmpErr := pdf.Output(&tmpBuf)
		if tmpErr == nil {
			f.Write(tmpBuf.Bytes())
			log.Printf("PDF DIAG: wrote %d bytes to /tmp/vecsys-debug.pdf", len(tmpBuf.Bytes()))
		} else {
			log.Printf("PDF DIAG: Output to debug file failed: %v", tmpErr)
		}
		f.Close()
	} else {
		log.Printf("PDF DIAG: could not create /tmp/vecsys-debug.pdf: %v", fErr)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("PDF DIAG: Output failed: %v", err)
		return nil, err
	}

	pdfBytes := buf.Bytes()
	log.Printf("PDF DIAG: generated %d bytes, first 8 hex = %x", len(pdfBytes), pdfBytes[:min(8, len(pdfBytes))])

	if !pdf.Ok() {
		log.Printf("PDF DIAG: final pdf.Ok()=false, err=%v", pdf.Err())
	}

	return pdfBytes, nil
}
