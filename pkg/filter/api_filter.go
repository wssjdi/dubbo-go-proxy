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

package filter

import (
	"net/http"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

func init() {
	extension.SetFilterFunc(constant.HttpApiFilter, ApiFilter())
}

// ApiFilter url match api
// 这里扩展动态从配置中心加载api配置数据，考虑使用zookeeper做为配置中心，可以监听上线下线
func ApiFilter() context.FilterFunc {
	return func(c context.Context) {
		url := c.GetUrl()
		method := c.GetMethod()

		if api, ok := model.EmptyApi.FindApi(url); ok {
			if !api.MatchMethod(method) {
				c.WriteWithStatus(http.StatusMethodNotAllowed, constant.Default405Body)
				c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
				c.Abort()
				return
			}

			if !api.IsOk(api.Name) {
				c.WriteWithStatus(http.StatusNotAcceptable, constant.Default406Body)
				c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
				c.Abort()
				return
			}

			c.Api(api)
			c.Next()
		} else {
			c.WriteWithStatus(http.StatusNotFound, constant.Default404Body)
			c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
			c.Abort()
		}
	}
}
