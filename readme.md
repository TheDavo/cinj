# About
Cinj is a command line tool that expands on the markdown syntax to make
code report generation easier. 

By modifying code block syntax in a markdown file,
automating code-heavy report generation is now an easier process, especially 
when paired with tools such as `pandoc` and `weasyprint` for PDF generation. 

Replace what would be a block of code with the Cinj command
in the markdown file as so:

```python

--- report.md
--- report content

cinj{./my_file.py}

--- more report content

```

Cinj will now look for that file and replace the contents of the parent file
with a code block of the file specified in the Cinj command.

The original Cinj file is not modified, and instead a new file with the same
base name but now with a `.md` Markdown extension. 

## Usage

To use Cinj call it from the terminal, specifying the file to be worked on. The
only allowed files are those with either a `.cinj` or `.cinj.md` extension.

```c

>> cinj ./my_report.cinj.md
>> ls
>> my_report.md my_report.cinj.md

```

A new name for the output file can be specified when calling Cinj by using the
`newname` flag. Note that since Go's `flag` package is used, all positional
arguements must come _after_ any other flags.

```c

>> cinj --newname="new_name.md" ./my_report.cinj.md 
>> ls
>> my_report.md new_name.md

```
# Language Support

## Generic Features

### Line Ranges - Not Implemented
The markdown file can call Cinj for any language to grab all the content within
a file, or grab specific line number ranges from a file. If another token is
provided, and the language is supported, such as a function or class definition
grab, the `--lines` argument is overruled.

```python

# Grab lines from line 10 to the end of the file 
cinj{./my_script.js --lines="10"}

# Grab lines from line 10 to 20
cinj{./server.c --lines="10-20"}

```

## Python

Cinj's commands can be extended to limit the scope of code copied into a
markdown file, such as a particular classes or functions

```python

# Grab a class from the file my_file.py
cinj{./my_file.py --class="ExampleClass"}

# Grab a function from the file
cinj{./my_file.py --function="example_function"}

```

Implemented:
- [x] class
- [x] functions
- [x] decorators

### Passing Both `Class` and `Function` Arguments

When both `class` and `function` arguments have a value, Cinj will look
for a `function` inside of the `class`. 
This is useful if the source file contains many classes
and only a particular `__init__` function needs to be copied over, for example.

## JavaScript

## HTML

## CSS

## Go

## Rust

## C
### C Header Files

# Error Handling - Not Implemented

Cinj will panic on by default on any error, but can be overridden with the
`--no-panic` flag. This

## File Not Found

Should the file that the Cinj command is used on is not found, an error is
thrown and not further action is taken. This error is not logged.

```c

>> cinj not_found.md
>> >> Error, file not found: not_found.md
>> 

```
## File Not Found (Inside Markdown File)

Should a file not be found during the markdown file scan, Cinj will `panic` by
default and exit execution. 

### Not Panicking - Not Implemented

When the `--no-panic` flag is set, Cinj will instead log to the console that
an error has occurred and that a `cinj.log` has been written with the latest
errors.

```c
>> cinj --no-panic my_report.md 
>> >> Error on line: 34 cinj{my_file.py --class="ExampleClass"}
>> >> Error: File my_file.py not found
>> >> Logged to cinj.log
>> ls
>> my_report.md my_report.cinj.md cinj.log
```

If the `--no-panic` flag is set, then the line with the Cinj command inside the
markdown file will be replaced with an empty line instead.
## Token Not Found - Not Implemented

If the desired function, class, etc. is not found, an error will be logged to
the `cinj.log` file and the user is notified about the error through a command
line message. Cinj will panic and stop any additional action.

If the `--no-panic` flag is set, then the line with the Cinj command inside the
markdown file will be replaced with an empty line instead.
