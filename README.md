pf
==

Simple Procfile parser in Go.

Usage
=====

```
procfileFile, _ := os.Open("/home/user/project/Procfile")
procfile, _ := ParseProcfile(procfileFile)
for _, e := range procfile.Entries {
  fmt.Printf("%s:%s\n", e.Type, e.Command)
}
```
