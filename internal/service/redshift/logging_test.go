// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redshift_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tfredshift "github.com/hashicorp/terraform-provider-aws/internal/service/redshift"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccRedshiftLogging_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var log redshift.DescribeLoggingStatusOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_logging.test"
	clusterResourceName := "aws_redshift_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.RedshiftEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.RedshiftServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoggingDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingConfig_basic(rName, string(types.LogDestinationTypeCloudwatch)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingExists(ctx, resourceName, &log),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_identifier", clusterResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "log_destination_type", string(types.LogDestinationTypeCloudwatch)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedshiftLogging_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var log redshift.DescribeLoggingStatusOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_logging.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.RedshiftEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.RedshiftServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoggingDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingConfig_basic(rName, string(types.LogDestinationTypeCloudwatch)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingExists(ctx, resourceName, &log),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfredshift.ResourceLogging, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRedshiftLogging_disappears_Cluster(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var log redshift.DescribeLoggingStatusOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_logging.test"
	clusterResourceName := "aws_redshift_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.RedshiftEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.RedshiftServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoggingDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingConfig_basic(rName, string(types.LogDestinationTypeCloudwatch)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingExists(ctx, resourceName, &log),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfredshift.ResourceCluster(), clusterResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckLoggingDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).RedshiftClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_redshift_logging" {
				continue
			}

			_, err := tfredshift.FindLoggingByID(ctx, conn, rs.Primary.ID)
			if errs.IsA[*retry.NotFoundError](err) {
				return nil
			}
			if err != nil {
				return create.Error(names.Redshift, create.ErrActionCheckingDestroyed, tfredshift.ResNameLogging, rs.Primary.ID, err)
			}

			return create.Error(names.Redshift, create.ErrActionCheckingDestroyed, tfredshift.ResNameLogging, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckLoggingExists(ctx context.Context, name string, log *redshift.DescribeLoggingStatusOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.Redshift, create.ErrActionCheckingExistence, tfredshift.ResNameLogging, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.Redshift, create.ErrActionCheckingExistence, tfredshift.ResNameLogging, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RedshiftClient(ctx)
		out, err := tfredshift.FindLoggingByID(ctx, conn, rs.Primary.ID)
		if err != nil {
			return create.Error(names.Redshift, create.ErrActionCheckingExistence, tfredshift.ResNameLogging, rs.Primary.ID, err)
		}

		*log = *out

		return nil
	}
}

func testAccLoggingConfigBase(rName string) string {
	return acctest.ConfigCompose(
		// "InvalidVPCNetworkStateFault: The requested AZ us-west-2a is not a valid AZ."
		acctest.ConfigAvailableAZsNoOptInExclude("usw2-az2"),
		fmt.Sprintf(`
resource "aws_redshift_cluster" "test" {
  cluster_identifier                  = %[1]q
  availability_zone                   = data.aws_availability_zones.available.names[0]
  database_name                       = "mydb"
  master_username                     = "foo_test"
  master_password                     = "Mustbe8characters"
  multi_az                            = false
  node_type                           = "dc2.large"
  automated_snapshot_retention_period = 0
  allow_version_upgrade               = false
  skip_final_snapshot                 = true

  lifecycle {
    ignore_changes = [logging]
  }
}
`, rName))
}

func testAccLoggingConfig_basic(rName, logDestinationType string) string {
	return acctest.ConfigCompose(
		testAccLoggingConfigBase(rName),
		fmt.Sprintf(`
resource "aws_redshift_logging" "test" {
  cluster_identifier   = aws_redshift_cluster.test.id
  log_destination_type = %[1]q
}
`, logDestinationType))
}
