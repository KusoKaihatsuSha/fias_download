package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"time"
)

const (
	logo = "files/Icon.png"
)

var configFile = "data.json"
var logFile = "logs.txt"

type Download struct {
	Information           Information
	InformationHost       string
	InformationParameters map[string]string
	InformationLimit      int
	FuncRun               func(interface{}) error
	FuncRes               func(interface{})
	Last                  int
	Current               int64
}

type Information []struct {
	VersionID          int    `json:"VersionId"`
	TextVersion        string `json:"TextVersion"`
	FiasCompleteDbfURL string `json:"FiasCompleteDbfUrl"`
	FiasCompleteXMLURL string `json:"FiasCompleteXmlUrl"`
	FiasDeltaDbfURL    string `json:"FiasDeltaDbfUrl"`
	FiasDeltaXMLURL    string `json:"FiasDeltaXmlUrl"`
	Kladr4ArjURL       string `json:"Kladr4ArjUrl"`
	Kladr47ZURL        string `json:"Kladr47ZUrl"`
	GarXMLFullURL      string `json:"GarXMLFullURL"`
	GarXMLDeltaURL     string `json:"GarXMLDeltaURL"`
	Date               string `json:"Date"`
}

type InformationOne struct {
	VersionID          int    `json:"VersionId"`
	TextVersion        string `json:"TextVersion"`
	FiasCompleteDbfURL string `json:"FiasCompleteDbfUrl"`
	FiasCompleteXMLURL string `json:"FiasCompleteXmlUrl"`
	FiasDeltaDbfURL    string `json:"FiasDeltaDbfUrl"`
	FiasDeltaXMLURL    string `json:"FiasDeltaXmlUrl"`
	Kladr4ArjURL       string `json:"Kladr4ArjUrl"`
	Kladr47ZURL        string `json:"Kladr47ZUrl"`
	GarXMLFullURL      string `json:"GarXMLFullURL"`
	GarXMLDeltaURL     string `json:"GarXMLDeltaURL"`
	Date               string `json:"Date"`
}

type Counter struct {
	Total   *int64
	Current *int64
}

type FabricWorkers struct {
	Jobs      chan (Job)
	Results   chan (Result)
	Done      chan bool
	Wg        *sync.WaitGroup
	startTime time.Time
	endTime   time.Time
}

type Job struct {
	id      frun
	idr     fres
	element interface{}
}

type Result struct {
	job  Job
	done bool
}

type ListVersions []struct {
	VersionID          int    `json:"VersionId"`
	TextVersion        string `json:"TextVersion"`
	FiasCompleteDbfURL string `json:"FiasCompleteDbfUrl"`
	FiasCompleteXMLURL string `json:"FiasCompleteXmlUrl"`
	FiasDeltaDbfURL    string `json:"FiasDeltaDbfUrl"`
	FiasDeltaXMLURL    string `json:"FiasDeltaXmlUrl"`
	Kladr4ArjURL       string `json:"Kladr4ArjUrl"`
	Kladr47ZURL        string `json:"Kladr47ZUrl"`
	GarXMLFullURL      string `json:"GarXMLFullURL"`
	GarXMLDeltaURL     string `json:"GarXMLDeltaURL"`
	Date               string `json:"Date"`
}

type Archive struct {
	Pg          *widget.ProgressBar
	Zip         *widget.ProgressBar
	UZip        *widget.ProgressBar
	Flfull      *widget.Label
	Fl          *widget.Label
	Name        string
	dSize       int64
	Version     string
	dURL        string
	Result      ArchiveResult
	Unziponly   bool
	Dlonly      bool
	FolderUnzip string
	Def         struct {
		VersionID          int    `json:"VersionId"`
		TextVersion        string `json:"TextVersion"`
		FiasCompleteDbfURL string `json:"FiasCompleteDbfUrl"`
		FiasCompleteXMLURL string `json:"FiasCompleteXmlUrl"`
		FiasDeltaDbfURL    string `json:"FiasDeltaDbfUrl"`
		FiasDeltaXMLURL    string `json:"FiasDeltaXmlUrl"`
		Kladr4ArjURL       string `json:"Kladr4ArjUrl"`
		Kladr47ZURL        string `json:"Kladr47ZUrl"`
		GarXMLFullURL      string `json:"GarXMLFullURL"`
		GarXMLDeltaURL     string `json:"GarXMLDeltaURL"`
		Date               string `json:"Date"`
	}
}

type ArchiveResult struct {
	Name    string
	dSize   int64
	dURL    string
	Version string
}

type Downloader struct {
	Go      *Gui
	Items   map[string]*Archive
	Count   int
	Type    int
	Gui     bool
	OnlyDl  bool
	FullCli bool
	Logs    []string
	FuncRun func(interface{}) bool
	FuncRes func(interface{})
}

