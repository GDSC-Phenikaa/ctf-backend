package routes

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-pdf/fpdf"
)

type CertificateRequest struct {
	Name string `json:"name"`
}

// GenerateCertificate handler creates a PDF certificate for the provided name.
func GenerateCertificate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CertificateRequest
		// You can also support GET requests by checking query parameters.
		// We'll support both body and query for convenience.
		if r.Method == http.MethodPost {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
		} else {
			req.Name = r.URL.Query().Get("name")
		}

		if req.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		// Create a landscape A4 PDF
		pdf := fpdf.New("L", "mm", "A4", "")
		pdf.AddPage()

		// Background image
		// A4 landscape dimensions: 297 x 210 mm
		pdf.ImageOptions("certificate_base.png", 0, 0, 297, 210, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

		// The specific location for the name
		// Let's use a nice large elegant font. We'll stick to Arial, bold italic for now.
		pdf.SetFont("Arial", "BI", 45)
		
		// Set Y roughly in the middle of the blank space.
		// A4 height is 210, midway is 105. Let's start placing text at Y=100.
		pdf.SetY(100)
		
		// The text color might need to match the theme. Let's use a dark gray or sort of brownish color, or just black. Let's go with a gold-ish or black color.
		// The original used gold/brown, let's set text color to that gold-like color: #c49a45 or just black #333333. Let's use #000000 black for the name since it's common.
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(297, 20, req.Name, "", 1, "C", false, 0, "")

		var buf bytes.Buffer
		if err := pdf.Output(&buf); err != nil {
			http.Error(w, "Failed to generate certificate", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\"certificate.pdf\"")
		w.Write(buf.Bytes())
	}
}

func CertificateRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/generate", GenerateCertificate())
	r.Post("/generate", GenerateCertificate())
	return r
}
