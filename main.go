package main

import (
	"fmt"

	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// "github.com/eiannone/keyboard"
	"github.com/gocolly/colly"
	"github.com/rodaine/table"
)
type PokemonProduct struct { 
	url, image, name, price string 
}

type RutorgObj struct { 
	name , size  ,url string 
}

func GetTorrent(url string) string{
  c := colly.NewCollector();
  link := "";
  c.OnXML("//a[@href]/@href", func(e *colly.XMLElement) { 
    
    if(strings.Contains(e.Text,"https://rutor.org.in/index.php?do=download&id")){
      // println(e.Text);
      link = e.Text;
    }
    
  	
  })

  c.Visit(url);

  return link
}


func Search(text string) []RutorgObj {

  table.DefaultHeaderFormatter = func(format string, vals ...interface{})   string {
    return strings.ToUpper(fmt.Sprintf(format, vals...))
  }
  tbl := table.New("Name", "Size", "Link")
  
  c := colly.NewCollector();
  var data []RutorgObj

  c.OnScraped(func(r *colly.Response) { 
  	fmt.Println(r.StatusCode, " scraped!") 
  })

  c.OnXML("html//tbody/tr", func(e *colly.XMLElement) { 
    
    obj := RutorgObj{}
    url := e.ChildAttr(".//a","href")
    obj.url = GetTorrent(url)
    
  	obj.name = e.ChildText(".//a/text()") 
  	obj.size = e.ChildText(".//td[@class='dftd3']/text()") 
    // println(obj.name,obj.size,obj.url)
    tbl.AddRow(obj.name, obj.size, obj.url)
  	data = append(data, obj)
  })
  
  m := make(map[string]string);
  m["do"]="search";
  m["subaction"]="search";
  m["x"]="0";
  m["y"]="0";
  m["story"]=text;
  
  
  c.Post("https://rutor.org.in",m);
  tbl.Print();
  return data

}

func parsePage(link string) []RutorgObj {
	c := colly.NewCollector()
  var data []RutorgObj

  c.OnRequest(func(r *colly.Request) { 
  	fmt.Println("Visiting: ", r.URL) 
  }) 
   

   
  c.OnScraped(func(r *colly.Response) { 
  	fmt.Println(r.Request.URL, " scraped!") 
  })

  c.OnXML("html//tbody/tr", func(e *colly.XMLElement) { 
    
    obj := RutorgObj{}
    url := e.ChildAttr(".//a","href")
    obj.url = GetTorrent(url)
    
  	obj.name = e.ChildText(".//a/text()") 
  	obj.size = e.ChildText(".//td[@class='dftd3']/text()") 
  	data = append(data, obj)
  })

  c.Visit(link)
  
  return data
}



func writeData(data []RutorgObj,filename string){
  // opening the CSV file 
	file, err := os.Create(filename+".csv") 
	if err != nil { 
		log.Fatalln("Failed to create output CSV file", err) 
	} 
	defer file.Close() 
 
	// initializing a file writer 
	writer := csv.NewWriter(file) 
 
	// writing the CSV headers 
	headers := []string{ 
		"name", 
		"url", 
		"size", 
	} 
	writer.Write(headers) 

  for _, d := range data { 
		// converting a PokemonProduct to an array of strings 
		record := []string{ 
			d.name, 
			d.url, 
			d.size, 
		} 
 
		// adding a CSV record to the output file 
		writer.Write(record) 
	} 
	defer writer.Flush() 
}


func SearchConsole()  {
  var text string;
  fmt.Print("#>")
  fmt.Scan(&text)
  Search(text);
  // writeData(tmp,"search")
  fmt.Println("======================================<<<<")
}


func main(){
  
  banner,_ := ioutil.ReadFile("banner.txt")
  fmt.Println(string(banner))
  // fmt.Println("q - quit\ns - search\nt - get last")
  SearchConsole()

  // for true{
  
  //     char, _, _ := keyboard.GetSingleKey()

  
  //     switch string(char){
  //       case "s":
  //         SearchConsole()
  //       case "t":
  //         getAll()
  //       case "q":
  //         os.Exit(0)

  //     }


  // }
  

}

func getAll(){
  var Rudata []RutorgObj
  var  ind = 1
  for true{
    
    var tmp = parsePage("https://rutor.org.in/page/"+fmt.Sprint(ind));
    
    fmt.Println(len(tmp))
    for _,val := range tmp{
      Rudata = append(Rudata, val)
    }
    
    ind++
    if ind>10{
      break
    }
  }
  fmt.Println("The end!!");

  writeData(Rudata,"games");

                        
}