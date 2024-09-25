package logic

import (
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/gravitl/netmaker/database"
	"github.com/gravitl/netmaker/models"
)

// CreateDefaultAclNetworkPolicies - create default acl network policies
func CreateDefaultAclNetworkPolicies(netID models.NetworkID) {
	defaultDeviceAcl := models.Acl{
		ID:        uuid.New(),
		Default:   true,
		Name:      "all-nodes",
		NetworkID: netID,
		RuleType:  models.DevicePolicy,
		Src: []models.AclPolicyTag{
			{
				ID:    models.DeviceAclID,
				Value: "*",
			}},
		Dst: []models.AclPolicyTag{
			{
				ID:    models.DeviceAclID,
				Value: "*",
			}},
		AllowedDirection: models.TrafficDirectionBi,
		Enabled:          true,
		CreatedBy:        "auto",
		CreatedAt:        time.Now().UTC(),
	}
	InsertAcl(defaultDeviceAcl)
	defaultUserAcl := models.Acl{
		ID:        uuid.New(),
		Default:   true,
		Name:      "all-users",
		NetworkID: netID,
		RuleType:  models.UserPolicy,
		Src: []models.AclPolicyTag{
			{
				ID:    models.UserAclID,
				Value: "*",
			},
			{
				ID:    models.UserGroupAclID,
				Value: "*",
			},
		},
		Dst: []models.AclPolicyTag{{
			ID:    models.DeviceAclID,
			Value: "*",
		}},
		AllowedDirection: models.TrafficDirectionUni,
		Enabled:          true,
		CreatedBy:        "auto",
		CreatedAt:        time.Now().UTC(),
	}
	InsertAcl(defaultUserAcl)
}

// DeleteDefaultNetworkPolicies - deletes all default network acl policies
func DeleteDefaultNetworkPolicies(netId models.NetworkID) {
	acls, _ := ListAcls(netId)
	for _, acl := range acls {
		if acl.NetworkID == netId && acl.Default {
			DeleteAcl(acl)
		}
	}
}

// InsertAcl - creates acl policy
func InsertAcl(a models.Acl) error {
	d, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return database.Insert(a.ID.String(), string(d), database.ACLS_TABLE_NAME)
}

// GetAcl - gets acl info by id
func GetAcl(aID string) (models.Acl, error) {
	a := models.Acl{}
	d, err := database.FetchRecord(database.ACLS_TABLE_NAME, aID)
	if err != nil {
		return a, err
	}
	err = json.Unmarshal([]byte(d), &a)
	if err != nil {
		return a, err
	}
	return a, nil
}

// IsAclPolicyValid - validates if acl policy is valid
func IsAclPolicyValid(acl models.Acl) bool {
	//check if src and dst are valid
	isValid := false
	switch acl.RuleType {
	case models.UserPolicy:
		// src list should only contain users
		for _, srcI := range acl.Src {

			if srcI.ID == "" || srcI.Value == "" {
				break
			}
			if srcI.ID != models.UserAclID &&
				srcI.ID != models.UserGroupAclID {
				break
			}
			// check if user group is valid
			if srcI.ID == models.UserAclID {
				_, err := GetUser(srcI.Value)
				if err != nil {
					break
				}
			} else if srcI.ID == models.UserGroupAclID {
				err := IsGroupValid(models.UserGroupID(srcI.Value))
				if err != nil {
					break
				}
			}

		}
		for _, dstI := range acl.Dst {

			if dstI.ID == "" || dstI.Value == "" {
				break
			}
			if dstI.ID == models.UserAclID ||
				dstI.ID == models.UserGroupAclID {
				break
			}
			if dstI.ID != models.DeviceAclID {
				break
			}
			// check if tag is valid
			_, err := GetTag(models.TagID(dstI.Value))
			if err != nil {
				break
			}
		}
		isValid = true
	case models.DevicePolicy:
		for _, srcI := range acl.Src {
			if srcI.ID == "" || srcI.Value == "" {
				break
			}
			if srcI.ID != models.DeviceAclID {
				break
			}
			// check if tag is valid
			_, err := GetTag(models.TagID(srcI.Value))
			if err != nil {
				break
			}
		}
		for _, dstI := range acl.Dst {

			if dstI.ID == "" || dstI.Value == "" {
				break
			}
			if dstI.ID != models.DeviceAclID {
				break
			}
			// check if tag is valid
			_, err := GetTag(models.TagID(dstI.Value))
			if err != nil {
				break
			}
		}
		isValid = true
	}
	return isValid
}

// UpdateAcl - updates allowed fields on acls and commits to DB
func UpdateAcl(newAcl, acl models.Acl) error {
	if newAcl.Name != "" {
		acl.Name = newAcl.Name
	}
	acl.Src = newAcl.Src
	acl.Dst = newAcl.Dst
	acl.AllowedDirection = newAcl.AllowedDirection
	acl.Enabled = newAcl.Enabled
	d, err := json.Marshal(acl)
	if err != nil {
		return err
	}
	return database.Insert(acl.ID.String(), string(d), database.ACLS_TABLE_NAME)
}

