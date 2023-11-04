'''
:Authors:         Shashank Agrawal
:Date:            5/2016
'''

from charm.toolbox.pairinggroup import PairingGroup, GT
import sys
import time
import timeit
sys.path.append("/data/System Performance/ABE2/")
from ABE.ac17 import AC17CPABE


def main():
    # instantiate a bilinear pairing map
    pairing_group = PairingGroup('MNT224')

    # AC17 CP-ABE under DLIN (2-linear)
    cpabe = AC17CPABE(pairing_group, 2)

    # run the set up
    (pk, msk) = cpabe.setup()

    # generate a key
    time1=time.time()
    attr_list = ["ATTR"+str(i) for i in range(60)]
    
    key = cpabe.keygen(pk, msk, attr_list)
    time2=time.time()
    
    elapsed_time = timeit.timeit(lambda: cpabe.keygen(pk, msk, attr_list), number=100) / 100 
    print(f'用时{(elapsed_time)*1000}')
    # choose a random message
    msg = pairing_group.random(GT)

    # generate a ciphertext
    policy_str = 'ATTR0'
    #policy_str='ATTR1'
    ctxt = cpabe.encrypt(pk, msg, policy_str)

    # decryption
    rec_msg = cpabe.decrypt(pk, ctxt, key)
    if debug:
        if rec_msg == msg:
            print ("Successful decryption.")
        else:
            print ("Decryption failed.")


if __name__ == "__main__":
    debug = True
    main()
