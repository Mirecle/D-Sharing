######################################################
## Author: Ming
## Created: 2023/11/12
## Discription
## This code is used for building the dissemination tree from the tested dataset
######################################################
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

# 获取当前文件的绝对路径
current_file_path = os.path.abspath(__file__)

# 获取当前文件所在的目录
current_dir = os.path.dirname(current_file_path)
logging.basicConfig(filename=current_dir+'\\example.log', level=logging.INFO, 
                    format='%(asctime)s - %(levelname)s - %(message)s')
# 拼接新的路径
def load_latest_data():
    all_files = sorted(glob.glob(os.path.join(current_dir, "tested_tree_*")), 
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
        self.totalUser={(None,user,0,1)}  # 修改为集合
        
    # 现在不需要isExist方法，因为集合的查找是O(1)复杂度

    def addNode(self,pre_node,user,length):
        self.nodeNumber += 1
        # 用add代替append添加到集合中
        self.totalUser.add((pre_node,user,0,length))
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
    sys.setrecursionlimit(20000) 
    filename = generate_filename("tested_tree_")
    with open(os.path.join(current_dir, filename), 'wb') as f:
        pickle.dump((trees, edge, users), f)

    # 保留最近10个文件
    all_files = sorted(glob.glob(os.path.join(current_dir, "tested_tree_*")), 
                       key=os.path.getmtime, reverse=True)
    for old_file in all_files[10:]:
        os.remove(old_file)
# def data_to_file(trees,edge,users):
#     with open(current_dir+'\\process_tree', 'wb') as f:

#         sys.setrecursionlimit(1000000) 
#         pickle.dump(trees,edge,users, f)

# 从文件中加载session对象
def load_tree_from_file():
    with open(current_dir+'\\tested_tree', 'rb') as f:
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


try:

    
    trees,edge,users=load_latest_data()

except Exception as e:
    with open(current_dir+'\\part_dataset\\result-4000.txt', 'r') as file:
        for line in file:
            # 跳过注释行
            if line.startswith('#'):
                continue

            parts = line.strip().split(' ')#'\t' for Email-EuAll.txt
            if len(parts) != 2:
                continue  # 忽略格式不正确的行

            from_node, to_node = parts
            add_email(int(from_node), int(to_node))

indexing=0
while len(edge)!=0:
    From, to = edge.pop()
    if From.id==to.id:
        continue
    else:
        edge.add((From,to))

    node=From
    prenode=None
    tree=Tree(node)
    used_user=set()
    
    waiting_user=[]
    waiting_user.append((None,node,1))
    length=1
    first=True
    if len(to.received_emails_from)==0: ##没有后续传播树，只是一次点对点共享
        edge.discard((node,to))
        tree.addNode(None,to,length+1)
        # used_user.add(to)
        print(f'传播树{len(trees)+1}颗构建完成！')

        trees.append(tree)
        if indexing%10==0:
            save_data(trees,edge,users)
            indexing=0
        gc.collect()
        indexing+=1
        continue
    while len(waiting_user)!=0: #选择一个待处理用户
        
        # for i in node.received_emails_from:
        #     if (node,i) not in edge:
        #         continue
        #     else: vaild_number+=1
        # if  vaild_number>0:
        valid_list=[]
        for i in node.received_emails_from:
            if (node,i) in edge and i not in used_user:
                valid_list.append(i)
        if len(valid_list)>0: 
            numbers=len(valid_list) if len(valid_list)<20 else 20 #随机选择data_num数量的节点作为该数据的传播节点
            data_num=random.randint(1,numbers) 
            
            tree.totalUser.discard((prenode,node,0,length))
            tree.totalUser.add((prenode,node,data_num,length))

        else:  #如果没有子节点，直接结束循环
            if len(waiting_user)!=0:
                waiting_user.pop(0)
            if len(waiting_user)!=0:
                    prenode,node,length=waiting_user[0]
                    continue
            else:
                break
        select_items=random.sample(valid_list, data_num)
        if first:
            select_items[0]=to #edge.pop 扔掉了(From,to)节点.这个节点将在第一次执行时被处理掉
            first=False
        for i in select_items: 
            try:
                edge.discard((node,i))
                tree.addNode(node,i,length+1)
                used_user.add(i)
                waiting_user.append((node,i,length+1))                    
                node.received_emails_from.remove(i)
                #logging.info(f'user{node.id} remove {i.id},current sub{node.received_emails_from}')
                if len(node.received_emails_from)==0:
                    zero_sub.append(node)
            except Exception as e:
                continue

        waiting_user.remove((prenode,node,length))
        if len(waiting_user)!=0:
            prenode,node,length=waiting_user[0]
        else:
            break
        print(f'传播树{len(trees)+1}，处理节点{node.id}完成，有效总子节点数{len(valid_list)}，选择{data_num}个，当前深度{tree.Maxlength}，当前剩余edge{len(edge)}')
    print(f'传播树{len(trees)+1}颗构建完成！')

    trees.append(tree)
    if indexing%10==0:
        save_data(trees,edge,users)
        indexing=0
    gc.collect()
    indexing+=1
# 举例输出用户0的信息

if __name__ == "__main__":
    print('all operation is done, analyze start')
    