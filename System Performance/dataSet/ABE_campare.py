######################################################
## Author: Ming
## Created: 2023/11/7
## Discription
## This code is used for drawing the comparsion figure
######################################################
import matplotlib.pyplot as plt
import numpy as np

# Data
involved_users = np.array([1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60])
ABE = np.array([28.607, 48.008, 74.552, 102.167, 148.194, 155.108, 181.979, 217.689, 241.808, 285.810, 294.059, 326.384, 348.112])
FMHKA = np.array([0.0013, 0.0078, 0.0136, 0.0198, 0.0285, 0.03540, 0.0439, 0.0528, 0.0532, 0.06131, 0.0648, 0.07505, 0.0859])

# Bar width
width = 2
ax.barh(tree_depths, tree_numbers, color='skyblue', edgecolor='black')
# Plotting
fig, ax = plt.subplots()
bar1 = ax.bar(involved_users , ABE, width, label='ABE')
bar2 = ax.bar(involved_users , FMHKA, width, label='FMHKA')

# Labels and titles
ax.set_xlabel('Involved Users')
ax.set_ylabel('Values')
ax.set_title('Bar Chart of ABE and FMHKA')
ax.set_xticks(involved_users)
ax.legend()

plt.show()
aaa=1