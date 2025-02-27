// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package multicloud

import (
	"strings"

	"yunion.io/x/pkg/errors"

	"yunion.io/x/onecloud/pkg/cloudprovider"
)

type STagBase struct {
}

func (self STagBase) GetSysTags() map[string]string {
	return nil
}

func (self STagBase) GetTags() (map[string]string, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetTags")
}

func (self STagBase) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type QcloudTags struct {
	TagSet []STag

	// Redis
	InstanceTags []STag
	// Elasticsearch
	TagList []STag
}

func (self *QcloudTags) GetTags() (map[string]string, error) {
	ret := map[string]string{}
	for _, tag := range self.TagSet {
		ret[tag.Key] = tag.Value
	}
	for _, tag := range self.InstanceTags {
		ret[tag.TagKey] = tag.TagValue
	}
	for _, tag := range self.TagList {
		ret[tag.TagKey] = tag.TagValue
	}
	return ret, nil
}

func (self *QcloudTags) GetSysTags() map[string]string {
	return nil
}

func (self *QcloudTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type STag struct {
	TagKey   string
	TagValue string

	Key   string
	Value string
}

type AliyunTags struct {
	Tags struct {
		Tag []STag
	}
}

func (self *AliyunTags) GetTags() (map[string]string, error) {
	ret := map[string]string{}
	for _, tag := range self.Tags.Tag {
		if strings.HasPrefix(tag.TagKey, "aliyun") || strings.HasPrefix(tag.TagKey, "acs:") ||
			strings.HasSuffix(tag.Key, "aliyun") || strings.HasPrefix(tag.Key, "acs:") {
			continue
		}
		if len(tag.TagKey) > 0 {
			ret[tag.TagKey] = tag.TagValue
		} else if len(tag.Key) > 0 {
			ret[tag.Key] = tag.Value
		}
	}
	return ret, nil
}

func (self *AliyunTags) GetSysTags() map[string]string {
	ret := map[string]string{}
	for _, tag := range self.Tags.Tag {
		if strings.HasPrefix(tag.TagKey, "aliyun") || strings.HasPrefix(tag.TagKey, "acs:") ||
			strings.HasPrefix(tag.Key, "aliyun") || strings.HasPrefix(tag.Key, "acs:") {
			if len(tag.TagKey) > 0 {
				ret[tag.TagKey] = tag.TagValue
			} else if len(tag.Key) > 0 {
				ret[tag.Key] = tag.Value
			}
		}
	}
	return ret
}

func (self *AliyunTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type ApsaraTags struct {
	Tags struct {
		Tag []STag
	}
}

func (self *ApsaraTags) GetTags() (map[string]string, error) {
	ret := map[string]string{}
	for _, tag := range self.Tags.Tag {
		if strings.HasPrefix(tag.TagKey, "aliyun") || strings.HasPrefix(tag.TagKey, "acs:") ||
			strings.HasSuffix(tag.Key, "aliyun") || strings.HasPrefix(tag.Key, "acs:") {
			continue
		}
		if len(tag.TagKey) > 0 {
			ret[tag.TagKey] = tag.TagValue
		} else if len(tag.Key) > 0 {
			ret[tag.Key] = tag.Value
		}
	}
	return ret, nil
}

func (self *ApsaraTags) GetSysTags() map[string]string {
	ret := map[string]string{}
	for _, tag := range self.Tags.Tag {
		if strings.HasPrefix(tag.TagKey, "aliyun") || strings.HasPrefix(tag.TagKey, "acs:") ||
			strings.HasPrefix(tag.Key, "aliyun") || strings.HasPrefix(tag.Key, "acs:") {
			if len(tag.TagKey) > 0 {
				ret[tag.TagKey] = tag.TagValue
			} else if len(tag.Key) > 0 {
				ret[tag.Key] = tag.Value
			}
		}
	}
	return ret
}

func (self *ApsaraTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type GoogleTags struct {
	Labels map[string]string
}

func (self *GoogleTags) GetTags() (map[string]string, error) {
	return self.Labels, nil
}

func (self *GoogleTags) GetSysTags() map[string]string {
	return nil
}

func (self *GoogleTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type AzureTags struct {
	Tags map[string]string
}

func (self *AzureTags) GetTags() (map[string]string, error) {
	return self.Tags, nil
}

func (self *AzureTags) GetSysTags() map[string]string {
	return nil
}

func (self *AzureTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type AwsTags struct {
	TagSet []STag
}

func (self *AwsTags) GetTags() (map[string]string, error) {
	ret := map[string]string{}
	for _, tag := range self.TagSet {
		if tag.Key == "Name" || tag.Key == "Description" {
			continue
		}
		ret[tag.Key] = tag.Value
	}
	return ret, nil
}

func (self *AwsTags) GetSysTags() map[string]string {
	return nil
}

func (self *AwsTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type CtyunTags struct {
}

func (self *CtyunTags) GetTags() (map[string]string, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetTags")
}

func (self *CtyunTags) GetSysTags() map[string]string {
	return nil
}

func (self *CtyunTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type EcloudTags struct {
}

func (self *EcloudTags) GetTags() (map[string]string, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetTags")
}

func (self *EcloudTags) GetSysTags() map[string]string {
	return nil
}

func (self *EcloudTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type JdcloudTags struct {
}

func (jt *JdcloudTags) GetTags() (map[string]string, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetTags")
}

func (self *JdcloudTags) GetSysTags() map[string]string {
	return nil
}

func (self *JdcloudTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type HuaweiTags struct {
	Tags []string
}

func (self *HuaweiTags) GetTags() (map[string]string, error) {
	tags := map[string]string{}
	for _, kv := range self.Tags {
		splited := strings.Split(kv, "=")
		if len(splited) == 2 {
			tags[splited[0]] = splited[1]
		}
	}
	return tags, nil
}

func (self *HuaweiTags) GetSysTags() map[string]string {
	return nil
}

func (self *HuaweiTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type HuaweiDiskTags struct {
	Tags map[string]string
}

func (self *HuaweiDiskTags) GetTags() (map[string]string, error) {
	return self.Tags, nil
}

func (self *HuaweiDiskTags) GetSysTags() map[string]string {
	return nil
}

func (self *HuaweiDiskTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type OpenStackTags struct {
	Metadata map[string]string
}

func (self *OpenStackTags) GetTags() (map[string]string, error) {
	return self.Metadata, nil
}

func (self *OpenStackTags) GetSysTags() map[string]string {
	return nil
}

func (self *OpenStackTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type UcloudTags struct {
}

func (self *UcloudTags) GetTags() (map[string]string, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetTags")
}

func (self *UcloudTags) GetSysTags() map[string]string {
	return nil
}

func (self *UcloudTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}

type ZStackTags struct {
}

func (self *ZStackTags) GetTags() (map[string]string, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetTags")
}

func (self *ZStackTags) GetSysTags() map[string]string {
	return nil
}

func (self *ZStackTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}
