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

import com.google.common.collect.Sets;
import java.util.Collection;
import java.util.Objects;
import java.util.Set;
import org.springframework.core.convert.converter.Converter;
import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.security.oauth2.server.resource.authentication.JwtGrantedAuthoritiesConverter;

public class CustomJwtAuthenticationConverter implements
    Converter<Jwt, AbstractAuthenticationToken> {

  private final Converter<Jwt, Collection<GrantedAuthority>> jwtGrantedAuthoritiesConverter = new JwtGrantedAuthoritiesConverter();

  @Override
  public AbstractAuthenticationToken convert(Jwt source) {
    Set<GrantedAuthority> authorities;
    String principal;
    if (source.getIssuer().toString().startsWith("https://cloud.google.com/iap")) {
      principal = (String) source.getClaimAsMap("gcip").get("email");
      // You can add additional specific authorities extraction here
      authorities = Sets.newHashSet(
          new SimpleGrantedAuthority("ROLE_IAP_USER")
      );
      authorities.addAll(Objects.requireNonNull(jwtGrantedAuthoritiesConverter.convert(source)));
    } else {
      principal = source.getClaimAsString("sub");
      // You can add additional specific authorities extraction here
      authorities = Sets.newHashSet(
          new SimpleGrantedAuthority("ROLE_API_USER")
      );
      authorities.addAll(Objects.requireNonNull(jwtGrantedAuthoritiesConverter.convert(source)));
    }
    return new JwtAuthenticationToken(source, authorities, principal);
  }
}
