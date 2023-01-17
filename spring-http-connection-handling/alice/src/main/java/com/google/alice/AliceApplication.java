// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package com.google.alice;

import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.impl.conn.PoolingHttpClientConnectionManager;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.http.client.HttpComponentsClientHttpRequestFactory;
import org.springframework.web.client.RestTemplate;

@SpringBootApplication
public class AliceApplication {

	public static void main(String[] args) {
		SpringApplication.run(AliceApplication.class, args);
	}
/* Standard resttemplate
	@Bean(name = "restTemplate")
	public RestTemplate prepareRestTemplate() {
		return new RestTemplate();
	}
*/

	@Bean(name = "restTemplate")
	public RestTemplate prepareRestTemplateWithPooling() {
		PoolingHttpClientConnectionManager poolingConnectionManager = new PoolingHttpClientConnectionManager();
		poolingConnectionManager.setMaxTotal(100);
		CloseableHttpClient client = HttpClientBuilder.create().setConnectionManager(poolingConnectionManager).disableAutomaticRetries().build();
		HttpComponentsClientHttpRequestFactory factory = new HttpComponentsClientHttpRequestFactory(client);
		factory.setConnectionRequestTimeout(300000);
		factory.setConnectTimeout(300000);
		factory.setReadTimeout(300000);
		return new RestTemplate(factory);
	}
}
