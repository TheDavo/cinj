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
