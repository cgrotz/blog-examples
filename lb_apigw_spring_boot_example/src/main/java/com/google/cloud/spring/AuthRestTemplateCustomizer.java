// Copyright 2022 Google LLC
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

package com.google.cloud.spring;

import com.google.auth.oauth2.GoogleCredentials;
import com.google.auth.oauth2.IdToken;
import com.google.auth.oauth2.IdTokenProvider;
import java.util.ArrayList;
import java.util.List;
import org.springframework.boot.web.client.RestTemplateCustomizer;
import org.springframework.http.HttpHeaders;
import org.springframework.http.client.ClientHttpRequestInterceptor;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.stereotype.Component;
import org.springframework.util.CollectionUtils;
import org.springframework.web.client.RestTemplate;

@Component
public class AuthRestTemplateCustomizer implements RestTemplateCustomizer {

  @Override
  public void customize(RestTemplate restTemplate) {
    List<ClientHttpRequestInterceptor> interceptors
        = restTemplate.getInterceptors();
    if (CollectionUtils.isEmpty(interceptors)) {
      interceptors = new ArrayList<>();
    }
    interceptors.add((request, body, execution) -> {
      // This behaviour could also be limited to the .a.run.app domain to only forward the Identity Token to other Cloud Run services
      GoogleCredentials adCredentials = GoogleCredentials.getApplicationDefault();
      if (adCredentials instanceof IdTokenProvider) {
        IdTokenProvider idTokenProvider = (IdTokenProvider) adCredentials;
        IdToken idToken = idTokenProvider.idTokenWithAudience(
            "https://" + request.getURI().getHost(), null);
        request.getHeaders().add(HttpHeaders.AUTHORIZATION, "Bearer " + idToken.getTokenValue());
      }
      return execution.execute(request, body);
    });

    interceptors.add((request, body, execution) -> {
      JwtAuthenticationToken auth = (JwtAuthenticationToken) SecurityContextHolder.getContext()
          .getAuthentication();
      request.getHeaders()
          .add("X-Forwarded-Authorization", "Bearer " + auth.getToken().getTokenValue());
      return execution.execute(request, body);
    });
    restTemplate.setInterceptors(interceptors);

  }
}