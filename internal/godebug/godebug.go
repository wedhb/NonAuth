package godebug

import (
    "os"
    "strings"
)

// Get returns the value of the named GODEBUG key, or "" if not present.
// GODEBUG is of the form "key=val,key2=val2,..."
func Get(name string) string {
    s := os.Getenv("GODEBUG")
    for s != "" {
        var key, val string
        i := strings.IndexByte(s, ',')
        if i < 0 {
            key = s
            s = ""
        } else {
            key = s[:i]
            s = s[i+1:]
        }
        if eq := strings.IndexByte(key, '='); eq >= 0 {
            key, val = key[:eq], key[eq+1:]
        }
        if key == name {
            return val
        }
    }
    return ""
}
