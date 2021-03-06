goconfig provides parsing for simple configuration file.

A configuration is expected to be in the format

```
[ sectionName ]
key1 = some value
key2 = some other value
# we want to explain the importance and great forethought
# in this next value.
key3 = unintuitive value
[ anotherSection ]
key1 = a value
key2 = yet another value
#...
```

Blank lines are skipped, and lines beginning with `#` are considered
comments to be skipped. It is an error to have a section marker ('[]')
without a section name. `key = ` lines will set the line to a blank
value. If no section is given, the default section (`default`).

Parsing a file can be done with the ParseFile function. It will return
a `map[string]map[string]string`. For example, if the section `foo` is
defined, and `foo = bar` is specified:

```
import config "github.com/screscent/goconfig"
func getFoo() []string {
        conf = config.ParseFile("config.conf")
        return conf["foo"]["bar"]
}
```