type Gui struct {
	Form      fyne.Window
	Widget    *fyne.Container
	W         int
	H         int
	Title     string
	Items     map[string]fyne.CanvasObject
	ItemsText map[string]interface{}
	PositionX int
	PositionY int
	Action    interface{}
	DefSize   fyne.Size
	Elements  *ElementsForm
	Config    string
}

type GuiElement struct {
	Element fyne.CanvasObject
	Title   string
	Parent  *Gui
	Bind    func(interface{})
}

type frun func(interface{}) bool

type fres func(interface{})

type ElementsFormValue struct {
	Page20  string
	Page21  bool
	Page22  bool
	Page23  string
	Page24  string
	Page25  string
	Page251 string
	Page26  string
	Page261 string
	Page27  float64
	Page28  string
}

// ElementsForm type
// contained information about elements on the form
type ElementsForm struct {
	Tab     *container.AppTabs `Name:"Tab" Page:"Tab"`
	Page1   *container.TabItem `Name:"Main" Page:"Main"`
	Page2   *container.TabItem `Name:"Settings" Page:"Settings"`
	Page3   *container.TabItem `Name:"..." Page:"..."`
	Button1 *widget.Button     `Name:"download delta and unpack" Page:"Page1"`
	Button2 *widget.Button     `Name:"download full and unpack" Page:"Page1"`
	Button3 *widget.Button     `Name:"only download delta" Page:"Page1"`
	Button4 *widget.Button     `Name:"only download full" Page:"Page1"`
	Button5 *widget.Button     `Name:"save" Page:"Page2"`
	Button6 *widget.Button     `Name:"default" Page:"Page2"`
	Page20  *container.TabItem `Name:"Address API" Page:"Main" Default:"https://fias.nalog.ru/WebServices/Public/GetAllDownloadFileInfo"`
	Page21  *container.TabItem `Name:"Full mode" Page:"Main" Default:"false"`
	Page22  *container.TabItem `Name:"Proxy" Page:"Settings" Default:"false"`
	Page23  *container.TabItem `Name:"Proxy IP" Page:"Settings" Default:"0.0.0.0"`
	Page24  *container.TabItem `Name:"Format" Page:"Settings" Default:"zip"`
	Page25  *container.TabItem `Name:"Path full(press Enter)" Page:"Settings" Default:"C:/1/"`
	Page26  *container.TabItem `Name:"Path delta(press Enter)" Page:"Settings" Default:"C:/2/"`
	Page27  *container.TabItem `Name:"Count last" Page:"Settings" Default:"14"`
	Page28  *container.TabItem `Name:"Regions(by - ;)" Page:"Settings" Default:"11;11;33"`
	Page251 *container.TabItem `Name:"full type" Page:"Settings" Default:"FiasCompleteDbfURL"`
	Page261 *container.TabItem `Name:"delta type" Page:"Settings" Default:"FiasDeltaDbfURL"`
}

var printf = fmt.Printf
var echo = fmt.Println

// init()
// initialization data mainPage
func (o *Gui) init() {
	o.Items = make(map[string]fyne.CanvasObject)
	o.ItemsText = make(map[string]interface{})
	o.Action = new(Downloader)
	o.Action.(*Downloader).init()

	a := app.New()
	ico, _ := fyne.LoadResourceFromPath(logo)
	a.SetIcon(ico)
	o.Form = a.NewWindow(o.Title)
	vbox := container.NewVBox()
	o.Widget = vbox
	o.mainPage()
}

// xTake(string, string) string
// Find element on Form
func xTake(type_, typetype_ string) string {
	tmp, _ := reflect.TypeOf(*new(ElementsForm)).FieldByName(type_)
	if val, ex := tmp.Tag.Lookup(typetype_); ex {
		return val
	}
	return ""
}

// addETo(*fyne.Container, fyne.CanvasObject) *fyne.CanvasObject
// add element on Form
func (o *Gui) addETo(obj *fyne.Container, item fyne.CanvasObject) *fyne.CanvasObject {
	obj.Add(item)
	return &item
}

// addETo(*fyne.Container, fyne.CanvasObject) *fyne.CanvasObject
// add element on Form on Tabs
func (o *Gui) addEToT(obj *container.AppTabs, item *container.TabItem) *container.TabItem {
	obj.Append(item)
	return item
}

