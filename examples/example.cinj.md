# Header 1
This is content before the code snippet!

command used: cinj{c:/users/david/documents/projects/cinj/examples/example.py}
```python
import math

print(math.sin(0.37*math.pi))


class Test:
    def __init__(self):
        self.content = "Nothing!"

    def print_content(self):
        print(self.content)
```
Use supported commands to grab pertinent content!

command used: cinj{./example.py --class=Test}
```python
class Test:
    def __init__(self):
        self.content = "Nothing!"

    def print_content(self):
        print(self.content)
```


command used: cinj{./example.py --function=print_content}
```python
    def print_content(self):
        print(self.content)
```
This is content after!
