package cloudrun

import (
	"github.com/pulumi/pulumi-gcp/sdk/v4/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v4/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v4/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreateCloudRunService(
	ctx *pulumi.Context,
	projectId string,
	resources []pulumi.Resource,
	serviceapi string,
	computeserviceapi string,
	serviceName string,
	serviceimagePath string,
	urlmapName string,
	iammemberName string) error {

	api, err := projects.NewService(ctx, "enable-"+serviceapi, &projects.ServiceArgs{
		DisableDependentServices: pulumi.Bool(true),
		Project:                  pulumi.String(projectId),
		Service:                  pulumi.String(serviceapi),
  }, pulumi.DependsOn(resources))
	if err != nil {
		return err
	}

	computeapi, err := projects.NewService(ctx, "enable2-"+serviceapi, &projects.ServiceArgs{
		DisableDependentServices: pulumi.Bool(true),
		Project:                  pulumi.String(projectId),
		Service:                  pulumi.String(computeserviceapi),
  }, pulumi.DependsOn(resources))
	if err != nil {
		return err
	}

	cloudrunNegService, err := cloudrun.NewService(ctx, serviceName, &cloudrun.ServiceArgs{
		Project: pulumi.String(projectId),
		Location:  pulumi.String("us-central1"),
		Metadata:  &cloudrun.ServiceMetadataArgs{
			Namespace: pulumi.String(projectId),
		},
		Template: &cloudrun.ServiceTemplateArgs{
			Spec: &cloudrun.ServiceTemplateSpecArgs{
				Containers: cloudrun.ServiceTemplateSpecContainerArray{
					&cloudrun.ServiceTemplateSpecContainerArgs{
						Image: pulumi.String(serviceimagePath),
          },
        },
      },
    },
		Traffics: cloudrun.ServiceTrafficArray{
      &cloudrun.ServiceTrafficArgs{
				Percent:        pulumi.Int(100),
        LatestRevision: pulumi.Bool(true),
      },
    },
  }, pulumi.DependsOn([]pulumi.Resource{api}))
  if err != nil {
    return err
  }

	//regionneg, err := compute.NewRegionNetworkEndpointGroup(ctx, "cloudrunregionneg", &compute.RegionNetworkEndpointGroupArgs{
	_, err = compute.NewRegionNetworkEndpointGroup(ctx, "cloudrunregionneg", &compute.RegionNetworkEndpointGroupArgs{
	//serverlessendpointgroup, err := compute.NewRegionNetworkEndpointGroup(ctx, "cloudrunNegRegionNetworkEndpointGroup", &compute.RegionNetworkEndpointGroupArgs{
		Project:							pulumi.String(projectId),
		NetworkEndpointType:	pulumi.String("SERVERLESS"),
		Region:								pulumi.String("us-central1"),
		CloudRun:	&compute.RegionNetworkEndpointGroupCloudRunArgs{
			Service: cloudrunNegService.Name,
		},
  }, pulumi.DependsOn([]pulumi.Resource{computeapi}))
  if err != nil {
    return err
  }
/*
	proxy, err := compute.NewGlobalNetworkEndpoint(ctx, "proxy", &compute.GlobalNetworkEndpointArgs{
    GlobalNetworkEndpointGroup: serverlessendpointgroup.ID(),
    Fqdn:                       pulumi.String("test.example.com"),
    Port:                       serverlessendpointgroup.DefaultPort,
	}, pulumi.Provider(google_beta))
	if err != nil {
    return err
	}
*/

  /*
	backendservice, err := compute.NewBackendService(ctx, "runbackendservice", &compute.BackendServiceArgs{
		Project:											pulumi.String(projectId),
		EnableCdn:                    pulumi.Bool(true),
    ConnectionDrainingTimeoutSec: pulumi.Int(10),
    Backends: compute.BackendServiceBackendArray{
			&compute.BackendServiceBackendArgs{
				Group: regionneg.ID(),
			},
		},
  }, pulumi.DependsOn([]pulumi.Resource{computeapi}))
	if err != nil {
    return err
  }

	_, err = compute.NewForwardingRule(ctx, "default", &compute.ForwardingRuleArgs{
		Project:				pulumi.String(projectId),
		Region:         pulumi.String("us-central1"),
		PortRange:      pulumi.String("8080"),
		BackendService: backendservice.ID(),
  }, pulumi.DependsOn([]pulumi.Resource{computeapi}))
	if err != nil {
		return err
	}

  //staticURLMap, err := compute.NewURLMap(ctx, urlmapName, &compute.URLMapArgs{
  _, err = compute.NewURLMap(ctx, urlmapName, &compute.URLMapArgs{
		Project:	pulumi.String(projectId),
    Description:    pulumi.String("a description"),
    DefaultService: backendservice.ID(),

    HostRules: compute.URLMapHostRuleArray{
      &compute.URLMapHostRuleArgs{
        Hosts: pulumi.StringArray{
          pulumi.String("*"),
        },
        PathMatcher: pulumi.String("mysite"),
      },
      &compute.URLMapHostRuleArgs{
        Hosts: pulumi.StringArray{
          pulumi.String("myothersite.com"),
        },
        PathMatcher: pulumi.String("otherpaths"),
      },
    },

    PathMatchers: compute.URLMapPathMatcherArray{
      &compute.URLMapPathMatcherArgs{
        Name: pulumi.String("mysite"),
        DefaultService: backendservice.ID(),
        PathRules: compute.URLMapPathMatcherPathRuleArray{
          &compute.URLMapPathMatcherPathRuleArgs{
            Paths: pulumi.StringArray{
              pulumi.String("/*"),
            },
            Service: backendservice.ID(),
          },
        },
      },
      &compute.URLMapPathMatcherArgs{
        Name:           pulumi.String("otherpaths"),
        DefaultService: backendservice.ID(),
      },
    },
  }, pulumi.DependsOn([]pulumi.Resource{api, backendservice}))
  if err != nil {
      return err
  }

	_, err = compute.NewTargetHttpProxy(ctx, "defaultTargetHttpProxy", &compute.TargetHttpProxyArgs{
		UrlMap: defaultURLMap.ID(),
  })
	if err != nil {
		return err
  }
	*/

	/*
  _, err = compute.NewBackendService(ctx, "cloudrun-defaultBackendService", &compute.BackendServiceArgs{
    HealthChecks: pulumi.String(pulumi.String{
	    defaultHttpHealthCheck.ID(),
	  }),
  })
  if err != nil {
    return err
  }
	*/

  _, err = cloudrun.NewIamMember(ctx, iammemberName, &cloudrun.IamMemberArgs{
    //Service: pulumi.String(cloudrunservice.Name),
		Project: pulumi.String(projectId),
    Service: cloudrunNegService.Name,
    Location:  pulumi.String("us-central1"),
    Role: pulumi.String("roles/run.invoker"),
    Member: pulumi.String("allUsers"),
  })

 	return nil
}
