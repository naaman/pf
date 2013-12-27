package pf

import (
  "os"
  "flag"
  "fmt"
)

var procfileFlag = flag.String("procfile", "", "Procfile path")

func main() {
  flag.Parse()

  procfileFile, err := os.Open(*procfileFlag)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  
  procfile, err := ParseProcfile(procfileFile)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  for _, e := range procfile.Entries {
    fmt.Printf("%s:%s\n", e.Type, e.Command)
  }
}
