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

import java.util.regex.Matcher;
import java.util.regex.Pattern;
import javax.servlet.http.HttpServletRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.server.resource.BearerTokenError;
import org.springframework.security.oauth2.server.resource.BearerTokenErrors;
import org.springframework.security.oauth2.server.resource.web.BearerTokenResolver;
import org.springframework.util.StringUtils;

public class GoogleCloudBearerTokenResolver implements BearerTokenResolver {

  private static final Logger logger = LoggerFactory.getLogger(
      GoogleCloudBearerTokenResolver.class);
  private static final Pattern authorizationPattern = Pattern.compile(
      "^Bearer \"(?<token>[a-zA-Z0-9-._~+/]+=*)\"$",
      Pattern.CASE_INSENSITIVE);

  @Override
  public String resolve(HttpServletRequest request) {
    String token = resolveFromXForwardedAuthorization(request);
    if (token != null) {
      return token;
    } else {
      return resolveFromXGoogIapJwtAssertion(request);
    }
  }

  private String resolveFromXGoogIapJwtAssertion(HttpServletRequest request) {
    String authorization = request.getHeader("X-Goog-Iap-Jwt-Assertion");
    if (StringUtils.hasLength(authorization)) {
      return authorization.replaceAll("^\"|\"$", "");
    } else {
      return null;
    }
  }

  private String resolveFromXForwardedAuthorization(HttpServletRequest request) {
    String authorization = request.getHeader("X-Forwarded-Authorization");
    if (StringUtils.hasLength(authorization)) {
      if (!StringUtils.startsWithIgnoreCase(authorization, "bearer")) {
        return null;
      } else {
        Matcher matcher = authorizationPattern.matcher(authorization);
        if (!matcher.matches()) {
          BearerTokenError error = BearerTokenErrors.invalidToken(
              "Bearer token in X-Forwarded-Authorization is malformed");
          throw new OAuth2AuthenticationException(error);
        }
        return matcher.group("token");
      }
    } else {
      return null;
    }
  }
}