// DeleteAcl - deletes acl policy
func DeleteAcl(a models.Acl) error {
	return database.DeleteRecord(database.ACLS_TABLE_NAME, a.ID.String())
}

// GetDefaultPolicy - fetches default policy in the network by ruleType
func GetDefaultPolicy(netID models.NetworkID, ruleType models.AclPolicyType) (models.Acl, error) {
	acls, _ := ListAcls(netID)
	for _, acl := range acls {
		if acl.Default && acl.RuleType == ruleType {
			return acl, nil
		}
	}
	return models.Acl{}, errors.New("default rule not found")
}

// ListUserPolicies - lists all acl policies enforced on an user
func ListUserPolicies(u models.User) []models.Acl {
	data, err := database.FetchRecords(database.TAG_TABLE_NAME)
	if err != nil && !database.IsEmptyRecord(err) {
		return []models.Acl{}
	}
	acls := []models.Acl{}
	for _, dataI := range data {
		acl := models.Acl{}
		err := json.Unmarshal([]byte(dataI), &acl)
		if err != nil {
			continue
		}

		if acl.RuleType == models.UserPolicy {
			srcMap := convAclTagToValueMap(acl.Src)
			if _, ok := srcMap[u.UserName]; ok {
				acls = append(acls, acl)
			} else {
				// check for user groups
				for gID := range u.UserGroups {
					if _, ok := srcMap[gID.String()]; ok {
						acls = append(acls, acl)
						break
					}
				}
			}

		}
	}
	return acls
}

// ListUserPoliciesByNetwork - lists all acl user policies in a network
func ListUserPoliciesByNetwork(netID models.NetworkID) []models.Acl {
	data, err := database.FetchRecords(database.TAG_TABLE_NAME)
	if err != nil && !database.IsEmptyRecord(err) {
		return []models.Acl{}
	}
	acls := []models.Acl{}
	for _, dataI := range data {
		acl := models.Acl{}
		err := json.Unmarshal([]byte(dataI), &acl)
		if err != nil {
			continue
		}
		if acl.NetworkID == netID && acl.RuleType == models.UserPolicy {
			acls = append(acls, acl)
		}
	}
	return acls
}

// listDevicePolicies - lists all device policies in a network
func listDevicePolicies(netID models.NetworkID) []models.Acl {
	data, err := database.FetchRecords(database.TAG_TABLE_NAME)
	if err != nil && !database.IsEmptyRecord(err) {
		return []models.Acl{}
	}
	acls := []models.Acl{}
	for _, dataI := range data {
		acl := models.Acl{}
		err := json.Unmarshal([]byte(dataI), &acl)
		if err != nil {
			continue
		}
		if acl.NetworkID == netID && acl.RuleType == models.DevicePolicy {
			acls = append(acls, acl)
		}
	}
	return acls
}

// ListAcls - lists all acl policies
func ListAcls(netID models.NetworkID) ([]models.Acl, error) {
	data, err := database.FetchRecords(database.TAG_TABLE_NAME)
	if err != nil && !database.IsEmptyRecord(err) {
		return []models.Acl{}, err
	}
	acls := []models.Acl{}
	for _, dataI := range data {
		acl := models.Acl{}
		err := json.Unmarshal([]byte(dataI), &acl)
		if err != nil {
			continue
		}
		if acl.NetworkID == netID {
			acls = append(acls, acl)
		}
	}
	return acls, nil
}

func convAclTagToValueMap(acltags []models.AclPolicyTag) map[string]struct{} {
	aclValueMap := make(map[string]struct{})
	for _, aclTagI := range acltags {
		aclValueMap[aclTagI.ID.String()] = struct{}{}
	}
	return aclValueMap
}

// IsNodeAllowedToCommunicate - check node is allowed to communicate with the peer
func IsNodeAllowedToCommunicate(node, peer models.Node) bool {
	// check default policy if all allowed return true
	defaultPolicy, err := GetDefaultPolicy(models.NetworkID(node.Network), models.DevicePolicy)
	if err == nil {
		if defaultPolicy.Enabled {
			return true
		}
	}
	// list device policies
	policies := listDevicePolicies(models.NetworkID(peer.Network))
	for _, policy := range policies {
		srcMap := convAclTagToValueMap(policy.Src)
		dstMap := convAclTagToValueMap(policy.Dst)
		for tagID := range peer.Tags {
			if _, ok := dstMap[tagID.String()]; ok {
				for tagID := range node.Tags {
					if _, ok := srcMap[tagID.String()]; ok {
						return true
					}
				}
			}
			if _, ok := srcMap[tagID.String()]; ok {
				for tagID := range node.Tags {
					if _, ok := dstMap[tagID.String()]; ok {
						return true
					}
				}
			}
		}
	}
	return false
}

// SortTagEntrys - Sorts slice of Tag entries by their id
func SortAclEntrys(acls []models.Acl) {
	sort.Slice(acls, func(i, j int) bool {
		return acls[i].Name < acls[j].Name
	})
}
