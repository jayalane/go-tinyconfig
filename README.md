
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


*Is it good to use?*

I like it.  

*What is it? *

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
