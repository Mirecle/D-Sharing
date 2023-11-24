######################################################
## Author: Ming
## Created: 2023/11/20
## Discription
## This code is used for drawing the distribution of Propagation Trees by Entropy (Log Scale) in case 3 and different entropy ranged by node number in case 2
######################################################
import numpy as np
import matplotlib.pyplot as plt
def show1(): #figure tree
    H_list={'0.00': 136, '1.00': 28, '1.58': 30, '2.00': 20, '2.32': 20, '2.58': 30, '2.81': 42, '3.00': 16, '3.17': 27, '3.32': 10, '3.46': 11, '3.58': 24, '3.70': 13, '3.81': 14, '3.91': 30, '4.00': 32, '4.17': 54, '4.46': 44, '4.52': 23, '4.58': 24, '4.91': 30, '5.13': 35, '5.17': 72, '5.36': 41, '5.39': 42, '5.43': 43, '5.46': 44, '5.58': 48, '5.86': 58, '5.93': 61, '5.98': 63, '6.07': 67, '6.11': 69, '6.39': 84, '6.61': 98, '6.73': 106, '6.82': 113, '7.12': 139, '7.23': 150, '7.64': 200, '7.79': 221, '7.89': 237, '8.40': 337, '8.42': 342, '8.91': 481, '9.64': 796, '10.61': 1561, '11.76': 3471, '13.59': 12316, '17.90': 244709}
    entropy = list(map(float, H_list.keys()))
    tree_counts = list(H_list.values())
    # 将信息熵的值转换为字符串，以便在x轴上等间距显示
    entropy_str = list(H_list.keys())

    # 计算传播树数量的对数
    tree_counts_log = np.log10(tree_counts)

    # 绘制柱状图
    plt.figure(figsize=(15, 7))
    plt.grid(color='lightgrey',linestyle='--',zorder=0) #启用网格
    bars = plt.bar(entropy_str, tree_counts_log, color='blue', linewidth=0.5,zorder=3)

    # 在每个柱子顶部添加数量标签
    for bar in bars:
        yval = bar.get_height()
        plt.text(bar.get_x() + bar.get_width()/2, yval, round(10**yval), ha='center', va='bottom', fontname='Times New Roman',fontsize=12)

    # 设置标题和坐标轴标签
    # plt.title('Distribution of Propagation Trees by Entropy (Log Scale)', fontname='Times New Roman', fontsize=14)
    plt.xlabel('Information Entropy', fontname='Times New Roman', fontsize=18)
    plt.ylabel('Number of Dissemination Trees (Log Scale)', fontname='Times New Roman', fontsize=18)

    # 设置坐标轴字体
    plt.xticks(rotation=90, fontname='Times New Roman', fontsize=14)
    plt.yticks(fontname='Times New Roman', fontsize=14)

    # 显示图表
    plt.tight_layout()
    plt.show()

def show2(): #figure Node
    entropy_list = [18.8365, 17.4811, 13.955, 12.3685, 11.4538, 10.7591, 10.2192, 9.8734, 9.66, 9.5137]
    num_list = [468120, 182952, 15881, 5288, 2805, 1733, 1192, 938, 809, 731]
    bar_width = 0.4
    # Number of children (from 0 to len(entropy_list) - 1)
    num_of_children = list(range(len(entropy_list)))
    node_count_bar_positions = [x + bar_width/2 for x in num_of_children]
    # Creating the figure and the first set of bars for entropy
    fig, ax1 = plt.subplots(figsize=(10, 6))
    ax1.bar(num_of_children, entropy_list, color='green', width=0.4, align='center',zorder=3)

    # Labelling the entropy values on top of the bars
    for i, v in enumerate(entropy_list):
        ax1.text(i, v + 0.3, str(v), ha='center', color='green',zorder=3)
    ax1.grid(color='lightgrey', linestyle='--', zorder=0)
    # Setting the labels and titles
    ax1.set_xlabel('Number of Child Nodes', fontname='Times New Roman', fontsize=18)
    ax1.set_ylabel('Information Entropy', fontname='Times New Roman', fontsize=18)

    # Setting the second y-axis for the number of nodes
    ax2 = ax1.twinx()
    ax2.bar([x + 0.4 for x in num_of_children], num_list, color='blue', width=0.4, align='center',zorder=3,log=True)

    # Labelling the node count values on top of the bars
    for i, v in enumerate(num_list):
        ax2.text(i+0.4, v + 0.3, str(v), ha='center', color='blue',zorder=3)

    ax2.set_ylabel('Number of Nodes', fontname='Times New Roman', fontsize=18)

    # Setting the grid
    
    plt.show()
    # lists = [18.8365, 17.4811, 13.955, 12.3685, 11.4538, 10.7591, 10.2192, 9.8734, 9.66, 9.5137]
    # num_list=[468120, 182952, 15881, 5288, 2805, 1733, 1192, 938, 809, 731]
    # # 子节点数量（从0到len(lists) - 1）
    # num_of_children = list(range(len(lists)))

    # # 绘制柱状图
    # plt.figure(figsize=(10, 6))
    # plt.grid(color='lightgrey',linestyle='--',zorder=0) #启用网格
    # bars=plt.bar(num_of_children, lists, color='green',zorder=3)
    # for bar in bars:
    #     yval = bar.get_height()
    #     plt.text(bar.get_x() + bar.get_width()/2, yval, yval, ha='center', va='bottom', fontname='Times New Roman',fontsize=14)
    # # 设置标题和坐标轴标签

    # ax2 = ax1.twinx()
    # ax2.bar([x + 0.4 for x in num_of_children], num_list, color='blue', width=0.4, align='center')
    # # plt.title('Entropy for Different Numbers of Child Nodes', fontsize=14)
    # plt.xlabel('Number of Child Nodes',fontname='Times New Roman', fontsize=18)
    # plt.ylabel('Information Entropy', fontname='Times New Roman',fontsize=18)

    # # 设置坐标轴字体大小
    # plt.xticks(num_of_children, fontsize=14)
    # plt.yticks(fontsize=14)

    # # 显示图表
    # plt.show()

if __name__=='__main__':
    show2()
    while True:
        continue