// mainPage()
// create mainPage
func (o *Gui) mainPage() {
	o.Widget.Objects = nil
	o.Action.(*Downloader).Go = o
	o.Elements = new(ElementsForm)
	o.Elements.Tab = container.NewAppTabs()
	o.Elements.Page1 = new(container.TabItem)
	o.Elements.Page1.Content = container.NewVBox()
	o.Elements.Page1.Text = xTake("Page1", "Name")
	o.Elements.Page2 = new(container.TabItem)
	o.Elements.Page2.Content = container.NewVBox()
	o.Elements.Page2.Text = xTake("Page2", "Name")
	o.Elements.Page3 = new(container.TabItem)
	o.Elements.Page3.Content = container.NewVBox()
	o.Elements.Page3.Text = xTake("Page3", "Name")
	o.Elements.Button1 = new(widget.Button)
	o.Elements.Button1.SetText(xTake("Button1", "Name"))
	o.Elements.Button1.OnTapped = o.Action.(*Downloader).click
	o.addETo(o.Elements.Page1.Content.(*fyne.Container), o.Elements.Button1)
	o.Elements.Button2 = new(widget.Button)
	o.Elements.Button2.SetText(xTake("Button2", "Name"))
	o.Elements.Button2.OnTapped = o.Action.(*Downloader).click
	o.addETo(o.Elements.Page1.Content.(*fyne.Container), o.Elements.Button2)
	o.Elements.Button3 = new(widget.Button)
	o.Elements.Button3.SetText(xTake("Button3", "Name"))
	o.Elements.Button3.OnTapped = o.Action.(*Downloader).clickdl
	o.addETo(o.Elements.Page1.Content.(*fyne.Container), o.Elements.Button3)
	o.Elements.Button4 = new(widget.Button)
	o.Elements.Button4.SetText(xTake("Button4", "Name"))
	o.Elements.Button4.OnTapped = o.Action.(*Downloader).clickdl
	o.addETo(o.Elements.Page1.Content.(*fyne.Container), o.Elements.Button4)
	o.Elements.Button5 = new(widget.Button)
	o.Elements.Button5.SetText(xTake("Button5", "Name"))
	o.Elements.Button5.OnTapped = o.Action.(*Downloader).clicksave
	o.addETo(o.Elements.Page2.Content.(*fyne.Container), o.Elements.Button5)
	o.Elements.Button6 = new(widget.Button)
	o.Elements.Button6.SetText(xTake("Button6", "Name"))
	o.Elements.Button6.OnTapped = o.Action.(*Downloader).clickdef
	o.addETo(o.Elements.Page2.Content.(*fyne.Container), o.Elements.Button6)
	o.addEToT(o.Elements.Tab, o.Elements.Page1)
	tabsS := container.NewAppTabs()
	o.Elements.Page20 = new(container.TabItem)
	o.Elements.Page20.Content = widget.NewEntry()
	o.Elements.Page20.Text = xTake("Page20", "Name")
	o.Elements.Page21 = new(container.TabItem)
	o.Elements.Page21.Content = widget.NewCheck(xTake("Page21", "Name"), nil)
	o.Elements.Page21.Text = xTake("Page21", "Name")
	o.Elements.Page22 = new(container.TabItem)
	o.Elements.Page22.Content = widget.NewCheck(xTake("Page22", "Name"), nil)
	o.Elements.Page22.Text = xTake("Page22", "Name")
	o.Elements.Page23 = new(container.TabItem)
	o.Elements.Page23.Content = widget.NewEntry()
	o.Elements.Page23.Text = xTake("Page23", "Name")
	o.Elements.Page24 = new(container.TabItem)
	o.Elements.Page24.Content = widget.NewSelect([]string{".zip", ".rar", ".tar"}, nil)
	o.Elements.Page24.Text = xTake("Page24", "Name")
	o.Elements.Page25 = new(container.TabItem)
	o.Elements.Page25.Content = widget.NewEntry()
	o.Elements.Page25.Text = xTake("Page25", "Name")
	o.Elements.Page26 = new(container.TabItem)
	o.Elements.Page26.Content = widget.NewEntry()
	o.Elements.Page26.Text = xTake("Page26", "Name")
	o.Elements.Page27 = new(container.TabItem)
	o.Elements.Page27.Content = widget.NewSlider(1, 90)
	o.Elements.Page27.Text = xTake("Page27", "Name")
	o.Elements.Page28 = new(container.TabItem)
	o.Elements.Page28.Content = widget.NewEntry()
	o.Elements.Page28.Text = xTake("Page28", "Name")
	o.Elements.Page251 = new(container.TabItem)
	o.Elements.Page251.Content = widget.NewEntry()
	o.Elements.Page251.Text = xTake("Page251", "Name")
	o.Elements.Page261 = new(container.TabItem)
	o.Elements.Page261.Content = widget.NewEntry()
	o.Elements.Page261.Text = xTake("Page261", "Name")
	o.addEToT(tabsS, o.Elements.Page20)
	o.addEToT(tabsS, o.Elements.Page21)
	o.addEToT(tabsS, o.Elements.Page22)
	o.addEToT(tabsS, o.Elements.Page23)
	o.addEToT(tabsS, o.Elements.Page24)
	o.addEToT(tabsS, o.Elements.Page25)
	o.addEToT(tabsS, o.Elements.Page251)
	o.addEToT(tabsS, o.Elements.Page26)
	o.addEToT(tabsS, o.Elements.Page261)
	o.addEToT(tabsS, o.Elements.Page27)
	o.addEToT(tabsS, o.Elements.Page28)
	o.addETo(o.Elements.Page2.Content.(*fyne.Container), tabsS)
	o.addEToT(o.Elements.Tab, o.Elements.Page2)
	o.addElement(o.Elements.Tab, "tabs", true)
	o.Elements.Tab.OnChanged = o.settingsTabPage
	o.addEToT(o.Elements.Tab, o.Elements.Page3)
	o.Elements.Tab.SelectTabIndex(0)
	o.Action.(*Downloader).clickload()
	o.Action.(*Downloader).turn(!o.Elements.Page21.Content.(*widget.Check).Checked)

}

