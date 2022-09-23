# Walk-CLI
A utility to parse, select and transform YAML/JSON strings in a *NIX pipe.

Sometimes your scripts will perform API calls or `cat` files containing a big amount of YAML or JSON and all you need
is to transform it to something more readable, or easier to digest in the next step of the pipe.

That's what Walk-CLI is all about.

## Basic usage
`walkcli -t "template"`

* `template`: is a template written in the [GoWalker](https://github.com/theirish81/gowalker) format. It's somewhat
similar to what you would use in JavaScript template literals, as it uses the `${...}` notation to open templates,
and the dot notation to navigate data structures.

* `stdin`: the program will wait for a standard input to digest. You can input the text after running the program or use
a pipe.

### Examples

* `cat myfile.json | walkcli -t '${user.id}'`
* `cat products.json | walkcli -t '${item[0]'`

### Pretty-print and conversion
The output of `walkcli` can be any string. If you, however, are in need to output JSON or YAML, there's a couple of
switches which may improve your experience.

* `-j`: parses the output string as a JSON and beautifies it
* `-c`: adds colors to the output JSON
* `-y`: parses the output string as a YAML

Obviously if you use one of the parsers, the process will only succeed if the output string is syntactically correct.

## External templates
Sometimes the desired output needs to be something more complex than a string or two. In this case inlining the template
won't do it. To solve this problem, you can reference external template files.

in the `--template` argument, provide, instead, a path to file preceded by `file://` as in `file://my_templates/main.templ`.

Requirements:
* Template files must have an extension. We recommend `.templ`
* The template format is the one described in [GoWalker](https://github.com/theirish81/gowalker)

**IMPORTANT:** the engine will load the provided file and use it as *main* template, but will also load all the other
files in the same directory as *sub-templates*. We therefore recommend that you use one directory for task.