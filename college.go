package college
/*
 ********************** WEB SCRAPING PROGRAM ***********************
 *******************************************************************
 *******************************************************************
 This program will scrape specific websites pulling information and
 exporting it to a json file in the folder it is running in.
 *****
 Currently, this program is set up to pull information from the site
 colleges.niche.com where there are several colleges listed.
 *****
 This continues to find the next page until it gets to a clear
 ending page. Then it proceeds to add each college to an array
 and then saves the array in a file.
 *****
 This program took a minute to write because of the bugs and trying
 to learn go, properly. This was written for practice and for a
 practical app that is in the works. It was a helluva lot better
 than sitting here and hardcoding things into a json array.
 *******************************************************************
 *******************************************************************
 *******************************************************************/
import (
   "os"
   "fmt"
   "strings"
   "net/http"
   "encoding/json"
   "golang.org/x/net/html"
)
// This is the data I wanted to represent the school.
type School struct {
   Name       string       `json:"name"`
   State      string       `json:"state"`
}
// JSON array
var j []School
// URL I'm scraping the information from
var baseUrl = "https://colleges.niche.com"
func Run() {
   pageGetter(baseUrl + "/all/?LetterGroup=3-A") // First Page
   finish()
}
// This gets the page, checks if an error occured or if the response is
// OK. It then turns it into a html.Tokenizer object and sents it to
// PageFilter via parameter
func pageGetter(url string) {
   fmt.Println(url)
   response, err := http.Get(url)
   if err != nil {
      panic(err)
   }
   defer response.Body.Close()
   if response.StatusCode != http.StatusOK {
      panic(response.StatusCode)
   }
   token := html.NewTokenizer(response.Body)
   pageFilter(token)
}
// This sifts through each html tag to find the next url until the
// specified end has been reached. Once that happens, it loops to
// find the start of the list of colleges. When it finds it, it is
// sent to the toJson function.
func pageFilter(token *html.Tokenizer) {
   var found bool             // This is if I found the <a> tag that I'm currently in
   var urlAcquired bool       // This is if I have already sent the URL to PageGetter function
   var i = 0
   var clear bool = true      // This is if I have sent the token to toJson
   for clear {
      f := token.Next()
      if f == html.StartTagToken {
         tag := token.Token()
         if tag.Data == "a" && !urlAcquired && !found {
            for _, a := range tag.Attr {
               if a.Key == "class" && a.Val == "selected" {
                  if i == 1 {
                     found = true
                     break
                  }
                  i++
                  break
               }
            }
         } else if tag.Data == "a" && !urlAcquired && found {
            for _, a := range tag.Attr {
               if a.Key == "href" && !strings.Contains(a.Val, "university") {
                  pageGetter(baseUrl + a.Val)
                  urlAcquired = true
                  break
               } else if a.Key == "href" && strings.Contains(a.Val, "university") {
                  urlAcquired = true
                  break
               }
            }
         } else if tag.Data == "div" && urlAcquired {
            for _, a := range tag.Attr {
               if a.Key == "class" && a.Val == "columns" {
                  clear = true
                  toJson(token)
                  break
               }
            }
         }
      }
      if f == html.ErrorToken { // If it reaches the end of the page, stop the loop
         break
      }
   }
}
// This function takes a html.Tokenizer and loops through to find the names
// and states these colleges are located in. It puts the information into
// school objects and then adds them to the JSON array
func toJson(token *html.Tokenizer) {
  var clear bool = true
  for clear {
    t := token.Next()
    if t == html.StartTagToken {
       tag := token.Token()
       if tag.Data == "div" {
          break
       }
       if tag.Data == "a" {
          for token.Next() != html.TextToken {
          }
          tag = token.Token()
          name := tag.String()
          t = token.Next()
          for t != html.TextToken {
             t = token.Next()
          }
          token.Next()
          token.Next()
          token.Next()
          tag = token.Token()
          s := strings.Split(tag.String(), ", ")
          var state string
          if len(s) > 1 {
             state = s[1]
          } else {
             state = s[0]
          }
          j = append(j, School{name, state})
       }
    }
  }
}
func makeG() map[string][]School {
  // JSON array by State
  g := make(map[string][]School)
  for _, i := range j {
    if i.State == "\n" && !strings.Contains(i.Name, "Online") {
      continue
    }
    g[i.State] = *new([]School)
  }
  for _, i := range j {
    if i.State == "\n" && !strings.Contains(i.Name, "Online") {
      continue
    }
    g[i.State] = append(g[i.State], i)
  }
  return g
}

// This finalizes the process. It puts the JSON array into a json file
func finish() {
  b, err := json.Marshal(makeG())
  if err != nil {
    panic(err)
  }
  p, err := os.Create("Schools.json")
  if err != nil {
    panic(err)
  }
  defer p.Close()
  _, err = p.Write(b)
  if err != nil {
    panic(err)
  }
}