// SaveJson()
// save data to config file
func (o *Downloader) SaveJson() {
	nn := ElementsFormValue{}
	nn.Page20 = o.Go.Elements.Page20.Content.(*widget.Entry).Text
	nn.Page21 = o.Go.Elements.Page21.Content.(*widget.Check).Checked
	nn.Page22 = o.Go.Elements.Page22.Content.(*widget.Check).Checked
	nn.Page23 = o.Go.Elements.Page23.Content.(*widget.Entry).Text
	nn.Page24 = o.Go.Elements.Page24.Content.(*widget.Select).Selected
	nn.Page25 = o.Go.Elements.Page25.Content.(*widget.Entry).Text
	os.MkdirAll(nn.Page25, 0644)
	nn.Page251 = o.Go.Elements.Page251.Content.(*widget.Entry).Text
	nn.Page26 = o.Go.Elements.Page26.Content.(*widget.Entry).Text
	os.MkdirAll(nn.Page26, 0644)
	nn.Page261 = o.Go.Elements.Page261.Content.(*widget.Entry).Text
	nn.Page27 = o.Go.Elements.Page27.Content.(*widget.Slider).Value
	nn.Page28 = o.Go.Elements.Page28.Content.(*widget.Entry).Text
	all, _ := json.MarshalIndent(nn, "", "  ")
	ioutil.WriteFile(configFile, all, 0775)
}

// LoadJson()
// get config file data
func (o *Downloader) LoadJson() {
	nn := ElementsFormValue{}
	b, _ := ioutil.ReadFile(configFile)
	json.Unmarshal(b, &nn)
	o.Go.Elements.Page20.Content.(*widget.Entry).SetText(nn.Page20)
	o.Go.Elements.Page21.Content.(*widget.Check).SetChecked(nn.Page21)
	o.Go.Elements.Page22.Content.(*widget.Check).SetChecked(nn.Page22)
	o.Go.Elements.Page23.Content.(*widget.Entry).SetText(nn.Page23)
	o.Go.Elements.Page24.Content.(*widget.Select).SetSelected(nn.Page24)
	o.Go.Elements.Page25.Content.(*widget.Entry).SetText(nn.Page25)
	o.Go.Elements.Page251.Content.(*widget.Entry).SetText(nn.Page251)
	o.Go.Elements.Page26.Content.(*widget.Entry).SetText(nn.Page26)
	o.Go.Elements.Page261.Content.(*widget.Entry).SetText(nn.Page261)
	o.Go.Elements.Page27.Content.(*widget.Slider).SetValue(nn.Page27)
	o.Go.Elements.Page28.Content.(*widget.Entry).SetText(nn.Page28)

}

// settingsTabPage(*container.TabItem)
// create settings tab
func (o *Gui) settingsTabPage(tab *container.TabItem) {
	o.Elements.Page27.Content.(*widget.Slider).OnChanged = o.opSli("")
	o.Elements.Page25.Content.(*widget.Entry).OnSubmitted = o.opFolder(o.Elements.Page25)
	o.Elements.Page26.Content.(*widget.Entry).OnSubmitted = o.opFolder(o.Elements.Page26)
	o.Action.(*Downloader).turn(!o.Elements.Page21.Content.(*widget.Check).Checked)
	o.Action.(*Downloader).clickload()
}

// opSli(string) func(float64)
// slider update
func (o *Gui) opSli(s string) func(float64) {
	return func(x float64) {
		o.Elements.Page27.Text = fmt.Sprintf("%g", x)
		o.Elements.Page2.Content.Refresh()
	}
}

