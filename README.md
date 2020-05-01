# go-boot-config 

go-boot-config lets you externalize your configuration so that you can work with the same application code in different environments.
You can use properties files, YAML files, environment variables, and command-line arguments to externalize configuration. 

Property values can be retrieved directly through go-boot-config module, accessed through go-boot-config ’s `Environment` abstraction

The idea behind go-boot-config is to standartized config retrieval in Go, to be able to build smart modules around.

## Simple Usage 


Before retrieving your parameter you need load the `Environment` 
```go
package main

import "github.com/furkilic/go-boot-config/pkg/go-boot-config"

func main() {
	gobootconfig.Load()
	myString, _ := gobootconfig.GetString("myString")
	myStringWithDefault := gobootconfig.GetStringWithDefault("myStringDef", "myDefault")
	// Play with confs....
}
```
```sh
my-app --myString=myValue --myStringDef=${myString}
```


You can easily retrieve few types of data (XXX being String, Bool, Int or Float): 
- GetXXX(key) : retrieve XXX type of key or error
```go
    myVal, err := gobootconfig.GetString("myString")
```
- GetXXXWithDefault(key, defaultValue) : retrieve XXX type of key or defaultValue
```go
    myValOrDefault, err := gobootconfig.GetStringWithDefault("myString", "myDefault")
```
- GetXXXSlice(key) : retrieve XXX type slice of key or error
```go
    myStringSlice, err := gobootconfig.GetStringSlice("myStringSlice")
```
- GetXXXSliceWithDefault(key, defaultValue) : retrieve XXX type slice of key or defaultValue
```go
    myStringSliceOrDefault, err := gobootconfig.GetStringSliceWithDefault("myStringSlice", []string{"default"})
```
- GetObject(key, theObjectAdress) : inject value of key in theObjectAdress or error if not found
```go
    var myObj MyObj
    err := gobootconfig.GetObj("myObject", &myObj)
```



## PropertySource Order

go-boot-config uses a very particular PropertySource order that is designed to allow sensible overriding of values.
Properties are considered in the following order:

1. Command line arguments.
2. OS environment variables.
3. RandomValuePropertySource that has properties only in random.*.
4. Profile-specific application properties  (application-{profile}.properties and YAML variants).
5. Application properties (application.properties and YAML variants).
6. Default properties 


```go
package main

import "github.com/furkilic/go-boot-config/pkg/go-boot-config"

func main() {
	gobootconfig.Load()
	name, err := gobootconfig.GetString("name")	
	// name has the value from the property sources
	// err if name does not exists
}
```

On your application config path (for example configs) you can have an `application.properties` file that provides 
a sensible default property value for `name`. 
For one-off testing, you can launch with a specific command line switch (for example, `my-app --name="GoBootConfig"`).


## Configuring Random Values

The `randomPropertySource` is useful for injecting random values (for example, into secrets or test cases). 
It can produce integers, floats, uuids, or strings, as shown in the following example:
```properties
my.secret=${random.value}
my.number=${random.int}
my.float=${random.float}
my.uuid=${random.uuid}
```


## Accessing Command Line Properties

By default, go-boot-config converts any command line option arguments 
(that is, arguments starting with `--`, such as `--server.port=9000` or `--server.port 9000` ) to a `property` and 
adds them to the go-boot-config `Environment`. 
As mentioned previously, command line properties always take precedence over other property sources.


## Application Property Files
go-boot-config loads properties from `application.properties` files in the following locations and adds them to the `Environment`:
   
1. A `/configs` subdirectory of the current directory
2. The current directory

The list is ordered by precedence (properties defined in locations higher in the list override those defined in lower locations).

:grey_exclamation: You can also use `YAML` ('.yml' or '.yaml') files as an alternative to '.properties'.

If you do not like `application.properties` as the configuration file name, you can switch to another file 
name by specifying a `go.config.name` environment property. 
You can also refer to an explicit location by using the `go.config.location` environment property 
(which is a comma-separated list of directory locations or file paths).
 The following example shows how to specify a different file name:
