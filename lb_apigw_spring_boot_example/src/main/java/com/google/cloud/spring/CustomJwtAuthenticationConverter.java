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
