/**
 * Tencent is pleased to support the open source community by making Polaris available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package pluginsdk

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/polaris-contrib/polaris-server-remote-plugin-common/api"
)

// MarshalRequest convert any proto to Request proto.
func MarshalRequest(msg proto.Message) (*api.Request, error) {
	anyPB, err := anypb.New(msg)
	if err != nil {
		return nil, err
	}
	req := &api.Request{
		Payload: anyPB,
	}
	return req, nil
}

// MarshalResponse convert any proto to Response proto.
func MarshalResponse(msg proto.Message) (*api.Response, error) {
	anyPB, err := anypb.New(msg)
	if err != nil {
		return nil, err
	}
	return &api.Response{
		Reply: anyPB,
	}, nil
}

// UnmarshalRequest convert request to plugin payload proto.
func UnmarshalRequest(req *api.Request, msg proto.Message) error {
	if req == nil {
		return nil
	}
	anyPB, err := anypb.New(msg)
	if err != nil {
		return err
	}
	payload := req.Payload
	if payload == nil {
		return nil
	}
	if payload.GetTypeUrl() != anyPB.GetTypeUrl() {
		return fmt.Errorf("got wrong payload type url: %s", payload.GetTypeUrl())
	}

	return anypb.UnmarshalTo(payload, msg, proto.UnmarshalOptions{})
}

// UnmarshalResponse convert response to plugin reply proto.
func UnmarshalResponse(resp *api.Response, msg proto.Message) error {
	if resp == nil {
		return nil
	}

	anyPB, err := anypb.New(msg)
	if err != nil {
		return err
	}

	reply := resp.Reply
	if reply == nil {
		return fmt.Errorf("nil plugin reply")
	}

	if anyPB.GetTypeUrl() != reply.GetTypeUrl() {
		return fmt.Errorf("got wrong reply type url: %s", reply.GetTypeUrl())
	}

	return anypb.UnmarshalTo(reply, msg, proto.UnmarshalOptions{})
}
