// Copyright Project Contour Authors
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

package featuretests

import (
	"testing"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	matcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher"
	"github.com/golang/protobuf/ptypes/wrappers"
	projcontour "github.com/projectcontour/contour/apis/projectcontour/v1"
	envoyv2 "github.com/projectcontour/contour/internal/envoy/v2"
	"github.com/projectcontour/contour/internal/fixture"
	"github.com/projectcontour/contour/internal/k8s"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestCorsPolicy(t *testing.T) {
	rh, c, done := setup(t)
	defer done()

	rh.OnAdd(fixture.NewService("svc1").
		WithPorts(v1.ServicePort{Port: 80, TargetPort: intstr.FromInt(8080)}),
	)

	// Allow origin
	rh.OnAdd(fixture.NewProxy("simple").WithSpec(
		projcontour.HTTPProxySpec{
			VirtualHost: &projcontour.VirtualHost{
				Fqdn: "hello.world",
				CORSPolicy: &projcontour.CORSPolicy{
					AllowOrigin: []string{"*"},
				},
			}, Routes: []projcontour.Route{{
				Services: []projcontour.Service{{
					Name: "svc1",
					Port: 80,
				}},
			}},
		}),
	)

	c.Request(routeType).Equals(&v2.DiscoveryResponse{
		Resources: resources(t,
			envoyv2.RouteConfiguration("ingress_http",
				envoyv2.CORSVirtualHost("hello.world",
					&envoy_api_v2_route.CorsPolicy{
						AllowCredentials: &wrappers.BoolValue{Value: false},
						AllowOriginStringMatch: []*matcher.StringMatcher{{
							MatchPattern: &matcher.StringMatcher_Exact{
								Exact: "*",
							},
							IgnoreCase: true,
						}},
					},
					&envoy_api_v2_route.Route{
						Match:  routePrefix("/"),
						Action: routecluster("default/svc1/80/da39a3ee5e"),
					}),
			),
		),
		TypeUrl: routeType,
	})

	// Allow credentials
	rh.OnAdd(fixture.NewProxy("simple").WithSpec(
		projcontour.HTTPProxySpec{
			VirtualHost: &projcontour.VirtualHost{
				Fqdn: "hello.world",
				CORSPolicy: &projcontour.CORSPolicy{
					AllowOrigin:      []string{"*"},
					AllowCredentials: true,
				},
			}, Routes: []projcontour.Route{{
				Services: []projcontour.Service{{
					Name: "svc1",
					Port: 80,
				}},
			}},
		}),
	)

	c.Request(routeType).Equals(&v2.DiscoveryResponse{
		Resources: resources(t,
			envoyv2.RouteConfiguration("ingress_http",
				envoyv2.CORSVirtualHost("hello.world",
					&envoy_api_v2_route.CorsPolicy{
						AllowOriginStringMatch: []*matcher.StringMatcher{{
							MatchPattern: &matcher.StringMatcher_Exact{
								Exact: "*",
							},
							IgnoreCase: true,
						}},
						AllowCredentials: &wrappers.BoolValue{Value: true},
					},
					&envoy_api_v2_route.Route{
						Match:  routePrefix("/"),
						Action: routecluster("default/svc1/80/da39a3ee5e"),
					}),
			),
		),
		TypeUrl: routeType,
	})

	// Allow methods
	rh.OnAdd(fixture.NewProxy("simple").WithSpec(
		projcontour.HTTPProxySpec{
			VirtualHost: &projcontour.VirtualHost{
				Fqdn: "hello.world",
				CORSPolicy: &projcontour.CORSPolicy{
					AllowOrigin:      []string{"*"},
					AllowCredentials: true,
					AllowMethods:     []projcontour.CORSHeaderValue{"GET", "POST", "OPTIONS"},
				},
			}, Routes: []projcontour.Route{{
				Services: []projcontour.Service{{
					Name: "svc1",
					Port: 80,
				}},
			}},
		}),
	)

	c.Request(routeType).Equals(&v2.DiscoveryResponse{
		Resources: resources(t,
			envoyv2.RouteConfiguration("ingress_http",
				envoyv2.CORSVirtualHost("hello.world",
					&envoy_api_v2_route.CorsPolicy{
						AllowOriginStringMatch: []*matcher.StringMatcher{{
							MatchPattern: &matcher.StringMatcher_Exact{
								Exact: "*",
							},
							IgnoreCase: true,
						}},
						AllowCredentials: &wrappers.BoolValue{Value: true},
						AllowMethods:     "GET,POST,OPTIONS",
					},
					&envoy_api_v2_route.Route{
						Match:  routePrefix("/"),
						Action: routecluster("default/svc1/80/da39a3ee5e"),
					}),
			),
		),
		TypeUrl: routeType,
	})

	// Allow headers
	rh.OnAdd(fixture.NewProxy("simple").WithSpec(
		projcontour.HTTPProxySpec{
			VirtualHost: &projcontour.VirtualHost{
				Fqdn: "hello.world",
				CORSPolicy: &projcontour.CORSPolicy{
					AllowOrigin:      []string{"*"},
					AllowCredentials: true,
					AllowHeaders:     []projcontour.CORSHeaderValue{"custom-header-1", "custom-header-2"},
				},
			}, Routes: []projcontour.Route{{
				Services: []projcontour.Service{{
					Name: "svc1",
					Port: 80,
				}},
			}},
		}),
	)

	c.Request(routeType).Equals(&v2.DiscoveryResponse{
		Resources: resources(t,
			envoyv2.RouteConfiguration("ingress_http",
				envoyv2.CORSVirtualHost("hello.world",
					&envoy_api_v2_route.CorsPolicy{
						AllowOriginStringMatch: []*matcher.StringMatcher{{
							MatchPattern: &matcher.StringMatcher_Exact{
								Exact: "*",
							},
							IgnoreCase: true,
						}},
						AllowCredentials: &wrappers.BoolValue{Value: true},
						AllowHeaders:     "custom-header-1,custom-header-2",
					},
					&envoy_api_v2_route.Route{
						Match:  routePrefix("/"),
						Action: routecluster("default/svc1/80/da39a3ee5e"),
					}),
			),
		),
		TypeUrl: routeType,
	})

	// Expose headers
	rh.OnAdd(fixture.NewProxy("simple").WithSpec(
		projcontour.HTTPProxySpec{
			VirtualHost: &projcontour.VirtualHost{
				Fqdn: "hello.world",
				CORSPolicy: &projcontour.CORSPolicy{
					AllowOrigin:      []string{"*"},
					AllowCredentials: true,
					ExposeHeaders:    []projcontour.CORSHeaderValue{"custom-header-1", "custom-header-2"},
				},
			}, Routes: []projcontour.Route{{
				Services: []projcontour.Service{{
					Name: "svc1",
					Port: 80,
				}},
			}},
		}),
	)

	c.Request(routeType).Equals(&v2.DiscoveryResponse{
		Resources: resources(t,
			envoyv2.RouteConfiguration("ingress_http",
				envoyv2.CORSVirtualHost("hello.world",
					&envoy_api_v2_route.CorsPolicy{
						AllowOriginStringMatch: []*matcher.StringMatcher{{
							MatchPattern: &matcher.StringMatcher_Exact{
								Exact: "*",
							},
							IgnoreCase: true,
						},
						},
						AllowCredentials: &wrappers.BoolValue{Value: true},
						ExposeHeaders:    "custom-header-1,custom-header-2",
					},
					&envoy_api_v2_route.Route{
						Match:  routePrefix("/"),
						Action: routecluster("default/svc1/80/da39a3ee5e"),
					}),
			),
		),
		TypeUrl: routeType,
	})

	// Max Age
	rh.OnAdd(fixture.NewProxy("simple").WithSpec(
		projcontour.HTTPProxySpec{
			VirtualHost: &projcontour.VirtualHost{
				Fqdn: "hello.world",
				CORSPolicy: &projcontour.CORSPolicy{
					AllowOrigin:      []string{"*"},
					AllowCredentials: true,
					MaxAge:           "10m",
				},
			}, Routes: []projcontour.Route{{
				Services: []projcontour.Service{{
					Name: "svc1",
					Port: 80,
				}},
			}},
		}),
	)

	c.Request(routeType).Equals(&v2.DiscoveryResponse{
		Resources: resources(t,
			envoyv2.RouteConfiguration("ingress_http",
				envoyv2.CORSVirtualHost("hello.world",
					&envoy_api_v2_route.CorsPolicy{
						AllowOriginStringMatch: []*matcher.StringMatcher{{
							MatchPattern: &matcher.StringMatcher_Exact{
								Exact: "*",
							},
							IgnoreCase: true,
						}},
						AllowCredentials: &wrappers.BoolValue{Value: true},
						MaxAge:           "600",
					},
					&envoy_api_v2_route.Route{
						Match:  routePrefix("/"),
						Action: routecluster("default/svc1/80/da39a3ee5e"),
					}),
			),
		),
		TypeUrl: routeType,
	})

	// Disable preflight request caching
	rh.OnAdd(fixture.NewProxy("simple").WithSpec(
		projcontour.HTTPProxySpec{
			VirtualHost: &projcontour.VirtualHost{
				Fqdn: "hello.world",
				CORSPolicy: &projcontour.CORSPolicy{
					AllowOrigin:      []string{"*"},
					AllowCredentials: true,
					MaxAge:           "0s",
				},
			}, Routes: []projcontour.Route{{
				Services: []projcontour.Service{{
					Name: "svc1",
					Port: 80,
				}},
			}},
		}),
	)

	c.Request(routeType).Equals(&v2.DiscoveryResponse{
		Resources: resources(t,
			envoyv2.RouteConfiguration("ingress_http",
				envoyv2.CORSVirtualHost("hello.world",
					&envoy_api_v2_route.CorsPolicy{
						AllowOriginStringMatch: []*matcher.StringMatcher{{
							MatchPattern: &matcher.StringMatcher_Exact{
								Exact: "*",
							},
							IgnoreCase: true,
						}},
						AllowCredentials: &wrappers.BoolValue{Value: true},
						MaxAge:           "0",
					},
					&envoy_api_v2_route.Route{
						Match:  routePrefix("/"),
						Action: routecluster("default/svc1/80/da39a3ee5e"),
					}),
			),
		),
		TypeUrl: routeType,
	})

	// Virtual hosts with an invalid max age in their policy are not added
	invvhost := &projcontour.HTTPProxy{
		ObjectMeta: fixture.ObjectMeta("simple"),
		Spec: projcontour.HTTPProxySpec{
			VirtualHost: &projcontour.VirtualHost{
				Fqdn: "hello.world",
				CORSPolicy: &projcontour.CORSPolicy{
					AllowOrigin:      []string{"*"},
					AllowCredentials: true,
					MaxAge:           "-10m",
				},
			}, Routes: []projcontour.Route{{
				Services: []projcontour.Service{{
					Name: "svc1",
					Port: 80,
				}},
			}},
		},
	}

	rh.OnAdd(invvhost)

	c.Request(routeType).Equals(&v2.DiscoveryResponse{
		Resources: resources(t,
			envoyv2.RouteConfiguration("ingress_http")),
		TypeUrl: routeType,
	}).Status(invvhost).Like(
		projcontour.HTTPProxyStatus{CurrentStatus: k8s.StatusInvalid},
	)

}
