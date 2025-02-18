package dal

import (
	"bufio"
	"errors"
	"github.com/XiaoYao-0/Contracts-source-code-collector/domain"
	"io"
	"os"
	"strings"
)

// Read contracts-list
func GetContracts(filepath string) ([]*domain.Contract, error) {
	var contracts []*domain.Contract
	f, err := os.Open(filepath)
	if err != nil {
		return contracts, err
	}
	defer f.Close()
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return contracts, err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		contracts = append(contracts, phaseContract(s))
	}
	return contracts, nil
}

// Phase contract information,
// valid example: ${Txhash},${ContractAddress},${ContractName}
// 0x0eb83dfc4bda0ee8bebb4827793adb24799419c23dd2a7cba173a5f719a7dd0f,0x1f9840a85d5af5bf1d1762f925bdaddc4201f984,UNI
func phaseContract(s string) *domain.Contract {
	var contract domain.Contract
	s = strings.ReplaceAll(s, "\"", "")
	if !strings.HasPrefix(s, "0x") {
		contract.Error = errors.New("non-contract")
		return &contract
	}
	split := strings.Split(s, ",")
	if len(split) != 3 {
		contract.Error = errors.New("non-contract")
		return &contract
	}
	if !isValid(split[1]) {
		contract.Error = errors.New("illegal address")
		return &contract
	}
	contract.Name = split[2]
	contract.Address = split[1]
	contract.ContractRepo = ContractRepo{}
	return &contract
}

// Check the validity of the address,
// valid example: 0x0eb83dfc4bda0ee8bebb4827793adb24799419c23dd2a7cba173a5f719a7dd0f
func isValid(address string) bool {
	if len(address) != 42 {
		return false
	}
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	return true
}
