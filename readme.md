
# About
Cinj is a command line tool that expands on the markdown syntax to make
code report generation easier. 

By modifying code block syntax a little bit,
automating report generation is now an easier process, especially when paired
with tools such as `pandoc` and `weasyprint` for PDF generation. 

Replace what would be the block of code with the Cinj command
in the markdown file as so:

```python

--- report.md
--- report content

cinj{./my_file.py}

--- more report content

```

Cinj will now look for that file and replace the contents of the parent file
with a code block of the content searched file.

The original Cinj file is not modified, and instead a new file with the same
name but without the Cinj extension is created.

```c

>> cinj my_report.cinj.md
>> ls
>> my_report.md my_report.cinj.md

```

A new name for the output file can be specified when calling Cinj

```c

>> cinj my_report.md --name="new_name.md"
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

## Python - Not Implemented

Cinj's commands can be extended to limit the scope of code copied into a
markdown file, such as a particular classes or functions

```python

# Grab a class from the file my_file.py
cinj{./my_file.py --class="ExampleClass"}

# Grab a function from the file
cinj{./my_file.py --func="example_function"}

# Opt out of including decorator names in the content grab, default is true
cinj{./my_files.py --decorators="false"}

```

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

## File Not Found - Not Implemented

Should the file that the Cinj command is used on is not found, an error is
thrown and not further action is taken. This error is not logged.

```c

>> cinj not_found.md
>> >> Error: file not_found.md not found
>> 

```
## File Not Found (Inside Markdown File) - Not Implemented

Should a file not be found during the markdown file scan, Cinj will `panic` by
default and exit execution. This can be overridden by using the `--no-panic`
flag.

### Not Panicking - Not Implemented

When the `--no-panic` flag is set, Cinj will instead log to the console that
an error has occurred and that a `cinj.log` has been written with the latest
errors.

```c
>> cinj my_report.md --no-panic
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
