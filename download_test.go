package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_create_test_config(t *testing.T) {
	fmt.Println(t.Name())
	configFile = "testdata.json"
	logFile = "testlogs.txt"
	filebody := `{
  "Page20": "https://fias.nalog.ru/WebServices/Public/GetAllDownloadFileInfo",
  "Page21": false,
  "Page22": false,
  "Page23": "0.0.0.0",
  "Page24": ".zip",
  "Page25": "delete_this_folder_manual\\full\\",
  "Page251": "GarXMLFullURL",
  "Page26": "delete_this_folder_manual\\delta\\",
  "Page261": "GarXMLDeltaURL",
  "Page27": 3,
  "Page28": "01;02"
}`
	ioutil.WriteFile(configFile, []byte(filebody), 0755)
}

func Test_downloading(t *testing.T) {
	fmt.Println(t.Name())
	gui := new(Gui)
	gui.Title = "-"
	gui.init()
	gui.Action.(*Downloader).clickodl()
}

func Test_files(t *testing.T) {
	fmt.Println(t.Name())
	searchFiles("delete_this_folder_manual")
}

func Test_delete_test_config(t *testing.T) {
	fmt.Println(t.Name())
	os.RemoveAll(configFile)
	os.RemoveAll(logFile)
}

func searchFiles(folder string) {
	fi, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println(err)
	}
	for i := range fi {
		if !fi[i].IsDir() {
			fmt.Println(fi[i].Name(), fi[i].Size())
		} else {
			searchFiles(folder + "/" + fi[i].Name())
		}
	}
}

func Test_delete_test_files(t *testing.T) {
	fmt.Println(t.Name())
	os.RemoveAll("delete_this_folder_manual")
}

func Test_create_example_config(t *testing.T) {
	fmt.Println(t.Name())
	configFile = "example_config_data.json"
	filebody := `{
  "Page20": "https://fias.nalog.ru/WebServices/Public/GetAllDownloadFileInfo",
  "Page21": false,
  "Page22": false,
  "Page23": "0.0.0.0",
  "Page24": ".zip",
  "Page25": "thisfolder\\full\\",
  "Page251": "GarXMLFullURL",
  "Page26": "thisfolder\\delta\\",
  "Page261": "GarXMLDeltaURL",
  "Page27": 7,
  "Page28": "01;02;03;04"
}`
	ioutil.WriteFile(configFile, []byte(filebody), 0755)
}
