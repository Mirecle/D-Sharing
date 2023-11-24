#########################################
##  This code is used for generating the number of dissemination trees at different depths
##  Created in 2023-11-18
##  Author Ming
##
#########################################



import os
import sys
import pickle
import random
import gc
import logging
import datetime
import glob

gc.disable()
#程序处理部分
class Tree:
    def __init__(self,user):
        self.root=user
        self.Maxlength=1
        self.MaxlenNode={user}  # 修改为集合
        self.nodeNumber=1
        self.totalUser={(user,1)}  # 修改为集合
        
    # 现在不需要isExist方法，因为集合的查找是O(1)复杂度

    def addNode(self,user,length):
        self.nodeNumber += 1
        # 用add代替append添加到集合中
        self.totalUser.add((user,length))
        if self.Maxlength<length:
            self.Maxlength=length
            self.MaxlenNode=[user]
        elif self.Maxlength==length:
            self.MaxlenNode.append(user)
current_dir='D:\\app\\SynologyDrive\\文档\\第三篇论文\\Transparent-Kites\\System Performance\\dataSet\\'


# 拼接新的路径
def load_latest_data():
    all_files = sorted(glob.glob(os.path.join(current_dir, "process_tree_*")), 
                       key=os.path.getmtime, reverse=True)
    for file in all_files:
        try:
            with open(file, 'rb') as f:
                trees, edge, users = pickle.load(f)
                print(f'读取文件{file}成功')
                return trees, edge, users
        except EOFError:
            continue
    return None  # 如果所有文件都无法读取



class Tree:
    def __init__(self,user):
        self.root=user
        self.Maxlength=1
        self.MaxlenNode={user}  # 修改为集合
        self.nodeNumber=1
        self.totalUser={(user,1)}  # 修改为集合
        
    # 现在不需要isExist方法，因为集合的查找是O(1)复杂度

    def addNode(self,user,length):
        self.nodeNumber += 1
        # 用add代替append添加到集合中
        self.totalUser.add((user,length))
        if self.Maxlength<length:
            self.Maxlength=length
            self.MaxlenNode=[user]
        elif self.Maxlength==length:
            self.MaxlenNode.append(user)
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
edge=set()
trees=[]
zero_sub=[]
def generate_filename(base_name):
    timestamp = datetime.datetime.now().strftime("%Y%m%d%H%M%S")
    return f"{base_name}_{timestamp}"

def save_data(trees, edge, users):
    sys.setrecursionlimit(10000) 
    filename = generate_filename("process_tree")
    with open(os.path.join(current_dir, filename), 'wb') as f:
        pickle.dump((trees, edge, users), f)

    # 保留最近10个文件
    all_files = sorted(glob.glob(os.path.join(current_dir, "process_tree_*")), 
                       key=os.path.getmtime, reverse=True)
    for old_file in all_files[10:]:
        os.remove(old_file)
# def data_to_file(trees,edge,users):
#     with open(current_dir+'\\process_tree', 'wb') as f:

#         sys.setrecursionlimit(1000000) 
#         pickle.dump(trees,edge,users, f)

# 从文件中加载session对象
def load_tree_from_file():
    with open(current_dir+'\\process_tree', 'rb') as f:
        trees,edge,users = pickle.load(f)
    return trees

def add_email(from_id, to_id):
    # 检查发件人和收件人是否已经存在于字典中
    if from_id not in users:
        users[from_id] = User(from_id)
    if to_id not in users:
        users[to_id] = User(to_id)

    # 将收件人添加到发件人的列表中
    users[from_id].add_sender(users[to_id])
    edge.add((users[from_id],users[to_id]))
# 读取数据集文件并处理每一行


def rebuild(tree):
    for i,length in tree.totalUser:
        tree.totalUser.discard((i,length))
        tree.totalUser.add((users[i.id],length))
    tree.totalUser.discard((tree.root,1))
    tree.totalUser.add((None,tree.root,len(users[tree.root.id].received_emails_from),1))
    current_nodes=[(None,tree.root,1)]
    # for length in range(tree.Maxlength):
    while len(current_nodes)!=0:
        pre_node,node,length=current_nodes[0]
        subNum=0
        for i in users[node.id].received_emails_from:
            if (i,length+1) in tree.totalUser:
                tree.totalUser.discard((i,length+1))
                tree.totalUser.add((node,i,length+1))
                current_nodes.append((node,i,length+1))
                subNum+=1
        tree.totalUser.discard((pre_node,node,length))
        tree.totalUser.add((pre_node,node,subNum,length))
        current_nodes.remove((pre_node,node,length))
trees,edge,_=load_latest_data() 



print('all operation is done, analyze start')
    
numberlist={}
Depth=[]
for i in trees:
    if i.Maxlength not in numberlist:
        numberlist[i.Maxlength]=1
        Depth.append(i.Maxlength)
    else:
        numberlist[i.Maxlength]+=1

tree_depths = sorted(Depth)
tree_numbers=[numberlist[i] for i in tree_depths]
print(tree_depths)

import matplotlib.pyplot as plt
from matplotlib import pyplot as plt

fig, ax = plt.subplots()



# Re-creating the horizontal bar plot with a logarithmic scale
plt.rcParams['font.family'] = 'Times New Roman'
plt.rcParams['font.size'] = 12  # Base font size; will scale up for labels

ax.spines['right'].set_visible(False)
ax.spines['top'].set_visible(False)

# Plotting with a logarithmic scale and adding a black edge color to the bars
ax.barh(tree_depths, tree_numbers, color='skyblue', edgecolor='black')

# Removing the vertical grid (dashed lines)
ax.grid(False)

# Setting the x-axis to a logarithmic scale
ax.set_xscale('log')
ax.set_yticks(tree_depths)
ax.set_ylim(1, 16)
# Scaling up the x and y axis labels by 1.5 times
ax.set_xlabel('Dissemination Tree Numbers (Log Scale)', fontsize=18)
ax.set_ylabel('Tree Depth', fontsize=18)

# Adding the actual numbers to the right of the bars with increased font size
for i in range(len(tree_depths)):
    ax.text(tree_numbers[i], tree_depths[i], f' {tree_numbers[i]}', va='center', ha='left', color='black', fontsize=14)

plt.show()
while True:
    pass