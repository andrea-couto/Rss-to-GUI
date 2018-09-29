//Andy Couto
//Project 1: Rust V GO
//The goal of this project was to make a GUI that would list data from
//StackOverflow's RSS feed within 50 miles of Bridgewater,MA

package main

import (
	"encoding/xml"
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

var jobs = returnJobs()

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	ChannelList []Channel `xml:"channel"`
}

type Channel struct {
	ItemList []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Category []string `xml:"category"`
	Location string `xml:"location"`
}

type modelHandler struct {

}

func returnM() *modelHandler {
	return new(modelHandler)
}

func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
	}
}

func returnJobs () []Item {
	resp, _ := http.Get("https://stackoverflow.com/jobs/feed?l=Bridgewater%2c+MA%2c+USA&u=Miles&d=50")
	bytes, _ := ioutil.ReadAll(resp.Body)

	var rss Rss

	xml.Unmarshal([]byte(bytes), &rss)

	rssChan := rss.ChannelList[0]
	jobs := rssChan.ItemList
	return jobs
}

func (mh *modelHandler) NumRows(m *ui.TableModel) int {
	return len(jobs)
}

func (mh *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {

	switch column {
	case 0:
		return ui.TableString(fmt.Sprintf("%d", row))
	case 1:
		return  ui.TableString(jobs[row].Title)
	case 2:
		var s []string
		for _, cat := range jobs[row].Category {
			s = append(s, cat)
		}
		return ui.TableString(strings.Join(s, ", "))
	case 3:
		return ui.TableString("Link")
	case 4:
		return ui.TableString(jobs[row].Location)
	}
	panic("unreachable")
}

func passToLink(row int) {
	openbrowser(jobs[row].Link)
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	if column == 3 {
		passToLink(row)
	}
}

func setupUI() {
	mainwin := ui.NewWindow("StackOverflow Bridgewater Jobs", 640, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	table := ui.NewTable(&ui.TableParams{
		Model: ui.NewTableModel(returnM()),
	})
	mainwin.SetChild(table)
	mainwin.SetMargined(true)

	table.AppendTextColumn("Num",
		0, ui.TableModelColumnNeverEditable, nil)

	table.AppendTextColumn("Title",
		1, ui.TableModelColumnNeverEditable, nil)

	table.AppendTextColumn("Categories",
		2, ui.TableModelColumnNeverEditable, nil)

	table.AppendTextColumn("Location",
		4, ui.TableModelColumnNeverEditable, nil)

	table.AppendButtonColumn("Links",
		3, ui.TableModelColumnAlwaysEditable)

	mainwin.Show()
}

//https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ui.Main(setupUI)
}