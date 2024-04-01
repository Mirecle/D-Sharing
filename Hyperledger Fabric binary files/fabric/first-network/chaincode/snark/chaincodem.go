package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	//"github.com/hyperledger/fabric-protos-go/peer"
)

// ShareInitiationChaincode 定义链码结构
type ShareInitiationChaincode struct {
}

// DataNode 表示数据节点的结构
type DataNode struct {
	Owner    string `json:"owner"`
	DataHash string `json:"dataHash"`
	RootKey  string `json:"rootKey"`  //加密数据的密钥
	ShareKey string `json:"shareKey"` //共享根密钥
	Hkey     string `json:"Hkey"`     //共享根密钥的哈希
	Token    string `json:"token"`
	Address  string `json:"address"`
	Status   string `json:"status"`
	Loc      string `json:"loc"`
}

// Init 链码初始化方法
func (s *ShareInitiationChaincode) Init(APIstub shim.ChaincodeStubInterface) peer.Response {

	fmt.Println("chaincode Init")
	var A, B, C string          // Entities
	var Aval, Bval, Cval string // Asset holdings

	A = "mprooft"
	Aval = "true"
	B = "mprooff"
	Bval = "false"
	C = "mproof"
	Cval = "true-false"
	err := APIstub.PutState(A, []byte(Aval))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(B, []byte(Bval))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(C, []byte(Cval))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// Invoke 链码调用方法
func (s *ShareInitiationChaincode) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "ShareInitiation" {
		return s.ShareInitiation(APIstub, args)
	} else if function == "ShareDissemination" {
		return s.ShareDissemination(APIstub, args)
	} else if function == "ShareUpdate" {
		return s.ShareUpdate(APIstub, args)
	} else if function == "ShareRevocation" {
		return s.ShareRevocation(APIstub, args)
	} else if function == "query" {
		return s.query(APIstub, args)
	} else if function == "querydata" {
		return s.querydata(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func mimcHash(data []byte, str string) string {
	length := 76
	rand.Seed(time.Now().UnixNano())
	digits := []rune("0123456789")
	result := make([]rune, length)
	for i := 0; i < length; i++ {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}

// ShareInitiation 实现Share Initiation的逻辑
func (s *ShareInitiationChaincode) ShareInitiation(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	owner := args[0]
	dataHash := args[1]
	rootKey := args[2]
	shareKey := args[3]
	token := args[4]
	address := args[5]
	loc := args[6]

	// 检查数据是否已经共享
	exists, err := s.checkExistence(APIstub, dataHash)
	if err != nil {
		return shim.Error(err.Error())
	}
	if exists {
		return shim.Error(fmt.Sprintf("Data with hash %s has already been shared", dataHash))
	}

	// 创建数据节点
	dataNode := DataNode{
		Owner:    owner,
		DataHash: dataHash,
		RootKey:  rootKey,
		ShareKey: shareKey,
		Token:    token,
		Address:  address,
		Status:   "Active",
		Loc:      loc,
	}
	preImage := []byte{0x01, 0x02, 0x03}
	dataNode.Hkey = mimcHash(preImage, dataNode.ShareKey)

	// 将数据节点保存到世界状态
	dataNodeJSON, err := json.Marshal(dataNode)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(dataHash, dataNodeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	DId := "Initiation#" + dataHash + "Owner#" + owner
	err = APIstub.PutState(DId, dataNodeJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// checkExistence 检查数据是否已经存在
func (s *ShareInitiationChaincode) checkExistence(APIstub shim.ChaincodeStubInterface, dataHash string) (bool, error) {
	//dataNodeJSON, err := APIstub.GetState(dataHash)
	_, err := APIstub.GetState(dataHash)
	if err != nil {
		return false, err
	}
	//return dataNodeJSON != nil, nil
	return false, nil
}

// ShareDissemination 实现Share Dissemination的逻辑
func (s *ShareInitiationChaincode) ShareDissemination(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	owner := args[0]
	dataHash := args[1]
	rootKey := args[2]
	token := args[3]
	address := args[4]
	loc := args[5]
	sowner := args[6]

	// 验证 SU1 的身份和权限,零知识证明和节点是否active
	isproof := s.Mproof(APIstub, args)
	fmt.Println("Mproof result", isproof)
	isactive := s.checkNodeStatus(APIstub, dataHash, sowner)
	fmt.Println("Active result", isactive)
	if isproof && isactive {
		// 创建新的传播节点
		dataNode := DataNode{
			Owner:    owner,
			DataHash: dataHash,
			RootKey:  rootKey,
			Token:    token,
			Address:  address,
			Status:   "Active",
			Loc:      loc,
		}

		// 计算新的共享密钥
		preImage := []byte{0x01, 0x02, 0x03}
		dataNode.ShareKey = mimcHash(preImage, dataNode.ShareKey+dataNode.Loc+dataNode.Token)
		preImage = []byte{0x01, 0x02, 0x03}
		dataNode.Hkey = mimcHash(preImage, dataNode.ShareKey)
		// 将数据节点保存到世界状态
		dataNodeJSON, err := json.Marshal(dataNode)
		if err != nil {
			return shim.Error(err.Error())
		}
		DId := "Initiation#" + dataHash + "Owner#" + owner
		err = APIstub.PutState(DId, dataNodeJSON)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	} else if isproof == false {
		return shim.Error("ShareDissemination failed --isproof false")
	} else if isactive == false {
		return shim.Error("ShareDissemination failed --isactive false")
	} else {
		return shim.Error("ShareDissemination failed")
	}
}

// checkNodeStatus 检查节点的状态
func (s *ShareInitiationChaincode) checkNodeStatus(APIstub shim.ChaincodeStubInterface, dataHash string, sowner string) bool {
	combinationKey := "Initiation#" + dataHash + "Owner#" + sowner
	result, _ := APIstub.GetState(combinationKey)
	var datanode DataNode
	if result != nil {
		json.Unmarshal(result, &datanode)
	} else {
		return false
	}
	if datanode.Status == "Active" {
		return true
	} else {
		return false
	}
}

// ShareUpdate 实现Share Update的逻辑
func (s *ShareInitiationChaincode) ShareUpdate(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	sowner := args[0]
	dataHash := args[1]

	// 检查节点的状态和零知识证明
	isproof := s.Mproof(APIstub, args)
	isactive := s.checkNodeStatus(APIstub, dataHash, sowner)
	if isproof && isactive {
		begin := "Initiation#" + dataHash + "Owner#"
		end := string(BytesPrefix([]byte(begin)))
		result, err := APIstub.GetStateByRange(begin, end)
		if err != nil {
			return shim.Error("Not spread yet")
		}
		if result != nil {
			defer func() {
				result.Close()
			}()
			for result.HasNext() {
				arecord, err := result.Next()
				if err != nil {
					return shim.Error("ShareUpdate failed")
				}
				var datanode DataNode
				json.Unmarshal(arecord.Value, &datanode)
				datanode.Token = newtoken(8)
				// 将数据节点保存到世界状态
				dataNodeJSON, err := json.Marshal(datanode)
				if err != nil {
					return shim.Error(err.Error())
				}
				DId := "Initiation#" + datanode.DataHash + "Owner#" + datanode.Owner
				err = APIstub.PutState(DId, dataNodeJSON)
				if err != nil {
					return shim.Error(err.Error())
				}
			}
		}
		return shim.Success([]byte("ShareUpdate finished"))
	} else if isproof == false {
		return shim.Error("ShareUpdate failed --isproof false")
	} else if isactive == false {
		return shim.Error("ShareUpdate failed --isactive false")
	} else {
		return shim.Error("ShareUpdate failed")
	}
}

// ShareRevocation 实现共享数据吊销功能
func (s *ShareInitiationChaincode) ShareRevocation(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}
	dataHash := args[0]
	sowner := args[1]
	// 检查节点的存在和状态
	isproof := s.Mproof(APIstub, args)
	isactive := s.checkNodeStatus(APIstub, dataHash, sowner)
	if isproof && isactive {
		begin := "Initiation#" + dataHash + "Owner#"
		end := string(BytesPrefix([]byte(begin)))
		result, err := APIstub.GetStateByRange(begin, end)
		if err != nil {
			return shim.Error("Not spread yet")
		}
		if result != nil {
			defer func() {
				result.Close()
			}()
			for result.HasNext() {
				arecord, err := result.Next()
				if err != nil {
					return shim.Error("ShareRevocation failed")
				}
				var datanode DataNode
				json.Unmarshal(arecord.Value, &datanode)
				datanode.Address = "initialization"
				// datanode.Status = "Inactive"
				// 将数据节点保存到世界状态
				dataNodeJSON, err := json.Marshal(datanode)
				if err != nil {
					return shim.Error(err.Error())
				}
				DId := "Initiation#" + datanode.DataHash + "Owner#" + datanode.Owner
				err = APIstub.PutState(DId, dataNodeJSON)
				if err != nil {
					return shim.Error(err.Error())
				}
			}
		}
		return shim.Success([]byte("ShareRevocation finished"))
	} else if isproof == false {
		return shim.Error("ShareRevocation failed --isproof false")
	} else if isactive == false {
		return shim.Error("ShareRevocation failed --isactive false")
	} else {
		return shim.Error("ShareRevocation failed")
	}
}

// BytesPrefix 前缀批查询
func BytesPrefix(prefix []byte) []byte {
	var limit []byte
	for i := len(prefix) - 1; i >= 0; i-- {
		c := prefix[i]
		if c < 0xff {
			limit = make([]byte, i+1)
			copy(limit, prefix)
			limit[i] = c + 1
			break
		}
	}
	return limit
}

func newtoken(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := "0123456789abcdefghijklmnopqrstuvwxyz"

	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

func (s *ShareInitiationChaincode) Mproof(APIstub shim.ChaincodeStubInterface, args []string) bool {
	Avalbytes, err := APIstub.GetState("mproof")
	if err != nil {
		return false
	}
	if string(Avalbytes) == "true" {
		return true
	} else {
		return false
	}
}

func (s *ShareInitiationChaincode) query(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := APIstub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}
	fmt.Printf("Query Response:%s\n", string(Avalbytes))
	return shim.Success(Avalbytes)
}

func (s *ShareInitiationChaincode) querydata(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	dataHash := args[0]
	var datanodelist []DataNode
	begin := "Initiation#" + dataHash + "Owner#"
	end := string(BytesPrefix([]byte(begin)))
	result, err := APIstub.GetStateByRange(begin, end)
	if err != nil {
		return shim.Error("Not spread yet")
	}
	if result != nil {
		defer func() {
			result.Close()
		}()
		for result.HasNext() {
			arecord, err := result.Next()
			if err != nil {
				return shim.Error("Acquisition failure!")
			}
			var datanode DataNode
			json.Unmarshal(arecord.Value, &datanode)
			datanodelist = append(datanodelist, DataNode{
				DataHash: datanode.DataHash,
				Owner:    datanode.Owner,
				Loc:      datanode.Loc,
				Token:    datanode.Token,
				Status:   datanode.Status,
				Address:  datanode.Address,
			})
		}
	}
	datanodelistData, _ := json.Marshal(datanodelist)
	return shim.Success(datanodelistData)
}

// 主函数，启动链码
func main() {
	err := shim.Start(new(ShareInitiationChaincode))
	if err != nil {
		fmt.Printf("Error starting ShareInitiationChaincode: %s", err)
	}
}
