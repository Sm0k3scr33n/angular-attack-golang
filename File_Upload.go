package main
import (
  "fmt"
  "io/ioutil"
  "net/http"
  "log"
)
//
// func fileUpload(w http.ResponseWriter, r *http.Request){
//   r.ParseMultipartForm(1024)
//   fileHeader := r.MultipartForm.File["uploaded"][0]
//   file, err := fileHeader.Open()
//   if err == nil {
//     data, err :=ioutil.ReadAll(file)
//     if err == nil {
//       fmt.Fprintln(w, string(data))
//     }
//   }
// }


func fileUpload(w http.ResponseWriter, r *http.Request){
  log.Println(r)
  log.Println("---")
  err := r.ParseMultipartForm(1024)
  log.Println(err)
  log.Println("---")
  file, _, err := r.FormFile("uploaded")
  log.Println(file)
  log.Println(err)

  log.Println(err)
  if err == nil {
    data, err := ioutil.ReadAll(file)
    log.Println("writing file")
    err = ioutil.WriteFile("/tmp/postfile.png", data, 0644)
    if err == nil {
      fmt.Fprintln(w, string(data))
    }
  }

}

const indexPage = `
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>Go File Upload</title>
  </head>
  <body>
    <form action="./fileupload?hello=world&thread=123" method="post" enctype="multipart/form-data">
      <input type="text" name="hello" value="ben gabbard" />
      <input type="text" name="post" value="599" />
      <input type = "file" name="uploaded">
      <input type="submit">
    </form>
  </body>
</html>
`
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, indexPage)
}
func main(){
  server := http.Server{
      Addr: "localhost:8080",
  }
  http.HandleFunc("/fileupload", fileUpload)
  http.HandleFunc("/", indexPageHandler)
  server.ListenAndServe()
}
