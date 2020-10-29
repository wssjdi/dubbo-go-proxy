/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package proxy

import (
	"encoding/json"
	"sync"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	_ "github.com/dubbogo/dubbo-go-proxy/pkg/filter"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/dubbogo/dubbo-go-proxy/pkg/service"
	_ "github.com/dubbogo/dubbo-go-proxy/pkg/service/api"
)

// Proxy
type Proxy struct {
	startWG sync.WaitGroup
}

// Start proxy start
func (p *Proxy) Start() {
	conf := config.GetBootstrap()

	p.startWG.Add(1)

	defer func() {
		if re := recover(); re != nil {
			logger.Error(re)
			// TODO stop
		}
	}()

	p.beforeStart()

	listeners := conf.GetListeners()

	for _, s := range listeners {
		ls := ListenerService{Listener: &s}
		go ls.Start()
	}
}

func (p *Proxy) beforeStart() {
	dubbo.SingleDubboClient().Init()

	// TODO mock api register
	ads := extension.GetMustApiDiscoveryService(constant.LocalMemoryApiDiscoveryService)

	a1 := &model.Api{
		Name:     "/api/v1/station/member/name",
		ITypeStr: "HTTP",
		OTypeStr: "DUBBO",
		Method:   "POST",
		Status:   1,
		Metadata: map[string]dubbo.DubboMetadata{
			"dubbo": {
				ApplicationName: "station",
				Group:           "",
				Version:         "0.0.1",
				Interface:       "com.chebada.station.service.MemberService",
				Method:          "getMemberByName",
				Retries:         "1",
				Types: []string{
					"java.lang.Long",
					"java.lang.String",
				},
				ClusterName: "test_dubbo",
			},
		},
	}

	j1, _ := json.Marshal(a1)
	ads.AddApi(*service.NewDiscoveryRequest(j1))

	a2 := &model.Api{
		Name:     "/api/v1/station/member",
		ITypeStr: "HTTP",
		OTypeStr: "DUBBO",
		Method:   "POST",
		Status:   1,
		Metadata: map[string]dubbo.DubboMetadata{
			"dubbo": {
				ApplicationName: "station",
				Group:           "",
				Version:         "0.0.1",
				Interface:       "com.chebada.station.service.MemberService",
				Method:          "getMember",
				Retries:         "1",
				Types: []string{
					"java.lang.Long",
				},
				ClusterName: "test_dubbo",
			},
		},
	}

	j2, _ := json.Marshal(a2)
	ads.AddApi(*service.NewDiscoveryRequest(j2))

	a3 := &model.Api{
		Name:     "/api/v1/station/airport",
		ITypeStr: "HTTP",
		OTypeStr: "DUBBO",
		Method:   "POST",
		Status:   1,
		Metadata: map[string]dubbo.DubboMetadata{
			"dubbo": {
				ApplicationName: "station",
				Group:           "",
				Version:         "0.0.1",
				Interface:       "com.chebada.station.service.AirportService",
				Method:          "findAll",
				Retries:         "1",
				Types: []string{
					"com.chebada.station.request.Area",
				},
				ClusterName: "test_dubbo",
			},
		},
	}

	j3, _ := json.Marshal(a3)
	ads.AddApi(*service.NewDiscoveryRequest(j3))

	a4 := &model.Api{
		Name:     "/api/v1/station/bus",
		ITypeStr: "HTTP",
		OTypeStr: "DUBBO",
		Method:   "POST",
		Status:   1,
		Metadata: map[string]dubbo.DubboMetadata{
			"dubbo": {
				ApplicationName: "station",
				Group:           "",
				Version:         "0.0.1",
				Interface:       "com.chebada.station.service.BusStationService",
				Method:          "findAll",
				Retries:         "1",
				Types: []string{
					"com.chebada.station.request.Area",
				},
				ClusterName: "test_dubbo",
			},
		},
	}

	j4, _ := json.Marshal(a4)
	ads.AddApi(*service.NewDiscoveryRequest(j4))

	a5 := &model.Api{
		Name:     "/api/v1/station/train",
		ITypeStr: "HTTP",
		OTypeStr: "DUBBO",
		Method:   "POST",
		Status:   1,
		Metadata: map[string]dubbo.DubboMetadata{
			"dubbo": {
				ApplicationName: "station",
				Group:           "",
				Version:         "0.0.1",
				Interface:       "com.chebada.station.service.TrainStationService",
				Method:          "findAll",
				Retries:         "1",
				Types: []string{
					"com.chebada.station.request.Area",
				},
				ClusterName: "test_dubbo",
			},
		},
	}

	j5, _ := json.Marshal(a5)
	ads.AddApi(*service.NewDiscoveryRequest(j5))
}

// NewProxy create proxy
func NewProxy() *Proxy {
	return &Proxy{
		startWG: sync.WaitGroup{},
	}
}

func Start(bs *model.Bootstrap) {
	logger.Infof("[dubboproxy go] start by config : %+v", bs)

	proxy := NewProxy()
	proxy.Start()

	proxy.startWG.Wait()
}
