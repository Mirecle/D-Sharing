import matplotlib.pyplot as plt
import numpy as np

# 创建示例数据
aps = ['300', '600', '900', '1200', '1500', '1800', '2100', '2400', '2700']
tps_values = [299.7 ,599.1 ,824.3 ,1189.2 ,1497.0 ,1534.1 ,1632.0 ,1668.1 ,1671.5 ]

# 最大时延数据
max_latency_values = [0.21,0.20,0.24,0.27,0.28,0.29,0.28,0.27,0.22]

# 平均时延数据
avg_latency_values = [0.04,0.04,0.04,0.04,0.04,0.05,0.08,0.05,0.05]

# 设置图形大小
fig, ax1 = plt.subplots(figsize=(18, 6))

# 画折线图 - 左边纵坐标
ax1.plot(aps, tps_values, color='red', marker='o', label='ShareRevocation')

# 设置左边纵坐标轴标签和标题
ax1.set_xlabel('Transaction Arrival Rate (aps)')
ax1.set_ylabel('Throughput (tps)', color='black')

# 创建第二个y轴，共享x轴
ax2 = ax1.twinx()

# 计算柱状图的位置
bar_width = 0.35
bar_positions_avg = np.arange(len(aps))
bar_positions_max = bar_positions_avg + bar_width

# 画柱状图 - 右边纵坐标

# 最大时延 - 浅蓝色
ax2.bar(bar_positions_avg, max_latency_values, width=bar_width, color='#7998B0', label='ShareRevocation maxLatency', alpha=0.7)
# 平均时延 - 深蓝色
ax2.bar(bar_positions_max, avg_latency_values, width=bar_width, color='#53618F', label='ShareRevocation avgLatency')

# 设置右边纵坐标轴标签
ax2.set_ylabel('Latency (s)', color='black')

# 将柱状图放置在顶层
ax1.set_zorder(1)
ax1.patch.set_visible(False)

# 设置x轴刻度位置和标签
ax1.set_xticks(bar_positions_avg + bar_width / 2)
ax1.set_xticklabels(aps)

# 显示网格线
ax1.grid(True, linestyle='--', alpha=0.5)

# 设置图例，使用bbox_to_anchor参数避免重叠
ax1.legend(loc='upper left', bbox_to_anchor=(0, 1))
ax2.legend(loc='upper left', bbox_to_anchor=(0, 0.95))

# 设置右边纵坐标刻度线为0.2的间隔
ax2.set_yticks(np.arange(0, 0.41, 0.1))

# 保存图为PDF文件
plt.savefig('ShareRevocation.pdf', format='pdf', bbox_inches='tight')


# 显示图形
plt.show()