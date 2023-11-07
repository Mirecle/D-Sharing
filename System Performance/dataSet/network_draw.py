######################################################
## Author: Ming
## Created: 2023/11/7
## Discription
##
######################################################
import os
import sys

import matplotlib.pyplot as plt
import networkx as nx
import numpy as np
import color_generation
# 获取当前文件的绝对路径
plt.figsize=(60, 60)

# 获取当前文件所在的目录
current_file_path = os.path.abspath(__file__)

# 获取当前文件所在的目录
current_dir = os.path.dirname(current_file_path)

# 拼接新的路径


class User:
    def __init__(self, id):
        self.id = id
        self.received_emails_from = []

    def add_sender(self, sender):
        self.received_emails_from.append(sender)

    def __repr__(self):
        return f"User({self.id})"


# 用于存储用户ID和用户对象的字典
users = {}

def add_email(from_id, to_id):
    # 检查发件人和收件人是否已经存在于字典中
    if from_id not in users:
        users[from_id] = User(from_id)
    if to_id not in users:
        users[to_id] = User(to_id)

    # 将收件人添加到发件人的列表中
    users[from_id].add_sender(users[to_id])

# 读取数据集文件并处理每一行

with open(current_dir+'\\Email-EuAll.txt', 'r') as file:
    for line in file:
        # 跳过注释行
        if line.startswith('#'):
            continue

        parts = line.strip().split('\t')
        if len(parts) != 2:
            continue  # 忽略格式不正确的行

        from_node, to_node = parts
        add_email(int(from_node), int(to_node))

plt.ion()  # 开启交互模式
# 假设 'users' 是一个包含了所有User对象的字典，键是用户的id
# users = {0: User(0), 1: User(1), ... }
class Node:
    def __init__(self,id,level):
        self.id = id
        self.level = level
        self.position=None
        self.color=None
    def a_position(self,position):
        self.position=position
    def set_color(self,color):
        self.color=color
G = nx.DiGraph()
def set_position(pre_node,node, parent_position, angle_span,index,total_valid_subnode_number):
        
        angle = (2 * np.pi / total_valid_subnode_number) * index if total_valid_subnode_number > 1 else 0
        # if node.level==3:
        #     random_radii=radii[pre_node.level]*np.random.rand()
        # elif node.level==2:
        #     sub_number=len(users[node.id].received_emails_from) if len(users[node.id].received_emails_from)>30 else np.random.rand()*80+len(users[node.id].received_emails_from)
        #     random_radii=radii[pre_node.level]*sub_number/400
        #     print(sub_number)
        random_radii=radii[pre_node.level]*np.random.rand()
        pos_x = pre_node.position[0]+random_radii * np.cos(angle)
        pos_y = pre_node.position[1]+random_radii * np.sin(angle)

        idx = index

        
        positions[node.id] = (pos_x, pos_y)
        node.a_position((pos_x, pos_y))
# 根据层级设置颜色
colors = {1: 'red', 2: 'blue', 3:'white'}
colors2 = {1: 'black', 2:'gray'}
radii = {1: 1.0, 2: 0.5, 3: 0.5}
node_colors = []
node_sizes=[]
edge_colors=[]
# 用于记录每个节点层级的字典
waiting_list=[]
waiting_list_id=[]
positions={}
drawed_id=[]
node_list=[]
edge_list=[]
def draw_node(node,pre_node):
    node_list.append(node.id)
    if node.level==2:
        choose_color=color_generation.get_random_color()
        node_colors.append(choose_color)
        node.set_color(choose_color)
        node_sizes.append(70)
    elif node.level==3:
        choose_color=color_generation.get_lighter_color(pre_node.color)
        node_colors.append(choose_color)
        node.set_color(choose_color) 
        node_sizes.append(50)   
    else:
        node_colors.append('black')# the origin node
        node.set_color('black')
        node_sizes.append(100)
edge_list=[]
origin_node=Node(0,1)
origin_node.a_position((0,0))
waiting_list.append(origin_node)
waiting_list_id.append(0)
draw_node(origin_node,Node)



parent_position=0
angle_span=360
positions[0] = (0, 0)
while len(waiting_list)!=0:
    node=waiting_list[0]
    user = users[node.id]
    index=1
    total_valid_subnode_number=0
    for recevier in user.received_emails_from:
        if recevier.id not in drawed_id and recevier.id not in waiting_list_id and node.level<3:
            total_valid_subnode_number=total_valid_subnode_number+1  #total_valid_subnode_number is used to recode users's (line 134) valid subnode number
    for recevier in user.received_emails_from:
        if recevier.id not in drawed_id and recevier.id not in waiting_list_id and node.level<3:
            index=index+1
            next_node=Node(recevier.id,node.level+1)
            set_position(node,next_node, parent_position, angle_span,index,total_valid_subnode_number)
            waiting_list.append(next_node)
            waiting_list_id.append(recevier.id)
            draw_node(next_node,node)
            edge_list.append({'from':node.id,'to':recevier.id,'color':colors2[node.level]})
            # G.add_edge(node.id,recevier.id)
            edge_colors.append(colors2[node.level])
    index=0
    waiting_list.pop(0)
    waiting_list_id.pop(0)
    drawed_id.append(node.id)


node_list.reverse()
node_colors.reverse()
edge_list.reverse()
edge_colors.reverse()
node_sizes.reverse()
for i in node_list:
    G.add_node(i)

for j in edge_list:
    G.add_edge(j['from'],j['to'])

# 画出网络图

# nx.draw(G, pos, node_color=node_colors, node_size=50, edge_color='gray', with_labels=False)
nx.draw_networkx_edges(G, positions, edge_color=edge_colors)
nx.draw_networkx_nodes(G, positions, node_color=node_colors, node_size=node_sizes, edgecolors='black', linewidths=0.5)
plt.savefig('network_diagram.pdf', format='pdf', dpi=1200)