```sh
./my-app --go.config.name=my-app
```
The following example shows how to specify two locations:
```sh
./my-app --go.config.location=configs/default.properties,configs/override.properties
```

:warning: `go.config.name` and `go.config.location` are used very early to determine which files have to be loaded. 
They must be defined as an environment property (typically an OS environment variable or a command-line argument).

If `go.config.location` contains directories (as opposed to files, and, at runtime, 
be appended with the names generated from `go.config.name` before being loaded, including profile-specific file names). 
Files specified in `go.config.location` are used as-is, with no support for profile-specific variants, and are overridden 
by any profile-specific properties.

When custom config locations are configured by using `go.config.location`, they replace the default locations


:grey_exclamation: If you use environment variables rather than properties, most operating systems disallow period-separated key names,
 but you can use underscores instead (for example, `GO_CONFIG_NAME` instead of `go.config.name`).
 
 
 ## Profile-specific Properties

In addition to `application.properties` files, profile-specific properties can also be defined by using the following 
naming convention: `application-{profile}.properties`. The `Environment` has a set of default profiles (by default, `default`) 
that are used if no active profiles are set. 
In other words, if no profiles are explicitly activated, then properties from `application-default.properties` are loaded.

Profile-specific properties are loaded from the same locations as standard `application.properties`, with profile-specific 
files always overriding the non-specific ones.

If several profiles are specified, a first-wins strategy applies. 
For example, profiles specified by the `go.profiles.active` property are added after those configured and therefore take precedence.

:grey_exclamation: If you have specified any files in `go.config.location`, profile-specific variants of those files are 
not considered. Use directories in `go.config.location` if you want to also use profile-specific properties.


## Placeholders in Properties
The values in `application.properties` are filtered through the existing `Environment` when they are used, 
so you can refer back to previously defined values (for example, from OS environment properties).

```properties
app.name=MyApp
app.description=${app.name} is using go-boot-config
app.nonexisting=${app.what:-mydefault} is not existing so 'mydefault' will be used
```

## Using YAML Instead of Properties

`YAML` is a superset of JSON and, as such, is a convenient format for specifying hierarchical configuration data. 
go-boot-config supports YAML as an alternative to properties by using `gopkg.in/yaml.v2`


## Using YAML Instead of Properties

`YAML` is a superset of JSON and, as such, is a convenient format for specifying hierarchical configuration data. 
go-boot-config supports YAML as an alternative to properties by using `gopkg.in/yaml.v2`

## Binding rules per property source

### Set PropertySource

Property Source | Simple | Slice 
--- | --- | ---
Command line arguments | nospacelowercase, camelCase, kebab-case, or underscore_notation | Comma-separated values or Mulitple declaration
| | `--myapp.mystring=myvalue` | `--myapp.myslice=myvalue1,myvalue2` or `--myapp.myslice=myvalue1 --myapp.myslice=myvalue2`
Environment Variables | Upper case format with underscore as the delimiter. `_` should not be used within a property name | Comma-separated values
| | `MYAPP_MYSTRING=myvalue` | `MYAPP_MYSLICE=myvalue1,myvalue2`
Properties Files | nospacelowercase, camelCase, kebab-case, or underscore_notation | Comma-separated values 
| | `myapp.mystring=myvalue` | `myapp.myslice=myvalue1,myvalue2`
YAML Files | nospacelowercase, camelCase, kebab-case, or underscore_notation | Standard YAML list syntax or comma-separated values
| | `myapp: {mystring: myvalue}` | `myapp: {myslice: [myvalue1, myvalue2]}` 


### Get PropertySource

To retrieve the property `nospacelowercase`is recommended
```go
gobootconfig.GetString("myapp.mystring")
// Works with : 
// myapp.mystring
// my-app.my-string
// myApp.myString
// my_app.my_string
// MYAPP_MYSTRING
```

If you want to retrieve a specific index of a list value, you can use standard list syntax using `[index]`
```go
// myapp.myslice=val1,val2,val3
gobootconfig.GetString("myapp.myslice[1]")
// Will return : 
// val2
```