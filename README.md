# Simple Properties


[![build](https://github.com/codesqueak/simpleProperties/actions/workflows/build.yml/badge.svg)](https://github.com/codesqueak/simpleProperties/actions/workflows/build.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)


A simple-to-use out of the box, properties loader


If you find this project useful, you may want to [__Buy me a
Coffee!__ :coffee:](https://www.buymeacoffee.com/codesqueak) Thanks :thumbsup:


## How to use

### What gets loaded

The default properties handler loads property values from files under the `/resource` project directory.  These can be a mix of basic properties (.properties), 
yaml (.yaml) or JSON files (.json).  Load order priority is .yaml least, .json middle to .properties highest.

#### File names

The first file(s) to check & load is `bootstrap.<yaml/json/properties>`.  This cannot contain expressions for evaluation, i.e. properties are just `key=value` type. Any 
expression found will be a treated as a fatal error.  This loading takes place when the `Properties` struct is created and a pointer is handed to the application. e.g.

```
properties := simpleProperties.DefaultProperties()

which is a shorthand way of saying

var properties *simpleProperties.Properties = simpleProperties.DefaultProperties()
```

The second stage of loading looks for files named `application.<yaml/json/properties>`.  This is triggered by the following operation.

```
properties.Load()
```

If any property loaded has the key of profile, then any associated values are used to load further files.  For example, if the following exists:

**Note:** The evaluation of property expression is carried out once the load operations have completed as part of the `Load()` call

`
profile=dev,xyzzy
`

then the following files are checked and loaded:

```
application_dev.<yaml/json/properties>
application_xyzzy.<yaml/json/properties>
```

In summary, loading file order is

```
boostrap.<yaml/json/properties>
application.<yaml/json/properties>
application_<profile_name>.<yaml/json/properties>
```

#### CLI Properties

Properties can also be added via the command line. For example, if the following where on the command line...

```-p1=abc something -=pqr -ghi= -p2=plugh```

then the three properties `p1=abc`, `ghi`= and `p2=plugh` would be loaded. Note that CLI properties are the
highest priority and will override anything loaded from files.

### Property Expressions and Default Values

Expressions can be used the RHS of property declarations. Each named value is delimited by `${}`. A default value can also be specified by adding a colon after
the property name and then stating the default. e.g. `${value:defaultValue}`

If we had the following properties:

```
p1=Hello ${name}
name=Fred
p2=Hello ${another:World}
```

Then this would evaluate to:

```
p1=Hello Fred
name=Fred
p2=Hello World
```

Note: related properties do not need to be in one file.  In the example above, each line could be 
in separate files and evaluation of expressions only occurs once all loading is completed (Including CLI properties)

### Get the library

Add this:

```
go get -u github.com/codesqueak/simpleProperties/pkg/simpleProperties@vlatest
```

### In Project Code to Import

Add this 

```
import (
	"github.com/codesqueak/simpleProperties/pkg"
)
```


### In project code to use

Add this

```
	properties := simpleProperties.DefaultProperties()
	properties.Load()
	var p1 = properties.GetProperty("p1")
```