// opFolder(*container.TabItem) func(string)
// dialog popup with change folder. Press 'Enter' for action
func (o *Gui) opFolder(s *container.TabItem) func(string) {
	return func(string) {
		dialog.ShowFolderOpen(func(t fyne.ListableURI, err error) {
			if t != nil {
				syml := `\`
				if strings.Contains(t.Path(), `/`) {
					syml = `/`
				}
				s.Content.(*widget.Entry).SetText(t.Path() + syml + strings.Split(s.Text, "(")[0] + syml)
			}
		}, o.Form)
	}
}

// ChooseFile(fyne.URIReadCloser, error)
// crop text URL to setting
func (o *Gui) ChooseFile(t fyne.URIReadCloser, err error) {
	if t != nil {
		o.Items["save_full"].(*widget.Entry).SetText(t.URI().String()[7:])
	}
}

// turn(bool)
// show/hide elements
func (o *Downloader) turn(val bool) {
	if val {
		o.Go.Elements.Button2.Hide()
		o.Go.Elements.Button4.Hide()
		o.Go.Elements.Button1.Show()
		o.Go.Elements.Button3.Show()

	} else {
		o.Go.Elements.Button2.Show()
		o.Go.Elements.Button4.Show()
		o.Go.Elements.Button1.Hide()
		o.Go.Elements.Button3.Hide()
	}
}

// xBool(interface{})
// wrapper around not Bool
func xBool(s interface{}) bool {
	switch type_ := s.(type) {
	case string:
		if type_ == "true" {
			return true
		} else {
			return false
		}
	case bool:
		return type_
	case int:
		if type_ == 1 {
			return true
		} else {
			return false
		}
	}
	return false
}

// clickdef()
// create tabs
func (o *Downloader) clickdef() {
	o.Go.Elements.Page20.Content.(*widget.Entry).SetText(xTake("Page20", "Default"))
	o.Go.Elements.Page21.Content.(*widget.Check).SetChecked(xBool(xTake("Page21", "Default")))
	o.Go.Elements.Page22.Content.(*widget.Check).SetChecked(xBool(xTake("Page22", "Default")))
	o.Go.Elements.Page23.Content.(*widget.Entry).SetText(xTake("Page23", "Default"))
	o.Go.Elements.Page24.Content.(*widget.Select).SetSelected(xTake("Page24", "Default"))
	o.Go.Elements.Page25.Content.(*widget.Entry).SetText(xTake("Page25", "Default"))
	o.Go.Elements.Page251.Content.(*widget.Entry).SetText(xTake("Page251", "Default"))
	o.Go.Elements.Page26.Content.(*widget.Entry).SetText(xTake("Page26", "Default"))
	o.Go.Elements.Page261.Content.(*widget.Entry).SetText(xTake("Page261", "Default"))
	feetFloat, _ := strconv.ParseFloat(xTake("Page27", "Default"), 64)
	o.Go.Elements.Page27.Content.(*widget.Slider).SetValue(feetFloat)
	o.Go.Elements.Page28.Content.(*widget.Entry).SetText(xTake("Page28", "Default"))
}

// clicksave()
// click save config
func (o *Downloader) clicksave() {
	o.SaveJson()
}

// clickload()
// click load config
func (o *Downloader) clickload() {
	o.LoadJson()
}

// show()
// show form
func (o *Gui) show() {
	o.Form.ShowAndRun()
}

// cl(bool, bool)
// start downloading click
func (o *Downloader) cl(dlonly, full bool) {
	o.FuncRun = o.Fill
	o.FuncRes = o.Statisticf
	o.GetList()
	o.FuncRun = o.DownloadFile
	o.FuncRes = o.Statisticd
	o.Download()
	o.FuncRun = o.Unzip
	o.FuncRes = o.Statisticu
	o.Unz(dlonly)
	o.SaveLogs()
}

// click()
// click download and unpack
func (o *Downloader) click() {
	o.OnlyDl = false
	o.clickodl()
}

// clickdl()
// click download
func (o *Downloader) clickdl() {
	o.OnlyDl = true
	o.clickodl()
}

// clickodl()
// click download and unpack considering o.OnlyDl
func (o *Downloader) clickodl() {
	o.Go.Elements.Page3.Content = container.NewVBox()
	o.cl(o.OnlyDl, o.Go.Elements.Page21.Content.(*widget.Check).Checked)
}

// addElement(fyne.CanvasObject, string, bool)
// add element on form
func (o *Gui) addElement(item fyne.CanvasObject, name string, show bool) {
	o.Items[name] = item
	o.Widget.Add(item)
	o.Form.SetContent(o.Widget)
	if show {
		o.Widget.Show()
	}
}

// init()
// init the runnable flags
func (o *Downloader) init() {
	flag.BoolVar(&o.Gui, "gui", true, "gui mode or hide cli mode")
	flag.BoolVar(&o.OnlyDl, "o", false, "if true, only download, not unpacked")
	flag.BoolVar(&o.FullCli, "full", false, "if true, download Complete Full package, else delta")
	flag.Parse()
}

// getWebSize(string) int64
// try get file size by url
func (o *Downloader) getWebSize(url_ string) int64 {
	client := &http.Client{}
	if o.Go.Elements.Page22.Content.(*widget.Check).Checked {
		proxy, _ := url.Parse(o.Go.Elements.Page23.Content.(*widget.Entry).Text)
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
	}
	for i := 1; i < 500; i++ {
		response, err := client.Get(url_)
		if err != nil {
			echo("Error while downloading", url_, "-", err)
			return 0
		}
		time.Sleep(3 * time.Millisecond)
		if response.ContentLength > 0 {
			return response.ContentLength
		}
		defer response.Body.Close()

	}
	return 0
}

// GetFullData() *ListVersions
// get list files from server
func (o *Downloader) GetFullData() *ListVersions {
	req, err := http.NewRequest("POST", o.Go.Elements.Page20.Content.(*widget.Entry).Text, nil)
	if err != nil {
		echo(err)
	}
	client := &http.Client{}
	if o.Go.Elements.Page22.Content.(*widget.Check).Checked {
		proxy, _ := url.Parse(o.Go.Elements.Page23.Content.(*widget.Entry).Text)
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
	}
	resp, err := client.Do(req)
	if err != nil {
		echo(err)
		time.Sleep(1 * time.Second)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		echo(err)
	}
	info := new(ListVersions)
	json.Unmarshal(b, info)
	return info
}

// Fill(interface{}) bool
// get data
func (o *Downloader) Fill(fileFias interface{}) bool {
	text := make(map[string]string)
	keys := reflect.ValueOf(&fileFias.(*Archive).Def).Elem()
	for i := 0; i < keys.NumField(); i++ {
		text[keys.Type().Field(i).Name] = keys.Field(i).String()
	}
	t, _ := time.Parse("02.01.2006", text["Date"])
	if o.Go.Elements.Page21.Content.(*widget.Check).Checked || o.FullCli {
		fileFias.(*Archive).Name = o.Go.Elements.Page251.Content.(*widget.Entry).Text + "[" + t.Format("2006-01-02") + "]"
	} else {
		fileFias.(*Archive).Name = o.Go.Elements.Page261.Content.(*widget.Entry).Text + "[" + t.Format("2006-01-02") + "]"
	}
	if o.Go.Elements.Page21.Content.(*widget.Check).Checked || o.FullCli {
		fileFias.(*Archive).dSize = o.getWebSize(text[o.Go.Elements.Page251.Content.(*widget.Entry).Text])
		fileFias.(*Archive).dURL = text[o.Go.Elements.Page251.Content.(*widget.Entry).Text]
		fileFias.(*Archive).Result.dURL = o.Go.Elements.Page25.Content.(*widget.Entry).Text + fileFias.(*Archive).Name + o.Go.Elements.Page24.Content.(*widget.Select).Selected
		fileFias.(*Archive).Result.dSize = fileSize(fileFias.(*Archive).Result.dURL)
	} else {
		fileFias.(*Archive).dSize = o.getWebSize(text[o.Go.Elements.Page261.Content.(*widget.Entry).Text])
		fileFias.(*Archive).dURL = text[o.Go.Elements.Page261.Content.(*widget.Entry).Text]
		fileFias.(*Archive).Result.dURL = o.Go.Elements.Page26.Content.(*widget.Entry).Text + fileFias.(*Archive).Name + o.Go.Elements.Page24.Content.(*widget.Select).Selected
		fileFias.(*Archive).Result.dSize = fileSize(fileFias.(*Archive).Result.dURL)
	}
	pbar := widget.NewProgressBar()
	upbarzip := widget.NewProgressBar()
	pbarzip := widget.NewProgressBar()
	filelabelfull := widget.NewLabel("")
	filelabel := widget.NewLabel("")
	filelabel.TextStyle.Italic = true
	filelabelfull.TextStyle.Italic = true
	url_, _ := url.Parse(fileFias.(*Archive).dURL)
	hl := widget.NewHyperlink(fileFias.(*Archive).Name, url_)

	o.Go.addETo(o.Go.Elements.Page3.Content.(*fyne.Container), container.NewHBox(hl, pbar, filelabelfull, upbarzip, filelabel, pbarzip))
	o.Go.Elements.Tab.SelectTabIndex(2)

	fileFias.(*Archive).Pg = pbar
	fileFias.(*Archive).Pg.Min = 0
	fileFias.(*Archive).Pg.Max = 100

	fileFias.(*Archive).Zip = pbarzip
	fileFias.(*Archive).Zip.Min = 0
	fileFias.(*Archive).Zip.Max = 100

	fileFias.(*Archive).UZip = upbarzip
	fileFias.(*Archive).Fl = filelabel
	fileFias.(*Archive).Flfull = filelabelfull

	return true
}

// GetList()
// fill list
func (o *Downloader) GetList() {
	i := 0
	o.Items = make(map[string]*Archive)
	for id, textdef := range *o.GetFullData() {
		count, _ := strconv.Atoi(fmt.Sprintf("%g", o.Go.Elements.Page27.Content.(*widget.Slider).Value))
		if id < count {
			fileFias := new(Archive)
			fileFias.Def = textdef
			o.Items[fmt.Sprintf("%d", i)] = fileFias
			i++
		}
		if o.Go.Elements.Page21.Content.(*widget.Check).Checked || o.FullCli {
			break
		}
	}
	o.FillWorkers()
}

// fileSize(string) int64
// get file size
func fileSize(path string) int64 {
	fi, err := os.Stat(path)
	var fileSize int64
	if err == nil {
		fileSize = fi.Size()
	}
	return fileSize
}

// WriteCounter type
// using Reader interface
type WriteCounter struct {
	Total   int64
	Current *widget.ProgressBar
	Parent  *Archive
}

// Write([]byte) (int, error)
// needed interface metod of Reader
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress()
// update downloading
func (wc WriteCounter) PrintProgress() {
	val := (wc.Total * 100 / wc.Parent.dSize)
	val2 := (wc.Total * 10 / wc.Parent.dSize)
	wc.Current.SetValue(float64(val))
	printf("\rDownloading.%s", strings.Repeat(".", int(val2))) //p("\rDownloading...%s", ToStr(int(val)))
}

// Statisticf(interface{})
// logging
func (o *Downloader) Statisticf(ar interface{}) {
	o.ToLogs("find " + ar.(Result).job.element.(*Archive).Name)
}

// Statisticd(interface{})
// logging
func (o *Downloader) Statisticd(ar interface{}) {
	o.ToLogs("download " + ar.(Result).job.element.(*Archive).Name)
}

// Statisticu(interface{})
// logging
func (o *Downloader) Statisticu(ar interface{}) {
	o.ToLogs("unzip " + ar.(Result).job.element.(*Archive).Name)
}

// Download()
// downloading
func (o *Downloader) Download() {
	for _, item := range o.Items {
		if strings.Trim(item.dURL, " ") == "" {
			continue
		}
		if o.Go.Elements.Page21.Content.(*widget.Check).Checked || o.FullCli {
			os.MkdirAll(o.Go.Elements.Page25.Content.(*widget.Entry).Text, 0644)
		} else {
			os.MkdirAll(o.Go.Elements.Page26.Content.(*widget.Entry).Text, 0644)
		}
		output, err := os.OpenFile(item.Result.dURL, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			continue
		}
		output.Close()
		r, _ := regexp.Compile(`\.(zip|rar)`)
		itog := r.ReplaceAllString(item.Result.dURL, "")
		os.RemoveAll(itog)
		if item.Result.dSize == item.dSize {
			item.Unziponly = true
		} else {
			item.Unziponly = false
		}
	}
	o.FillWorkers()
}

// Unz(bool)
// unzipped
func (o *Downloader) Unz(dlonly bool) {
	for _, item := range o.Items {
		item.Dlonly = dlonly
	}
	o.FillWorkers()
}

// DownloadFile(interface{}) error
// query and download, init pipe
func DownloadFile(ar interface{}) error {
	if ar.(*Archive).dURL == "" {
		return errors.New("Empty URL")
	}
	out, err := os.Create(ar.(*Archive).Result.dURL)
	if err != nil {
		return err
	}
	defer out.Close()
	counter := &WriteCounter{}
	counter.Current = ar.(*Archive).Pg
	counter.Parent = ar.(*Archive)
	pr, pw := io.Pipe()
	go load(ar.(*Archive).dURL, pw, int64(0), counter)
	_, err = io.Copy(out, pr)
	if err != nil {
		return err
	}
	return err
}

// load(string, *io.PipeWriter, int64, *WriteCounter)
// chunking downloading 25mb
func load(url string, w *io.PipeWriter, begin int64, counter *WriteCounter) {
	defer w.Close()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err.Error())
	}
	end := begin + 25000000
	header := fmt.Sprintf("bytes=%v-%v", begin, end)
	req.Header.Set("Range", header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusPartialContent { // OR you can use code 206
		current, err := io.Copy(w, io.TeeReader(resp.Body, counter))
		if current == int64(0) || err != nil {
			return
		}
		load(url, w, begin+current, counter)
	} else {
		if resp.StatusCode != 416 {

		}
		return
	}
}

// DownloadFile(interface{}) bool
// wrap for downloading
func (o *Downloader) DownloadFile(ar interface{}) bool {
	if !ar.(*Archive).Unziponly {
		DownloadFile(ar)
	}
	ar.(*Archive).Pg.SetValue(100)
	r, _ := regexp.Compile(`\.(zip|rar)`)
	itog := r.ReplaceAllString(ar.(*Archive).Result.dURL, "")
	ar.(*Archive).Flfull.SetText(ar.(*Archive).Result.dURL)
	ar.(*Archive).FolderUnzip = itog
	return true
}

// Unzip(interface{}) bool
// unzip
func (o *Downloader) Unzip(ar interface{}) bool {
	arr := ar.(*Archive)
	if strings.Trim(arr.dURL, " ") == "" {
		arr.Flfull.SetText("EMPTY")
		arr.Fl.SetText("ERROR! EMPTY URL")
		arr.UZip.SetValue(0)
		arr.Pg.SetValue(0)
	}
	if !arr.Dlonly {
		reader, err := zip.OpenReader(arr.Result.dURL)
		if err != nil {
			arr.Fl.SetText("ERROR. REPEAT DOWNLOADING.")
			return false
		}
		defer reader.Close()
		if err := os.MkdirAll(arr.FolderUnzip, 0755); err != nil {
			return false
		}
		arr.UZip.Max = float64(len(reader.File))
		re := regexp.MustCompile(`^\d{2,3}`)
		for _, file := range reader.File {
			pathext := ""
			newFilename := file.Name
			arr.Fl.SetText(newFilename)
			arr.dSize = file.FileInfo().Size()
			arr.UZip.SetValue(arr.UZip.Value + 1)
			regions := strings.Split(o.Go.Elements.Page28.Content.(*widget.Entry).Text, ";")
			if re.Match([]byte(newFilename)) {
				pathext = re.FindString(newFilename)
				mkd := filepath.Join(arr.FolderUnzip, pathext)
				count_ := len(regions)
				for _, val := range regions {
					if pathext != val {
						count_--
					}
				}
				if count_ == 0 {
					arr.Fl.SetText("")
					continue
				}
				os.MkdirAll(mkd, 644)
				newFilename = re.ReplaceAllString(newFilename, "")
			}
			path := filepath.Join(arr.FolderUnzip, pathext, newFilename)
			if file.FileInfo().IsDir() {
				os.MkdirAll(path, file.Mode())
				continue
			}
			fileReader, err := file.Open()
			if err != nil {
				return false
			}
			defer fileReader.Close()
			out, err := os.Create(path)
			if err != nil {
				return false
			}
			defer out.Close()
			counter := &WriteCounter{}
			counter.Current = arr.Zip
			counter.Parent = arr
			_, err = io.Copy(out, io.TeeReader(fileReader, counter))
			if err != nil {
				return false
			}
		}
		arr.Fl.SetText("UNPACK")
	}
	return true
}

// RemoveFile(string, bool) bool
// remove file
func RemoveFile(path string, isserver bool) bool {
	err := os.Remove(path)
	return err == nil
}

// CloseFile(interface{})
// close file
func CloseFile(file interface{}) {
	switch file_ := file.(type) {
	case *os.File:
		file_.Close()
	default:
		echo("Error! close")
	}
}

// ToLogs(string)
// logging
func (obj *Downloader) ToLogs(text string) {
	obj.Logs = append(obj.Logs, text)
}

// SaveLogs()
// logging
func (obj *Downloader) SaveLogs() {
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()
	log.SetOutput(f)
	for _, text := range obj.Logs {
		log.Println(text)
	}
}

// FillWorkers()
// create task
func (o *Downloader) FillWorkers() {
	work := new(FabricWorkers)
	work.startTime = time.Now()
	work.Init()
	go work.FillWorkers(o.FuncRun, o.FuncRes, o.Items)
	go work.Result()
	go work.CreateWorkerPool()
	work.End()
	work.endTime = time.Now()
}

// Init()
// init worker
func (obj *FabricWorkers) Init() {
	obj.Jobs = make(chan Job, runtime.NumCPU())
	obj.Results = make(chan Result, runtime.NumCPU())
	var wg sync.WaitGroup
	obj.Wg = &wg
	obj.Done = make(chan bool)
}

// End()
// task complete
func (obj *FabricWorkers) End() {
	<-obj.Done
}

// Worker()
// worker action
func (obj *FabricWorkers) Worker() {
	for job := range obj.Jobs {
		output := Result{job, job.id(job.element)}
		obj.Results <- output
	}
	obj.Wg.Done()
}

// CreateWorkerPool()
// run workers
func (obj *FabricWorkers) CreateWorkerPool() {
	for i := 0; i < runtime.NumCPU(); i++ {
		obj.Wg.Add(1)
		go obj.Worker()
	}
	obj.Wg.Wait()
	close(obj.Results)
}

// FillWorkers(frun, fres, interface{})
// work with task
func (obj *FabricWorkers) FillWorkers(run frun, res fres, elements interface{}) {
	if reflect.ValueOf(elements).Kind() == reflect.Map {
		v := reflect.ValueOf(elements).MapRange()
		for v.Next() {
			f := v.Value()
			job := Job{run, res, f.Interface()}
			obj.Jobs <- job
		}
	}
	close(obj.Jobs)
}

// Result()
// get result
func (obj *FabricWorkers) Result() {
	for result := range obj.Results {
		result.job.idr(result)
	}
	obj.Done <- true
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()
	runtime.Gosched()
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	gui := new(Gui)
	gui.Title = "FIAS DOWNLOADER"
	gui.init()
	if gui.Action.(*Downloader).Gui {
		gui.show()
	} else {
		gui.Action.(*Downloader).clickodl()
	}
}
