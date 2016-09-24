package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/andlabs/ui"
)

type Mot struct {
	unMot   string
	tableau []string
	essais  int
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func (m *Mot) ChargerCSV(fichier string) {
	var tableauCSV []string
	csvfile, err := os.Open(fichier)
	Check(err)
	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 // over 9000
	reader.Comma = ','
	rawCSVdata, err := reader.ReadAll()
	Check(err)
	for _, each := range rawCSVdata {
		c := each[0]
		tableauCSV = append(tableauCSV, c)
	}
	m.tableau = tableauCSV
}

func (m *Mot) Tirage() {
	m.unMot = m.tableau[rand.Intn(len(m.tableau))]
	fmt.Println(m.unMot)
}

func (m *Mot) Dire() {
	binary, err := exec.LookPath("picospeaker")
	Check(err)
	cmd := exec.Command(binary, "-l", "fr-FR", m.unMot)
	cmd.Start()
	cmd.Wait()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	nomFichier := os.Args[1]
	var m, ok, nok Mot
	m.ChargerCSV(nomFichier)
	ok.ChargerCSV("ok.txt")
	nok.ChargerCSV("nok.txt")
	//
	err := ui.Main(func() {
		le_mot := ui.NewLabel("")
		input := ui.NewEntry()
		essaie := ui.NewButton("essayer")
		box := ui.NewVerticalBox()
		box.Append(input, false)
		box.Append(le_mot, false)
		box.Append(essaie, false)
		window := ui.NewWindow("Dictée Magique", 300, 100, false)
		window.SetChild(box)
		m.Tirage()
		m.Dire()
		essaie.OnClicked(func(*ui.Button) {
			if m.unMot == input.Text() {
				ok.Tirage()
				ok.Dire()
				m.Tirage()
				m.Dire()
			} else {
				nok.Tirage()
				nok.Dire()
			}
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	Check(err)
}
