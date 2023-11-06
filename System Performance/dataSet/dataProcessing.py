import os
import sys
# 获取当前文件的绝对路径
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

# 举例输出用户0的信息
print(users[0])
# 打印用户0收到邮件的用户列表
print(users[0].received_emails_from)