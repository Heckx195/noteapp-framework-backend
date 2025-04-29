package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

// ExportNotebook exports a notebook as a PDF
func ExportNotebook(c *gin.Context) {
	var request struct {
		NotebookName string `json:"notebook_name"`
		Notes        []struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"notes"`
	}

	// Bind the JSON payload
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create a new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Notebook: "+request.NotebookName)

	// Add notes to the PDF
	pdf.SetFont("Arial", "", 12)
	for _, note := range request.Notes {
		pdf.Ln(10)
		pdf.Cell(40, 10, "Note: "+note.Title)
		pdf.Ln(5)
		pdf.MultiCell(0, 10, note.Content, "", "", false)
	}

	// Output the PDF
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=notebook.pdf")
	err := pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
	}
}

// ExportNote exports a single note as a PDF
func ExportNote(c *gin.Context) {
	var request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	// Bind the JSON payload
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create a new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Note: "+request.Title)

	// Add note content to the PDF
	pdf.SetFont("Arial", "", 12)
	pdf.Ln(10)
	pdf.MultiCell(0, 10, request.Content, "", "", false)

	// Output the PDF
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=note.pdf")
	err := pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
	}
}
