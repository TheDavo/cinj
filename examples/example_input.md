# Header 1
This is content before the code snippet!

command used: cinj{c:/users/david/documents/projects/cinj/examples/example.py}
```python
import math
from abc import ABC
print(math.sin(0.37*math.pi))


class Test:
    def __init__(self):
        self.content = "Nothing!"

    def print_content(self):
        print(self.content)

    @decorator1
    @decorator2
    def todo(self):
        pass


class MyABC(ABC):
    def __init__(self):
        pass

    @decorator3
    @decorator4
    def abc_func(self):
        pass

    def todo(self):
        pass
```

Use supported commands to grab pertinent content!


command used: cinj{./example.py --class=Test}
```python
# in file example.py
class Test:
    def __init__(self):
        self.content = "Nothing!"

    def print_content(self):
        print(self.content)

    @decorator1
    @decorator2
    def todo(self):
        pass


```

command used: cinj{./example.py --function=print_content}
```python
# in file example.py
# in class Test
def print_content(self):
    print(self.content)

```

command used: cinj{./example.py --function=abc_func}
```python
# in file example.py
# in class MyABC(ABC)
@decorator3
@decorator4
def abc_func(self):
    pass

```

command used: cinj{./example.py --function=abc_func --decorator=false}
```python
# in file example.py
# in class MyABC(ABC)
def abc_func(self):
    pass

```

command used: cinj{./example.py --class=Test --function=todo}
```python
# in file example.py
# in class Test
@decorator1
@decorator2
def todo(self):
    pass


```

command used: cinj{./example.py --class=MyABC --function=todo}
```python
# in file example.py
# in class MyABC(ABC)
def todo(self):
    pass
```

This is content after!
