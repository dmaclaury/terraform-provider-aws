// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package elasticache

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	elasticache_sdkv2 "github.com/aws/aws-sdk-go-v2/service/elasticache"
	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	elasticache_sdkv1 "github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newServerlessCacheResource,
			Name:    "Serverless Cache",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceCluster,
			TypeName: "aws_elasticache_cluster",
			Name:     "Cluster",
		},
		{
			Factory:  dataSourceReplicationGroup,
			TypeName: "aws_elasticache_replication_group",
			Name:     "Replication Group",
		},
		{
			Factory:  dataSourceSubnetGroup,
			TypeName: "aws_elasticache_subnet_group",
			Name:     "Subnet Group",
		},
		{
			Factory:  dataSourceUser,
			TypeName: "aws_elasticache_user",
			Name:     "User",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceCluster,
			TypeName: "aws_elasticache_cluster",
			Name:     "Cluster",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceGlobalReplicationGroup,
			TypeName: "aws_elasticache_global_replication_group",
			Name:     "Global Replication Group",
		},
		{
			Factory:  resourceParameterGroup,
			TypeName: "aws_elasticache_parameter_group",
			Name:     "Parameter Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceReplicationGroup,
			TypeName: "aws_elasticache_replication_group",
			Name:     "Replication Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceSubnetGroup,
			TypeName: "aws_elasticache_subnet_group",
			Name:     "Subnet Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceUser,
			TypeName: "aws_elasticache_user",
			Name:     "User",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceUserGroup,
			TypeName: "aws_elasticache_user_group",
			Name:     "User Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceUserGroupAssociation,
			TypeName: "aws_elasticache_user_group_association",
			Name:     "User Group Association",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.ElastiCache
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*elasticache_sdkv1.ElastiCache, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return elasticache_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config[names.AttrEndpoint].(string))})), nil
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*elasticache_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return elasticache_sdkv2.NewFromConfig(cfg, func(o *elasticache_sdkv2.Options) {
		if endpoint := config[names.AttrEndpoint].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
