import math
from abc import ABC
print(math.sin(0.37*math.pi))


class Test:
    def __init__(self):
        self.content = "Nothing!"

    def print_content(self):
        print(self.content)


class MyABC(ABC):
    def __init__(self):
        pass

    def abc_func(self):
        pass
