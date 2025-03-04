package dal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/XiaoYao-0/Contracts-source-code-collector/conf"
	"github.com/XiaoYao-0/Contracts-source-code-collector/domain"
	"io/ioutil"
	"net/http"
	"strings"
)

type ApiJson struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		Sourcecode           string `json:"SourceCode"`
		Abi                  string `json:"ABI"`
		ContractName         string `json:"ContractName"`
		CompilerVersion      string `json:"CompilerVersion"`
		OptimizationUsed     string `json:"OptimizationUsed"`
		Runs                 string `json:"Runs"`
		ConstructorArguments string `json:"ConstructorArguments"`
		EvmVersion           string `json:"EVMVersion"`
		Library              string `json:"Library"`
		LicenseType          string `json:"LicenseType"`
		Proxy                string `json:"Proxy"`
		Implementation       string `json:"Implementation"`
		SwarmSource          string `json:"SwarmSource"`
	} `json:"result"`
}

type ContractRepo struct{}

// Collect sourceCode json of every contract address
func (contractRepo ContractRepo) Collect(contract *domain.Contract) error {
	apiKey, err := conf.ApiKey()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s",
		contract.Address, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	apiJson := ApiJson{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	err = json.Unmarshal([]byte(body), &apiJson)
	if err != nil {
		return err
	}
	if len(apiJson.Result) == 0 {
		err = errors.New("collect result empty error")
		return err
	}
	if len(apiJson.Result) > 1 {
		fmt.Printf("Warning: Result of {%s} has 2 or more srcs, address = \"%s\"\n", contract.Name, contract.Address)
	}
	contract.Json = apiJson.Result[0].Sourcecode
	err = unmarshal(contract)
	if err != nil {
		return err
	}
	return nil
}

// Unmarshal the SourceCode json string of a contract
func unmarshal(contract *domain.Contract) error {
	var src interface{}
	if strings.HasPrefix(contract.Json, "{{") {
		contract.Json = contract.Json[1 : len(contract.Json)-1]
		err := json.Unmarshal([]byte(contract.Json), &src)
		if err != nil {
			return err
		}
		src0, ok := src.(map[string]interface{})
		if ok {
			for k, v := range src0 {
				if k == "sources" {
					src1, ok1 := v.(map[string]interface{})
					if ok1 {
						for k1, v1 := range src1 {
							src2, ok2 := v1.(map[string]interface{})
							if ok2 {
								for _, v2 := range src2 {
									v3, ok3 := v2.(string)
									if ok3 {
										sol := &domain.Sol{Name: k1, Src: v3}
										contract.Sols = append(contract.Sols, sol)
									} else {
										return errors.New("unmarshal error")
									}
								}
							} else {
								return errors.New("unmarshal error")
							}
						}
						return nil
					} else {
						return errors.New("unmarshal error")
					}
				}
			}
		} else {
			return errors.New("unmarshal error")
		}
	}
	err := json.Unmarshal([]byte(contract.Json), &src)
	if err != nil {
		solName := fmt.Sprintf("%s.sol", contract.Name)
		sol := &domain.Sol{Name: solName, Src: contract.Json}
		contract.Sols = append(contract.Sols, sol)
		return nil
	}
	src0, ok := src.(map[string]interface{})
	if !ok {
		return errors.New("unmarshal error")
	}
	for k, v := range src0 {
		v0, ok0 := v.(map[string]interface{})
		if !ok0 {
			return errors.New("unmarshal error")
		}
		if ok0 {
			for _, v1 := range v0 {
				if v2, ok := v1.(string); ok {
					sol := &domain.Sol{Name: k, Src: v2}
					contract.Sols = append(contract.Sols, sol)
				}
			}
		}
	}
	return nil
}
