import matplotlib.pyplot as plt
import networkx as nx
import numpy as np
from random import randint
from PIL import ImageColor

# ...（您的其他代码）

# 生成随机颜色的函数
def get_random_color():
    return "#%06x" % randint(0, 0xFFFFFF)

# 生成浅色版本的函数
def get_lighter_color(color):
    # 将hex颜色转换为RGB
    rgb = ImageColor.getcolor(color, "RGB")
    # 创建浅色版本
    lighter_rgb = tuple(min(255, int(c + (255 - c) * 0.7)) for c in rgb)
    # 将RGB转换回hex
    return '#%02x%02x%02x' % lighter_rgb

print(get_random_color())