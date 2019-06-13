## Service configuration

Describes strategy, runtime parameters and other options for one executable.

```yaml
services:
	# list of services definition
# plugins definition
```

Example of service config

```yaml
label: "service name"
command: "path to executable"
args:
	- arguments
	- to
	- executable
environment:
	SOME_PARAM: some value
	ANOTHER_PARAM: another value
envFiles:
	- source environment file1
	- source environment file2
workdir: "working directory"
stop_timeout: 3s
restart_delay: 5s
restart: -1
logFile: "path to log file"
raw: false
```

### label  
  
> string, not required, default is random name

Readable name of service. Will be presented in all logs.
By default human-readable random name will be generated.

*example*: `service-1`

### command

> string, required

Path to executable binary. For non-absolute path binary has to be available through `PATH` environment variable.

*example*: `cat`

### args

> array of string, not required, default is empty

Arguments that will be passed to the `command`. 


Hi @thim81 ! Thanks for you question. Let me answer in a reverse order:

### args

> array of string, not required, 

Each argument should be place as separated array element. There is no way how to start single command joined with arguments in one line due to security issue: when all arguments are passing separately, it's impossible to make an insecure call.

A some theoretical example: assume that `command` accepts a single-line command that should be invoked by the system. In some cases (not only by hackers, but just by mistake) someone may write command like: `echo aa bb cc; ls /`. That will interpreted by the shell as a two commands: `echo aa bb cc` and `ls /` that probably not what you really want to do.

In case of separated arguments, all arguments passed not to shell but to the command it self. So it doesn't matter what symbols, escaping characters or anything else someone will put.

Anyway there is still a hack: to set `command: /usr/bin/bash`  and `args: ['-c', 'any complex command'].

However, writing manually all parameters can make a hassle. To help with it, monexec has a special CLI argument: `--generate` that will make most parameters automatically. For example:
`monexec run --generate -- nc -l 9001` creates:
```yaml
services:
- label: crystal-spangle
  command: nc
  args:
  - -l
  - "9001"
  stop_timeout: 5s
  restart_delay: 5s
  restart: -1
```

*example*: 
```yaml
args:
 - "-l"
 - "9001"
```

### environment

>  map of string=>string, not required, default is empty

Additional environment variables that will be appended over system for the service.

*example*:

```yaml
environment:
	API_TOKEN: "xx-yy-zz"
	ENDPOINT: "api.example.com"
```

### envFiles

> list of string, not required, default is empty
  
Load environment variables from file.

General environment variables processing:  
  
1. system environments  
2. environments defined by `--env, -e` (CLI) or in `environment` (config)  
3. (since `0.1.14`), environments files in order as they defined by `--env-file, -E` (CLI) or in `envFile` (config)   

File path are absolute or relative to the 
If environment file couldn't be read, the application **can still be launched** - only warning in log presented.  
  
Environment file format:  
  
* pair: KEY=VALUE    
* comment: line started with #    
* empty lines or invalid (without = symbol) ignored    
* there is no way to escape new line symbol  
  
Example of file:  
  
```env  
# some comment    
    
 it was empty lines and this is a broken record MAIL=will be replaced MAIL=owner@reddec.net SUBJECT=multiple = are supported  
```