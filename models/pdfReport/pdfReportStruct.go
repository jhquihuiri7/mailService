package pdfReport

import (
	"fmt"
	"github.com/benoitmasson/plotters/piechart"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"image/color"
	"log"
	"os"
)

type PdfReport struct {
	ClientName string
	Limits     []int
	ErrorCount int
	ErrorMail  []string
}

func (p *PdfReport) GenerateChart() *os.File {
	plt := plot.New()
	pie1, err := piechart.NewPieChart(plotter.Values{float64(p.Limits[1] - p.Limits[0] - p.ErrorCount)})
	if err != nil {
		log.Fatal(err)
	}
	pie1.Total = float64(p.Limits[1] - p.Limits[0])
	pie1.Labels.Nominal = []string{"Enviados"}
	pie1.Labels.Values.Show = true
	pie1.Labels.Values.Percentage = true
	pie1.Color = color.RGBA{0, 255, 0, 255}
	pie2, err := piechart.NewPieChart(plotter.Values{float64(p.ErrorCount)})
	if err != nil {
		log.Fatal(err)
	}
	pie2.Offset.Value = float64(p.Limits[1] - p.Limits[0] - p.ErrorCount)
	pie2.Total = float64(p.Limits[1] - p.Limits[0])
	pie2.Labels.Nominal = []string{"No enviados"}
	pie2.Labels.Values.Show = true
	pie2.Labels.Values.Percentage = true
	pie2.Color = color.RGBA{200, 51, 73, 1}

	plt.HideAxes()
	plt.Title.Text = "Registro de correos enviados"
	plt.Title.TextStyle.Font.Size = 20
	plt.Title.Padding = -50

	plt.Legend.Add(fmt.Sprintf("%s: %d", "Correos enviados con éxito", p.Limits[1]-p.Limits[0]-p.ErrorCount), pie1)
	plt.Legend.Add(fmt.Sprintf("%s: %d", "Correos enviados fallidos", p.ErrorCount), pie2)
	plt.Add(pie1, pie2)
	nf, err := os.Create("reportPlot.png")
	if err != nil {
		log.Fatal(err)
	}
	w, err := plt.WriterTo(500, 500, "png")
	if err != nil {
		log.Fatal(err)
	}
	w.WriteTo(nf)
	return nf
}
func (p *PdfReport) GenerateBulkReport() *os.File {
	plotFile := p.GenerateChart()
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 20, 20)
	m.Row(12, func() {
		m.Text(fmt.Sprintf("%s %s",
			"Reporte de correos masivos enviador por", p.ClientName),
			props.Text{
				Align: consts.Center,
				Size:  15,
				Style: consts.Bold,
			},
		)
	})
	m.Row(12, func() {
		m.Text(fmt.Sprintf("%s %d - %d",
			"Correos entre los números", p.Limits[0], p.Limits[1]),
			props.Text{
				Align: consts.Center,
				Size:  12,
				Style: consts.Italic,
			},
		)
	})
	m.Row(100, func() {
		m.Col(12, func() {
			m.FileImage(plotFile.Name(), props.Rect{Center: true})
		})
	})
	if p.ErrorMail != nil {
		m.Row(10, func() {})
		m.Row(10, func() {
			m.Text("Corros no enviados",
				props.Text{
					Align: consts.Center,
					Size:  12,
					Style: consts.BoldItalic,
				},
			)
		})
		for i, v := range p.ErrorMail {
			m.Row(10, func() {
				m.Text(fmt.Sprintf("%d) %s", i+1, v),
					props.Text{
						Align: consts.Left,
						Size:  12,
						Style: consts.Italic,
					},
				)
			})
		}
	}
	buffer, err := m.Output()
	if err != nil {
		log.Fatal(err)
	}
	nf, err := os.CreateTemp(os.TempDir(), "BulkReport*.pdf")
	if err != nil {
		log.Fatal(err)
	}
	_, err = nf.Write(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return nf
}
