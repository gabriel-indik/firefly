// Copyright © 2024 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
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

package apiserver

import (
	"net/http"
	"strconv"

	"github.com/hyperledger/firefly-common/pkg/ffapi"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/hyperledger/firefly/pkg/database"
)

var getSubscriptionEventsFiltered = &ffapi.Route{
	Name:   "getSubscriptionEventsFiltered",
	Path:   "subscriptions/{subid}/events",
	Method: http.MethodGet,
	PathParams: []*ffapi.PathParam{
		{Name: "subid", Description: coremsgs.APIParamsSubscriptionID},
	},
	QueryParams: []*ffapi.QueryParam{
		{Name: "startsequence", IsBool: false, Description: coremsgs.APIFetchDataDesc},
		{Name: "endsequence", IsBool: false, Description: coremsgs.APIFetchDataDesc},
	},
	FilterFactory:   database.EventQueryFactory,
	Description:     coremsgs.APIEndpointsGetSubscriptionEventsFiltered,
	JSONInputValue:  nil,
	JSONOutputValue: func() interface{} { return []*core.Event{} },
	JSONOutputCodes: []int{http.StatusOK},
	Extensions: &coreExtensions{
		CoreJSONHandler: func(r *ffapi.APIRequest, cr *coreRequest) (output interface{}, err error) {
			subscription, _ := cr.or.GetSubscriptionByID(cr.ctx, r.PP["subid"])

			startSeq, err := strconv.Atoi(r.QP["startsequence"])
			if err != nil {
				return nil, err
			}

			endSeq, err := strconv.Atoi(r.QP["endsequence"])
			if err != nil {
				return nil, err
			}

			return r.FilterResult(cr.or.GetSubscriptionEventsHistorical(cr.ctx, subscription, r.Filter, startSeq, endSeq))
		},
	},
}
