
GoLang Config Utility
=====================

*Why a config utility?*

The thing that makes most Java frameworks impossible to work with is
the initialization process needed to support the configuration
mechanisms.  Most everything needs a little bit of config, whether
embedded libraries or the main application, but it would be nice to
have something that allows for that without forcing opinions on every
bit of code in the call stack.  That's my diatribe.  This library is a
super short bit of code that builds a map of strings to strings with
default values and an easy way to find out from the command line what
are the options and defaults.

First little feature creep: I've added a key "configEnvVar" which if
set in config.txt will make the code load a file "config_$configEnvVar.txt"

e.g.:

config.txt has this line:

```
configEnvVar = PROD
```

then config_PROD.txt will be loaded after and overriding config.txt.
This enables docker images to be generated, one image that runs with different
configs in different environments (e.g. QA vs. PROD).  

Second little feature creep: For each key in a file it is parsing, it
will check any environment variables named like
"TINYCONFIG_OVERRIDE_*" and set a key, equal to whatever * matches,
set equal to the value of the env variable.  It does this check after
everything else happens (so it won't cause a reload of the files in
feature creep number one.  


*Is it good to use?*

I like it.  

*What is it?*

```
var theConfig config.Config
var defaultConfig = `#
numParsedMessageListener=4
metaPort=9093
metaHost=localhost
portList=9090,9092,9094
maxConnToUse=5
messagesBuffer=10000
parsedMsgBuffer=1000000  // these now work
syslogMode1isFront2IsBack3IsBoth=3
randomString=RandomString
#  Hey there
`

if len(os.Args) > 1 && os.Args[1] == "--dumpConfig" {
	log.Println(defaultConfig)
	return
}
var err error
theConfig, err = config.ReadConfig("config.txt", defaultConfig)
log.Println("Config", theConfig)
if err != nil {
	log.Println("Error opening config.txt", err.Error())
	if theConfig == nil {
		os.Exit(11)
	}
}
parsedMessagesChannel = make(chan parsedMessage, theConfig["parsedMsgBuffer"].IntVal)

```

*Who owns this code?*

Chris Lane http://github.com/jayalane

*Adivce for starting out*

If you integrate, please let me or them know of your experience and
any suggestions for improvement.

The current API can best be seen in the _test files probably.  

*Requirements*

None at present.  
