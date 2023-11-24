######################################################
## Author: Ming
## Created: 2023/11/7
## Discription
## This code is for simulating the key generation of FMHKA
######################################################
import hashlib
import random
import time
import timeit
def sha256_hash(data):
    # 创建一个SHA256哈希对象
    sha256 = hashlib.sha256()

    # 更新哈希对象中的数据
    sha256.update(data.encode('utf-8'))

    # 返回十六进制哈希值
    return sha256.hexdigest()

def key_compute(shared_key,loc,token):
    return sha256_hash(str(shared_key)+str(loc)+str(token))

def generate_random_hex(length=8):
    return ''.join([hex(random.randint(0, 15))[2:] for _ in range(length)])



def recover_key(shared_key,path_vecs):
    for path_vec in path_vecs:
        shared_key=key_compute(shared_key,path_vec['loc'],path_vec['token'])
    return shared_key



if __name__ == '__main__':
    ######### single key generation #########
    test_loc='001'
    shared_key=sha256_hash(generate_random_hex()) #randomly generate a shared_key
    test_token=generate_random_hex() #randomly generate a token
    print('当前loc：'+test_loc,'当前token：'+test_token,'当前shared_key：'+shared_key)
    elapsed_time = timeit.timeit(lambda: key_compute(shared_key, test_loc, test_token), number=10000) / 10000 #将key_compute 运行了10000次
    print(f'新生成的shared key：{key_compute(shared_key, test_loc, test_token)}，平均用时：{elapsed_time*1000:.7f}ms')

    ######### key recovering ##########
    for r in range(13):
        path_length=r*5
        path_vecs=[]
        for i in range(path_length):
            path_vec={}
            path_vec['loc']='00'+str(i)
            path_vec['token']=generate_random_hex()
            path_vecs.append(path_vec)
    # print('当前path_vecs：',path_vecs,'\n当前shared_key：'+shared_key)
    
        elapsed_time = timeit.timeit(lambda: recover_key(shared_key,path_vecs), number=10000) / 10000 #将key_compute 运行了10000次
        print(f'r={r*5}, 新生成的shared key：{recover_key(shared_key,path_vecs)}，平均用时：{elapsed_time*1000:.7f} ms \n')
    recovered_key=recover_key(shared_key,path_vecs)