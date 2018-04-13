package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleVpnGateway_basic(t *testing.T) {
	var v, v2 fcu.VpnGateway

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_vpn_gateway.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists(
						"outscale_vpn_gateway.foo", &v),
				),
			},

			resource.TestStep{
				Config: testAccVpnGatewayConfigChangeVPC,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists(
						"outscale_vpn_gateway.foo", &v2),
				),
			},
		},
	})
}
func TestAccOutscaleVpnGateway_delete(t *testing.T) {
	var vpnGateway fcu.VpnGateway

	testDeleted := func(r string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			_, ok := s.RootModule().Resources[r]
			if ok {
				return fmt.Errorf("VPN Gateway %q should have been deleted.", r)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_vpn_gateway.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("outscale_vpn_gateway.foo", &vpnGateway)),
			},
			resource.TestStep{
				Config: testAccNoVpnGatewayConfig,
				Check:  resource.ComposeTestCheckFunc(testDeleted("outscale_vpn_gateway.foo")),
			},
		},
	})
}

func TestAccOutscaleVpnGateway_tags(t *testing.T) {
	var v fcu.VpnGateway

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_vpn_gateway.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckVpnGatewayConfigTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("outscale_vpn_gateway.foo", &v),
					testAccCheckTags(&v.Tags, "foo", "bar"),
				),
			},
		},
	})
}

func testAccOutscaleVpnGatewayDisappears(gateway *fcu.VpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		var err error

		opts := &fcu.DeleteVpnGatewayInput{
			VpnGatewayId: gateway.VpnGatewayId,
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.VM.DeleteVpnGateway(opts)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err != nil {
			return err
		}

		return resource.Retry(40*time.Minute, func() *resource.RetryError {
			opts := &fcu.DescribeVpnGatewaysInput{
				VpnGatewayIds: []*string{gateway.VpnGatewayId},
			}

			var resp *fcu.DescribeVpnGatewaysOutput
			var err error

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				resp, err = conn.VM.DescribeVpnGateways(opts)
				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return resource.NonRetryableError(err)
			})
			if err != nil {
				cgw, ok := err.(awserr.Error)
				if ok && cgw.Code() == "InvalidVpnGatewayID.NotFound" {
					return nil
				}
				if ok && cgw.Code() == "IncorrectState" {
					return resource.RetryableError(fmt.Errorf(
						"Waiting for VPN Gateway to be in the correct state: %v", gateway.VpnGatewayId))
				}
				return resource.NonRetryableError(
					fmt.Errorf("Error retrieving VPN Gateway: %s", err))
			}
			if *resp.VpnGateways[0].State == "deleted" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf(
				"Waiting for VPN Gateway: %v", gateway.VpnGatewayId))
		})
	}
}

func testAccCheckVpnGatewayDestroy(s *terraform.State) error {
	FCU := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vpn_gateway" {
			continue
		}

		var resp *fcu.DescribeVpnGatewaysOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = FCU.VM.DescribeVpnGateways(&fcu.DescribeVpnGatewaysInput{
				VpnGatewayIds: []*string{aws.String(rs.Primary.ID)},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})
		if err == nil {
			var v *fcu.VpnGateway
			for _, g := range resp.VpnGateways {
				if *g.VpnGatewayId == rs.Primary.ID {
					v = g
				}
			}

			if v == nil {
				// wasn't found
				return nil
			}

			if *v.State != "deleted" {
				return fmt.Errorf("Expected VPN Gateway to be in deleted state, but was not: %s", v)
			}
			return nil
		}

		// Verify the error is what we want
		ec2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if ec2err.Code() != "InvalidVpnGatewayID.NotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckVpnGatewayExists(n string, ig *fcu.VpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		FCU := testAccProvider.Meta().(*OutscaleClient).FCU

		var resp *fcu.DescribeVpnGatewaysOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = FCU.VM.DescribeVpnGateways(&fcu.DescribeVpnGatewaysInput{
				VpnGatewayIds: []*string{aws.String(rs.Primary.ID)},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})
		if err != nil {
			return err
		}
		if len(resp.VpnGateways) == 0 {
			return fmt.Errorf("VPN Gateway not found")
		}

		*ig = *resp.VpnGateways[0]

		return nil
	}
}

const testAccNoVpnGatewayConfig = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}
`

const testAccVpnGatewayConfig = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_vpn_gateway" "foo" {
}
`

const testAccVpnGatewayConfigChangeVPC = `
resource "outscale_lin" "bar" {
	cidr_block = "10.2.0.0/16"
}

resource "outscale_vpn_gateway" "foo" {
}
`

const testAccCheckVpnGatewayConfigTags = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_vpn_gateway" "foo" {
	tag {
		foo = "bar"
	}
}
`

const testAccCheckVpnGatewayConfigTagsUpdate = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_vpn_gateway" "foo" {
	tag {
		bar = "baz"
	}
}
`

const testAccCheckVpnGatewayConfigReattach = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin" "bar" {
	cidr_block = "10.2.0.0/16"
}

resource "outscale_vpn_gateway" "foo" {
}

resource "outscale_vpn_gateway" "bar" {
}
`

const testAccCheckVpnGatewayConfigReattachChange = `
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin" "bar" {
	cidr_block = "10.2.0.0/16"
}

resource "outscale_vpn_gateway" "foo" {
}

resource "outscale_vpn_gateway" "bar" {
}
`
