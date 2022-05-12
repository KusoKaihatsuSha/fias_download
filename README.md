# Downloader FIAS/GAR
App for helping downloading archives with **FIAS/GAR** files from https://fias.nalog.ru

`* App can be using with GUI mode or CLI mode`

`* App using configuration file, which could be create on tab 'settings' in GUI mode or you can create this file manually, for example seeing inside test files.`

### **Available flags**

> Run GUI/CLI mode

`-gui=true` or `-gui=false`  

> Only download mode

`-o=true` or `-o=false`

> Download full archive or delta archive 

`-full=true` or `-full=false`

### **Build**

`go build -ldflags "-s -w -H=windowsgui"`

### **Description of config file:**

```ini
Address API             - Address fias API
Full mode               - Full  base download mode
Proxy                   - proxy using(NOT kerberos)
Proxy IP                - proxy ip
Format                  - Archive type (!use zip, other not worked)
Path full(press Enter)  - folder path full base
Path delta(press Enter) - folder path delta base
Count last              - How many delta bases want download
Regions(by              - ;) - use format "01;02" for num region filter
full type               - type full base in API
delta type              - type delta base in API
```

**Workable config example:**

```json
{
  "Page20": "https://fias.nalog.ru/WebServices/Public/GetAllDownloadFileInfo",
  "Page21": false,
  "Page22": false,
  "Page23": "0.0.0.0",
  "Page24": ".zip",
  "Page25": "fias_gar\\full\\",
  "Page251": "GarXMLFullURL",
  "Page26": "fias_gar\\delta\\",
  "Page261": "GarXMLDeltaURL",
  "Page27": 7,
  "Page28": "01;02"
}
```

Screenshots:


<div style="width:50%">
<img src="/pictures/001.png" >
</div>
