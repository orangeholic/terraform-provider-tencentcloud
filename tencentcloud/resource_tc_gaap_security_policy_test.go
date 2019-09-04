package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudGaapSecurityPolicy_basic(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapSecurityPolicyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityPolicyExists("tencentcloud_gaap_security_policy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_security_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityPolicy_disable(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapSecurityPolicyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityPolicyExists("tencentcloud_gaap_security_policy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "true"),
				),
			},
			{
				Config: testAccGaapSecurityPolicyDisable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityPolicyExists("tencentcloud_gaap_security_policy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "false"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityPolicy_drop(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapSecurityPolicyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityPolicyDrop,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityPolicyExists("tencentcloud_gaap_security_policy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "action", "DROP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "true"),
				),
			},
		},
	})
}

func testAccCheckGaapSecurityPolicyExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no listener ID is set")
		}

		service := GaapService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, _, _, exist, err := service.DescribeSecurityPolicy(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if !exist {
			return fmt.Errorf("security policy not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckGaapSecurityPolicyDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := GaapService{client: client}

		_, _, _, exist, err := service.DescribeSecurityPolicy(context.TODO(), *id)
		if err != nil {
			return err
		}

		if exist {
			return errors.New("security policy still exists")
		}

		return nil
	}
}

const testAccGaapSecurityPolicyBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  action   = "ACCEPT"
}
`

const testAccGaapSecurityPolicyDisable = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  action   = "ACCEPT"
  enable   = false
}
`

const testAccGaapSecurityPolicyDrop = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  action   = "DROP"
}
`