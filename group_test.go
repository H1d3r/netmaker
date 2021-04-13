package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gravitl/netmaker/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var Networks []models.Network

func TestCreateNetwork(t *testing.T) {
	network := models.Network{}
	network.NetID = "skynet"
	network.AddressRange = "10.71.0.0/16"
	t.Run("CreateNetwork", func(t *testing.T) {
		response, err := api(t, network, http.MethodPost, "http://localhost:8081/api/networks", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("InvalidToken", func(t *testing.T) {
		response, err := api(t, network, http.MethodPost, "http://localhost:8081/api/networks", "badkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, message.Code)
		assert.Equal(t, "W1R3: You are unauthorized to access this endpoint.", message.Message)
	})
	t.Run("BadName", func(t *testing.T) {
		//issue #42
		t.Skip()
	})
	t.Run("BadAddress", func(t *testing.T) {
		//issue #42
		t.Skip()
	})
	t.Run("DuplicateNetwork", func(t *testing.T) {
		//issue #42
		t.Skip()
	})
}

func TestGetNetworks(t *testing.T) {
	t.Run("ValidToken", func(t *testing.T) {
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		err = json.NewDecoder(response.Body).Decode(&Networks)
		assert.Nil(t, err, err)
	})
	t.Run("InvalidToken", func(t *testing.T) {
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks", "badkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.Equal(t, http.StatusUnauthorized, message.Code)
		assert.Equal(t, "W1R3: You are unauthorized to access this endpoint.", message.Message)
	})
}

func TestGetNetwork(t *testing.T) {
	t.Run("ValidToken", func(t *testing.T) {
		var network models.Network
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		err = json.NewDecoder(response.Body).Decode(&network)
		assert.Nil(t, err, err)
		assert.Equal(t, "skynet", network.DisplayName)
	})
	t.Run("InvalidToken", func(t *testing.T) {
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks/skynet", "badkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.Equal(t, http.StatusUnauthorized, message.Code)
		assert.Equal(t, "W1R3: You are unauthorized to access this endpoint.", message.Message)
	})
	t.Run("InvalidNetwork", func(t *testing.T) {
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks/badnetwork", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, "W1R3: This network does not exist.", message.Message)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}

func TestGetNetworkNodeNumber(t *testing.T) {
	t.Run("ValidKey", func(t *testing.T) {
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks/skynet/numnodes", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message int
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		//assert.Equal(t, "W1R3: This network does not exist.", message.Message)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("InvalidKey", func(t *testing.T) {
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks/skynet/numnodes", "badkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.Equal(t, http.StatusUnauthorized, message.Code)
		assert.Equal(t, "W1R3: You are unauthorized to access this endpoint.", message.Message)
	})
	t.Run("BadNetwork", func(t *testing.T) {
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks/badnetwork/numnodes", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, "W1R3: This network does not exist.", message.Message)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}

func TestDeleteNetwork(t *testing.T) {
	t.Run("InvalidKey", func(t *testing.T) {
		response, err := api(t, "", http.MethodDelete, "http://localhost:8081/api/networks/skynet", "badkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.Equal(t, http.StatusUnauthorized, message.Code)
		assert.Equal(t, "W1R3: You are unauthorized to access this endpoint.", message.Message)
	})
	t.Run("ValidKey", func(t *testing.T) {
		response, err := api(t, "", http.MethodDelete, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message mongo.DeleteResult
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, int64(1), message.DeletedCount)

	})
	t.Run("BadNetwork", func(t *testing.T) {
		response, err := api(t, "", http.MethodDelete, "http://localhost:8081/api/networks/badnetwork", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, "W1R3: This network does not exist.", message.Message)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
	t.Run("NodesExist", func(t *testing.T) {
		t.Skip()
	})
	//Create Network for follow-on tests
	createNetwork(t)
}

func TestCreateAccessKey(t *testing.T) {
	key := models.AccessKey{}
	key.Name = "skynet"
	key.Uses = 10
	t.Run("MultiUse", func(t *testing.T) {
		response, err := api(t, key, http.MethodPost, "http://localhost:8081/api/networks/skynet/keys", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		message, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err, err)
		assert.NotNil(t, message, message)
		returnedkey := getKey(t, key.Name)
		assert.Equal(t, key.Name, returnedkey.Name)
		assert.Equal(t, key.Uses, returnedkey.Uses)
	})
	deleteKey(t, "skynet", "skynet")
	t.Run("ZeroUse", func(t *testing.T) {
		//t.Skip()
		key.Uses = 0
		response, err := api(t, key, http.MethodPost, "http://localhost:8081/api/networks/skynet/keys", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		message, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err, err)
		assert.NotNil(t, message, message)
		returnedkey := getKey(t, key.Name)
		assert.Equal(t, key.Name, returnedkey.Name)
		assert.Equal(t, 1, returnedkey.Uses)
	})
	t.Run("DuplicateAccessKey", func(t *testing.T) {
		//t.Skip()
		//this will fail
		response, err := api(t, key, http.MethodPost, "http://localhost:8081/api/networks/skynet/keys", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
		deleteKey(t, key.Name, "skynet")
	})

	t.Run("InvalidToken", func(t *testing.T) {
		response, err := api(t, key, http.MethodPost, "http://localhost:8081/api/networks/skynet/keys", "badkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, message.Code)
		assert.Equal(t, "W1R3: You are unauthorized to access this endpoint.", message.Message)
	})
	t.Run("BadNetwork", func(t *testing.T) {
		response, err := api(t, key, http.MethodPost, "http://localhost:8081/api/networks/badnetwork/keys", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, "W1R3: This network does not exist.", message.Message)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}

func TestDeleteKey(t *testing.T) {
	t.Run("KeyValid", func(t *testing.T) {
		//fails -- deletecount not returned
		response, err := api(t, "", http.MethodDelete, "http://localhost:8081/api/networks/skynet/keys/skynet", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message mongo.DeleteResult
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, int64(1), message.DeletedCount)
	})
	t.Run("InValidKey", func(t *testing.T) {
		//fails -- status message  not returned
		response, err := api(t, "", http.MethodDelete, "http://localhost:8081/api/networks/skynet/keys/badkey", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, "W1R3: This key does not exist.", message.Message)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
	t.Run("KeyInValidNetwork", func(t *testing.T) {
		response, err := api(t, "", http.MethodDelete, "http://localhost:8081/api/networks/badnetwork/keys/skynet", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, "W1R3: This network does not exist.", message.Message)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
	t.Run("InvalidCredentials", func(t *testing.T) {
		response, err := api(t, "", http.MethodDelete, "http://localhost:8081/api/networks/skynet/keys/skynet", "badkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, message.Code)
		assert.Equal(t, "W1R3: You are unauthorized to access this endpoint.", message.Message)
	})
}

func TestGetKeys(t *testing.T) {
	createKey(t)
	t.Run("Valid", func(t *testing.T) {
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks/skynet/keys", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		var keys []models.AccessKey
		err = json.NewDecoder(response.Body).Decode(&keys)
		assert.Nil(t, err, err)
	})
	//deletekeys
	t.Run("InvalidNetwork", func(t *testing.T) {
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks/badnetwork/keys", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, "W1R3: This network does not exist.", message.Message)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
	t.Run("InvalidCredentials", func(t *testing.T) {
		response, err := api(t, "", http.MethodGet, "http://localhost:8081/api/networks/skynet/keys", "badkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, message.Code)
		assert.Equal(t, "W1R3: You are unauthorized to access this endpoint.", message.Message)
	})
}

func TestUpdateNetwork(t *testing.T) {
	var returnedNetwork models.Network
	t.Run("UpdateNetID", func(t *testing.T) {
		type Network struct {
			NetID string
		}
		var network Network
		network.NetID = "wirecat"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, network.NetID, returnedNetwork.NetID)
	})
	t.Run("NetIDInvalidCredentials", func(t *testing.T) {
		type Network struct {
			NetID string
		}
		var network Network
		network.NetID = "wirecat"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "badkey")
		assert.Nil(t, err, err)
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnauthorized, message.Code)
		assert.Equal(t, "W1R3: You are unauthorized to access this endpoint.", message.Message)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
	t.Run("InvalidNetwork", func(t *testing.T) {
		type Network struct {
			NetID string
		}
		var network Network
		network.NetID = "wirecat"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/badnetwork", "secretkey")
		assert.Nil(t, err, err)
		defer response.Body.Close()
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusNotFound, message.Code)
		assert.Equal(t, "W1R3: This network does not exist.", message.Message)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
	t.Run("UpdateNetIDTooLong", func(t *testing.T) {
		type Network struct {
			NetID string
		}
		var network Network
		network.NetID = "wirecat-skynet"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})
	t.Run("UpdateAddress", func(t *testing.T) {
		type Network struct {
			AddressRange string
		}
		var network Network
		network.AddressRange = "10.0.0.1/24"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, network.AddressRange, returnedNetwork.AddressRange)
	})
	t.Run("UpdateAddressInvalid", func(t *testing.T) {
		type Network struct {
			AddressRange string
		}
		var network Network
		network.AddressRange = "10.0.0.1/36"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})
	t.Run("UpdateDisplayName", func(t *testing.T) {
		type Network struct {
			DisplayName string
		}
		var network Network
		network.DisplayName = "wirecat"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, network.DisplayName, returnedNetwork.DisplayName)

	})
	t.Run("UpdateDisplayNameInvalidName", func(t *testing.T) {
		type Network struct {
			DisplayName string
		}
		var network Network
		//create name that is longer than 100 chars
		name := ""
		for i := 0; i < 101; i++ {
			name = name + "a"
		}
		network.DisplayName = name
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnprocessableEntity, message.Code)
		assert.Equal(t, "W1R3: Field validation for 'DisplayName' failed.", message.Message)
		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})
	t.Run("UpdateInterface", func(t *testing.T) {
		type Network struct {
			DefaultInterface string
		}
		var network Network
		network.DefaultInterface = "netmaker"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, network.DefaultInterface, returnedNetwork.DefaultInterface)

	})
	t.Run("UpdateListenPort", func(t *testing.T) {
		type Network struct {
			DefaultListenPort int32
		}
		var network Network
		network.DefaultListenPort = 6000
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, network.DefaultListenPort, returnedNetwork.DefaultListenPort)
	})
	t.Run("UpdateListenPortInvalidPort", func(t *testing.T) {
		type Network struct {
			DefaultListenPort int32
		}
		var network Network
		network.DefaultListenPort = 1023
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnprocessableEntity, message.Code)
		assert.Equal(t, "W1R3: Field validation for 'DefaultListenPort' failed.", message.Message)
		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})
	t.Run("UpdatePostUP", func(t *testing.T) {
		type Network struct {
			DefaultPostUp string
		}
		var network Network
		network.DefaultPostUp = "sudo wg add-conf wc-netmaker /etc/wireguard/peers/conf"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, network.DefaultPostUp, returnedNetwork.DefaultPostUp)
	})
	t.Run("UpdatePreUP", func(t *testing.T) {
		type Network struct {
			DefaultPreUp string
		}
		var network Network
		network.DefaultPreUp = "test string"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, network.DefaultPreUp, returnedNetwork.DefaultPreUp)
	})
	t.Run("UpdateKeepAlive", func(t *testing.T) {
		type Network struct {
			DefaultKeepalive int32
		}
		var network Network
		network.DefaultKeepalive = 60
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, network.DefaultKeepalive, returnedNetwork.DefaultKeepalive)
	})
	t.Run("UpdateKeepAliveTooBig", func(t *testing.T) {
		type Network struct {
			DefaultKeepAlive int32
		}
		var network Network
		network.DefaultKeepAlive = 1001
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnprocessableEntity, message.Code)
		assert.Equal(t, "W1R3: Field validation for 'DefaultKeepAlive' failed.", message.Message)
		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})
	t.Run("UpdateSaveConfig", func(t *testing.T) {
		//causes panic
		t.Skip()
		type Network struct {
			DefaultSaveConfig *bool
		}
		var network Network
		value := false
		network.DefaultSaveConfig = &value
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, *network.DefaultSaveConfig, *returnedNetwork.DefaultSaveConfig)
	})
	t.Run("UpdateManualSignUP", func(t *testing.T) {
		t.Skip()
		type Network struct {
			AllowManualSignUp *bool
		}
		var network Network
		value := true
		network.AllowManualSignUp = &value
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, *network.AllowManualSignUp, *returnedNetwork.AllowManualSignUp)
	})
	t.Run("DefaultCheckInterval", func(t *testing.T) {
		type Network struct {
			DefaultCheckInInterval int32
		}
		var network Network
		network.DefaultCheckInInterval = 6000
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, network.DefaultCheckInInterval, returnedNetwork.DefaultCheckInInterval)
	})
	t.Run("DefaultCheckIntervalTooBig", func(t *testing.T) {
		type Network struct {
			DefaultCheckInInterval int32
		}
		var network Network
		network.DefaultCheckInInterval = 100001
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		var message models.ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&message)
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusUnprocessableEntity, message.Code)
		assert.Equal(t, "W1R3: Field validation for 'DefaultCheckInInterval' failed.", message.Message)
		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})
	t.Run("MultipleFields", func(t *testing.T) {
		type Network struct {
			DisplayName       string
			DefaultListenPort int32
		}
		var network Network
		network.DefaultListenPort = 7777
		network.DisplayName = "multi"
		response, err := api(t, network, http.MethodPut, "http://localhost:8081/api/networks/skynet", "secretkey")
		assert.Nil(t, err, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&returnedNetwork)
		assert.Nil(t, err, err)
		assert.Equal(t, network.DisplayName, returnedNetwork.DisplayName)
		assert.Equal(t, network.DefaultListenPort, returnedNetwork.DefaultListenPort)
	})
}
